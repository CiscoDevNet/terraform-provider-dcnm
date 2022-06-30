package dcnm

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var switchDeployMutexMap = make(map[string]*sync.Mutex, 0)

func resourceDCNMInventroy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMInventroyCreate,
		Update: resourceDCNMInventroyUpdate,
		Read:   resourceDCNMInventroyRead,
		Delete: resourceDCNMInventroyDelete,

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auth_protocol": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"max_hops": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"second_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"preserve_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},

			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"switch_config": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"switch_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"role": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"leaf",
								"spine",
								"border",
								"border_spine",
								"border_gateway",
								"border_gateway_spine",
								"super_spine",
								"border_super_spine",
								"border_gateway_super_spine",
							}, false),
						},

						"switch_db_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"serial_number": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"model": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: resourceDCNMSwitchConfigHash,
			},

			"config_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},

			"deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ips": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func roleMappingFunc(role string) string {
	roleMapping := map[string]string{
		"leaf":                       "leaf",
		"spine":                      "spine",
		"border":                     "border",
		"border_spine":               "border spine",
		"border_gateway":             "border gateway",
		"border_gateway_spine":       "border gateway spine",
		"super_spine":                "super spine",
		"border_super_spine":         "border super spine",
		"border_gateway_super_spine": "border gateway super spine",
	}

	return roleMapping[role]
}

func extractFabricID(dcnmClient *client.Client, fabricName string) (int, error) {
	durl := fmt.Sprintf("/rest/control/fabrics/%s", fabricName)

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(stripQuotes(cont.S("id").String()))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func extractSwitchinfo(contList *container.Container) models.Switch {
	s := models.Switch{}

	cont := contList.Index(0)

	s.Reachable = stripQuotes(cont.S("reachable").String())
	s.Auth = stripQuotes(cont.S("auth").String())
	s.Known = stripQuotes(cont.S("known").String())
	s.Valid = stripQuotes(cont.S("valid").String())
	s.Selectable = stripQuotes(cont.S("selectable").String())
	s.SysName = stripQuotes(cont.S("sysName").String())
	s.IP = stripQuotes(cont.S("ipaddr").String())
	s.Platform = stripQuotes(cont.S("platform").String())
	s.Version = stripQuotes(cont.S("version").String())
	s.LastChange = stripQuotes(cont.S("lastChange").String())
	s.Hops, _ = strconv.Atoi(stripQuotes(cont.S("hopCount").String()))
	s.DeviceIndex = stripQuotes(cont.S("deviceIndex").String())
	s.StatReason = stripQuotes(cont.S("statusReason").String())

	return s
}

func extractSerialNumber(cont *container.Container, ip string) (string, error) {
	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		infoCont := cont.Index(i)

		ipGet := stripQuotes(infoCont.S("ipAddress").String())
		if ipGet == ip {
			return stripQuotes(infoCont.S("serialNumber").String()), nil
		}
	}

	return "", fmt.Errorf("No inventory found for given ip address")
}

func getRemoteSwitch(dcnmClient *client.Client, fabric, ip, serialNum string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabric)

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		switchCont := cont.Index(i)
		if ip != "" {
			ipGet := stripQuotes(switchCont.S("ipAddress").String())
			if ipGet == ip {
				return switchCont, nil
			}
		} else {
			serialGet := stripQuotes(switchCont.S("serialNumber").String())
			if serialGet == serialNum {
				return switchCont, nil
			}
		}
	}
	return nil, fmt.Errorf("Desired switch not found")
}

func getSwitchInfo(cont *container.Container) map[string]interface{} {

	sInfo := make(map[string]interface{})
	sInfo["ip"] = stripQuotes(cont.S("ipAddress").String())
	sInfo["switch_name"] = stripQuotes(cont.S("logicalName").String())
	sInfo["switch_db_id"] = stripQuotes(cont.S("switchDbID").String())
	sInfo["serial_number"] = stripQuotes(cont.S("serialNumber").String())
	sInfo["model"] = stripQuotes(cont.S("model").String())
	sInfo["mode"] = stripQuotes(cont.S("mode").String())

	return sInfo
}

func resourceDCNMInventroyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	delSwtiches := make([]string, 0, 1)

	ips := make([]string, 0, 1)
	discoveredIps := make([]string, 0, 1)

	inv := models.Inventory{}
	switchObjs := make([]*models.Switch, 0, 1)
	switchInfos := d.Get("switch_config").(*schema.Set).List()

	inv.Username = d.Get("username").(string)
	inv.Password = d.Get("password").(string)

	if auth, ok := d.GetOk("auth_protocol"); ok {
		inv.V3auth = auth.(int)
	}
	if hops, ok := d.GetOk("max_hops"); ok {
		inv.MaxHops = hops.(int)
	}
	if secTime, ok := d.GetOk("second_timeout"); ok {
		inv.SecondTimeout = secTime.(int)
	}
	if preConfig, ok := d.GetOk("preserve_config"); ok {
		inv.PreserveConfig = preConfig.(string)
	}
	if platform, ok := d.GetOk("platform"); ok {
		inv.Platform = platform.(string)
	}

	for _, val := range switchInfos {
		sInfo := val.(map[string]interface{})
		ip := sInfo["ip"].(string)

		inv.SeedIP = ip

		fabricID, err := extractFabricID(dcnmClient, fabricName)
		if err != nil {
			return err
		}

		dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/test-reachability", strconv.Itoa(fabricID))
		cont, err := dcnmClient.Save(dUrl, &inv)
		if err != nil {
			log.Println("error at test reachability for switch : ", ip, err)
			continue
		}

		switchM := extractSwitchinfo(cont)

		if switchM.Selectable != "true" || switchM.Reachable != "true" {
			log.Println("Desired switch is not reachable or not selectable or invalid user/password or bad authentication protocol", ip)
			continue
		}

		switchObjs = append(switchObjs, &switchM)
		discoveredIps = append(discoveredIps, ip)
	}

	inv.SeedIP = strings.Join(discoveredIps, ",")

	invModel := models.NewSwitch(&inv, switchObjs)

	dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/discover", fabricName)
	_, err := dcnmClient.Save(dUrl, invModel)
	if err != nil {
		return fmt.Errorf("error at discovery for switches %s", err)
	}

	var delFlag bool
	deployedIps := make([]string, 0, 1)
	deployedSerial := make([]string, 0, 1)
	for _, ip := range discoveredIps {
		var serialNum string
		configTimeout := (d.Get("config_timeout").(int)) * 60
		migrate := true

		for configTimeout > 0 {
			cont, err := getRemoteSwitch(dcnmClient, fabricName, ip, "")
			if err != nil {
				log.Println("error at get call for switch in creation :", ip, err)
				continue
			}
			serialNum = stripQuotes(cont.S("serialNumber").String())

			if stripQuotes(cont.S("mode").String()) != "Migration" {
				time.Sleep(10 * time.Second)
				configTimeout = configTimeout - 10
				migrate = false
				break
			}
			time.Sleep(5 * time.Second)
			configTimeout = configTimeout - 5
		}
		if migrate {
			log.Println("Timeout occurs before going into normal mode. Hence removing it!", ip)
			delSwtiches = append(delSwtiches, serialNum)
			delFlag = true
			continue
		}

		err = deployswitch(dcnmClient, fabricName, serialNum, configTimeout)
		if err != nil {
			delSwtiches = append(delSwtiches, serialNum)
			log.Println("error at switch deployment :", ip, err)
			delFlag = true
			continue
		}

		deployedIps = append(deployedIps, ip)
		deployedSerial = append(deployedSerial, serialNum)
	}

	if delFlag {
		for _, serial := range delSwtiches {
			_, err := getRemoteSwitch(dcnmClient, fabricName, "", serial)
			if err == nil {
				durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serial)
				_, delerr := dcnmClient.Delete(durl)
				if delerr != nil {
					log.Println("error at deletion of switch ", err)
				}
			}
		}
		log.Println("Some switches failed to discover and deploy, resuming procedure for successfully discovered switches")
	}

	for _, ip := range deployedIps {
		for _, val := range switchInfos {
			sInfo := val.(map[string]interface{})

			if sInfo["ip"].(string) == ip {
				if sInfo["role"] != "" {
					cont, err := getRemoteSwitch(dcnmClient, fabricName, ip, "")
					if err != nil {
						log.Println("error at get call for switch in creation :", ip, err)
						continue
					}
					serialNum := stripQuotes(cont.S("serialNumber").String())

					durl := fmt.Sprintf("/rest/control/switches/roles")
					sRole := models.SwitchRole{}
					sRole.Role = roleMappingFunc(sInfo["role"].(string))
					sRole.SerialNumber = serialNum

					_, err = dcnmClient.SaveForAttachment(durl, &sRole)
					if err != nil {
						log.Println("error at switch role assignment :", ip, err)
					}
				}
			}
		}
	}

	delFlag = false
	err = deployFabric(dcnmClient, fabricName)
	if err != nil {
		log.Println("error at fabric deployment :", err)
		delFlag = true
	}

	if delFlag {
		for _, serial := range deployedSerial {
			_, err := getRemoteSwitch(dcnmClient, fabricName, "", serial)
			if err == nil {
				durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serial)
				_, delerr := dcnmClient.Delete(durl)
				if delerr != nil {
					log.Println("error at deletion of switch ", err)
				}
			}
		}
	} else {
		ips = append(ips, deployedIps...)
	}

	d.Set("ips", ips)

	if len(ips) == 0 {
		return fmt.Errorf("none of the switches are discovered and deployed on the fabric, some internal issue in switches")
	}

	d.SetId(strings.Join(ips, ","))

	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMInventroyRead(d, m)
}

func resourceDCNMInventroyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	if d.HasChange("deploy") && d.Get("deploy").(bool) == false {
		d.Set("deploy", true)
		return fmt.Errorf("Deployed switches can not be undeployed")
	}

	delSwtiches := make([]string, 0, 1)
	var delFlag bool

	ipDns := d.Get("ips").([]interface{})

	ips := make([]string, 0, 1)
	switchInfosOld, switchInfosNew := d.GetChange("switch_config")

	deleteSwitches := getSerialsForDelete(switchInfosOld.(*schema.Set).List(), switchInfosNew.(*schema.Set).List())
	err := deleteSpecificSwitches(dcnmClient, fabricName, deleteSwitches)
	if err != nil {
		return err
	}

	inv := models.Inventory{}
	inv.Username = d.Get("username").(string)
	inv.Password = d.Get("password").(string)

	if auth, ok := d.GetOk("auth_protocol"); ok {
		inv.V3auth = auth.(int)
	}
	if hops, ok := d.GetOk("max_hops"); ok {
		inv.MaxHops = hops.(int)
	}
	if secTime, ok := d.GetOk("second_timeout"); ok {
		inv.SecondTimeout = secTime.(int)
	}
	if preConfig, ok := d.GetOk("preserve_config"); ok {
		inv.PreserveConfig = preConfig.(string)
	}
	if platform, ok := d.GetOk("platform"); ok {
		inv.Platform = platform.(string)
	}

	newSwitchFlag := false
	discoveredIps := make([]string, 0, 1)
	switchObjs := make([]*models.Switch, 0, 1)
	switchInfos := switchInfosNew.(*schema.Set).List()
	for _, val := range switchInfos {

		sInfo := val.(map[string]interface{})
		ip := sInfo["ip"].(string)

		if contains(ipDns, ip) {
			auth := d.Get("auth_protocol").(int)

			cont, err := getRemoteSwitch(dcnmClient, fabricName, ip, "")
			if err != nil {
				return err
			}

			switchDbID := stripQuotes(cont.S("switchDbID").String())

			body := []byte(fmt.Sprintf("switchIds=%s&userName=%s&password=%s&v3protocol=%s", switchDbID, inv.Username, inv.Password, strconv.Itoa(auth)))

			durl := fmt.Sprintf("/fm/fmrest/lanConfig/saveSwitchCredentials")
			if dcnmClient.GetPlatform() == "nd" {
				durl = fmt.Sprintf("/rest/lanConfig/saveSwitchCredentials")
			}
			cont, err = dcnmClient.UpdateCred(durl, body)
			if err != nil {
				log.Println("error at updation of switch :", ip, err)
			}

			ips = append(ips, ip)

		} else {
			newSwitchFlag = true

			inv.SeedIP = ip

			fabricID, err := extractFabricID(dcnmClient, fabricName)
			if err != nil {
				return err
			}

			dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/test-reachability", strconv.Itoa(fabricID))
			cont, err := dcnmClient.Save(dUrl, &inv)
			if err != nil {
				log.Println("error at test reachability for switch :", ip, err)
				continue
			}

			switchM := extractSwitchinfo(cont)

			if switchM.Selectable != "true" || switchM.Reachable != "true" {
				log.Println("Desired switch: is not reachable or not selectable or invalid user/password or bad authentication protocol", ip)
				continue
			}

			discoveredIps = append(discoveredIps, ip)
			switchObjs = append(switchObjs, &switchM)
		}
	}

	if newSwitchFlag {
		inv.SeedIP = strings.Join(discoveredIps, ",")
		invModel := models.NewSwitch(&inv, switchObjs)

		dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/discover", fabricName)
		_, err = dcnmClient.Save(dUrl, invModel)
		if err != nil {
			return fmt.Errorf("error at discovery for switch : %s", err)
		}

		deployedIps := make([]string, 0, 1)
		deployedSerial := make([]string, 0, 1)
		for _, ip := range discoveredIps {
			var serialNum string
			configTimeout := (d.Get("config_timeout").(int)) * 60
			migrate := true

			for configTimeout > 0 {
				cont, err := getRemoteSwitch(dcnmClient, fabricName, ip, "")
				if err != nil {
					log.Println("error at get call for switch in updation :", ip, err)
					continue
				}
				serialNum = stripQuotes(cont.S("serialNumber").String())

				if stripQuotes(cont.S("mode").String()) != "Migration" {
					time.Sleep(10 * time.Second)
					configTimeout = configTimeout - 10
					migrate = false
					break
				}
				time.Sleep(5 * time.Second)
				configTimeout = configTimeout - 5
			}
			if migrate {
				log.Println("Timeout occurs before going into normal mode. Hence removing it!", ip)
				delSwtiches = append(delSwtiches, serialNum)
				delFlag = true
				continue
			}

			err := deployswitch(dcnmClient, fabricName, serialNum, configTimeout)
			if err != nil {
				delSwtiches = append(delSwtiches, serialNum)
				log.Println("error at switch deployment :", ip, err)
				delFlag = true
				continue
			}

			deployedIps = append(deployedIps, ip)
			deployedSerial = append(deployedSerial, serialNum)
		}

		if delFlag {
			for _, serial := range delSwtiches {
				_, err := getRemoteSwitch(dcnmClient, fabricName, "", serial)
				if err == nil {
					durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serial)
					_, delerr := dcnmClient.Delete(durl)
					if delerr != nil {
						log.Println("error at deletion of switch ", err)
					}
				}
			}
			log.Println("Some switches failed to discover and deploy, resuming procedure for successfully discovered switches")
		}

		delFlag = false
		err = deployFabric(dcnmClient, fabricName)
		if err != nil {
			log.Println("error at fabric deployment :", err)
			delFlag = true
		}

		if delFlag {
			for _, serial := range deployedSerial {
				_, err := getRemoteSwitch(dcnmClient, fabricName, "", serial)
				if err == nil {
					durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serial)
					_, delerr := dcnmClient.Delete(durl)
					if delerr != nil {
						log.Println("error at deletion of switch ", err)
					}
				}
			}
		} else {
			ips = append(ips, deployedIps...)
		}

	}

	for _, ip := range ips {
		for _, val := range switchInfos {
			sInfo := val.(map[string]interface{})

			if sInfo["ip"].(string) == ip {
				if sInfo["role"] != "" {
					cont, err := getRemoteSwitch(dcnmClient, fabricName, ip, "")
					if err != nil {
						log.Println("error at get call for switch in creation :", ip, err)
						continue
					}
					serialNum := stripQuotes(cont.S("serialNumber").String())

					durl := fmt.Sprintf("/rest/control/switches/roles")
					sRole := models.SwitchRole{}
					sRole.Role = roleMappingFunc(sInfo["role"].(string))
					sRole.SerialNumber = serialNum

					_, err = dcnmClient.SaveForAttachment(durl, &sRole)
					if err != nil {
						log.Println("error at switch role assignment :", ip, err)
					}
				}
			}
		}
	}

	err = deployFabric(dcnmClient, fabricName)
	if err != nil {
		log.Println("error at fabric deployment after role assignment: ", err)
	}

	d.Set("ips", ips)
	d.SetId(strings.Join(ips, ","))

	log.Println("[DEBUG] End of Update method ", d.Id())
	return resourceDCNMInventroyRead(d, m)
}

func resourceDCNMInventroyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	switchIps := d.Get("ips").([]interface{})

	switchConfigGet := make([]interface{}, 0, 1)
	switchSerial := make([]string, 0, 1)
	ips := make([]string, 0, 1)

	for _, ip := range switchIps {
		cont, err := getRemoteSwitch(dcnmClient, fabricName, ip.(string), "")
		if err == nil {
			switchMap := getSwitchInfo(cont)

			ips = append(ips, ip.(string))

			role, err := getSwitchRole(dcnmClient, switchMap["serial_number"].(string))
			if err == nil {
				switchMap["role"] = strings.ReplaceAll(strings.Trim(role, " "), " ", "_")
			} else {
				log.Println("error in read at fetching switch role :", ip, err)
			}

			switchSerial = append(switchSerial, switchMap["serial_number"].(string))

			switchConfigGet = append(switchConfigGet, switchMap)

		}
	}
	d.Set("switch_config", switchConfigGet)

	deployFlag := true
	for _, serial := range switchSerial {
		flag, err := checkDeploy(dcnmClient, fabricName, serial)
		if err == nil {
			if !flag {
				deployFlag = false
			}
		}
	}

	if deployFlag {
		d.Set("deploy", true)
	} else {
		d.Set("deploy", false)
	}

	d.SetId(strings.Join(ips, ","))

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMInventroyDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	dn := strings.Split(d.Id(), ",")

	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabricName)
	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return err
	}

	delErr := false
	deletedIps := make([]string, 0, 1)

	for _, ip := range dn {
		serialNumber, err := extractSerialNumber(cont, strings.Trim(ip, " "))
		if err != nil {
			return err
		}

		durl = fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serialNumber)
		_, err = dcnmClient.Delete(durl)
		if err != nil {
			log.Println("error at deletion of switch :", ip, err)
			delErr = true
		} else {
			deletedIps = append(deletedIps, ip)
		}
	}

	if delErr {
		leftIps := make([]string, 0, 1)
		diffIps := difference(dn, deletedIps)
		for _, ip := range diffIps {
			leftIps = append(leftIps, ip.(string))
		}
		d.SetId(strings.Join(leftIps, ","))

		return fmt.Errorf("All switches are not deleted properly")
	}
	d.SetId("")

	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}

func checkDeploy(client *client.Client, fabric, serialNum string) (bool, error) {
	durl := fmt.Sprintf("rest/control/fabrics/%s/config-preview/%s", fabric, serialNum)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return false, err
	}

	totalSwitch := len(cont.Data().([]interface{}))
	for i := 0; i < totalSwitch; i++ {
		switchCont := cont.Index(i)
		if stripQuotes(switchCont.S("switchId").String()) == serialNum {
			if stripQuotes(switchCont.S("status").String()) == "In-Sync" {
				return true, nil
			}
		}
	}

	return false, nil
}

func deployswitch(client *client.Client, fabric, serialNum string, configTime int) error {
	log.Println("[DEBUG] Begining Deployment of switch ", serialNum)

	// Step 1 switch configuration
	configDone := false
	timeLeft := configTime
	durl := fmt.Sprintf("rest/control/fabrics/%s/config-preview", fabric)
	for timeLeft > 0 {
		cont, err := client.GetviaURL(durl)
		if err != nil {
			return err
		}

		var flag bool
		totalSwitch := len(cont.Data().([]interface{}))
		for i := 0; i < totalSwitch; i++ {
			switchCont := cont.Index(i)
			if stripQuotes(switchCont.S("switchId").String()) == serialNum {
				if stripQuotes(switchCont.S("status").String()) == "Out-of-Sync" {
					flag = true
					configDone = true
				} else if stripQuotes(switchCont.S("status").String()) == "In-Sync" {
					return nil
				}
				break
			}
		}

		if flag {
			break
		}

		timeLeft = timeLeft / 2
		time.Sleep(time.Duration(timeLeft) * time.Second)
	}
	if !configDone {
		return fmt.Errorf("timeout occurs before completion of switch configuration")
	}

	//Step 2 deploy switch into fabric
	durl = fmt.Sprintf("rest/control/fabrics/%s/config-deploy/%s", fabric, serialNum)
	_, err := client.SaveAndDeploy(durl)
	if err != nil {
		return err
	}

	return nil
}

func deployFabric(client *client.Client, fabric string) error {
	//Step 3 deploy fabric
	durl := fmt.Sprintf("rest/control/fabrics/%s/config-deploy", fabric)
	_, err := client.SaveAndDeploy(durl)
	if err != nil {
		return err
	}

	//Step 4 Save configuration
	durl = fmt.Sprintf("rest/control/fabrics/%s/config-save", fabric)
	_, err = client.SaveAndDeploy(durl)
	if err != nil {
		return err
	}

	//Step 5 deploy fabric
	durl = fmt.Sprintf("rest/control/fabrics/%s/config-deploy", fabric)
	_, err = client.SaveAndDeploy(durl)
	if err != nil {
		return err
	}

	return nil
}

func getSwitchRole(client *client.Client, serial string) (string, error) {
	durl := fmt.Sprintf("/rest/control/switches/roles?serialNumber=%s", serial)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return "", err
	}

	return stripQuotes(cont.Index(0).S("role").String()), nil
}

func resourceDCNMSwitchConfigHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["ip"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["role"].(string)))

	return hashString(buf.String())
}

func getSerialsForDelete(old []interface{}, new []interface{}) []string {
	oldIps := make([]string, 0, 1)
	for _, val := range old {
		info := val.(map[string]interface{})

		oldIps = append(oldIps, info["ip"].(string))
	}

	newIps := make([]string, 0, 1)
	for _, val := range new {
		info := val.(map[string]interface{})

		newIps = append(newIps, info["ip"].(string))
	}

	diff := setDifference(oldIps, newIps)

	return diff
}

func deleteSpecificSwitches(client *client.Client, fabricName string, ips []string) error {
	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabricName)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		serialNumber, err := extractSerialNumber(cont, strings.Trim(ip, " "))
		if err != nil {
			return err
		}

		durl = fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serialNumber)
		_, err = client.Delete(durl)
		if err != nil {
			return err
		}
	}
	return nil
}
