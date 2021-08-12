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
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 2 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	fabricName := importInfo[0]
	name := importInfo[1]

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes", fabricName)
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", fabricName,name)
	}

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		return nil,err
	}

	stateImport := setServiceNodeAttributes(d, cont)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return []*schema.ResourceData{stateImport},nil
}


func setServiceNodeAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {

	d.Set("name", stripQuotes(cont.S("name").String()))
	d.Set("service_fabric", stripQuotes(cont.S("fabricName").String()))
	d.Set("node_type", stripQuotes(cont.S("type").String()))
	d.Set("form_factor", stripQuotes(cont.S("formFactor").String()))
	d.Set("interface_name", stripQuotes(cont.S("interfaceName").String()))
	d.Set("link_template_name", stripQuotes(cont.S("linkTemplateName").String()))

	switchList := strings.Split(stripQuotes(cont.S("attachedSwitchSn").String()),",")
	if switcheSn,ok := d.GetOk("switches"); ok {
		tfSwitches := toStringList(switcheSn.(*schema.Set).List())
		if !reflect.DeepEqual(tfSwitches,switchList) {
			d.Set("switches", switchList)
		} else {
			d.Set("switches", tfSwitches)
		}
	} else {
		d.Set("switches", switchList)
	}

	d.Set("attached_switch_interface_name", stripQuotes(cont.S("attachedSwitchinterfaceName").String()))
	d.Set("attached_fabric", stripQuotes(cont.S("attachedFabricName").String()))
	d.Set("form_factor", stripQuotes(cont.S("formFactor").String()))
	d.Set("speed", stripQuotes(cont.S("nvPair","SPEED").String()))
	d.Set("mtu", stripQuotes(cont.S("nvPair","MTU").String()))
	d.Set("allowed_vlans", stripQuotes(cont.S("nvPair","ALLOWED_VLANS").String()))
	d.Set("bpduguard_enabled", stripQuotes(cont.S("nvPair","BPDUGUARD_ENABLED").String()))
	d.Set("porttype_fast_enabled", stripQuotes(cont.S("nvPair","PORTTYPE_FAST_ENABLED").String()))
	d.Set("admin_state", stripQuotes(cont.S("nvPair","ADMIN_STATE").String()))
	d.Set("source_if_name", stripQuotes(cont.S("nvPair","SOURCE_IF_NAME").String()))
	d.Set("source_fabric_name", stripQuotes(cont.S("nvPair","SOURCE_FABRIC_NAME").String()))
	d.Set("source_switch_name", stripQuotes(cont.S("nvPair","SOURCE_SWITCH_NAME").String()))
	d.Set("link_uuid", stripQuotes(cont.S("nvPair","LINK_UUID").String()))
	d.Set("priority", stripQuotes(cont.S("nvPair","PRIORITY").String()))
	d.Set("dest_fabric_name", stripQuotes(cont.S("nvPair","DEST_FABRIC_NAME").String()))
	d.Set("policy_id", stripQuotes(cont.S("nvPair","POLICY_ID").String()))
	d.Set("dest_switch_name", stripQuotes(cont.S("nvPair","DEST_SWITCH_NAME").String()))
	d.Set("is_metaswitch", stripQuotes(cont.S("nvPair","IS_METASWITCH").String()))
	d.Set("dest_if_name", stripQuotes(cont.S("nvPair","DEST_IF_NAME").String()))
	d.Set("dest_serial_number", stripQuotes(cont.S("nvPair","DEST_SERIAL_NUMBER").String()))
	d.Set("source_serial_number", stripQuotes(cont.S("nvPair","SOURCE_SERIAL_NUMBER").String()))
	d.Set("policy_description", stripQuotes(cont.S("nvPair","POLICY_DESC").String()))

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
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/?attached-fabric=%s", serviceNode.AttachedFabricName)
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes", serviceNode.FabricName)
	}

	_, err := dcnmClient.Save(durl, &serviceNode)
	if err != nil {
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
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/?attached-fabric=%s", serviceNode.AttachedFabricName)
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", serviceNode.FabricName,serviceNode.Name)
	}

	_, err := dcnmClient.Update(durl, &serviceNode)
	if err != nil {
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

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes", fabricName)
	} else {
		durl = fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", fabricName,name)
	}

	cont, err := dcnmClient.GetviaURL(durl)
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

	durl := fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s", serviceFabric, d.Id())
	_, err := dcnmClient.Delete(durl)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}
