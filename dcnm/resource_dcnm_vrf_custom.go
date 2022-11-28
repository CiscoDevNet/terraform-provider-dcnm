package dcnm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDCNMVRFCustom() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMVRFCustomCreate,
		Read:   resourceDCNMVRFCustomRead,
		Update: resourceDCNMVRFCustomUpdate,
		Delete: resourceDCNMVRFCustomDelete,

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"segment_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_VRF_Universal",
			},
			"extension_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_VRF_Extension_Universal",
			},

			"service_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_props": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"deploy_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},

			"attachments": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"serial_number": {
							Type:     schema.TypeString,
							Required: true,
						},

						"vlan_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"attach": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"free_form_config": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"extension_values": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"loopback_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"loopback_ipv4": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"loopback_ipv6": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vrf_lite": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"peer_vrf_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"interface_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"dot1q_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ip_mask": {
										Type:     schema.TypeString,
										Optional: true,
										// Computed: true,
									},
									"neighbor_ip": {
										Type:     schema.TypeString,
										Optional: true,

										// Computed: true,
									},
									"neighbor_asn": {
										Type:     schema.TypeString,
										Optional: true,

										// Computed: true,
									},
									"ipv6_mask": {
										Type:     schema.TypeString,
										Optional: true,

										// Computed: true,
									},
									"ipv6_neighbor": {
										Type:     schema.TypeString,
										Optional: true,

										// Computed: true,
									},
									"auto_vrf_lite_flag": {
										Type:     schema.TypeString,
										Optional: true,

										// Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func setVRFCustomAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	if cont.Exists("fabric") {
		d.Set("fabric_name", stripQuotes(cont.S("fabric").String()))
	}
	d.Set("name", stripQuotes(cont.S("vrfName").String()))
	d.Set("template", stripQuotes(cont.S("vrfTemplate").String()))
	d.Set("extension_template", stripQuotes(cont.S("vrfExtensionTemplate").String()))
	d.Set("segment_id", stripQuotes(cont.S("vrfId").String()))

	if cont.Exists("serviceVrfTemplate") && stripQuotes(cont.S("serviceVrfTemplate").String()) != "null" {
		d.Set("service_template", stripQuotes(cont.S("serviceVrfTemplate").String()))
	}
	if cont.Exists("source") && stripQuotes(cont.S("source").String()) != "null" {
		d.Set("source", stripQuotes(cont.S("source").String()))
	}

	var strByte []byte
	if cont.Exists("vrfTemplateConfig") {
		strJson := models.G(cont, "vrfTemplateConfig")
		strJson = strings.ReplaceAll(strJson, "\\", "")
		strByte = []byte(strJson)
		var vrfTemplateConfig map[string]string
		json.Unmarshal(strByte, &vrfTemplateConfig)
		props, ok := d.GetOk("template_props")

		map2 := make(map[string]interface{})
		for k := range props.(map[string]interface{}) {
			map2[k] = vrfTemplateConfig[k]
		}
		if !ok {
			d.Set("template_props", vrfTemplateConfig)
		} else {

			d.Set("template_props", map2)
		}
	}

	d.SetId(stripQuotes(cont.S("vrfName").String()))
	return d
}

func resourceDCNMVRFCustomCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	vrf := models.VRF{}
	vrf.Name = d.Get("name").(string)
	vrf.Fabric = d.Get("fabric_name").(string)
	vrf.Template = d.Get("template").(string)
	vrf.ExtensionTemplate = d.Get("extension_template").(string)

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); !ok {
			return fmt.Errorf("attachments must be configured if deploy=true")
		}
	}

	if segmentId, ok := d.GetOk("segment_id"); ok {
		vrf.Id = segmentId.(string)
	} else {
		//request to get the next vrf segment id
		if dcnmClient.GetPlatform() == "nd" {
			cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/top-down/fabrics/%s/vrfinfo", vrf.Fabric))
			if err != nil {
				return err
			}
			vrf.Id = cont.S("l3vni").String()
		} else {
			cont, err := dcnmClient.GetSegID(fmt.Sprintf("/rest/managed-pool/fabrics/%s/partitions/ids", vrf.Fabric))
			if err != nil {
				return err
			}
			vrf.Id = cont.S("partitionSegmentId").String()
		}
	}

	if srcTemp, ok := d.GetOk("service_template"); ok {
		vrf.ServiceVRFTemplate = srcTemp.(string)
	}

	if src, ok := d.GetOk("source"); ok {
		vrf.Source = src.(string)
	}

	vrfConfig := d.Get("template_props").(map[string]interface{})

	confStr, err := json.Marshal(vrfConfig)
	if err != nil {
		return err
	}
	vrf.Config = string(confStr)

	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs", vrf.Fabric)
	_, err = dcnmClient.Save(durl, &vrf)
	if err != nil {
		return err
	}

	d.SetId(vrf.Name)

	//VRF attachment
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); ok {
			attachList := make([]map[string]interface{}, 0, 1)
			for _, val := range d.Get("attachments").(*schema.Set).List() {
				attachment := val.(map[string]interface{})

				attachMap := make(map[string]interface{})

				durl := fmt.Sprintf("/rest/control/switches/%s/fabric-name", attachment["serial_number"].(string))
				cont, err := dcnmClient.GetviaURL(durl)
				if err != nil {
					return err
				}
				attachmentFabricName := stripQuotes(cont.S("fabricName").String())

				attachMap["fabric"] = attachmentFabricName
				attachMap["vrfName"] = vrf.Name
				attachMap["deployment"] = attachment["attach"].(bool)
				attachMap["serialNumber"] = attachment["serial_number"].(string)

				if attachment["vlan_id"].(int) != 0 {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				} else {
					attachMap["vlan"] = vrfConfig["vrfVlanId"]
				}
				if attachment["free_form_config"] != nil {
					attachMap["freeformConfig"] = attachment["free_form_config"].(string)
				}
				if attachment["extension_values"] != nil {
					attachMap["extensionValues"] = attachment["extension_values"].(string)
				}

				flag := false
				instance := models.VRFInstance{}
				if attachment["loopback_id"] != nil {
					instance.LookbackID = attachment["loopback_id"].(int)
					flag = true
				}
				if attachment["loopback_ipv4"] != nil {
					instance.LoopbackIpv4 = attachment["loopback_ipv4"].(string)
					flag = true
				}
				if attachment["loopback_ipv6"] != nil {
					instance.LoopbackIpv6 = attachment["loopback_ipv6"].(string)
					flag = true
				}
				if flag {
					instStr, err := json.Marshal(instance)
					if err != nil {
						return err
					}
					attachMap["instanceValues"] = string(instStr)
				}

				if attachment["vrf_lite"] != nil {
					vrfLiteList := make([]map[string]interface{}, 0, 1)
					ifNameNotFoundList := make([]string, 0, 1)
					for _, val := range attachment["vrf_lite"].(*schema.Set).List() {
						log.Println("vrf_lite enter")
						vrfLite := val.(map[string]interface{})
						vrfLiteMap := make(map[string]interface{})

						durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/switches?vrf-names=%s&serial-numbers=%s", attachmentFabricName, vrf.Name, attachMap["serialNumber"].(string))
						cont, err := dcnmClient.GetviaURL(durl)
						if err != nil {
							return err
						}
						ifNameFound := false
						extensionProtValues := cont.Index(0).S("switchDetailsList").Index(0).S("extensionPrototypeValues")
						for i := 0; i < len(extensionProtValues.Data().([]interface{})); i++ {
							extensionProtVal := extensionProtValues.Index(i)
							if ifName := stripQuotes(extensionProtVal.S("interfaceName").String()); ifName == vrfLite["interface_name"] {
								ifNameFound = true
								extensionValueString := stripQuotes(extensionProtVal.S("extensionValues").String())

								var extensionValues map[string]interface{}
								extensionValueString = strings.Replace(extensionValueString, "\\", "", -1)
								err = json.Unmarshal([]byte(extensionValueString), &extensionValues)
								if err != nil {
									return err
								}
								if len(extensionValues) != 0 {
									vrfLiteMap["PEER_VRF_NAME"] = vrfLite["peer_vrf_name"]
									vrfLiteMap["IF_NAME"] = vrfLite["interface_name"]

									if vrfLite["dot1q_id"] != "" {
										vrfLiteMap["DOT1Q_ID"] = vrfLite["dot1q_id"]
									} else {
										durl := "/rest/resource-manager/reserve-id"
										dot1q := models.VRFDot1qID{
											ScopeType:    "DeviceInterface",
											UsageType:    "TOP_DOWN_L3_DOT1Q",
											AllocatedTo:  vrf.Name,
											SerialNumber: attachMap["serialNumber"].(string),
											IfName:       ifName,
										}
										cont, err := dcnmClient.Save(durl, &dot1q)
										if err != nil {
											return err
										}
										vrfLiteMap["DOT1Q_ID"] = stripQuotes(cont.String())
									}

									if vrfLite["ip_mask"] != "" {
										vrfLiteMap["IP_MASK"] = vrfLite["ip_mask"].(string)
									} else if extensionValues["IP_MASK"] != nil {
										vrfLiteMap["IP_MASK"] = extensionValues["IP_MASK"].(string)
									}

									if vrfLite["neighbor_ip"] != "" {
										vrfLiteMap["NEIGHBOR_IP"] = vrfLite["neighbor_ip"].(string)
									} else if extensionValues["NEIGHBOR_IP"] != nil {
										vrfLiteMap["NEIGHBOR_IP"] = extensionValues["NEIGHBOR_IP"].(string)
									}

									if vrfLite["neighbor_asn"] != "" {
										vrfLiteMap["NEIGHBOR_ASN"] = vrfLite["neighbor_asn"].(string)
									} else if extensionValues["NEIGHBOR_ASN"] != nil {
										vrfLiteMap["NEIGHBOR_ASN"] = extensionValues["NEIGHBOR_ASN"].(string)
									}

									if vrfLite["ipv6_mask"] != "" {
										vrfLiteMap["IPV6_MASK"] = vrfLite["ipv6_mask"].(string)
									} else if extensionValues["IPV6_MASK"] != nil {
										vrfLiteMap["IPV6_MASK"] = extensionValues["IPV6_MASK"].(string)
									}

									if vrfLite["ipv6_neighbor"] != "" {
										vrfLiteMap["IPV6_NEIGHBOR"] = vrfLite["ipv6_neighbor"].(string)
									} else if extensionValues["IPV6_NEIGHBOR"] != nil {
										vrfLiteMap["IPV6_NEIGHBOR"] = extensionValues["IPV6_NEIGHBOR"].(string)
									}

									if vrfLite["auto_vrf_lite_flag"] != "" {
										vrfLiteMap["AUTO_VRF_LITE_FLAG"] = vrfLite["auto_vrf_lite_flag"].(string)
									} else if extensionValues["AUTO_VRF_LITE_FLAG"] != nil {
										vrfLiteMap["AUTO_VRF_LITE_FLAG"] = extensionValues["AUTO_VRF_LITE_FLAG"].(string)
									}
									vrfLiteMap["VRF_LITE_JYTHON_TEMPLATE"] = extensionValues["VRF_LITE_JYTHON_TEMPLATE"].(string)

								} else {
									return fmt.Errorf("No VRF_LITE Data found for switch %s", attachMap["serialNumber"].(string))
								}
								contMap := make(map[string]interface{})
								vrfLiteStr, err := json.Marshal(map[string]interface{}{
									"VRF_LITE_CONN": vrfLiteList,
								})
								if err != nil {
									return err
								}
								contMap["VRF_LITE_CONN"] = string(vrfLiteStr)
								vrfLiteStr, err = json.Marshal(contMap)
								if err != nil {
									return err
								}
								attachMap["extensionValues"] = string(vrfLiteStr)
							}
						}
						if !ifNameFound {
							ifNameNotFoundList = append(ifNameNotFoundList, vrfLite["interface_name"].(string))
						}

						vrfLiteList = append(vrfLiteList, vrfLiteMap)
					}

					if len(ifNameNotFoundList) > 0 {
						return fmt.Errorf("VRF LITE Config not found for attachment:%s", ifNameNotFoundList)
					}

					contMap := make(map[string]interface{})
					vrfLiteStr, err := json.Marshal(map[string]interface{}{
						"VRF_LITE_CONN": vrfLiteList,
					})
					if err != nil {
						return err
					}
					contMap["VRF_LITE_CONN"] = string(vrfLiteStr)
					vrfLiteStr, err = json.Marshal(contMap)
					if err != nil {
						return err
					}
					attachMap["extensionValues"] = string(vrfLiteStr)
				} else {
					attachMap["extensionValues"] = ""
				}
				attachList = append(attachList, attachMap)
			}

			vrfAttach := models.NewVRFAttachment(vrf.Name, attachList)
			durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/attachments", vrf.Fabric)
			cont, err := dcnmClient.SaveForAttachment(durl, vrfAttach)
			if err != nil {
				return err
			}

			// VRF Deployment
			for _, v := range cont.Data().(map[string]interface{}) {
				if v != "SUCCESS" && v != "SUCCESS Peer attach Reponse :  SUCCESS" {
					return fmt.Errorf("VRF record is created but not deployed yet. Error while attachment : %s", v)
				}
			}
			vrfD := models.VRFDeploy{}
			vrfD.Name = vrf.Name
			durl = fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/deployments", vrf.Fabric)
			_, err = dcnmClient.Save(durl, &vrfD)
			if err != nil {
				d.Set("deploy", false)
			}

			deployFlag := false
			deployTimeout := d.Get("deploy_timeout").(int)
			for j := 0; j < (deployTimeout / 5); j++ {
				deployStatus, err := getVRFDeploymentStatus(dcnmClient, vrf.Fabric, vrf.Name)
				if err != nil {
					return err
				}
				deployFlag = deployStatus == "DEPLOYED"
				if !deployFlag {
					time.Sleep(5 * time.Second)
				} else {
					deployFlag = true
					break
				}
			}
			if !deployFlag {
				return fmt.Errorf("VRF record is created but not deployed yet. deployment timeout occured")
			}

		} else {
			d.Set("deploy", false)
			d.Set("attachments", make([]interface{}, 0, 1))
			return fmt.Errorf("VRF record is created but not deployed yet. Either make deploy=false or provide attachments")
		}
	}

	d.SetId(vrf.Name)
	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMVRFRead(d, m)
}

func resourceDCNMVRFCustomRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	fabricName := d.Get("fabric_name").(string)

	cont, err := getRemoteVRF(dcnmClient, fabricName, dn)
	if err != nil {
		return err
	}

	setVRFCustomAttributes(d, cont)

	flag, err := checkvrfDeploy(dcnmClient, fabricName, dn)
	if err != nil {
		d.Set("deploy", false)
		return err
	}
	d.Set("deploy", flag)

	if attaches, ok := d.GetOk("attachments"); ok {
		attachGet := make([]interface{}, 0, 1)

		for _, val := range attaches.(*schema.Set).List() {
			attachMap := val.(map[string]interface{})
			serialNum := attachMap["serial_number"].(string)

			attachStatus, vlan, err := getSwitchAttachStatus(dcnmClient, fabricName, dn, serialNum)
			if err == nil {
				attachMap["attach"] = attachStatus
				if attachMap["vlan_id"].(int) != 0 {
					attachMap["vlan_id"] = vlan
				}
			}
			if attachMap["vrf_lite"] != nil {
				lites := attachMap["vrf_lite"].(*schema.Set).List()

				liteGet := make([]interface{}, 0, 1)

				durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/switches?vrf-names=%s&serial-numbers=%s", fabricName, dn, attachMap["serial_number"].(string))
				cont, err := dcnmClient.GetviaURL(durl)
				if err != nil {
					return err
				}
				extensionValueString := stripQuotes(cont.Index(0).S("switchDetailsList").Index(0).S("extensionValues").String())

				if extensionValueString != "null" {
					var extensionValues map[string]interface{}
					extensionValueString = strings.Replace(extensionValueString, "\\\"", "\"", -1)
					extensionValueString = strings.Replace(extensionValueString, "\\\"", "\"", -1)

					var extensionValuesList []interface{}

					err = json.Unmarshal([]byte(extensionValueString), &extensionValues)
					if err != nil {
						fmt.Println(err)
					}

					_ = json.Unmarshal([]byte(extensionValues["VRF_LITE_CONN"].(string)), &extensionValues)

					extensionValuesList = extensionValues["VRF_LITE_CONN"].([]interface{})

					for i, _ := range extensionValuesList {
						vrfLiteMap := make(map[string]interface{}, 0)

						for _, val := range lites {
							vrfLite := val.(map[string]interface{})

							extensionValues = extensionValuesList[i].(map[string]interface{})

							if extensionValues["IF_NAME"].(string) == vrfLite["interface_name"] {

								if extensionValues["PEER_VRF_NAME"] != nil {
									vrfLiteMap["peer_vrf_name"] = extensionValues["PEER_VRF_NAME"].(string)
								}

								if extensionValues["IF_NAME"] != nil {
									vrfLiteMap["interface_name"] = extensionValues["IF_NAME"].(string)
								}

								if len(extensionValues) != 0 {

									if vrfLite["dot1q_id"] != "" {
										vrfLiteMap["dot1q_id"] = extensionValues["DOT1Q_ID"].(string)
									}
									if vrfLite["neighbor_ip"] != "" {
										vrfLiteMap["neighbor_ip"] = extensionValues["NEIGHBOR_IP"].(string)
									}
									if vrfLite["ip_mask"] != "" {
										vrfLiteMap["ip_mask"] = extensionValues["IP_MASK"].(string)
									}

									if vrfLite["neighbor_asn"] != "" {
										vrfLiteMap["neighbor_asn"] = extensionValues["NEIGHBOR_ASN"].(string)
									}
									if vrfLite["ipv6_mask"] != "" {
										vrfLiteMap["ipv6_mask"] = extensionValues["IPV6_MASK"].(string)
									}
									if vrfLite["ipv6_neighbor"] != "" {
										vrfLiteMap["ipv6_neighbor"] = extensionValues["IPV6_NEIGHBOR"].(string)
									}
									if vrfLite["auto_vrf_lite_flag"] != "" {
										vrfLiteMap["auto_vrf_lite_flag"] = extensionValues["AUTO_VRF_LITE_FLAG"].(string)
									}
								}
								liteGet = append(liteGet, vrfLiteMap)
							} else {
								if extensionValues["PEER_VRF_NAME"] != nil {
									vrfLiteMap["peer_vrf_name"] = extensionValues["PEER_VRF_NAME"].(string)
								}
								if extensionValues["IF_NAME"] != nil {
									vrfLiteMap["interface_name"] = extensionValues["IF_NAME"].(string)
								}
							}
						}

						liteGet = append(liteGet, vrfLiteMap)
					}
				}

				attachMap["vrf_lite"] = liteGet
			}
			attachGet = append(attachGet, attachMap)
		}
		log.Printf("attachGet: %v\n", attachGet)
		d.Set("attachments", attachGet)
	}

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMVRFCustomUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	dcnmClient := m.(*client.Client)

	vrf := models.VRF{}
	vrf.Name = d.Get("name").(string)
	vrf.Fabric = d.Get("fabric_name").(string)
	vrf.Template = d.Get("template").(string)
	vrf.ExtensionTemplate = d.Get("extension_template").(string)
	vrf.Id = d.Get("segment_id").(string)

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); !ok {
			return fmt.Errorf("attachments must be configured if deploy=true")
		}
	}

	if srcTemp, ok := d.GetOk("service_template"); ok {
		vrf.ServiceVRFTemplate = srcTemp.(string)
	}

	if src, ok := d.GetOk("source"); ok {
		vrf.Source = src.(string)
	}

	vrfConfig := d.Get("template_props").(map[string]interface{})

	confStr, err := json.Marshal(vrfConfig)
	if err != nil {
		return err
	}
	vrf.Config = string(confStr)

	dn := d.Id()
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/%s", vrf.Fabric, dn)
	_, err = dcnmClient.Update(durl, &vrf)
	if err != nil {
		return err
	}

	d.SetId(vrf.Name)

	//VRF Attachment
	if d.HasChange("deploy") && d.Get("deploy").(bool) == false {
		return fmt.Errorf("Deployed VRF can not be undeployed")
	}

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); ok {
			attachList := make([]map[string]interface{}, 0, 1)
			for _, val := range d.Get("attachments").(*schema.Set).List() {
				attachment := val.(map[string]interface{})

				attachMap := make(map[string]interface{})

				durl := fmt.Sprintf("/rest/control/switches/%s/fabric-name", attachment["serial_number"].(string))
				cont, err := dcnmClient.GetviaURL(durl)
				if err != nil {
					return err
				}
				attachmentFabricName := stripQuotes(cont.S("fabricName").String())

				attachMap["fabric"] = attachmentFabricName
				attachMap["vrfName"] = vrf.Name
				attachMap["deployment"] = attachment["attach"].(bool)
				attachMap["serialNumber"] = attachment["serial_number"].(string)

				if attachment["vlan_id"].(int) != 0 {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				} else {
					attachMap["vlan"] = vrfConfig["vrfVlanId"]
				}
				if attachment["free_form_config"] != nil {
					attachMap["freeformConfig"] = attachment["free_form_config"].(string)
				}
				if attachment["extension_values"] != nil {
					attachMap["extensionValues"] = attachment["extension_values"].(string)
				}

				flag := false
				instance := models.VRFInstance{}
				if attachment["loopback_id"] != nil {
					instance.LookbackID = attachment["loopback_id"].(int)
					flag = true
				}
				if attachment["loopback_ipv4"] != nil {
					instance.LoopbackIpv4 = attachment["loopback_ipv4"].(string)
					flag = true
				}
				if attachment["loopback_ipv6"] != nil {
					instance.LoopbackIpv6 = attachment["loopback_ipv6"].(string)
					flag = true
				}
				if flag {
					instStr, err := json.Marshal(instance)
					if err != nil {
						return err
					}
					attachMap["instanceValues"] = string(instStr)
				}

				if attachment["vrf_lite"] != nil {
					vrfLiteList := make([]map[string]interface{}, 0, 1)
					ifNameNotFoundList := make([]string, 0, 1)
					for _, val := range attachment["vrf_lite"].(*schema.Set).List() {
						log.Println("vrf_lite enter")
						vrfLite := val.(map[string]interface{})
						vrfLiteMap := make(map[string]interface{})

						durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/switches?vrf-names=%s&serial-numbers=%s", attachmentFabricName, vrf.Name, attachMap["serialNumber"].(string))
						cont, err := dcnmClient.GetviaURL(durl)
						if err != nil {
							return err
						}
						ifNameFound := false
						extensionProtValues := cont.Index(0).S("switchDetailsList").Index(0).S("extensionPrototypeValues")
						for i := 0; i < len(extensionProtValues.Data().([]interface{})); i++ {
							extensionProtVal := extensionProtValues.Index(i)
							if ifName := stripQuotes(extensionProtVal.S("interfaceName").String()); ifName == vrfLite["interface_name"] {
								ifNameFound = true
								extensionValueString := stripQuotes(extensionProtVal.S("extensionValues").String())

								var extensionValues map[string]interface{}
								extensionValueString = strings.Replace(extensionValueString, "\\", "", -1)
								err = json.Unmarshal([]byte(extensionValueString), &extensionValues)
								if err != nil {
									return err
								}
								if len(extensionValues) != 0 {
									vrfLiteMap["PEER_VRF_NAME"] = vrfLite["peer_vrf_name"]
									vrfLiteMap["IF_NAME"] = vrfLite["interface_name"]

									if vrfLite["dot1q_id"] != "" {
										vrfLiteMap["DOT1Q_ID"] = vrfLite["dot1q_id"]
									} else {
										durl := "/rest/resource-manager/reserve-id"
										dot1q := models.VRFDot1qID{
											ScopeType:    "DeviceInterface",
											UsageType:    "TOP_DOWN_L3_DOT1Q",
											AllocatedTo:  vrf.Name,
											SerialNumber: attachMap["serialNumber"].(string),
											IfName:       ifName,
										}
										cont, err := dcnmClient.Save(durl, &dot1q)
										if err != nil {
											return err
										}
										vrfLiteMap["DOT1Q_ID"] = stripQuotes(cont.String())
									}

									if vrfLite["ip_mask"] != "" {
										vrfLiteMap["IP_MASK"] = vrfLite["ip_mask"].(string)
									} else if extensionValues["IP_MASK"] != nil {
										vrfLiteMap["IP_MASK"] = extensionValues["IP_MASK"].(string)
									}

									if vrfLite["neighbor_ip"] != "" {
										vrfLiteMap["NEIGHBOR_IP"] = vrfLite["neighbor_ip"].(string)
									} else if extensionValues["NEIGHBOR_IP"] != nil {
										vrfLiteMap["NEIGHBOR_IP"] = extensionValues["NEIGHBOR_IP"].(string)
									}

									if vrfLite["neighbor_ip"] != "" {
										vrfLiteMap["NEIGHBOR_ASN"] = vrfLite["neighbor_asn"].(string)
									} else if extensionValues["NEIGHBOR_ASN"] != nil {
										vrfLiteMap["NEIGHBOR_ASN"] = extensionValues["NEIGHBOR_ASN"].(string)
									}

									if vrfLite["ipv6_mask"] != "" {
										vrfLiteMap["IPV6_MASK"] = vrfLite["ipv6_mask"].(string)
									} else if extensionValues["IPV6_MASK"] != nil {
										vrfLiteMap["IPV6_MASK"] = extensionValues["IPV6_MASK"].(string)
									}

									if vrfLite["ipv6_neighbor"] != "" {
										vrfLiteMap["IPV6_NEIGHBOR"] = vrfLite["ipv6_neighbor"].(string)
									} else if extensionValues["IPV6_NEIGHBOR"] != nil {
										vrfLiteMap["IPV6_NEIGHBOR"] = extensionValues["IPV6_NEIGHBOR"].(string)
									}

									if vrfLite["auto_vrf_lite_flag"] != "" {
										vrfLiteMap["AUTO_VRF_LITE_FLAG"] = vrfLite["auto_vrf_lite_flag"].(string)
									} else if extensionValues["AUTO_VRF_LITE_FLAG"] != nil {
										vrfLiteMap["AUTO_VRF_LITE_FLAG"] = extensionValues["AUTO_VRF_LITE_FLAG"].(string)
									}
									vrfLiteMap["VRF_LITE_JYTHON_TEMPLATE"] = extensionValues["VRF_LITE_JYTHON_TEMPLATE"].(string)

								} else {
									return fmt.Errorf("No VRF_LITE Data found for switch %s", attachMap["serialNumber"].(string))
								}
								contMap := make(map[string]interface{})
								vrfLiteStr, err := json.Marshal(map[string]interface{}{
									"VRF_LITE_CONN": vrfLiteList,
								})
								if err != nil {
									return err
								}
								contMap["VRF_LITE_CONN"] = string(vrfLiteStr)
								vrfLiteStr, err = json.Marshal(contMap)
								if err != nil {
									return err
								}
								attachMap["extensionValues"] = string(vrfLiteStr)
							}
						}
						if !ifNameFound {
							ifNameNotFoundList = append(ifNameNotFoundList, vrfLite["interface_name"].(string))
						}

						vrfLiteList = append(vrfLiteList, vrfLiteMap)
					}

					if len(ifNameNotFoundList) > 0 {
						return fmt.Errorf("VRF LITE Config not found for attachment:%s", ifNameNotFoundList)
					}

					contMap := make(map[string]interface{})
					vrfLiteStr, err := json.Marshal(map[string]interface{}{
						"VRF_LITE_CONN": vrfLiteList,
					})
					if err != nil {
						return err
					}
					contMap["VRF_LITE_CONN"] = string(vrfLiteStr)
					vrfLiteStr, err = json.Marshal(contMap)
					if err != nil {
						return err
					}
					attachMap["extensionValues"] = string(vrfLiteStr)
				} else {
					attachMap["extensionValues"] = ""
				}

				attachList = append(attachList, attachMap)
			}

			vrfAttach := models.NewVRFAttachment(vrf.Name, attachList)
			durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/attachments", vrf.Fabric)
			cont, err := dcnmClient.SaveForAttachment(durl, vrfAttach)
			if err != nil {
				return err
			}

			// VRF Deployment
			for _, v := range cont.Data().(map[string]interface{}) {
				if v != "SUCCESS" && v != "SUCCESS Peer attach Reponse :  SUCCESS" {
					return fmt.Errorf("VRF record is created but not deployed yet. Error while attachment : %s", v)
				}
			}
			vrfD := models.VRFDeploy{}
			vrfD.Name = vrf.Name
			durl = fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/deployments", vrf.Fabric)
			_, err = dcnmClient.Save(durl, &vrfD)
			if err != nil {
				d.Set("deploy", false)
			}

			deployFlag := false
			deployTimeout := d.Get("deploy_timeout").(int)
			for j := 0; j < (deployTimeout / 5); j++ {
				deployStatus, err := getVRFDeploymentStatus(dcnmClient, vrf.Fabric, vrf.Name)
				if err != nil {
					return err
				}
				deployFlag = deployStatus == "DEPLOYED"
				if !deployFlag {
					time.Sleep(5 * time.Second)
				} else {
					deployFlag = true
					break
				}
			}
			if !deployFlag {
				d.Set("deploy", false)
				return fmt.Errorf("VRF record is updated and deployment is initialised, but deployment timeout occured before completion of the deployment process")
			}

		} else {
			d.Set("deploy", false)
			d.Set("attachments", make([]interface{}, 0, 1))
			return fmt.Errorf("VRF record is not deployed yet. Either make deploy=false or provide attachments")
		}
	}

	d.SetId(vrf.Name)
	log.Println("[DEBUG] End of Update method ", d.Id())
	return resourceDCNMVRFRead(d, m)
}

func resourceDCNMVRFCustomDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)

	dn := d.Id()
	fabricName := d.Get("fabric_name").(string)

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); ok {
			attachList := make([]map[string]interface{}, 0, 1)
			for _, val := range d.Get("attachments").(*schema.Set).List() {
				attachment := val.(map[string]interface{})

				attachMap := make(map[string]interface{})
				durl := fmt.Sprintf("/rest/control/switches/%s/fabric-name", attachment["serial_number"].(string))
				cont, err := dcnmClient.GetviaURL(durl)
				if err != nil {
					return err
				}
				attachmentFabricName := stripQuotes(cont.S("fabricName").String())
				attachMap["fabric"] = attachmentFabricName
				attachMap["vrfName"] = dn
				attachMap["deployment"] = false
				attachMap["serialNumber"] = attachment["serial_number"].(string)
				if attachment["vlan_id"].(int) == 0 {
					attachMap["vlan"] = d.Get("template_props").(map[string]interface{})["vrfVlanId"]
				} else {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				}
				attachList = append(attachList, attachMap)
			}

			vrfAttach := models.NewVRFAttachment(dn, attachList)
			durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/attachments", fabricName)
			cont, err := dcnmClient.SaveForAttachment(durl, vrfAttach)
			if err != nil {
				return err
			}

			// VRF Deployment
			for _, v := range cont.Data().(map[string]interface{}) {
				if v != "SUCCESS" && v != "SUCCESS Peer attach Reponse :  SUCCESS" {
					return fmt.Errorf("failure at the time of detachment : %s", v)
				}
			}
			vrfD := models.VRFDeploy{}
			vrfD.Name = dn
			durl = fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/deployments", fabricName)
			_, err = dcnmClient.Save(durl, &vrfD)
			if err != nil {
				d.Set("deploy", false)
			}

			deployFlag := false
			deployTimeout := d.Get("deploy_timeout").(int)
			for j := 0; j < int(deployTimeout/5); j++ {
				deployStatus, err := getVRFDeploymentStatus(dcnmClient, fabricName, dn)
				if err != nil {
					return err
				}
				deployFlag = deployStatus == "NA"
				if !deployFlag {
					time.Sleep(5 * time.Second)
				} else {
					deployFlag = true
					break
				}
			}
			if !deployFlag {
				return fmt.Errorf("VRF record can not be deleted. deployment timeout occured")
			}
		}
	}

	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/%s", fabricName, dn)
	_, err := dcnmClient.Delete(durl)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}
