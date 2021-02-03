package dcnm

import (
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func datasourceDCNMInterface() *schema.Resource {
	return &schema.Resource{
		Read: datasourceDCNMInterfaceRead,

		Schema: map[string]*schema.Schema{
			"serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"port-channel",
					"vpc",
					"loopback",
					"sub-interface",
					"ethernet",
				}, false),
			},

			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"switch_name_1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"bpdu_gaurd_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"configuration": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"admin_state": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func datasourceDCNMInterfaceRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ")

	dcnmClient := m.(*client.Client)

	var serialNum1 string
	var serialNum2 string
	name := d.Get("name").(string)
	intfType := d.Get("type").(string)
	serialNum := d.Get("serial_number").(string)
	if intfType == "vpc" {
		vpcSerialNums := strings.Split(serialNum, "~")
		if len(vpcSerialNums) != 2 {
			return fmt.Errorf("serial number is not valid for vpc interface")
		}
		serialNum1 = vpcSerialNums[0]
		serialNum2 = vpcSerialNums[1]
	} else {
		serialNum1 = serialNum
	}

	cont, err := getRemoteInterface(dcnmClient, serialNum1, name)
	if err != nil {
		return err
	}

	setInterfaceAttributes(d, cont.Index(0), intfType)

	d.SetId(name)

	fabName := d.Get("fabric_name").(string)
	if intfType == "vpc" {
		swName1, err := getSwitchName(dcnmClient, fabName, serialNum1)
		if err == nil {
			d.Set("switch_name_1", swName1)
		}
		swName2, err := getSwitchName(dcnmClient, fabName, serialNum2)
		if err == nil {
			d.Set("switch_name_2", swName2)
		}
	} else {
		swName1, err := getSwitchName(dcnmClient, fabName, serialNum1)
		if err == nil {
			d.Set("switch_name_1", swName1)
		}
	}

	if flag, err := checkIntfDeploy(dcnmClient, serialNum, d.Get("name").(string), intfType); err != nil {
		d.Set("deploy", flag)
	} else {
		d.Set("deploy", flag)
	}

	log.Println("[DEBUG] Begining Read method ", d.Id())
	return nil
}
