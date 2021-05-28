package dcnm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDCNMInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMInterfaceCreate,
		Update: resourceDCNMInterfaceUpdate,
		Read:   resourceDCNMInterfaceRead,
		Delete: resourceDCNMInterfaceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDCNMInterfaceImporter,
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
				ForceNew: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port-channel",
					"vpc",
					"loopback",
					"sub-interface",
					"ethernet",
				}, false),
			},

			"policy": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"switch_name_1": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"switch_name_2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vrf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_routing_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_ls_routing": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_router_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_replication_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer1_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"vpc_peer2_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"vpc_peer1_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"vpc_peer2_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"on",
					"active",
					"passive",
				}, false),
			},

			"bpdu_gaurd_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true",
					"false",
					"no",
				}, false),
			},

			"port_fast_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer1_allowed_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer2_allowed_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer1_access_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer2_access_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer1_desc": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer2_desc": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer1_conf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vpc_peer2_conf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pc_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"access_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"allowed_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subinterface_vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"ipv4_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"subinterface_mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ethernet_speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Auto",
					"100Mb",
					"1Gb",
					"10Gb",
					"25Gb",
					"40Gb",
					"100Gb",
				}, false),
			},

			"serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"configuration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"admin_state": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func getRemoteInterface(client *client.Client, serialNum, name string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/interface?serialNumber=%s&ifName=%s", serialNum, name)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return cont, err
	}
	return cont, nil
}

func setInterfaceAttributes(d *schema.ResourceData, cont *container.Container, intftype string) *schema.ResourceData {
	d.Set("policy", stripQuotes(cont.S("policy").String()))

	interfaces := cont.S("interfaces").Index(0)
	d.Set("serial_number", stripQuotes(interfaces.S("serialNumber").String()))
	d.Set("type", intftype)
	// d.Set("fabric_name", stripQuotes(interfaces.S("nvPairs", "FABRIC_NAME").String()))
	d.Set("name", (stripQuotes(interfaces.S("ifName").String())))
	if state, err := strconv.ParseBool(stripQuotes(interfaces.S("nvPairs", "ADMIN_STATE").String())); err == nil {
		d.Set("admin_state", state)
	}

	if intftype == "loopback" {

		d.Set("ipv4", stripQuotes(interfaces.S("nvPairs", "IP").String()))
		d.Set("ipv6", stripQuotes(interfaces.S("nvPairs", "V6IP").String()))
		d.Set("loopback_tag", stripQuotes(interfaces.S("nvPairs", "ROUTE_MAP_TAG").String()))
		d.Set("loopback_router_id", stripQuotes(interfaces.S("nvPairs", "routerId").String()))
		d.Set("loopback_ls_routing", stripQuotes(interfaces.S("nvPairs", "LINK_STATE_ROUTING").String()))
		d.Set("loopback_routing_tag", stripQuotes(interfaces.S("nvPairs", "ROUTING_TAG").String()))
		d.Set("loopback_replication_mode", stripQuotes(interfaces.S("nvPairs", "REPLICATION_MODE").String()))
		d.Set("configuration", stripQuotes(interfaces.S("nvPairs", "CONF").String()))
		d.Set("description", stripQuotes(interfaces.S("nvPairs", "DESC").String()))
		d.Set("vrf", stripQuotes(interfaces.S("nvPairs", "INTF_VRF").String()))

		d.Set("pc_interface", make([]interface{}, 0, 1))
		d.Set("vpc_peer1_interface", make([]interface{}, 0, 1))
		d.Set("vpc_peer2_interface", make([]interface{}, 0, 1))

	} else if intftype == "vpc" {

		if p1ID, err := strconv.Atoi(stripQuotes(interfaces.S("nvPairs", "PEER1_PCID").String())); err == nil {
			d.Set("vpc_peer1_id", p1ID)
		}
		if p2ID, err := strconv.Atoi(stripQuotes(interfaces.S("nvPairs", "PEER2_PCID").String())); err == nil {
			d.Set("vpc_peer2_id", p2ID)
		}
		p1intfAct := interfaceToStrList(d.Get("vpc_peer1_interface"))
		p2intfAct := interfaceToStrList(d.Get("vpc_peer2_interface"))
		p1intfGet := stringToList(stripQuotes(interfaces.S("nvPairs", "PEER1_MEMBER_INTERFACES").String()))
		p2intfGet := stringToList(stripQuotes(interfaces.S("nvPairs", "PEER2_MEMBER_INTERFACES").String()))
		if compareStrLists(p1intfAct, p1intfGet) {
			d.Set("vpc_peer1_interface", d.Get("vpc_peer1_interface"))
		} else {
			d.Set("vpc_peer1_interface", stringToList(stripQuotes(interfaces.S("nvPairs", "PEER1_MEMBER_INTERFACES").String())))
		}
		if compareStrLists(p2intfAct, p2intfGet) {
			d.Set("vpc_peer2_interface", d.Get("vpc_peer2_interface"))
		} else {
			d.Set("vpc_peer2_interface", stringToList(stripQuotes(interfaces.S("nvPairs", "PEER2_MEMBER_INTERFACES").String())))
		}
		d.Set("mode", stripQuotes(interfaces.S("nvPairs", "PC_MODE").String()))
		d.Set("bpdu_gaurd_flag", stripQuotes(interfaces.S("nvPairs", "BPDUGUARD_ENABLED").String()))
		if ppf, err := strconv.ParseBool(stripQuotes(interfaces.S("nvPairs", "PORTTYPE_FAST_ENABLED").String())); err == nil {
			d.Set("port_fast_flag", ppf)
		}
		d.Set("mtu", stripQuotes(interfaces.S("nvPairs", "MTU").String()))
		d.Set("vpc_peer1_allowed_vlans", stripQuotes(interfaces.S("nvPairs", "PEER1_ALLOWED_VLANS").String()))
		d.Set("vpc_peer2_allowed_vlans", stripQuotes(interfaces.S("nvPairs", "PEER2_ALLOWED_VLANS").String()))
		d.Set("vpc_peer1_access_vlans", stripQuotes(interfaces.S("nvPairs", "PEER1_ACCESS_VLAN").String()))
		d.Set("vpc_peer2_access_vlans", stripQuotes(interfaces.S("nvPairs", "PEER2_ACCESS_VLAN").String()))
		d.Set("vpc_peer1_desc", stripQuotes(interfaces.S("nvPairs", "PEER1_PO_DESC").String()))
		d.Set("vpc_peer2_desc", stripQuotes(interfaces.S("nvPairs", "PEER2_PO_DESC").String()))
		d.Set("vpc_peer1_conf", stripQuotes(interfaces.S("nvPairs", "PEER1_PO_CONF").String()))
		d.Set("vpc_peer2_conf", stripQuotes(interfaces.S("nvPairs", "PEER2_PO_CONF").String()))

		d.Set("pc_interface", make([]interface{}, 0, 1))

	} else if intftype == "port-channel" {

		pcIntfAcc := interfaceToStrList(d.Get("pc_interface"))
		pcIntfGet := stringToList(stripQuotes(interfaces.S("nvPairs", "MEMBER_INTERFACES").String()))
		if compareStrLists(pcIntfAcc, pcIntfGet) {
			d.Set("pc_interface", d.Get("pc_interface"))
		} else {
			d.Set("pc_interface", stringToList(stripQuotes(interfaces.S("nvPairs", "MEMBER_INTERFACES").String())))
		}
		d.Set("mode", stripQuotes(interfaces.S("nvPairs", "PC_MODE").String()))
		d.Set("bpdu_gaurd_flag", stripQuotes(interfaces.S("nvPairs", "BPDUGUARD_ENABLED").String()))
		if ppf, err := strconv.ParseBool(stripQuotes(interfaces.S("nvPairs", "PORTTYPE_FAST_ENABLED").String())); err == nil {
			d.Set("port_fast_flag", ppf)
		}
		d.Set("mtu", stripQuotes(interfaces.S("nvPairs", "MTU").String()))
		d.Set("allowed_vlans", stripQuotes(interfaces.S("nvPairs", "ALLOWED_VLANS").String()))
		d.Set("access_vlans", stripQuotes(interfaces.S("nvPairs", "ACCESS_VLAN").String()))
		d.Set("configuration", stripQuotes(interfaces.S("nvPairs", "CONF").String()))
		d.Set("description", stripQuotes(interfaces.S("nvPairs", "DESC").String()))

		d.Set("vpc_peer1_interface", make([]interface{}, 0, 1))
		d.Set("vpc_peer2_interface", make([]interface{}, 0, 1))

	} else if intftype == "sub-interface" {

		if vlan, err := strconv.Atoi(stripQuotes(interfaces.S("nvPairs", "VLAN").String())); err == nil {
			d.Set("subinterface_vlan", vlan)
		}
		d.Set("vrf", stripQuotes(interfaces.S("nvPairs", "INTF_VRF").String()))
		d.Set("ipv4", stripQuotes(interfaces.S("nvPairs", "IP").String()))
		d.Set("ipv6", stripQuotes(interfaces.S("nvPairs", "IPv6").String()))
		d.Set("ipv6_prefix", stripQuotes(interfaces.S("nvPairs", "IPv6_PREFIX").String()))
		d.Set("ipv4_prefix", stripQuotes(interfaces.S("nvPairs", "PREFIX").String()))
		d.Set("subinterface_mtu", stripQuotes(interfaces.S("nvPairs", "MTU").String()))
		d.Set("configuration", stripQuotes(interfaces.S("nvPairs", "CONF").String()))
		d.Set("description", stripQuotes(interfaces.S("nvPairs", "DESC").String()))

		d.Set("pc_interface", make([]interface{}, 0, 1))
		d.Set("vpc_peer1_interface", make([]interface{}, 0, 1))
		d.Set("vpc_peer2_interface", make([]interface{}, 0, 1))

	} else if intftype == "ethernet" {

		d.Set("bpdu_gaurd_flag", stripQuotes(interfaces.S("nvPairs", "BPDUGUARD_ENABLED").String()))
		if ppf, err := strconv.ParseBool(stripQuotes(interfaces.S("nvPairs", "PORTTYPE_FAST_ENABLED").String())); err == nil {
			d.Set("port_fast_flag", ppf)
		}
		d.Set("mtu", stripQuotes(interfaces.S("nvPairs", "MTU").String()))
		d.Set("ipv4", stripQuotes(interfaces.S("nvPairs", "IP").String()))
		d.Set("ipv6", stripQuotes(interfaces.S("nvPairs", "IPv6").String()))
		d.Set("access_vlans", stripQuotes(interfaces.S("nvPairs", "ACCESS_VLAN").String()))
		d.Set("ipv6_prefix", stripQuotes(interfaces.S("nvPairs", "IPv6_PREFIX").String()))
		d.Set("ipv4_prefix", stripQuotes(interfaces.S("nvPairs", "PREFIX").String()))
		d.Set("ethernet_speed", stripQuotes(interfaces.S("nvPairs", "SPEED").String()))
		d.Set("allowed_vlans", stripQuotes(interfaces.S("nvPairs", "ALLOWED_VLANS").String()))
		d.Set("configuration", stripQuotes(interfaces.S("nvPairs", "CONF").String()))
		d.Set("description", stripQuotes(interfaces.S("nvPairs", "DESC").String()))
		d.Set("vrf", stripQuotes(interfaces.S("nvPairs", "INTF_VRF").String()))

		d.Set("pc_interface", make([]interface{}, 0, 1))
		d.Set("vpc_peer1_interface", make([]interface{}, 0, 1))
		d.Set("vpc_peer2_interface", make([]interface{}, 0, 1))

	}

	return d
}

func resourceDCNMInterfaceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Importer ", d.Id())

	dcnmClient := m.(*client.Client)

	var serialNum1 string
	var serialNum2 string
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 4 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	intfType := importInfo[0]
	serialNum := importInfo[1]
	name := importInfo[2]
	fabricName := importInfo[3]
	if intfType == "vpc" {
		vpcSerialNums := strings.Split(serialNum, "~")
		if len(vpcSerialNums) != 2 {
			return nil, fmt.Errorf("serial number is not valid for vpc interface")
		}
		serialNum1 = vpcSerialNums[0]
		serialNum2 = vpcSerialNums[1]
	} else {
		serialNum1 = serialNum
	}

	cont, err := getRemoteInterface(dcnmClient, serialNum1, name)
	if err != nil {
		errorMsg, flag := checkIntfErrors(cont)
		if flag {
			return nil, fmt.Errorf(errorMsg)
		}
	}

	setInterfaceAttributes(d, cont.Index(0), intfType)
	d.SetId(name)
	d.Set("fabric_name", fabricName)

	if intfType == "vpc" {
		swName1, err := getSwitchName(dcnmClient, fabricName, serialNum1)
		if err == nil {
			d.Set("switch_name_1", swName1)
		}
		swName2, err := getSwitchName(dcnmClient, fabricName, serialNum2)
		if err == nil {
			d.Set("switch_name_2", swName2)
		}
	} else {
		swName1, err := getSwitchName(dcnmClient, fabricName, serialNum1)
		if err == nil {
			d.Set("switch_name_1", swName1)
		}
	}

	flag, err := checkIntfDeploy(dcnmClient, serialNum, d.Get("name").(string), intfType)
	if err != nil {
		return nil, err
	}
	d.Set("deploy", flag)

	importState := d

	log.Println("[DEBUG] End of Importer ", d.Id())
	return []*schema.ResourceData{importState}, nil
}

func resourceDCNMInterfaceCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)
	policy := d.Get("policy").(string)
	name := d.Get("name").(string)
	intfType := d.Get("type").(string)

	intf := models.Interface{}
	intf.Policy = policy

	intfConfig := models.InterfaceConfig{}
	intfConfig.Fabric = fabricName
	intfConfig.InterfaceName = name

	switch1 := d.Get("switch_name_1")
	switchCont, err := getRemoteSwitchforDS(dcnmClient, fabricName, switch1.(string))
	if err != nil {
		return err
	}
	serial1 := stripQuotes(switchCont.S("serialNumber").String())

	nvPairMap := make(map[string]interface{})
	nvPairMap["INTF_NAME"] = name

	if intfType == "loopback" {
		intf.Type = "INTERFACE_LOOPBACK"
		intfConfig.InterfaceType = "INTERFACE_LOOPBACK"
		intfConfig.SerialNumber = serial1

		if vrf, ok := d.GetOk("vrf"); ok {
			nvPairMap["INTF_VRF"] = vrf.(string)
		} else {
			nvPairMap["INTF_VRF"] = ""
		}
		if ip, ok := d.GetOk("ipv4"); ok {
			nvPairMap["IP"] = ip.(string)
		} else {
			nvPairMap["IP"] = ""
		}
		if ipv6, ok := d.GetOk("ipv6"); ok {
			nvPairMap["V6IP"] = ipv6.(string)
		} else {
			nvPairMap["V6IP"] = ""
		}
		if tag, ok := d.GetOk("loopback_tag"); ok {
			nvPairMap["ROUTE_MAP_TAG"] = tag.(string)
		} else {
			nvPairMap["ROUTE_MAP_TAG"] = ""
		}
		if tag, ok := d.GetOk("loopback_router_id"); ok {
			nvPairMap["routerId"] = tag.(string)
		} else {
			nvPairMap["routerId"] = ""
		}
		if tag, ok := d.GetOk("loopback_ls_routing"); ok {
			nvPairMap["LINK_STATE_ROUTING"] = tag.(string)
		} else {
			nvPairMap["LINK_STATE_ROUTING"] = ""
		}
		if tag, ok := d.GetOk("loopback_routing_tag"); ok {
			nvPairMap["ROUTING_TAG"] = tag.(string)
		} else {
			nvPairMap["ROUTING_TAG"] = ""
		}
		if tag, ok := d.GetOk("loopback_replication_mode"); ok {
			nvPairMap["REPLICATION_MODE"] = tag.(string)
		} else {
			nvPairMap["REPLICATION_MODE"] = ""
		}
		if desc, ok := d.GetOk("description"); ok {
			nvPairMap["DESC"] = desc.(string)
		} else {
			nvPairMap["DESC"] = ""
		}
		if conf, ok := d.GetOk("configuration"); ok {
			nvPairMap["CONF"] = conf.(string)
		} else {
			nvPairMap["CONF"] = ""
		}

	} else if intfType == "vpc" {
		var serial2 string
		if switch2, ok := d.GetOk("switch_name_2"); ok {
			switchCont, err := getRemoteSwitchforDS(dcnmClient, fabricName, switch2.(string))
			if err != nil {
				return err
			}
			serial2 = stripQuotes(switchCont.S("serialNumber").String())
		} else {
			return fmt.Errorf("switch_name_2 field is required for vpc interface")
		}

		intf.Type = "INTERFACE_VPC"
		intfConfig.InterfaceType = "INTERFACE_VPC"
		intfConfig.SerialNumber = fmt.Sprintf("%s~%s", serial1, serial2)

		if p1ID, ok := d.GetOk("vpc_peer1_id"); ok {
			nvPairMap["PEER1_PCID"] = p1ID.(int)
		} else {
			nvPairMap["PEER1_PCID"] = ""
		}
		if p2ID, ok := d.GetOk("vpc_peer2_id"); ok {
			nvPairMap["PEER2_PCID"] = p2ID.(int)
		} else {
			nvPairMap["PEER2_PCID"] = ""
		}
		if p1intf, ok := d.GetOk("vpc_peer1_interface"); ok {
			nvPairMap["PEER1_MEMBER_INTERFACES"] = listToString(p1intf)
		} else {
			nvPairMap["PEER1_MEMBER_INTERFACES"] = ""
		}
		if p2intf, ok := d.GetOk("vpc_peer2_interface"); ok {
			nvPairMap["PEER2_MEMBER_INTERFACES"] = listToString(p2intf)
		} else {
			nvPairMap["PEER2_MEMBER_INTERFACES"] = ""
		}
		if mode, ok := d.GetOk("mode"); ok {
			nvPairMap["PC_MODE"] = mode.(string)
		} else {
			nvPairMap["PC_MODE"] = ""
		}
		if bpduF, ok := d.GetOk("bpdu_gaurd_flag"); ok {
			nvPairMap["BPDUGUARD_ENABLED"] = bpduF.(string)
		} else {
			nvPairMap["BPDUGUARD_ENABLED"] = ""
		}
		if pff, ok := d.GetOk("port_fast_flag"); ok {
			nvPairMap["PORTTYPE_FAST_ENABLED"] = pff.(bool)
		}
		if mtu, ok := d.GetOk("mtu"); ok {
			nvPairMap["MTU"] = mtu.(string)
		} else {
			nvPairMap["MTU"] = ""
		}
		if p1vlan, ok := d.GetOk("vpc_peer1_allowed_vlans"); ok {
			nvPairMap["PEER1_ALLOWED_VLANS"] = p1vlan.(string)
		} else {
			nvPairMap["PEER1_ALLOWED_VLANS"] = ""
		}
		if p2vlan, ok := d.GetOk("vpc_peer2_allowed_vlans"); ok {
			nvPairMap["PEER2_ALLOWED_VLANS"] = p2vlan.(string)
		} else {
			nvPairMap["PEER2_ALLOWED_VLANS"] = ""
		}
		if p1Avlan, ok := d.GetOk("vpc_peer1_access_vlans"); ok {
			nvPairMap["PEER1_ACCESS_VLAN"] = p1Avlan.(string)
		} else {
			nvPairMap["PEER1_ACCESS_VLAN"] = ""
		}
		if p2Avlan, ok := d.GetOk("vpc_peer2_access_vlans"); ok {
			nvPairMap["PEER2_ACCESS_VLAN"] = p2Avlan.(string)
		} else {
			nvPairMap["PEER2_ACCESS_VLAN"] = ""
		}
		if p1desc, ok := d.GetOk("vpc_peer1_desc"); ok {
			nvPairMap["PEER1_PO_DESC"] = p1desc.(string)
		} else {
			nvPairMap["PEER1_PO_DESC"] = ""
		}
		if p2desc, ok := d.GetOk("vpc_peer2_desc"); ok {
			nvPairMap["PEER2_PO_DESC"] = p2desc.(string)
		} else {
			nvPairMap["PEER2_PO_DESC"] = ""
		}
		if p1conf, ok := d.GetOk("vpc_peer1_conf"); ok {
			nvPairMap["PEER1_PO_CONF"] = p1conf.(string)
		} else {
			nvPairMap["PEER1_PO_CONF"] = ""
		}
		if p2conf, ok := d.GetOk("vpc_peer2_conf"); ok {
			nvPairMap["PEER2_PO_CONF"] = p2conf.(string)
		} else {
			nvPairMap["PEER2_PO_CONF"] = ""
		}

	} else if intfType == "port-channel" {
		intf.Type = "INTERFACE_PORT_CHANNEL"
		intfConfig.InterfaceType = "INTERFACE_PORT_CHANNEL"
		intfConfig.SerialNumber = serial1

		nvPairMap["PO_ID"] = name
		if intf, ok := d.GetOk("pc_interface"); ok {
			nvPairMap["MEMBER_INTERFACES"] = listToString(intf)
		} else {
			nvPairMap["MEMBER_INTERFACES"] = ""
		}
		if mode, ok := d.GetOk("mode"); ok {
			nvPairMap["PC_MODE"] = mode.(string)
		} else {
			nvPairMap["PC_MODE"] = ""
		}
		if bpduF, ok := d.GetOk("bpdu_gaurd_flag"); ok {
			nvPairMap["BPDUGUARD_ENABLED"] = bpduF.(string)
		} else {
			nvPairMap["BPDUGUARD_ENABLED"] = ""
		}
		if pff, ok := d.GetOk("port_fast_flag"); ok {
			nvPairMap["PORTTYPE_FAST_ENABLED"] = pff.(bool)
		}
		if mtu, ok := d.GetOk("mtu"); ok {
			nvPairMap["MTU"] = mtu.(string)
		} else {
			nvPairMap["MTU"] = ""
		}
		if vlans, ok := d.GetOk("allowed_vlans"); ok {
			nvPairMap["ALLOWED_VLANS"] = vlans.(string)
		} else {
			nvPairMap["ALLOWED_VLANS"] = ""
		}
		if accVlans, ok := d.GetOk("access_vlans"); ok {
			nvPairMap["ACCESS_VLAN"] = accVlans.(string)
		} else {
			nvPairMap["ACCESS_VLAN"] = ""
		}
		if desc, ok := d.GetOk("description"); ok {
			nvPairMap["DESC"] = desc.(string)
		} else {
			nvPairMap["DESC"] = ""
		}
		if conf, ok := d.GetOk("configuration"); ok {
			nvPairMap["CONF"] = conf.(string)
		} else {
			nvPairMap["CONF"] = ""
		}

	} else if intfType == "sub-interface" {
		intf.Type = "SUBINTERFACE"
		intfConfig.InterfaceType = "SUBINTERFACE"
		intfConfig.SerialNumber = serial1

		if vlan, ok := d.GetOk("subinterface_vlan"); ok {
			nvPairMap["VLAN"] = vlan.(int)
		} else {
			nvPairMap["VLAN"] = ""
		}
		if vrf, ok := d.GetOk("vrf"); ok {
			nvPairMap["INTF_VRF"] = vrf.(string)
		} else {
			nvPairMap["INTF_VRF"] = ""
		}
		if ip, ok := d.GetOk("ipv4"); ok {
			nvPairMap["IP"] = ip.(string)
		} else {
			nvPairMap["IP"] = ""
		}
		if ipv4Pre, ok := d.GetOk("ipv4_prefix"); ok {
			nvPairMap["PREFIX"] = ipv4Pre.(string)
		} else {
			nvPairMap["PREFIX"] = ""
		}
		if ipv6, ok := d.GetOk("ipv6"); ok {
			nvPairMap["IPv6"] = ipv6.(string)
		} else {
			nvPairMap["IPv6"] = ""
		}
		if ipv6Pre, ok := d.GetOk("ipv6_prefix"); ok {
			nvPairMap["IPv6_PREFIX"] = ipv6Pre.(string)
		} else {
			nvPairMap["IPv6_PREFIX"] = ""
		}
		if mtu, ok := d.GetOk("subinterface_mtu"); ok {
			nvPairMap["MTU"] = mtu.(string)
		} else {
			nvPairMap["MTU"] = ""
		}
		if desc, ok := d.GetOk("description"); ok {
			nvPairMap["DESC"] = desc.(string)
		} else {
			nvPairMap["DESC"] = ""
		}
		if conf, ok := d.GetOk("configuration"); ok {
			nvPairMap["CONF"] = conf.(string)
		} else {
			nvPairMap["CONF"] = ""
		}

	} else if intfType == "ethernet" {
		return fmt.Errorf("Ethernet interface can only be modified")

	}

	if state, ok := d.GetOk("admin_state"); ok {
		nvPairMap["ADMIN_STATE"] = state.(bool)
	} else {
		nvPairMap["ADMIN_STATE"] = ""
	}

	intfModel := models.NewInterface(&intf, &intfConfig, nvPairMap)

	cont, err := dcnmClient.Save("/rest/interface", intfModel)
	if err != nil {
		if cont != nil {
			errorMsg, flag := checkIntfErrors(cont)
			if flag {
				return fmt.Errorf(errorMsg)
			}
		} else {
			return err
		}
	}

	d.Set("serial_number", intfConfig.SerialNumber)
	d.Set("type", intfType)
	d.SetId(intfConfig.InterfaceName)

	//Deployment of interface
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		log.Println("[DEBUG] Begining Deployment ", d.Id())

		intfDeploy := models.InterfaceDelete{}
		intfDeploy.SerialNumber = intfConfig.SerialNumber
		intfDeploy.Name = intfConfig.InterfaceName
		cont, err = dcnmClient.SaveForAttachment("/rest/interface/deploy", &intfDeploy)
		if err != nil {
			errorMsg, flag := checkIntfErrors(cont)
			if flag {
				d.Set("deploy", false)
				return fmt.Errorf("interface is created but failed to deploy with error : %s", errorMsg)
			}
		}

		log.Println("[DEBUG] End of Deployment ", d.Id())
	}

	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMInterfaceRead(d, m)
}

func resourceDCNMInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)
	policy := d.Get("policy").(string)
	name := d.Id()
	intfType := d.Get("type").(string)
	serialnum := d.Get("serial_number").(string)

	if d.HasChange("switch_name_1") {
		switch1Old, switch1New := d.GetChange("switch_name_1")
		switchCont, err := getRemoteSwitchforDS(dcnmClient, fabricName, switch1New.(string))
		if err != nil {
			return err
		}
		serial1 := stripQuotes(switchCont.S("serialNumber").String())
		if intfType == "vpc" && d.HasChange("switch_name_2") {
			switch2Old, switch2New := d.GetChange("switch_name_2")
			switchCont, err := getRemoteSwitchforDS(dcnmClient, fabricName, switch2New.(string))
			if err != nil {
				return err
			}
			serial2 := stripQuotes(switchCont.S("serialNumber").String())
			serial := fmt.Sprintf("%s~%s", serial1, serial2)
			if serial != serialnum {
				d.Set("switch_name_1", switch1Old)
				d.Set("switch_name_2", switch2Old)
				return fmt.Errorf("switch names should not be updated")
			}
		} else if serial1 != serialnum {
			d.Set("switch_name_1", switch1Old)
			return fmt.Errorf("switch names should not be updated")
		}
	}

	intf := models.Interface{}
	intf.Policy = policy

	intfConfig := models.InterfaceConfig{}
	intfConfig.Fabric = fabricName
	intfConfig.InterfaceName = name

	nvPairMap := make(map[string]interface{})
	nvPairMap["INTF_NAME"] = name

	if intfType == "loopback" {
		intfConfig.SerialNumber = serialnum

		if vrf, ok := d.GetOk("vrf"); ok {
			nvPairMap["INTF_VRF"] = vrf.(string)
		} else {
			nvPairMap["INTF_VRF"] = ""
		}
		if ip, ok := d.GetOk("ipv4"); ok {
			nvPairMap["IP"] = ip.(string)
		} else {
			nvPairMap["IP"] = ""
		}
		if ipv6, ok := d.GetOk("ipv6"); ok {
			nvPairMap["V6IP"] = ipv6.(string)
		} else {
			nvPairMap["V6IP"] = ""
		}
		if tag, ok := d.GetOk("loopback_tag"); ok {
			nvPairMap["ROUTE_MAP_TAG"] = tag.(string)
		} else {
			nvPairMap["ROUTE_MAP_TAG"] = ""
		}
		if tag, ok := d.GetOk("loopback_router_id"); ok {
			nvPairMap["routerId"] = tag.(string)
		} else {
			nvPairMap["routerId"] = ""
		}
		if tag, ok := d.GetOk("loopback_ls_routing"); ok {
			nvPairMap["LINK_STATE_ROUTING"] = tag.(string)
		} else {
			nvPairMap["LINK_STATE_ROUTING"] = ""
		}
		if tag, ok := d.GetOk("loopback_routing_tag"); ok {
			nvPairMap["ROUTING_TAG"] = tag.(string)
		} else {
			nvPairMap["ROUTING_TAG"] = ""
		}
		if tag, ok := d.GetOk("loopback_replication_mode"); ok {
			nvPairMap["REPLICATION_MODE"] = tag.(string)
		} else {
			nvPairMap["REPLICATION_MODE"] = ""
		}
		if desc, ok := d.GetOk("description"); ok {
			nvPairMap["DESC"] = desc.(string)
		} else {
			nvPairMap["DESC"] = ""
		}
		if conf, ok := d.GetOk("configuration"); ok {
			nvPairMap["CONF"] = conf.(string)
		} else {
			nvPairMap["CONF"] = ""
		}

	} else if intfType == "vpc" {
		intfConfig.SerialNumber = serialnum

		if p1ID, ok := d.GetOk("vpc_peer1_id"); ok {
			nvPairMap["PEER1_PCID"] = p1ID.(int)
		} else {
			nvPairMap["PEER1_PCID"] = ""
		}
		if p2ID, ok := d.GetOk("vpc_peer2_id"); ok {
			nvPairMap["PEER2_PCID"] = p2ID.(int)
		} else {
			nvPairMap["PEER2_PCID"] = ""
		}
		if p1intf, ok := d.GetOk("vpc_peer1_interface"); ok {
			nvPairMap["PEER1_MEMBER_INTERFACES"] = listToString(p1intf)
		} else {
			nvPairMap["PEER1_MEMBER_INTERFACES"] = ""
		}
		if p2intf, ok := d.GetOk("vpc_peer2_interface"); ok {
			nvPairMap["PEER2_MEMBER_INTERFACES"] = listToString(p2intf)
		} else {
			nvPairMap["PEER2_MEMBER_INTERFACES"] = ""
		}
		if mode, ok := d.GetOk("mode"); ok {
			nvPairMap["PC_MODE"] = mode.(string)
		} else {
			nvPairMap["PC_MODE"] = ""
		}
		if bpduF, ok := d.GetOk("bpdu_gaurd_flag"); ok {
			nvPairMap["BPDUGUARD_ENABLED"] = bpduF.(string)
		} else {
			nvPairMap["BPDUGUARD_ENABLED"] = ""
		}
		if pff, ok := d.GetOk("port_fast_flag"); ok {
			nvPairMap["PORTTYPE_FAST_ENABLED"] = pff.(bool)
		}
		if mtu, ok := d.GetOk("mtu"); ok {
			nvPairMap["MTU"] = mtu.(string)
		} else {
			nvPairMap["MTU"] = ""
		}
		if p1vlan, ok := d.GetOk("vpc_peer1_allowed_vlans"); ok {
			nvPairMap["PEER1_ALLOWED_VLANS"] = p1vlan.(string)
		} else {
			nvPairMap["PEER1_ALLOWED_VLANS"] = ""
		}
		if p2vlan, ok := d.GetOk("vpc_peer2_allowed_vlans"); ok {
			nvPairMap["PEER2_ALLOWED_VLANS"] = p2vlan.(string)
		} else {
			nvPairMap["PEER2_ALLOWED_VLANS"] = ""
		}
		if p1Avlan, ok := d.GetOk("vpc_peer1_access_vlans"); ok {
			nvPairMap["PEER1_ACCESS_VLAN"] = p1Avlan.(string)
		} else {
			nvPairMap["PEER1_ACCESS_VLAN"] = ""
		}
		if p2Avlan, ok := d.GetOk("vpc_peer2_access_vlans"); ok {
			nvPairMap["PEER2_ACCESS_VLAN"] = p2Avlan.(string)
		} else {
			nvPairMap["PEER2_ACCESS_VLAN"] = ""
		}
		if p1desc, ok := d.GetOk("vpc_peer1_desc"); ok {
			nvPairMap["PEER1_PO_DESC"] = p1desc.(string)
		} else {
			nvPairMap["PEER1_PO_DESC"] = ""
		}
		if p2desc, ok := d.GetOk("vpc_peer2_desc"); ok {
			nvPairMap["PEER2_PO_DESC"] = p2desc.(string)
		} else {
			nvPairMap["PEER2_PO_DESC"] = ""
		}
		if p1conf, ok := d.GetOk("vpc_peer1_conf"); ok {
			nvPairMap["PEER1_PO_CONF"] = p1conf.(string)
		} else {
			nvPairMap["PEER1_PO_CONF"] = ""
		}
		if p2conf, ok := d.GetOk("vpc_peer2_conf"); ok {
			nvPairMap["PEER2_PO_CONF"] = p2conf.(string)
		} else {
			nvPairMap["PEER2_PO_CONF"] = ""
		}

	} else if intfType == "port-channel" {
		intfConfig.SerialNumber = serialnum

		nvPairMap["PO_ID"] = name
		if intf, ok := d.GetOk("pc_interface"); ok {
			nvPairMap["MEMBER_INTERFACES"] = listToString(intf)
		} else {
			nvPairMap["MEMBER_INTERFACES"] = ""
		}
		if mode, ok := d.GetOk("mode"); ok {
			nvPairMap["PC_MODE"] = mode.(string)
		} else {
			nvPairMap["PC_MODE"] = ""
		}
		if bpduF, ok := d.GetOk("bpdu_gaurd_flag"); ok {
			nvPairMap["BPDUGUARD_ENABLED"] = bpduF.(string)
		} else {
			nvPairMap["BPDUGUARD_ENABLED"] = ""
		}
		if pff, ok := d.GetOk("port_fast_flag"); ok {
			nvPairMap["PORTTYPE_FAST_ENABLED"] = pff.(bool)
		}
		if mtu, ok := d.GetOk("mtu"); ok {
			nvPairMap["MTU"] = mtu.(string)
		} else {
			nvPairMap["MTU"] = ""
		}
		if vlans, ok := d.GetOk("allowed_vlans"); ok {
			nvPairMap["ALLOWED_VLANS"] = vlans.(string)
		} else {
			nvPairMap["ALLOWED_VLANS"] = ""
		}
		if accVlans, ok := d.GetOk("access_vlans"); ok {
			nvPairMap["ACCESS_VLAN"] = accVlans.(string)
		} else {
			nvPairMap["ACCESS_VLAN"] = ""
		}
		if desc, ok := d.GetOk("description"); ok {
			nvPairMap["DESC"] = desc.(string)
		} else {
			nvPairMap["DESC"] = ""
		}
		if conf, ok := d.GetOk("configuration"); ok {
			nvPairMap["CONF"] = conf.(string)
		} else {
			nvPairMap["CONF"] = ""
		}

	} else if intfType == "sub-interface" {
		intfConfig.SerialNumber = serialnum

		if vlan, ok := d.GetOk("subinterface_vlan"); ok {
			nvPairMap["VLAN"] = vlan.(int)
		} else {
			nvPairMap["VLAN"] = ""
		}
		if vrf, ok := d.GetOk("vrf"); ok {
			nvPairMap["INTF_VRF"] = vrf.(string)
		} else {
			nvPairMap["INTF_VRF"] = ""
		}
		if ip, ok := d.GetOk("ipv4"); ok {
			nvPairMap["IP"] = ip.(string)
		} else {
			nvPairMap["IP"] = ""
		}
		if ipv4Pre, ok := d.GetOk("ipv4_prefix"); ok {
			nvPairMap["PREFIX"] = ipv4Pre.(string)
		} else {
			nvPairMap["PREFIX"] = ""
		}
		if ipv6, ok := d.GetOk("ipv6"); ok {
			nvPairMap["IPv6"] = ipv6.(string)
		} else {
			nvPairMap["IPv6"] = ""
		}
		if ipv6Pre, ok := d.GetOk("ipv6_prefix"); ok {
			nvPairMap["IPv6_PREFIX"] = ipv6Pre.(string)
		} else {
			nvPairMap["IPv6_PREFIX"] = ""
		}
		if mtu, ok := d.GetOk("subinterface_mtu"); ok {
			nvPairMap["MTU"] = mtu.(string)
		} else {
			nvPairMap["MTU"] = ""
		}
		if desc, ok := d.GetOk("description"); ok {
			nvPairMap["DESC"] = desc.(string)
		} else {
			nvPairMap["DESC"] = ""
		}
		if conf, ok := d.GetOk("configuration"); ok {
			nvPairMap["CONF"] = conf.(string)
		} else {
			nvPairMap["CONF"] = ""
		}

	} else if intfType == "ethernet" {
		intfConfig.SerialNumber = serialnum

		if vrf, ok := d.GetOk("vrf"); ok {
			nvPairMap["INTF_VRF"] = vrf.(string)
		} else {
			nvPairMap["INTF_VRF"] = ""
		}
		if bpduF, ok := d.GetOk("bpdu_gaurd_flag"); ok {
			nvPairMap["BPDUGUARD_ENABLED"] = bpduF.(string)
		} else {
			nvPairMap["BPDUGUARD_ENABLED"] = ""
		}
		if pff, ok := d.GetOk("port_fast_flag"); ok {
			nvPairMap["PORTTYPE_FAST_ENABLED"] = pff.(bool)
		}
		if mtu, ok := d.GetOk("mtu"); ok {
			nvPairMap["MTU"] = mtu.(string)
		} else {
			nvPairMap["MTU"] = ""
		}
		if speed, ok := d.GetOk("ethernet_speed"); ok {
			nvPairMap["SPEED"] = speed.(string)
		} else {
			nvPairMap["SPEED"] = ""
		}
		if vlans, ok := d.GetOk("allowed_vlans"); ok {
			nvPairMap["ALLOWED_VLANS"] = vlans.(string)
		} else {
			nvPairMap["ALLOWED_VLANS"] = ""
		}
		if ip, ok := d.GetOk("ipv4"); ok {
			nvPairMap["IP"] = ip.(string)
		} else {
			nvPairMap["IP"] = ""
		}
		if ipv4Pre, ok := d.GetOk("ipv4_prefix"); ok {
			nvPairMap["PREFIX"] = ipv4Pre.(string)
		} else {
			nvPairMap["PREFIX"] = ""
		}
		if ipv6, ok := d.GetOk("ipv6"); ok {
			nvPairMap["IPv6"] = ipv6.(string)
		} else {
			nvPairMap["IPv6"] = ""
		}
		if ipv6Pre, ok := d.GetOk("ipv6_prefix"); ok {
			nvPairMap["IPv6_PREFIX"] = ipv6Pre.(string)
		} else {
			nvPairMap["IPv6_PREFIX"] = ""
		}
		if accVlans, ok := d.GetOk("access_vlans"); ok {
			nvPairMap["ACCESS_VLAN"] = accVlans.(string)
		} else {
			nvPairMap["ACCESS_VLAN"] = ""
		}
		if desc, ok := d.GetOk("description"); ok {
			nvPairMap["DESC"] = desc.(string)
		} else {
			nvPairMap["DESC"] = ""
		}
		if conf, ok := d.GetOk("configuration"); ok {
			nvPairMap["CONF"] = conf.(string)
		} else {
			nvPairMap["CONF"] = ""
		}

	}

	if state, ok := d.GetOk("admin_state"); ok {
		nvPairMap["ADMIN_STATE"] = state.(bool)
	} else {
		nvPairMap["ADMIN_STATE"] = ""
	}

	intfModel := models.NewInterface(&intf, &intfConfig, nvPairMap)

	cont, err := dcnmClient.Update("/rest/interface", intfModel)
	if err != nil {
		if cont != nil {
			errorMsg, flag := checkIntfErrors(cont)
			if flag {
				return fmt.Errorf(errorMsg)
			}
		} else {
			return err
		}
	}

	d.Set("serial_number", intfConfig.SerialNumber)
	d.Set("type", intfType)
	d.SetId(intfConfig.InterfaceName)

	if d.HasChange("deploy") && d.Get("deploy").(bool) == false {
		return fmt.Errorf("Deployed interface can not be undeployed")
	}

	//Deployment of interface
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		log.Println("[DEBUG] Begining Deployment ", d.Id())

		intfDeploy := models.InterfaceDelete{}
		intfDeploy.SerialNumber = intfConfig.SerialNumber
		intfDeploy.Name = intfConfig.InterfaceName
		cont, err = dcnmClient.SaveForAttachment("/rest/interface/deploy", &intfDeploy)
		if err != nil {
			errorMsg, flag := checkIntfErrors(cont)
			if flag {
				d.Set("deploy", false)
				return fmt.Errorf("interface is created but failed to deploy with error : %s", errorMsg)
			}
		}

		log.Println("[DEBUG] End of Deployment ", d.Id())
	}

	log.Println("[DEBUG] End of Update method ", d.Id())
	return nil
}

func resourceDCNMInterfaceRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	var serialNum1 string
	dn := d.Id()
	intfType := d.Get("type").(string)
	serialNum := d.Get("serial_number").(string)
	if intfType == "vpc" {
		serialNum1 = (strings.Split(d.Get("serial_number").(string), "~"))[0]
	} else {
		serialNum1 = d.Get("serial_number").(string)
	}

	cont, err := getRemoteInterface(dcnmClient, serialNum1, dn)
	if err != nil {
		if cont != nil {
			errorMsg, flag := checkIntfErrors(cont)
			if flag {
				return fmt.Errorf(errorMsg)
			}
		} else {
			return err
		}
	}

	setInterfaceAttributes(d, cont.Index(0), intfType)
	d.SetId(dn)

	flag, err := checkIntfDeploy(dcnmClient, serialNum, d.Get("name").(string), intfType)
	if err != nil {
		return err
	}
	d.Set("deploy", flag)

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	serialNum := d.Get("serial_number").(string)
	intfType := d.Get("type").(string)

	if intfType == "ethernet" {
		d.SetId("")

		log.Println("[DEBUG] End of Delete method ")
		return nil
	}

	intfDel := models.InterfaceDelete{}
	intfDel.SerialNumber = serialNum
	intfDel.Name = dn
	cont, err := dcnmClient.DeleteWithPayload("/rest/interface", &intfDel)
	if err != nil {
		if cont != nil {
			errorMsg, flag := checkIntfErrors(cont)
			if flag {
				return fmt.Errorf(errorMsg)
			}
		} else {
			return err
		}
	}

	d.SetId("")

	log.Println("[DEBUG] End of Delete method ")
	return nil
}

func checkIntfErrors(cont *container.Container) (string, bool) {
	totalMsg := len(cont.Data().([]interface{}))
	flag := false
	errMsg := ""
	for i := 0; i < totalMsg; i++ {
		if stripQuotes(cont.Index(i).S("reportItemType").String()) == "ERROR" {
			flag = true
			errMsg = fmt.Sprintf("%s %s", errMsg, stripQuotes(cont.Index(i).S("message").String()))
		}
	}
	return errMsg, flag
}

func checkIntfDeploy(client *client.Client, serialnum, name, intftype string) (bool, error) {
	flag := false
	intfStr := fmt.Sprintf("%s~%s", serialnum, name)

	var serial1 string
	if intftype == "vpc" {
		serial1 = (strings.Split(serialnum, "~"))[0]
	} else {
		serial1 = serialnum
	}

	cont, err := client.GetviaURL(fmt.Sprintf("/rest/interface/detail?serialNumber=%s", serial1))
	if err != nil {
		if cont != nil {
			errorMsg, flag := checkIntfErrors(cont)
			if flag {
				return flag, fmt.Errorf(errorMsg)
			}
		} else {
			return flag, err
		}
	}

	totalIntf := len(cont.Data().([]interface{}))
	for i := 0; i < totalIntf; i++ {
		if stripQuotes(cont.Index(i).S("entityId").String()) == intfStr {
			if stripQuotes(cont.Index(i).S("complianceStatus").String()) == "In-Sync" {
				return true, nil
			}
		}
	}

	return flag, nil
}

func getSwitchName(client *client.Client, fabric, serialNum string) (string, error) {
	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabric)

	cont, err := client.GetviaURL(durl)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		switchCont := cont.Index(i)

		serial := stripQuotes(switchCont.S("serialNumber").String())
		if serial == serialNum {
			return stripQuotes(switchCont.S("logicalName").String()), nil
		}
	}
	return "", fmt.Errorf("no switch found for given serial-number in given fabric")
}
