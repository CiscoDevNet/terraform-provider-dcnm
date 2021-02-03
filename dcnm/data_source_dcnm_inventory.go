package dcnm

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDCNMInventory() *schema.Resource {
	return &schema.Resource{
		Read: datasourceDCNMInventoryRead,

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"switch_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"switch_db_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"model": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func getRemoteSwitchforDS(dcnmClient *client.Client, fabric, name string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/control/fabrics/%s/inventory", fabric)

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		switchCont := cont.Index(i)

		nameGet := stripQuotes(switchCont.S("logicalName").String())
		if nameGet == name {
			return switchCont, nil
		}
	}
	return nil, fmt.Errorf("Desired switch not found")
}

func datasourceDCNMInventoryRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ")

	dcnmClient := m.(*client.Client)

	name := d.Get("switch_name").(string)

	fabricName := d.Get("fabric_name").(string)

	cont, err := getRemoteSwitchforDS(dcnmClient, fabricName, name)
	if err != nil {
		return err
	}

	setSwitchAttributes(d, cont)

	flag, err := checkDeploy(dcnmClient, fabricName, d.Get("serial_number").(string))
	if err != nil {
		return err
	}
	if flag {
		d.Set("deploy", true)
	} else {
		d.Set("deploy", false)
	}

	role, err := getSwitchRole(dcnmClient, d.Get("serial_number").(string))
	if err == nil {
		d.Set("role", role)
	}

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}
