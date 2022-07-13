package dcnm

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDCNMNetwork() *schema.Resource {
	return &schema.Resource{
		Read: datasourceDCNMNetworkRead,

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"extension_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"l2_only_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"vlan_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv4_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"secondary_gw_1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"secondary_gw_2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"arp_supp_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"ir_enable_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"mcast_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_vrf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"trm_enable_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"rt_both_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"l3_gateway_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},

			"attachments": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"serial_number": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"attach": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},

						"switch_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"vlan_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"dot1_qvlan": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"untagged": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
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

						"instance_values": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceDCNMNetworkRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ")

	dcnmClient := m.(*client.Client)

	name := d.Get("name").(string)
	fabricName := d.Get("fabric_name").(string)

	cont, err := getRemoteNetwork(dcnmClient, fabricName, name)
	if err != nil {
		return err
	}

	setNetworkAttributes(d, cont)

	deployed, err := checkNetworkDeploy(dcnmClient, fabricName, name)
	if err != nil {
		d.Set("deploy", false)
		return err
	}
	d.Set("deploy", deployed)

	attachments, err := getNetworkAttachmentList(dcnmClient, fabricName, name)
	if err == nil {
		d.Set("attachments", attachments)
	} else {
		d.Set("attachments", make([]interface{}, 0, 1))
	}

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func getNetworkAttachmentList(client *client.Client, fabric, network string) ([]map[string]interface{}, error) {
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s/attachments", fabric, network)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	attachments := make([]map[string]interface{}, 0, 1)
	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		attachment := cont.Index(i)

		attachMap := make(map[string]interface{})
		attachMap["serial_number"] = stripQuotes(attachment.S("switchSerialNo").String())

		if stripQuotes(attachment.S("isLanAttached").String()) == "true" {
			attachMap["attach"] = true
		} else {
			attachMap["attach"] = false
		}

		if stripQuotes(attachment.S("vlanId").String()) != "null" {
			attachMap["vlan_id"] = int((attachment.S("vlanId").Data()).(float64))
		} else {
			attachMap["vlan_id"] = 0
		}

		if stripQuotes(attachment.S("portNames").String()) != "null" {
			attachMap["switch_ports"] = stringToList(stripQuotes(attachment.S("portNames").String()))
		} else {
			attachMap["switch_ports"] = make([]string, 0, 1)
		}

		attachMap["dot1_qvlan"] = 0
		attachMap["untagged"] = false
		attachMap["free_form_config"] = ""
		attachMap["extension_values"] = ""
		attachMap["instance_values"] = ""

		attachments = append(attachments, attachMap)
	}

	return attachments, nil
}
