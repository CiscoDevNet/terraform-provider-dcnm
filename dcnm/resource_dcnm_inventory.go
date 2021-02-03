package dcnm

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDCNMInventroy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMInventroyCreate,
		Update: resourceDCNMInventroyUpdate,
		Read:   resourceDCNMInventroyRead,
		Delete: resourceDCNMInventroyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDCNMInventoryImporter,
		},

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"password": &schema.Schema{
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

			"auth_protocol": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"max_hops": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"second_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"preserve_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
				ForceNew: true,
			},

			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"config_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
		},
	}
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

func getRemoteSwitch(dcnmClient *client.Client, fabric, ip string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabric)

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		switchCont := cont.Index(i)

		ipGet := stripQuotes(switchCont.S("ipAddress").String())
		if ipGet == ip {
			return switchCont, nil
		}
	}
	return nil, fmt.Errorf("Desired switch not found")
}

func setSwitchAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {

	d.Set("ip", stripQuotes(cont.S("ipAddress").String()))
	d.Set("fabric_name", stripQuotes(cont.S("fabricName").String()))
	d.Set("switch_name", stripQuotes(cont.S("logicalName").String()))
	d.Set("switch_db_id", stripQuotes(cont.S("switchDbID").String()))
	d.Set("serial_number", stripQuotes(cont.S("serialNumber").String()))
	d.Set("model", stripQuotes(cont.S("model").String()))
	d.Set("mode", stripQuotes(cont.S("mode").String()))

	d.SetId(stripQuotes(cont.S("ipAddress").String()))

	return d
}

func resourceDCNMInventoryImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Importer ", d.Id())

	dcnmClient := m.(*client.Client)

	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 2 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	fabricName := importInfo[0]
	name := importInfo[1]

	cont, err := getRemoteSwitchforDS(dcnmClient, fabricName, name)
	if err != nil {
		return nil, err
	}

	importState := setSwitchAttributes(d, cont)

	d.Set("preserve_config", "false")

	role, err := getSwitchRole(dcnmClient, d.Get("serial_number").(string))
	if err != nil {
		return nil, err
	}
	d.Set("role", role)

	log.Println("[DEBUG] End of Importer ", d.Id())
	return []*schema.ResourceData{importState}, nil
}

func resourceDCNMInventroyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)
	ip := d.Get("ip").(string)
	user := d.Get("username").(string)
	pass := d.Get("password").(string)

	inv := models.Inventory{}
	inv.SeedIP = ip
	inv.Username = user
	inv.Password = pass

	if auth, ok := d.GetOk("auth_protocol"); ok {
		inv.V3auth = auth.(int)
	}

	if maxHop, ok := d.GetOk("max_hops"); ok {
		inv.MaxHops = maxHop.(int)
	}

	if secTime, ok := d.GetOk("second_timeout"); ok {
		inv.SecondTimeout = secTime.(int)
	}

	if preConf, ok := d.GetOk("preserve_config"); ok {
		inv.PreserveConfig = preConf.(string)
	}

	if platform, ok := d.GetOk("platform"); ok {
		inv.Platform = platform.(string)
	}

	fabricID, err := extractFabricID(dcnmClient, fabricName)
	if err != nil {
		return err
	}

	dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/test-reachability", strconv.Itoa(fabricID))
	cont, err := dcnmClient.Save(dUrl, &inv)
	if err != nil {
		return err
	}

	switchM := extractSwitchinfo(cont)

	if switchM.Selectable != "true" || switchM.Reachable != "true" {
		return fmt.Errorf("Desired switch is not reachable or not selectable or invalid user/password or bad authentication protocol")
	}

	invModel := models.NewSwitch(&inv, &switchM)

	dUrl = fmt.Sprintf("/rest/control/fabrics/%s/inventory/discover", fabricName)
	_, err = dcnmClient.Save(dUrl, invModel)
	if err != nil {
		return err
	}

	d.SetId(ip)

	var serialNum string
	if d.Get("deploy").(bool) == true {
		var configTimeout int
		for i := 0; i < 3; i++ {
			cont, err = getRemoteSwitch(dcnmClient, fabricName, ip)
			if err != nil {
				return err
			}
			serialNum = stripQuotes(cont.S("serialNumber").String())
			configTimeout = d.Get("config_timeout").(int)
			log.Println(stripQuotes(cont.S("mode").String()))
			if stripQuotes(cont.S("mode").String()) != "Migration" {
				time.Sleep(10 * time.Second)
				break
			}
			time.Sleep(5 * time.Second)
		}

		err = deployswitch(dcnmClient, fabricName, serialNum, configTimeout)
		if err != nil {
			durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serialNum)
			_, delerr := dcnmClient.Delete(durl)
			if delerr != nil {
				return delerr
			}
			d.SetId("")
			return err
		}
	}

	if role, ok := d.GetOk("role"); ok {
		durl := fmt.Sprintf("/rest/control/switches/roles")
		sRole := models.SwitchRole{}
		sRole.Role = role.(string)
		sRole.SerialNumber = serialNum

		_, err := dcnmClient.SaveForAttachment(durl, &sRole)
		if err != nil {
			return err
		}
	}

	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMInventroyRead(d, m)
}

func resourceDCNMInventroyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	ip := d.Get("ip").(string)

	if d.HasChange("username") || d.HasChange("password") || d.HasChange("auth_protocol") {
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		auth := d.Get("auth_protocol").(int)

		cont, err := getRemoteSwitch(dcnmClient, fabricName, ip)
		if err != nil {
			return err
		}

		switchDbID := stripQuotes(cont.S("switchDbID").String())

		body := []byte(fmt.Sprintf("switchIds=%s&userName=%s&password=%s&v3protocol=%s", switchDbID, username, password, strconv.Itoa(auth)))

		durl := fmt.Sprintf("/fm/fmrest/lanConfig/saveSwitchCredentials")
		cont, err = dcnmClient.UpdateCred(durl, body)
		if err != nil {
			return err
		}
	}

	d.SetId(ip)

	var serialNum string
	if d.HasChange("deploy") && d.Get("deploy").(bool) == false {
		d.Set("deploy", true)
		return fmt.Errorf("Deployed switch can not be undeployed")
	} else {
		var configTimeout int
		for i := 0; i < 3; i++ {
			cont, err := getRemoteSwitch(dcnmClient, fabricName, ip)
			if err != nil {
				return err
			}
			serialNum = stripQuotes(cont.S("serialNumber").String())
			configTimeout = d.Get("config_timeout").(int)
			if stripQuotes(cont.S("mode").String()) != "Migration" {
				break
			}
			time.Sleep(5 * time.Second)
		}

		err := deployswitch(dcnmClient, fabricName, serialNum, configTimeout)
		if err != nil {
			d.Set("deploy", false)
			return err
		}
	}

	if d.HasChange("role") {
		durl := fmt.Sprintf("/rest/control/switches/roles")
		sRole := models.SwitchRole{}
		sRole.Role = d.Get("role").(string)
		sRole.SerialNumber = serialNum

		_, err := dcnmClient.SaveForAttachment(durl, &sRole)
		if err != nil {
			return err
		}
	}

	log.Println("[DEBUG] End of Update method ", d.Id())
	return resourceDCNMInventroyRead(d, m)
}

func resourceDCNMInventroyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)
	dn := d.Id()

	cont, err := getRemoteSwitch(dcnmClient, fabricName, dn)
	if err != nil {
		return err
	}

	setSwitchAttributes(d, cont)

	flag, err := checkDeploy(dcnmClient, fabricName, d.Get("serial_number").(string))
	if err != nil {
		return err
	}
	if flag {
		d.Set("deploy", true)
	} else {
		d.Set("deploy", false)
	}

	role, err := getSwitchRole(dcnmClient, d.Get("serial_number").(string))
	if err != nil {
		return err
	}
	d.Set("role", role)

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMInventroyDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	dn := d.Id()

	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabricName)
	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return err
	}

	serialNumber, err := extractSerialNumber(cont, dn)
	if err != nil {
		return err
	}

	durl = fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serialNumber)
	_, err = dcnmClient.Delete(durl)
	if err != nil {
		return err
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
		time.Sleep(time.Duration(timeLeft) * time.Minute)
	}
	if !configDone {
		return fmt.Errorf("Timeout occurs before completion of switch configuration")
	}

	//Step 2 deploy switch into fabric
	durl = fmt.Sprintf("rest/control/fabrics/%s/config-deploy/%s", fabric, serialNum)
	_, err := client.SaveAndDeploy(durl)
	if err != nil {
		return err
	}

	//Step 3 deploy fabric
	durl = fmt.Sprintf("rest/control/fabrics/%s/config-deploy", fabric)
	_, err = client.SaveAndDeploy(durl)
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

	//Step 6 check deployment
	durl = fmt.Sprintf("rest/control/fabrics/%s/config-preview", fabric)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return err
	}
	flag := false
	totalSwitch := len(cont.Data().([]interface{}))
	for i := 0; i < totalSwitch; i++ {
		switchCont := cont.Index(i)
		if stripQuotes(switchCont.S("switchId").String()) == serialNum {
			if stripQuotes(switchCont.S("status").String()) == "In-Sync" {
				flag = true
			}
			break
		}
	}
	if !flag {
		return fmt.Errorf("Switch deployment is not in sync")
	}

	log.Println("[DEBUG] End of Deployment of switch ", serialNum)
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
