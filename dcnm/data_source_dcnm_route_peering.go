package dcnm

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func datasourceDCNMRoutePeering() *schema.Resource {
	return &schema.Resource{
		Read: datasourceRoutePeeringRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attached_fabric": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deployment_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IntraTenantFW",
					"InterTenantFW",
					"OneArmADC",
					"TwoArmADC",
					"OneArmVNF",
				}, false),
			},
			"service_fabric": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"option": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"StaticPeering",
					"EBGPDynamicPeering",
					"None",
				}, false),
			},
			"reverse_next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_networks": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"network_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"InsideNetworkFW",
								"OutsideNetworkFW",
								"ArmOneADC",
								"ArmTwoADC",
								"ArmOneVNF",
							}, false),
						},
						"template_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vlan_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"vrf_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"gateway_ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"service_node_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_node_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Firewall",
					"ADC",
					"VNF",
				}, false),
			},
			"routes": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"vrf_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"route_parmas": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
		},
	}
}

func datasourceRoutePeeringRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method", d.Id())
	dcnmClient := m.(*client.Client)

	AttachedFabricName := d.Get("attached_fabric").(string)
	extFabric := d.Get("service_fabric").(string)
	node := d.Get("service_node_name").(string)
	name := d.Get("name").(string)
	cont, err := getRoutePeering(dcnmClient, AttachedFabricName, extFabric, node, name)
	if err != nil {
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}
	setPeeringAttributes(d, cont)
	d.SetId(name)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}
