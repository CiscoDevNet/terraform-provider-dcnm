package dcnm

import (
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDCNMTemplate() *schema.Resource {
	return &schema.Resource{
		Read: datasourceDCNMTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"supported_platforms": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_sub_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_content_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func datasourceDCNMTemplateRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read Method ", d.Id())

	dcnmClient := m.(*client.Client)

	name := d.Get("name").(string)

	cont, err := getTemplate(dcnmClient, name)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}
	setTemplateAttribute(d, cont)
	d.SetId(name)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}
