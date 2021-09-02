package dcnm

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDCNMServiceNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDCNMServiceNodeRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"form_factor": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"service_fabric": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"interface_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"link_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"switches": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},

			"attached_switch_interface_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"attached_fabric": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"allowed_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"bpdu_guard_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"porttype_fast_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"admin_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source_if_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source_fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source_switch_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"priority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dest_fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dest_switch_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_metaswitch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dest_if_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dest_serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source_serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"policy_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDCNMServiceNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Get("name").(string))

	dcnmClient := m.(*client.Client)

	serviceNodeName := d.Get("name").(string)
	fabricName := d.Get("service_fabric").(string)
	attachedFabricName := d.Get("attached_fabric").(string)

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/ndfc/api/v1/elastic-service/fabrics/%s/service-nodes/%s", fabricName, serviceNodeName)
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", fabricName, serviceNodeName)
	}

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return err
	}

	setServiceNodeAttributes(d, cont)
	d.Set("source_if_name", stripQuotes(cont.S("nvPairs", "SOURCE_IF_NAME").String()))
	d.Set("source_fabric_name", stripQuotes(cont.S("nvPairs", "SOURCE_FABRIC_NAME").String()))
	d.Set("source_switch_name", stripQuotes(cont.S("nvPairs", "SOURCE_SWITCH_NAME").String()))
	d.Set("priority", stripQuotes(cont.S("nvPairs", "PRIORITY").String()))
	d.Set("dest_fabric_name", stripQuotes(cont.S("nvPairs", "DEST_FABRIC_NAME").String()))
	d.Set("policy_id", stripQuotes(cont.S("nvPairs", "POLICY_ID").String()))
	d.Set("dest_switch_name", stripQuotes(cont.S("nvPairs", "DEST_SWITCH_NAME").String()))
	d.Set("is_metaswitch", stripQuotes(cont.S("nvPairs", "IS_METASWITCH").String()))
	d.Set("dest_if_name", stripQuotes(cont.S("nvPairs", "DEST_IF_NAME").String()))
	d.Set("dest_serial_number", stripQuotes(cont.S("nvPairs", "DEST_SERIAL_NUMBER").String()))
	d.Set("source_serial_number", stripQuotes(cont.S("nvPairs", "SOURCE_SERIAL_NUMBER").String()))
	d.Set("policy_description", stripQuotes(cont.S("nvPairs", "POLICY_DESC").String()))

	d.SetId(fmt.Sprintf("%s/%s/%s", fabricName, attachedFabricName, serviceNodeName))
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}
