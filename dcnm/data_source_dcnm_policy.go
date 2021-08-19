package dcnm

import (
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDCNMPolicy() *schema.Resource {
	return &schema.Resource{
		Read: datasourceDCNMPolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_content_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_props": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
func datasourceDCNMPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read Method ", d.Id())
	dcnmClient := m.(*client.Client)

	policyId := d.Get("policy_id").(string)

	cont, err := getAllPolicy(dcnmClient, policyId)

	if err != nil {
		return err
	}
	setPolicyAttributes(d, cont)
	d.SetId(policyId)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}
