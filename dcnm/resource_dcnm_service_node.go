package dcnm

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDCNMServiceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMServiceNodeCreate,
		Read:   resourceDCNMServiceNodeRead,
		Update: resourceDCNMServiceNodeUpdate,
		Delete: resourceDCNMServiceNodeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDCNMServiceNodeImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"node_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Firewall",
					"ADC",
					"VNF",
				}, false),
			},

			"form_factor": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Physical",
					"Vitual",
				}, false),
			},

			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"interface_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"link_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"attached_switch_sn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"attached_switch_interface_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"attached_fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"speed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Auto",
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "jumbo",
			},

			"allowed_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none",
			},

			"bpdu_guard_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "no",
			},

			"porttype_fast_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
			},

			"admin_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
			},

			"source_if_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Ethernet1/8",
			},

			"source_fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "macross",
			},

			"source_switch_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FDO23420QS7",
			},

			"link_uuid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "link-UUID-327240",
			},

			"priority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "500",
			},

			"dest_fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "service_fabric",
			},

			"policy_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "POLICY-327250",
			},

			"dest_switch_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LB",
			},

			"is_metaswitch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "true",
			},

			"dest_if_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "lb_interface",
			},

			"dest_serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "LB-service_fabric",
			},

			"source_serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FDO23420QS7",
			},

			"policy_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"force_deletion": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"retain_switch": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDCNMServiceNodeImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return nil,nil
}

func resourceDCNMServiceNodeCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	serviceNode := models.ServiceNode{
		Name:                        d.Get("name").(string),
		Type:                        d.Get("node_type").(string),
		FabricName:                  d.Get("fabric_name").(string),
		InterfaceName:               d.Get("interfaceName").(string),
		LinkTemplateName:            d.Get("linkTemplateName").(string),
		AttachedSwitchSn:            d.Get("attachedSwitchSn").(string),
		AttachedSwitchInterfaceName: d.Get("attachedSwitchinterfaceName").(string),
		AttachedFabricName:          d.Get("attachedFabricName").(string),
	}

	if formFactor, ok := d.GetOk("form_factor"); ok {
		serviceNode.FormFactor = formFactor.(string)
	}
	nvPairMap := make(map[string]interface{})

	if speed, ok := d.GetOk("speed"); ok {
		nvPairMap["SPEED"] = speed.(string)
	} else {
		nvPairMap["SPEED"] = ""
	}
	if MTU, ok := d.GetOk("mtu"); ok {
		nvPairMap["MTU"] = MTU.(string)
	} else {
		nvPairMap["MTU"] = ""
	}
	if allowedVlans, ok := d.GetOk("allowed_vlans"); ok {
		nvPairMap["ALLOWED_VLANS"] = allowedVlans.(string)
	} else {
		nvPairMap["ALLOWED_VLANS"] = ""
	}
	if bpduguardEnabled, ok := d.GetOk("bpduguard_enabled"); ok {
		nvPairMap["BPDUGUARD_ENABLED"] = bpduguardEnabled.(string)
	} else {
		nvPairMap["BPDUGUARD_ENABLED"] = ""
	}
	if porttypeFastEnabled, ok := d.GetOk("porttype_fast_enabled"); ok {
		nvPairMap["PORTTYPE_FAST_ENABLED"] = porttypeFastEnabled.(string)
	} else {
		nvPairMap["PORTTYPE_FAST_ENABLED"] = ""
	}
	if adminState, ok := d.GetOk("admin_state"); ok {
		nvPairMap["ADMIN_STATE"] = adminState.(string)
	} else {
		nvPairMap["ADMIN_STATE"] = ""
	}
	if sourceIfName, ok := d.GetOk("source_if_name"); ok {
		nvPairMap["SOURCE_IF_NAME"] = sourceIfName.(string)
	}
	if sourceFabricName, ok := d.GetOk("source_fabric_name"); ok {
		nvPairMap["SOURCE_FABRIC_NAME"] = sourceFabricName.(string)
	}
	if sourceSwitchName, ok := d.GetOk("source_switch_name"); ok {
		nvPairMap["SOURCE_SWITCH_NAME"] = sourceSwitchName.(string)
	}
	if linkUUID, ok := d.GetOk("link_uuid"); ok {
		nvPairMap["LINK_UUID"] = linkUUID.(string)
	}
	if prio, ok := d.GetOk("priority"); ok {
		nvPairMap["PRIORITY"] = prio.(string)
	}
	if destFabricName, ok := d.GetOk("dest_fabric_name"); ok {
		nvPairMap["DEST_FABRIC_NAME"] = destFabricName.(string)
	}
	if policyID, ok := d.GetOk("policy_id"); ok {
		nvPairMap["POLICY_ID"] = policyID.(string)
	}
	if destSwitchName, ok := d.GetOk("dest_switch_name"); ok {
		nvPairMap["DEST_SWITCH_NAME"] = destSwitchName.(string)
	}
	if isMetaswitch, ok := d.GetOk("is_metaswitch"); ok {
		nvPairMap["IS_METASWITCH"] = isMetaswitch.(string)
	}
	if destIfName, ok := d.GetOk("dest_if_name"); ok {
		nvPairMap["DEST_IF_NAME"] = destIfName.(string)
	}
	if destSerialNumber, ok := d.GetOk("dest_serial_number"); ok {
		nvPairMap["DEST_SERIAL_NUMBER"] = destSerialNumber.(string)
	}
	if sourceSerialNumber, ok := d.GetOk("source_serial_number"); ok {
		nvPairMap["SOURCE_SERIAL_NUMBER"] = sourceSerialNumber.(string)
	}
	if policyDesc, ok := d.GetOk("policy_description"); ok {
		nvPairMap["POLICY_DESC"] = policyDesc.(string)
	}
	if nvPairMap != nil {
		serviceNode.NVPairs = nvPairMap
	}

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes", serviceNode.FabricName)
	} else {
		durl = fmt.Sprintf("/rest/fabrics/%s/service-node", serviceNode.FabricName)
	}

	_, err := dcnmClient.Save(durl, &serviceNode)
	if err != nil {
		return err
	}

	d.SetId(serviceNode.Name)
	return resourceDCNMServiceNodeRead(d, m)
}

func resourceDCNMServiceNodeUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDCNMServiceNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	fabricName := d.Get("fabric_name").(string)

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes", fabricName)
	} else {
		durl = fmt.Sprintf("/rest/fabrics/%s/service-node", fabricName)
	}

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return err
	}

	setServiceNodeAttributes(d, cont)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMServiceNodeDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func getRemoteServiceNode(client *client.Client, fabricName, nodeName string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/fabrics/%s/service-nodes/%s", fabricName, nodeName)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return cont, err
	}
	return cont, nil
}

func setServiceNodeAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {

	d.Set("name", stripQuotes(cont.S("name").String()))
	d.Set("fabricName", stripQuotes(cont.S("fabricName").String()))
	d.Set("form_factor", stripQuotes(cont.S("formFactor").String()))
	d.Set("interfaceName", stripQuotes(cont.S("interfaceName").String()))
	d.Set("linkTemplateName", stripQuotes(cont.S("linkTemplateName").String()))
	d.Set("attachedSwitchSn", stripQuotes(cont.S("attachedSwitchSn").String()))
	d.Set("attachedSwitchinterfaceName", stripQuotes(cont.S("attachedSwitchinterfaceName").String()))
	d.Set("attachedFabricName", stripQuotes(cont.S("attachedFabricName").String()))

	return d
}
