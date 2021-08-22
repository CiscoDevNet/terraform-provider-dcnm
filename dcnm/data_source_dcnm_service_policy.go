package dcnm

import (
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDCNMServicePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDCNMServicePolicyRead,

		Schema: map[string]*schema.Schema{
			"policy_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"attached_fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dest_network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dest_vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"peering_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"policy_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"reverse_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"reverse_next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"service_node_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"service_node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source_network": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source_vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"src_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dest_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"next_hop_action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"fwd_direction": &schema.Schema{
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

func dataSourceDCNMServicePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Get("policy_name").(string))

	dcnmClient := m.(*client.Client)

	policyName := d.Get("policy_name").(string)
	fabricName := d.Get("fabric_name").(string)
	attachedFabricName := d.Get("attached_fabric_name").(string)
	serviceNodeName := d.Get("service_node_name").(string)

	cont, err := getServicePolicy(dcnmClient, attachedFabricName, fabricName, serviceNodeName, policyName)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}
	setServicePolicyAttributes(d, cont)
	d.Set("reverse_next_hop_ip", stripQuotes(cont.S("reverseNextHopIp").String()))
	d.Set("service_node_type", stripQuotes(cont.S("serviceNodeType").String()))

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}
