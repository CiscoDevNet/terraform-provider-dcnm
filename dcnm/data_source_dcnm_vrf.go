package dcnm

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDCNMVRF() *schema.Resource {
	return &schema.Resource{
		Read: datasourceDCNMVRFRead,

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"segment_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},

			"max_bgp_path": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"max_ibgp_path": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"trm_enable": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rp_external_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},

			"trm_bgw_msite_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"advertise_host_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"advertise_default_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"static_default_route": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"extension_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"service_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_props": {
				Type:     schema.TypeMap,
				Optional: true,
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

						"vlan_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceDCNMVRFRead(d *schema.ResourceData, m interface{}) error {
	dcnmClient := m.(*client.Client)

	dn := d.Get("name").(string)

	fabricName := d.Get("fabric_name").(string)

	cont, err := getRemoteVRF(dcnmClient, fabricName, dn)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("template_props"); ok {
		setVRFCustomTemplateAttributes(d, cont)
	} else {
		setVRFAttributes(d, cont)
	}

	flag, err := checkvrfDeploy(dcnmClient, fabricName, dn)
	if err != nil {
		d.Set("deploy", false)
		return err
	}
	d.Set("deploy", flag)

	attachments, err := getAttachmentList(dcnmClient, fabricName, dn)
	if err == nil {
		d.Set("attachments", attachments)
	} else {
		d.Set("attachments", make([]interface{}, 0, 1))
	}

	return nil
}

func getAttachmentList(client *client.Client, fabric, name string) ([]interface{}, error) {
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/attachments?vrf-names=%s", fabric, name)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	attachList := cont.Index(0).S("lanAttachList")

	attachmentList := make([]interface{}, 0, 1)
	for i := 0; i < len(attachList.Data().([]interface{})); i++ {
		attachMap := make(map[string]interface{})

		attachMap["serial_number"] = stripQuotes(attachList.Index(i).S("switchSerialNo").String())

		if attachList.Index(i).Exists("vlanId") && stripQuotes(attachList.Index(i).S("vlanId").String()) != "" {
			if vlan, err := strconv.Atoi(stripQuotes(attachList.Index(i).S("vlanId").String())); err == nil {
				attachMap["vlan_id"] = vlan
			}
		}

		if attachList.Index(i).Exists("isLanAttached") && stripQuotes(attachList.Index(i).S("isLanAttached").String()) != "" {
			if attach, err := strconv.ParseBool(stripQuotes(attachList.Index(i).S("isLanAttached").String())); err == nil {
				attachMap["attach"] = attach
			}
		}
		attachmentList = append(attachmentList, attachMap)
	}
	return attachmentList, nil
}
