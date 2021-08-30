package dcnm

import (
	"encoding/json"
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

func resourceDCNMVRF() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMVRFCreate,
		Read:   resourceDCNMVRFRead,
		Update: resourceDCNMVRFUpdate,
		Delete: resourceDCNMVRFDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDCNMVRFImporter,
		},

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

			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  9216,
			},

			"vlan_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"intf_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "12345",
			},

			"max_bgp_path": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"max_ibgp_path": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"trm_enable": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
			},

			"rp_external_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
			},

			"rp_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"mutlicast_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mutlicast_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_link_local_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
			},

			"trm_bgw_msite_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
			},

			"advertise_host_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
			},

			"advertise_default_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
			},

			"static_default_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
				}, false),
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"peer_vrf_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"dot1q_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ip_mask": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"neighbor_ip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"neighbor_asn": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ipv6_mask": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ipv6_neighbor": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"auto_vrf_lite_flag": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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

func getRemoteVRF(client *client.Client, fabricName, vrfName string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/%s", fabricName, vrfName)

	cont, err := client.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	return cont, nil
}

func setVRFAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	d.Set("fabric_name", stripQuotes(cont.S("fabric").String()))
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

	cont, err := cleanJsonString(stripQuotes(cont.S("vrfTemplateConfig").String()))
	if err == nil {
		if cont.Exists("mtu") {
			if mtu, err := strconv.Atoi(stripQuotes(cont.S("mtu").String())); err == nil {
				d.Set("mtu", mtu)
			}
		}
		if cont.Exists("vrfVlanId") {
			if vlan, err := strconv.Atoi(stripQuotes(cont.S("vrfVlanId").String())); err == nil {
				d.Set("vlan_id", vlan)
			}
		}
		if cont.Exists("tag") {
			d.Set("tag", stripQuotes(cont.S("tag").String()))
		}
		if cont.Exists("vrfVlanName") {
			d.Set("vlan_name", stripQuotes(cont.S("vrfVlanName").String()))
		}
		if cont.Exists("vrfDescription") {
			d.Set("description", stripQuotes(cont.S("vrfDescription").String()))
		}
		if cont.Exists("vrfIntfDescription") {
			d.Set("intf_description", stripQuotes(cont.S("vrfIntfDescription").String()))
		}
		if cont.Exists("maxBgpPaths") {
			if bgp, err := strconv.Atoi(stripQuotes(cont.S("maxBgpPaths").String())); err == nil {
				d.Set("max_bgp_path", bgp)
			}
		}
		if cont.Exists("maxIbgpPaths") {
			if ibgp, err := strconv.Atoi(stripQuotes(cont.S("maxIbgpPaths").String())); err == nil {
				d.Set("max_ibgp_path", ibgp)
			}
		}
		if cont.Exists("trmEnabled") {
			d.Set("trm_enable", stripQuotes(cont.S("trmEnabled").String()))
		}
		if cont.Exists("isRPExternal") {
			d.Set("rp_external_flag", stripQuotes(cont.S("isRPExternal").String()))
		}
		if cont.Exists("loopbackNumber") {
			if loopback, err := strconv.Atoi(stripQuotes(cont.S("loopbackNumber").String())); err == nil {
				d.Set("loopback_id", loopback)
			}
		}
		if cont.Exists("rpAddress") {
			d.Set("rp_address", stripQuotes(cont.S("rpAddress").String()))
		}
		if cont.Exists("L3VniMcastGroup") {
			d.Set("mutlicast_address", stripQuotes(cont.S("L3VniMcastGroup").String()))
		}
		if cont.Exists("ipv6LinkLocalFlag") {
			d.Set("ipv6_link_local_flag", stripQuotes(cont.S("ipv6LinkLocalFlag").String()))
		}
		if cont.Exists("multicastGroup") {
			d.Set("mutlicast_group", stripQuotes(cont.S("multicastGroup").String()))
		}
		if cont.Exists("trmBGWMSiteEnabled") {
			d.Set("trm_bgw_msite_flag", stripQuotes(cont.S("trmBGWMSiteEnabled").String()))
		}
		if cont.Exists("advertiseHostRouteFlag") {
			d.Set("advertise_host_route", stripQuotes(cont.S("advertiseHostRouteFlag").String()))
		}
		if cont.Exists("advertiseDefaultRouteFlag") {
			d.Set("advertise_default_route", stripQuotes(cont.S("advertiseDefaultRouteFlag").String()))
		}
		if cont.Exists("configureStaticDefaultRouteFlag") {
			d.Set("static_default_route", stripQuotes(cont.S("configureStaticDefaultRouteFlag").String()))
		}
	}

	d.SetId(stripQuotes(cont.S("vrfName").String()))
	return d
}

func resourceDCNMVRFImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Importer ", d.Id())

	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 2 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	fabricName := importInfo[0]
	vrf := importInfo[1]

	cont, err := getRemoteVRF(dcnmClient, fabricName, vrf)
	if err != nil {
		return nil, err
	}

	flag, err := checkvrfDeploy(dcnmClient, fabricName, vrf)
	if err != nil {
		d.Set("deploy", false)
		return nil, err
	}
	d.Set("deploy", flag)

	stateImport := setVRFAttributes(d, cont)

	log.Println("[DEBUG] End of Importer ", d.Id())
	return []*schema.ResourceData{stateImport}, nil
}

func resourceDCNMVRFCreate(d *schema.ResourceData, m interface{}) error {
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

	configMap := models.VRFProfileConfig{}
	if vlan, ok := d.GetOk("vlan_id"); ok {
		configMap.Vlan = vlan.(int)
	} else {
		durl := fmt.Sprintf("/rest/resource-manager/vlan/%s?vlanUsageType=TOP_DOWN_VRF_VLAN", d.Get("fabric_name").(string))
		cont, err := dcnmClient.GetviaURL(durl)
		if err != nil {
			return err
		}
		vlan, err := strconv.Atoi(cont.String())
		if err == nil {
			configMap.Vlan = vlan
		}
	}
	if mtu, ok := d.GetOk("mtu"); ok {
		configMap.Mtu = mtu.(int)
	}
	if vlanName, ok := d.GetOk("vlan_name"); ok {
		configMap.VlanName = vlanName.(string)
	}
	if desc, ok := d.GetOk("description"); ok {
		configMap.Description = desc.(string)
	}
	if intfDesc, ok := d.GetOk("intf_description"); ok {
		configMap.IntfDescription = intfDesc.(string)
	}
	if tag, ok := d.GetOk("tag"); ok {
		configMap.Tag = tag.(string)
	}
	if bgp, ok := d.GetOk("max_bgp_path"); ok {
		configMap.BGP = bgp.(int)
	}
	if ibgp, ok := d.GetOk("max_ibgp_path"); ok {
		configMap.IBGP = ibgp.(int)
	}
	if trm, ok := d.GetOk("trm_enable"); ok {
		configMap.TRM = trm.(string)
	}
	if rpExtr, ok := d.GetOk("rp_external_flag"); ok {
		configMap.RPexternal = rpExtr.(string)
	}
	if rpAddr, ok := d.GetOk("rp_address"); ok {
		configMap.RPaddress = rpAddr.(string)
	}
	if loopback, ok := d.GetOk("loopback_id"); ok {
		configMap.Lookback = loopback.(int)
	}
	if mcastAddr, ok := d.GetOk("mutlicast_address"); ok {
		configMap.Mcastaddr = mcastAddr.(string)
	}
	if mcastGrp, ok := d.GetOk("mutlicast_group"); ok {
		configMap.Mcastgroup = mcastGrp.(string)
	}
	if ipv6, ok := d.GetOk("ipv6_link_local_flag"); ok {
		configMap.IPv6Link = ipv6.(string)
	}
	if trmbgw, ok := d.GetOk("trm_bgw_msite_flag"); ok {
		configMap.TRMBGW = trmbgw.(string)
	}
	if hostR, ok := d.GetOk("advertise_host_route"); ok {
		configMap.AdhostRoute = hostR.(string)
	}
	if defaultR, ok := d.GetOk("advertise_default_route"); ok {
		configMap.AdDefaultRoute = defaultR.(string)
	}
	if staticR, ok := d.GetOk("static_default_route"); ok {
		configMap.StaticRoute = staticR.(string)
	}
	configMap.SegmentID = vrf.Id
	configMap.VrfName = vrf.Name

	confStr, err := json.Marshal(configMap)
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

				attachMap["fabric"] = vrf.Fabric
				attachMap["vrfName"] = vrf.Name
				attachMap["deployment"] = attachment["attach"].(bool)
				attachMap["serialNumber"] = attachment["serial_number"].(string)

				if attachment["vlan_id"].(int) != 0 {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				} else {
					attachMap["vlan"] = configMap.Vlan
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
					for _, val := range attachment["vrf_lite"].(*schema.Set).List() {
						log.Println("vrf_lite enter")
						vrfLite := val.(map[string]interface{})
						vrfLiteMap := make(map[string]interface{})

						durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/switches?vrf-names=%s&serial-numbers=%s", vrf.Fabric, vrf.Name, attachMap["serialNumber"].(string))
						cont, err := dcnmClient.GetviaURL(durl)
						if err != nil {
							return err
						}
						extensionProtValues := cont.Index(0).S("switchDetailsList").Index(0).S("extensionPrototypeValues").Index(0)
						ifName := stripQuotes(extensionProtValues.S("interfaceName").String())
						extensionValueString := stripQuotes(extensionProtValues.S("extensionValues").String())

						if len(extensionValueString) == 0 {
							return fmt.Errorf("No VRF_LITE Data found for switch %s", attachMap["serialNumber"].(string))
						}

						var extensionValues map[string]interface{}
						extensionValueString = strings.Replace(extensionValueString, "\\", "", -1)
						err = json.Unmarshal([]byte(extensionValueString), &extensionValues)
						if err != nil {
							return err
						}
						vrfLiteMap["PEER_VRF_NAME"] = vrfLite["peer_vrf_name"]
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

						log.Printf("extvals: %v", extensionValues)

						if vrfLite["ip_mask"] != "" {
							vrfLiteMap["IP_MASK"] = vrfLite["ip_mask"].(string)
						} else {
							vrfLiteMap["IP_MASK"] = extensionValues["IP_MASK"].(string)
						}
						log.Printf("extvals: %s", extensionValues["IP_MASK"])
						log.Printf("vrfLiteMap: %s", vrfLiteMap["IP_MASK"])
						log.Printf("vrfLite[\"ip_mask\"]: %v\n", vrfLite["ip_mask"])

						if vrfLite["neighbor_ip"] != "" {
							vrfLiteMap["NEIGHBOR_IP"] = vrfLite["neighbor_ip"].(string)
						} else {
							vrfLiteMap["NEIGHBOR_IP"] = extensionValues["NEIGHBOR_IP"].(string)
						}

						if vrfLite["neighbor_ip"] != "" {
							vrfLiteMap["NEIGHBOR_ASN"] = vrfLite["neighbor_asn"].(string)
						} else {
							vrfLiteMap["NEIGHBOR_ASN"] = extensionValues["NEIGHBOR_ASN"].(string)
						}

						if vrfLite["ipv6_mask"] != "" {
							vrfLiteMap["IPV6_MASK"] = vrfLite["ipv6_mask"].(string)
						} else {
							vrfLiteMap["IPV6_MASK"] = extensionValues["IPV6_MASK"].(string)
						}

						if vrfLite["ipv6_neighbor"] != "" {
							vrfLiteMap["IPV6_NEIGHBOR"] = vrfLite["ipv6_neighbor"].(string)
						} else {
							vrfLiteMap["IPV6_NEIGHBOR"] = extensionValues["IPV6_NEIGHBOR"].(string)
						}

						if vrfLite["auto_vrf_lite_flag"] != "" {
							vrfLiteMap["AUTO_VRF_LITE_FLAG"] = vrfLite["auto_vrf_lite_flag"].(string)
						} else {
							vrfLiteMap["AUTO_VRF_LITE_FLAG"] = extensionValues["AUTO_VRF_LITE_FLAG"].(string)
						}
						vrfLiteMap["IF_NAME"] = extensionValues["IF_NAME"].(string)

						vrfLiteList = append(vrfLiteList, vrfLiteMap)
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

func resourceDCNMVRFUpdate(d *schema.ResourceData, m interface{}) error {
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

	configMap := models.VRFProfileConfig{}
	if vlan, ok := d.GetOk("vlan_id"); ok {
		configMap.Vlan = vlan.(int)
	}
	if mtu, ok := d.GetOk("mtu"); ok {
		configMap.Mtu = mtu.(int)
	}
	if vlanName, ok := d.GetOk("vlan_name"); ok {
		configMap.VlanName = vlanName.(string)
	}
	if desc, ok := d.GetOk("description"); ok {
		configMap.Description = desc.(string)
	}
	if intfDesc, ok := d.GetOk("intf_description"); ok {
		configMap.IntfDescription = intfDesc.(string)
	}
	if tag, ok := d.GetOk("tag"); ok {
		configMap.Tag = tag.(string)
	}
	if bgp, ok := d.GetOk("max_bgp_path"); ok {
		configMap.BGP = bgp.(int)
	}
	if ibgp, ok := d.GetOk("max_ibgp_path"); ok {
		configMap.IBGP = ibgp.(int)
	}
	if trm, ok := d.GetOk("trm_enable"); ok {
		configMap.TRM = trm.(string)
	}
	if rpExtr, ok := d.GetOk("rp_external_flag"); ok {
		configMap.RPexternal = rpExtr.(string)
	}
	if rpAddr, ok := d.GetOk("rp_address"); ok {
		configMap.RPaddress = rpAddr.(string)
	}
	if loopback, ok := d.GetOk("loopback_id"); ok {
		configMap.Lookback = loopback.(int)
	}
	if mcastAddr, ok := d.GetOk("mutlicast_address"); ok {
		configMap.Mcastaddr = mcastAddr.(string)
	}
	if mcastGrp, ok := d.GetOk("mutlicast_group"); ok {
		configMap.Mcastgroup = mcastGrp.(string)
	}
	if ipv6, ok := d.GetOk("ipv6_link_local_flag"); ok {
		configMap.IPv6Link = ipv6.(string)
	}
	if trmbgw, ok := d.GetOk("trm_bgw_msite_flag"); ok {
		configMap.TRMBGW = trmbgw.(string)
	}
	if hostR, ok := d.GetOk("advertise_host_route"); ok {
		configMap.AdhostRoute = hostR.(string)
	}
	if defaultR, ok := d.GetOk("advertise_default_route"); ok {
		configMap.AdDefaultRoute = defaultR.(string)
	}
	if staticR, ok := d.GetOk("static_default_route"); ok {
		configMap.StaticRoute = staticR.(string)
	}
	configMap.SegmentID = vrf.Id
	configMap.VrfName = vrf.Name

	confStr, err := json.Marshal(configMap)
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

				attachMap["fabric"] = vrf.Fabric
				attachMap["vrfName"] = vrf.Name
				attachMap["deployment"] = attachment["attach"].(bool)
				attachMap["serialNumber"] = attachment["serial_number"].(string)

				if attachment["vlan_id"].(int) != 0 {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				} else {
					attachMap["vlan"] = configMap.Vlan
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
					for _, val := range attachment["vrf_lite"].(*schema.Set).List() {
						log.Println("vrf_lite enter")
						vrfLite := val.(map[string]interface{})
						vrfLiteMap := make(map[string]interface{})

						durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/switches?vrf-names=%s&serial-numbers=%s", vrf.Fabric, vrf.Name, attachMap["serialNumber"].(string))
						cont, err := dcnmClient.GetviaURL(durl)
						if err != nil {
							return err
						}
						extensionProtValues := cont.Index(0).S("switchDetailsList").Index(0).S("extensionPrototypeValues")
						// if len(extensionProtValues) < 0 {
						// 	return fmt.Errorf("")
						// }
						ifName := stripQuotes(extensionProtValues.S("interfaceName").String())
						extensionValueString := stripQuotes(extensionProtValues.S("extensionValues").String())

						var extensionValues map[string]interface{}

						extensionValueString = strings.Replace(extensionValueString, "\\", "", -1)
						err = json.Unmarshal([]byte(extensionValueString), &extensionValues)
						if err != nil {
							return err
						}
						if len(extensionValues) == 0 {
							return fmt.Errorf("No VRF_LITE Data found for switch %s", attachMap["serialNumber"].(string))
						}
						vrfLiteMap["PEER_VRF_NAME"] = vrfLite["peer_vrf_name"]

						if vrfLite["dot1q_id"] != nil {
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

						log.Printf("extvals: %v", extensionValues)

						if vrfLite["ip_mask"] != "" {
							vrfLiteMap["IP_MASK"] = vrfLite["ip_mask"].(string)
						} else {
							vrfLiteMap["IP_MASK"] = extensionValues["IP_MASK"].(string)
						}
						log.Printf("extvals: %s", extensionValues["IP_MASK"])
						log.Printf("vrfLiteMap: %s", vrfLiteMap["IP_MASK"])
						log.Printf("vrfLite[\"ip_mask\"]: %v\n", vrfLite["ip_mask"])

						if vrfLite["neighbor_ip"] != "" {
							vrfLiteMap["NEIGHBOR_IP"] = vrfLite["neighbor_ip"].(string)
						} else {
							vrfLiteMap["NEIGHBOR_IP"] = extensionValues["NEIGHBOR_IP"].(string)
						}

						if vrfLite["neighbor_ip"] != "" {
							vrfLiteMap["NEIGHBOR_ASN"] = vrfLite["neighbor_asn"].(string)
						} else {
							vrfLiteMap["NEIGHBOR_ASN"] = extensionValues["NEIGHBOR_ASN"].(string)
						}

						if vrfLite["ipv6_mask"] != "" {
							vrfLiteMap["IPV6_MASK"] = vrfLite["ipv6_mask"].(string)
						} else {
							vrfLiteMap["IPV6_MASK"] = extensionValues["IPV6_MASK"].(string)
						}

						if vrfLite["ipv6_neighbor"] != "" {
							vrfLiteMap["IPV6_NEIGHBOR"] = vrfLite["ipv6_neighbor"].(string)
						} else {
							vrfLiteMap["IPV6_NEIGHBOR"] = extensionValues["IPV6_NEIGHBOR"].(string)
						}

						if vrfLite["auto_vrf_lite_flag"] != "" {
							vrfLiteMap["AUTO_VRF_LITE_FLAG"] = vrfLite["auto_vrf_lite_flag"].(string)
						} else {
							vrfLiteMap["AUTO_VRF_LITE_FLAG"] = extensionValues["AUTO_VRF_LITE_FLAG"].(string)
						}
						vrfLiteMap["IF_NAME"] = extensionValues["IF_NAME"].(string)

						vrfLiteList = append(vrfLiteList, vrfLiteMap)
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

func resourceDCNMVRFRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	fabricName := d.Get("fabric_name").(string)

	cont, err := getRemoteVRF(dcnmClient, fabricName, dn)
	if err != nil {
		return err
	}

	setVRFAttributes(d, cont)

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
				for _, val := range lites {
					vrfLiteMap := val.(map[string]interface{})

					durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/switches?vrf-names=%s&serial-numbers=%s", fabricName, dn, attachMap["serial_number"].(string))
					cont, err := dcnmClient.GetviaURL(durl)
					if err != nil {
						return err
					}
					extensionProtValues := cont.Index(0).S("switchDetailsList").Index(0).S("extensionPrototypeValues").Index(0)
					extensionValueString := stripQuotes(extensionProtValues.S("extensionValues").String())
					var extensionValues map[string]interface{}

					extensionValueString = strings.Replace(extensionValueString, "\\", "", -1)
					err = json.Unmarshal([]byte(extensionValueString), &extensionValues)
					if err != nil {
						return err
					}

					if len(extensionValues) == 0 {
						d.SetId("")
						return nil
					}

					//vrfLiteMap["peer_vrf_name"] = extensionValues["PEER_VRF_NAME"].(string)
					vrfLiteMap["dot1q_id"] = extensionValues["DOT1Q_ID"].(string)
					vrfLiteMap["neighbor_ip"] = extensionValues["NEIGHBOR_IP"].(string)
					vrfLiteMap["neighbor_asn"] = extensionValues["NEIGHBOR_ASN"].(string)
					vrfLiteMap["ipv6_mask"] = extensionValues["IPV6_MASK"].(string)
					vrfLiteMap["ipv6_neighbor"] = extensionValues["IPV6_NEIGHBOR"].(string)
					vrfLiteMap["auto_vrf_lite_flag"] = extensionValues["AUTO_VRF_LITE_FLAG"].(string)

					liteGet = append(liteGet, vrfLiteMap)
				}
				d.Set("attachments", attachGet)
			}

			attachGet = append(attachGet, attachMap)
		}

		d.Set("attachments", attachGet)
	}

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMVRFDelete(d *schema.ResourceData, m interface{}) error {
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

				attachMap["fabric"] = fabricName
				attachMap["vrfName"] = dn
				attachMap["deployment"] = false
				attachMap["serialNumber"] = attachment["serial_number"].(string)
				if attachment["vlan_id"].(int) == 0 {
					attachMap["vlan"] = d.Get("vlan_id").(int)
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
			for j := 0; j < (deployTimeout / 5); j++ {
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

func checkvrfDeploy(client *client.Client, fabric, vrf string) (bool, error) {
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/attachments?vrf-names=%s", fabric, vrf)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return false, err
	}

	attachList := cont.Index(0).S("lanAttachList")

	flag := false
	for i := 0; i < len(attachList.Data().([]interface{})); i++ {
		if stripQuotes(attachList.Index(i).S("lanAttachState").String()) == "DEPLOYED" {
			flag = true
			break
		}
	}

	return flag, nil
}

func getSwitchAttachStatus(client *client.Client, fabric, vrf, switchNum string) (bool, int, error) {
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/attachments?vrf-names=%s", fabric, vrf)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return false, 0, err
	}

	attachList := cont.Index(0).S("lanAttachList")

	for i := 0; i < len(attachList.Data().([]interface{})); i++ {
		if stripQuotes(attachList.Index(i).S("switchSerialNo").String()) == switchNum {
			if stripQuotes(attachList.Index(i).S("isLanAttached").String()) == "true" {
				if stripQuotes(attachList.Index(i).S("vlanId").String()) != "null" {
					vlan, err := strconv.Atoi(stripQuotes(attachList.Index(i).S("vlanId").String()))
					if err == nil {
						return true, vlan, nil
					}
				}
				return true, 0, nil
			}
			return false, 0, nil
		}
	}
	return false, 0, nil
}

func getVRFDeploymentStatus(client *client.Client, fabricName, vrfName string) (string, error) {

	dURL := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs", fabricName)
	vrfCont, err := client.GetviaURL(dURL)
	if err != nil {
		return "", err
	}

	status := ""
	for i := 0; i < len(vrfCont.Data().([]interface{})); i++ {
		if stripQuotes(vrfCont.Index(i).S("vrfName").String()) == vrfName {
			status = stripQuotes(vrfCont.Index(i).S("vrfStatus").String())
			break
		}
	}

	return status, nil
}
