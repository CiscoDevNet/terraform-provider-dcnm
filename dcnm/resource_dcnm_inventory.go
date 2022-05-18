package dcnm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDCNMInventroy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDCNMInventroyCreate,
		UpdateContext: resourceDCNMInventroyUpdate,
		ReadContext:   resourceDCNMInventroyRead,
		DeleteContext: resourceDCNMInventroyDelete,

		Schema: map[string]*schema.Schema{
			"fabric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"auth_protocol": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"max_hops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"second_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"preserve_config": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "false",
			},

			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"switch_config": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},

						"switch_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"role": {
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

						"switch_db_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"serial_number": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"model": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: resourceDCNMSwitchConfigHash,
			},

			"config_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},

			"deploy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ips": {
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

	id, err := strconv.Atoi(models.G(cont, "id"))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func extractSwitchinfo(contList *container.Container) models.Switch {
	s := models.Switch{}

	cont := contList.Index(0)

	s.Reachable = models.G(cont, "reachable")
	s.Auth = models.G(cont, "auth")
	s.Known = models.G(cont, "known")
	s.Valid = models.G(cont, "valid")
	s.Selectable = models.G(cont, "selectable")
	s.SysName = models.G(cont, "sysName")
	s.IP = models.G(cont, "ipaddr")
	s.Platform = models.G(cont, "platform")
	s.Version = models.G(cont, "version")
	s.LastChange = models.G(cont, "lastChange")
	s.Hops, _ = strconv.Atoi(models.G(cont, "hopCount"))
	s.DeviceIndex = models.G(cont, "deviceIndex")
	s.StatReason = models.G(cont, "statusReason")

	return s
}

func extractSerialNumber(cont *container.Container, ip string) (string, error) {

	switchCont, err := cont.SearchInObjectList(
		func(infoCont *container.Container) bool {
			return ip == models.G(infoCont, "ipAddress")
		},
	)
	if err != nil {
		return "", fmt.Errorf("no inventory found for given ip address")
	}

	return models.G(switchCont, "serialNumber"), nil
}

func getRemoteSwitch(dcnmClient *client.Client, fabric, ip, serialNum string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabric)

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	switchCont, err := cont.SearchInObjectList(
		func(infoCont *container.Container) bool {
			if ip != "" {
				return ip == models.G(infoCont, "ipAddress")
			}
			return serialNum == models.G(infoCont, "serialNumber")
		},
	)
	if err != nil {
		return nil, fmt.Errorf("desired switch not found")
	}

	return switchCont, nil
}

func getSwitchInfo(cont *container.Container) map[string]interface{} {

	sInfo := make(map[string]interface{})
	sInfo["ip"] = models.G(cont, "ipAddress")
	sInfo["switch_name"] = models.G(cont, "logicalName")
	sInfo["switch_db_id"] = models.G(cont, "switchDbID")
	sInfo["serial_number"] = models.G(cont, "serialNumber")
	sInfo["model"] = models.G(cont, "model")
	sInfo["mode"] = models.G(cont, "mode")

	return sInfo
}

func resourceDCNMInventroyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Create method ")

	var diags diag.Diagnostics
	dcnmClient := m.(*client.Client)

	// Get attribute values from Terraform Config
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

	// Test reachability of desired switches
	for _, val := range switchInfos {
		sInfo := val.(map[string]interface{})
		ip := sInfo["ip"].(string)

		inv.SeedIP = ip

		fabricID, err := extractFabricID(dcnmClient, fabricName)
		if err != nil {
			return diag.FromErr(err)
		}

		dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/test-reachability", strconv.Itoa(fabricID))
		cont, err := dcnmClient.Save(dUrl, &inv)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Detail:   fmt.Sprintf("error at test reachability for switch %s: %s", ip, err),
			})
			continue
		}

		switchM := extractSwitchinfo(cont)

		if switchM.Selectable != "true" || switchM.Reachable != "true" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Detail:   fmt.Sprintf("desired switch is not reachable or not selectable or invalid user/password or bad authentication protocol %s", ip),
			})
			continue
		}

		switchObjs = append(switchObjs, &switchM)
		discoveredIps = append(discoveredIps, ip)
	}

	inv.SeedIP = strings.Join(discoveredIps, ",")

	invModel := models.NewSwitch(&inv, switchObjs)

	// Discover switches
	dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/discover", fabricName)
	_, err := dcnmClient.Save(dUrl, invModel)
	if err != nil {
		return diag.Errorf("error at discovery for switches: %s", err)
	}

	var delFlag bool
	deployedIps := make([]string, 0, 1)
	deployedSerial := make([]string, 0, 1)
	for _, ip := range discoveredIps {
		var serialNum string
		configTimeout := (d.Get("config_timeout").(int)) * 60
		migrate := true
		statusCheckInit := time.Now()
		configTime := time.Now()
		for configTime.Sub(statusCheckInit) <= (time.Duration(configTimeout) * time.Second) {
			cont, err := getRemoteSwitch(dcnmClient, fabricName, ip, "")
			if err != nil {
				log.Println("error at get call for switch in creation :", ip, err)
				continue
			}
			serialNum = models.G(cont, "serialNumber")

			if models.G(cont, "mode") != "Migration" {
				time.Sleep(10 * time.Second)
				migrate = false
				break
			}
			time.Sleep(5 * time.Second)
			configTime = time.Now()
		}
		if migrate {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Detail:   fmt.Sprintf("timeout occurs before going into normal mode. Hence removing it! %s", ip),
			})
			delSwtiches = append(delSwtiches, serialNum)
			delFlag = true
			continue
		}

		// err = deployswitch(dcnmClient, fabricName, serialNum, configTimeout)
		// if err != nil {
		// 	delSwtiches = append(delSwtiches, serialNum)
		// 	diags = append(diags, diag.Diagnostic{
		// 		Severity: diag.Warning,
		// 		Detail:   fmt.Sprintf("error at switch deployment %s: %s", ip, err),
		// 	})
		// 	delFlag = true
		// 	continue
		// }

		deployedIps = append(deployedIps, ip)
		deployedSerial = append(deployedSerial, serialNum)
	}

	// if delFlag {
	// 	for _, serial := range delSwtiches {
	// 		_, err := getRemoteSwitch(dcnmClient, fabricName, "", serial)
	// 		if err == nil {
	// 			durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serial)
	// 			_, delerr := dcnmClient.Delete(durl)
	// 			if delerr != nil {
	// 				log.Println()
	// 				diags = append(diags, diag.Diagnostic{
	// 					Severity: diag.Warning,
	// 					Detail:   fmt.Sprintf("error at deletion of switch %s", err),
	// 				})
	// 			}
	// 		}
	// 	}
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Detail:   "some switches failed to discover and deploy, resuming procedure for successfully discovered switches",
	// 	})
	// }

	for _, ip := range deployedIps {
		for _, val := range switchInfos {
			sInfo := val.(map[string]interface{})

			if sInfo["ip"].(string) == ip && sInfo["role"] != "" {
				cont, err := getRemoteSwitch(dcnmClient, fabricName, ip, "")
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Detail:   fmt.Sprintf("error at get call for switch in creation %s: %s", ip, err),
					})
					continue
				}
				serialNum := models.G(cont,"serialNumber")
				durl := "/rest/control/switches/roles"
				sRole := models.SwitchRole{}
				sRole.Role = roleMappingFunc(sInfo["role"].(string))
				sRole.SerialNumber = serialNum

				_, err = dcnmClient.SaveForAttachment(durl, &sRole)
				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Detail:   fmt.Sprintf("error at switch role assignment %s: %s", ip, err),
					})
				}
			}
		}
	}

	delFlag = false
	err = deployFabric(dcnmClient, fabricName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Detail:   fmt.Sprintf("error at fabric deployment: %s", err),
		})
		delFlag = true
	}

	if delFlag {
		for _, serial := range deployedSerial {
			_, err := getRemoteSwitch(dcnmClient, fabricName, "", serial)
			if err == nil {
				durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serial)
				_, delerr := dcnmClient.Delete(durl)
				if delerr != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Detail:   fmt.Sprintf("error at deletion of switch: %s", err),
					})
				}
			}
		}
	} else {
		ips = append(ips, deployedIps...)
	}

	d.Set("ips", ips)

	if len(ips) == 0 {
		return append(diags, diag.Errorf("none of the switches are discovered and deployed on the fabric, some internal issue in switches")...)
	}

	d.SetId(strings.Join(ips, ","))

	log.Println("[DEBUG] End of Create method ", d.Id())
	return append(diags, resourceDCNMInventroyRead(ctx, d, m)...)
}

func resourceDCNMInventroyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	var diags diag.Diagnostics
	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	if d.HasChange("deploy") && !d.Get("deploy").(bool) {
		d.Set("deploy", true)
		return append(diags, diag.Errorf("deployed switches can not be undeployed")...)
	}
	delSwtiches := make([]string, 0, 1)
	var delFlag bool

	ipDns := d.Get("ips").([]interface{})

	ips := make([]string, 0, 1)
	switchInfosOld, switchInfosNew := d.GetChange("switch_config")

	deleteSwitches := getSerialsForDelete(switchInfosOld.(*schema.Set).List(), switchInfosNew.(*schema.Set).List())
	err := deleteSpecificSwitches(dcnmClient, fabricName, deleteSwitches)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
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
				return append(diags, diag.FromErr(err)...)
			}

			switchDbID := stripQuotes(cont.S("switchDbID").String())

			body := []byte(fmt.Sprintf("switchIds=%s&userName=%s&password=%s&v3protocol=%s", switchDbID, inv.Username, inv.Password, strconv.Itoa(auth)))

			durl := "/fm/fmrest/lanConfig/saveSwitchCredentials"
			if dcnmClient.GetPlatform() == "nd" {
				durl = "/rest/lanConfig/saveSwitchCredentials"
			}
			_, err = dcnmClient.UpdateCred(durl, body)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Detail:   fmt.Sprintf("error at updation of switch %s: %s", ip, err),
				})
			}

			ips = append(ips, ip)

		} else {
			newSwitchFlag = true

			inv.SeedIP = ip

			fabricID, err := extractFabricID(dcnmClient, fabricName)
			if err != nil {
				return append(diags, diag.FromErr(err)...)
			}

			dUrl := fmt.Sprintf("/rest/control/fabrics/%s/inventory/test-reachability", strconv.Itoa(fabricID))
			cont, err := dcnmClient.Save(dUrl, &inv)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Detail:   fmt.Sprintf("error at test reachability for switch %s: %s", ip, err),
				})
				continue
			}

			switchM := extractSwitchinfo(cont)

			if switchM.Selectable != "true" || switchM.Reachable != "true" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Detail:   fmt.Sprintf("desired switch: is not reachable or not selectable or invalid user/password or bad authentication protocol %s", ip),
				})
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
			return append(diags, diag.Errorf("error at discovery for switch : %s", err)...)
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
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Warning,
						Detail:   fmt.Sprintf("error at get call for switch in updation %s: %s", ip, err),
					})
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
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Detail:   fmt.Sprintf("timeout occurs before going into normal mode. Hence removing it! %s", ip),
				})
				delSwtiches = append(delSwtiches, serialNum)
				delFlag = true
				continue
			}

			err := deployswitch(dcnmClient, fabricName, serialNum, configTimeout)
			if err != nil {
				delSwtiches = append(delSwtiches, serialNum)
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Detail:   fmt.Sprintf("error at switch deployment %s: %s", ip, err),
				})
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
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Detail:   fmt.Sprintf("error at deletion of switch %s", err),
						})
					}
				}
			}
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Detail:   "some switches failed to discover and deploy, resuming procedure for successfully discovered switches",
			})
		}

		delFlag = false
		err = deployFabric(dcnmClient, fabricName)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Detail:   fmt.Sprintf("error at fabric deployment: %s", err),
			})
			delFlag = true
		}

		if delFlag {
			for _, serial := range deployedSerial {
				_, err := getRemoteSwitch(dcnmClient, fabricName, "", serial)
				if err == nil {
					durl := fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serial)
					_, delerr := dcnmClient.Delete(durl)
					if delerr != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Detail:   fmt.Sprintf("error at deletion of switch %s", err),
						})
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
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Detail:   fmt.Sprintf("error at get call for switch in creation %s: %s", ip, err),
						})
						continue
					}
					serialNum := stripQuotes(cont.S("serialNumber").String())

					durl := "/rest/control/switches/roles"
					sRole := models.SwitchRole{}
					sRole.Role = roleMappingFunc(sInfo["role"].(string))
					sRole.SerialNumber = serialNum

					_, err = dcnmClient.SaveForAttachment(durl, &sRole)
					if err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Detail:   fmt.Sprintf("error at switch role assignment %s: %s", ip, err),
						})
					}
				}
			}
		}
	}

	err = deployFabric(dcnmClient, fabricName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Detail:   fmt.Sprintf("error at fabric deployment after role assignment: %s", err),
		})
	}

	d.Set("ips", ips)
	d.SetId(strings.Join(ips, ","))

	log.Println("[DEBUG] End of Update method ", d.Id())
	return append(diags, resourceDCNMInventroyRead(ctx, d, m)...)
}

func resourceDCNMInventroyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceDCNMInventroyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	var diags diag.Diagnostics
	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	dn := strings.Split(d.Id(), ",")

	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabricName)
	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return diag.FromErr(err)
	}

	delErr := false
	deletedIps := make([]string, 0, 1)

	for _, ip := range dn {
		serialNumber, err := extractSerialNumber(cont, strings.Trim(ip, " "))
		if err != nil {
			return diag.FromErr(err)
		}

		durl = fmt.Sprintf("/rest/control/fabrics/%s/switches/%s", fabricName, serialNumber)
		_, err = dcnmClient.Delete(durl)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Detail:   fmt.Sprintf("error at deletion of switch %s: %s", ip, err),
			})
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

		return append(diags, diag.Errorf("all switches are not deleted properly")...)
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
	// //Step 3 deploy fabric
	// durl := fmt.Sprintf("rest/control/fabrics/%s/config-deploy", fabric)
	// _, err := client.SaveAndDeploy(durl)
	// if err != nil {
	// 	return err
	// }

	//Step 4 Save configuration
	durl := fmt.Sprintf("rest/control/fabrics/%s/config-save", fabric)
	_, err := client.SaveAndDeploy(durl)
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
