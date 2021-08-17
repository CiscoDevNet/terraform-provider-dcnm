package dcnm

import (
	"fmt"
	"log"
	"reflect"
	"strings"

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
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice([]string{
					"Physical",
					"Virtual",
				}, false),
				Default: "Virtual",
			},

			"service_fabric": &schema.Schema{
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
				Optional: true,
				Default:  "service_link_trunk",
			},

			"switches": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
				Set:      schema.HashString,
			},

			"attached_switch_interface_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"attached_fabric": &schema.Schema{
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

			"bpdu_guard_flag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes", "no",
				}, false),
				Default: "no",
			},

			"porttype_fast_enabled": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
				Default: "true",
			},

			"admin_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
				Default: "true",
			},

			"policy_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDCNMServiceNodeImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 2 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	fabricName := importInfo[0]
	name := importInfo[1]

	cont, err := getServiceNodeAttributes(dcnmClient,fabricName, name)
	if err != nil {
		return nil, err
	}

	stateImport := setServiceNodeAttributes(d, cont)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return []*schema.ResourceData{stateImport}, nil
}

func getServiceNodeAttributes(dcnmClient *client.Client,fabricName, name string) (*container.Container, error) {
	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s", fabricName, name)
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", fabricName, name)
	}

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return nil, err
	}
	return cont,nil
}
func setServiceNodeAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {

	d.Set("name", stripQuotes(cont.S("name").String()))
	d.Set("service_fabric", stripQuotes(cont.S("fabricName").String()))
	d.Set("node_type", stripQuotes(cont.S("type").String()))
	d.Set("form_factor", stripQuotes(cont.S("formFactor").String()))
	d.Set("interface_name", stripQuotes(cont.S("interfaceName").String()))
	d.Set("link_template_name", stripQuotes(cont.S("linkTemplateName").String()))

	switchList := strings.Split(stripQuotes(cont.S("attachedSwitchSn").String()), ",")
	if switcheSn, ok := d.GetOk("switches"); ok {
		tfSwitches := toStringList(switcheSn.(*schema.Set).List())
		if !reflect.DeepEqual(tfSwitches, switchList) {
			d.Set("switches", switchList)
		} else {
			d.Set("switches", tfSwitches)
		}
	} else {
		d.Set("switches", switchList)
	}

	d.Set("attached_switch_interface_name", stripQuotes(cont.S("attachedSwitchInterfaceName").String()))
	d.Set("attached_fabric", stripQuotes(cont.S("attachedFabricName").String()))
	d.Set("form_factor", stripQuotes(cont.S("formFactor").String()))
	d.Set("speed", stripQuotes(cont.S("nvPairs", "SPEED").String()))
	d.Set("mtu", stripQuotes(cont.S("nvPairs", "MTU").String()))
	d.Set("allowed_vlans", stripQuotes(cont.S("nvPairs", "ALLOWED_VLANS").String()))
	d.Set("bpdu_guard_flag", stripQuotes(cont.S("nvPairs", "BPDUGUARD_ENABLED").String()))
	d.Set("porttype_fast_enabled", stripQuotes(cont.S("nvPairs", "PORTTYPE_FAST_ENABLED").String()))
	d.Set("admin_state", stripQuotes(cont.S("nvPairs", "ADMIN_STATE").String()))

	d.SetId(stripQuotes(cont.S("name").String()))
	return d
}

func resourceDCNMServiceNodeCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	attachedFabric := d.Get("attached_fabric").(string)
	switches := toStringList(d.Get("switches").(*schema.Set).List())
	if len(switches) > 2 {
		return fmt.Errorf("Fabric: %s - Upto 2 switches only allowed", attachedFabric)
	}
	attachedSwitchSn := strings.Join(switches[:], ",")

	serviceNode := models.ServiceNode{
		Name:                        d.Get("name").(string),
		Type:                        d.Get("node_type").(string),
		FabricName:                  d.Get("service_fabric").(string),
		InterfaceName:               d.Get("interface_name").(string),
		LinkTemplateName:            d.Get("link_template_name").(string),
		AttachedSwitchSn:            attachedSwitchSn,
		AttachedSwitchInterfaceName: d.Get("attached_switch_interface_name").(string),
		AttachedFabricName:          attachedFabric,
	}

	if formFactor, ok := d.GetOk("form_factor"); ok {
		serviceNode.FormFactor = formFactor.(string)
	}
	nvPairMap := make(map[string]interface{})

	if speed, ok := d.GetOk("speed"); ok {
		nvPairMap["SPEED"] = speed.(string)
	}
	if MTU, ok := d.GetOk("mtu"); ok {
		nvPairMap["MTU"] = MTU.(string)
	}
	if allowedVlans, ok := d.GetOk("allowed_vlans"); ok {
		nvPairMap["ALLOWED_VLANS"] = allowedVlans.(string)
	}
	if bpduguardEnabled, ok := d.GetOk("bpdu_guard_flag"); ok {
		nvPairMap["BPDUGUARD_ENABLED"] = bpduguardEnabled.(string)
	}
	if porttypeFastEnabled, ok := d.GetOk("porttype_fast_enabled"); ok {
		nvPairMap["PORTTYPE_FAST_ENABLED"] = porttypeFastEnabled.(string)
	}
	if adminState, ok := d.GetOk("admin_state"); ok {
		nvPairMap["ADMIN_STATE"] = adminState.(string)
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
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes", serviceNode.FabricName)
	}

	cont, err := dcnmClient.Save(durl, &serviceNode)
	if err != nil {
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}

	d.SetId(serviceNode.Name)
	log.Println("[DEBUG] End of Create ", d.Id())
	return resourceDCNMServiceNodeRead(d, m)
}

func resourceDCNMServiceNodeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	attachedFabric := d.Get("attached_fabric").(string)
	switches := toStringList(d.Get("switches").(*schema.Set).List())
	if len(switches) > 2 {
		return fmt.Errorf("Fabric: %s - Upto 2 switches only allowed", attachedFabric)
	}
	attachedSwitchSn := strings.Join(switches[:], ",")

	serviceNode := models.ServiceNode{
		Name:                        d.Get("name").(string),
		Type:                        d.Get("node_type").(string),
		FabricName:                  d.Get("service_fabric").(string),
		InterfaceName:               d.Get("interface_name").(string),
		LinkTemplateName:            d.Get("link_template_name").(string),
		AttachedSwitchSn:            attachedSwitchSn,
		AttachedSwitchInterfaceName: d.Get("attached_switch_interface_name").(string),
		AttachedFabricName:          attachedFabric,
	}

	if formFactor, ok := d.GetOk("form_factor"); ok {
		serviceNode.FormFactor = formFactor.(string)
	}
	nvPairMap := make(map[string]interface{})

	if speed, ok := d.GetOk("speed"); ok {
		nvPairMap["SPEED"] = speed.(string)
	}
	if MTU, ok := d.GetOk("mtu"); ok {
		nvPairMap["MTU"] = MTU.(string)
	}
	if allowedVlans, ok := d.GetOk("allowed_vlans"); ok {
		nvPairMap["ALLOWED_VLANS"] = allowedVlans.(string)
	}
	if bpduguardEnabled, ok := d.GetOk("bpdu_guard_flag"); ok {
		nvPairMap["BPDUGUARD_ENABLED"] = bpduguardEnabled.(string)
	}
	if porttypeFastEnabled, ok := d.GetOk("porttype_fast_enabled"); ok {
		nvPairMap["PORTTYPE_FAST_ENABLED"] = porttypeFastEnabled.(string)
	}
	if adminState, ok := d.GetOk("admin_state"); ok {
		nvPairMap["ADMIN_STATE"] = adminState.(string)
	}
	if policyDesc, ok := d.GetOk("policy_description"); ok {
		nvPairMap["POLICY_DESC"] = policyDesc.(string)
	}
	if nvPairMap != nil {
		serviceNode.NVPairs = nvPairMap
	}

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s", serviceNode.FabricName, serviceNode.Name)
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", serviceNode.FabricName, serviceNode.Name)
	}

	cont, err := dcnmClient.Update(durl, &serviceNode)
	if err != nil {
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}

	d.SetId(serviceNode.Name)
	log.Println("[DEBUG] End of Update ", d.Id())
	return resourceDCNMServiceNodeRead(d, m)
}

func resourceDCNMServiceNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	name := d.Get("name").(string)
	fabricName := d.Get("service_fabric").(string)

	cont, err := getServiceNodeAttributes(fabricName, name)
	if err != nil {
		d.SetId("")
		return nil
	}

	setServiceNodeAttributes(d, cont)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMServiceNodeDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)
	serviceFabric := d.Get("service_fabric").(string)

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s", serviceFabric, d.Id())
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", serviceFabric, d.Id())
	}
	_, err := dcnmClient.Delete(durl)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}
