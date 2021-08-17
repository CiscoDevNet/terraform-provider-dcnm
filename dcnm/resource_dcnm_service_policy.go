package dcnm

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDCNMServicePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMServicePolicyCreate,
		Read:   resourceDCNMServicePolicyRead,
		Update: resourceDCNMServicePolicyUpdate,
		Delete: resourceDCNMServicePolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDCNMServicePolicyImporter,
		},

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

			"destination_network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_network_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"last_update": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"peering_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policy_template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"reverse_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},

			"reverse_next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"route_map_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"service_node_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"service_node_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_network_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"status_details": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
			},

			"attach_details": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
			},

			"destination_interfaces": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
			},

			"source_interfaces": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
			},

			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"src_port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dest_port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"route_map_action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"next_hop_action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"reverse": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"reverse_next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"fwd_direction": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func setServicePolicyAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {

	d.Set("policy_name", stripQuotes(cont.S("policyName").String()))
	d.Set("fabric_name", stripQuotes(cont.S("fabricName").String()))
	d.Set("attached_fabric_name", stripQuotes(cont.S("attachedFabricName").String()))
	d.Set("destination_network", stripQuotes(cont.S("destinationNetwork").String()))
	d.Set("destination_vrf_name", stripQuotes(cont.S("destinationVRFName").String()))
	d.Set("enabled", stripQuotes(cont.S("enabled").String()))
	d.Set("last_update", stripQuotes(cont.S("lastUpdate").String()))
	d.Set("next_hop_ip", stripQuotes(cont.S("nextHopIP").String()))
	d.Set("peering_name", stripQuotes(cont.S("peeringName").String()))
	d.Set("policy_template_name", stripQuotes(cont.S("policyTemplateName").String()))
	d.Set("reverse_enabled", stripQuotes(cont.S("reverseEnabled").String()))
	d.Set("reverse_next_hop_ip", stripQuotes(cont.S("reverseNextHopIP").String()))
	d.Set("route_map_name", stripQuotes(cont.S("routeMapName").String()))
	d.Set("service_node_name", stripQuotes(cont.S("serviceNodeName").String()))
	d.Set("service_node_type", stripQuotes(cont.S("serviceNodeType").String()))
	d.Set("source_network", stripQuotes(cont.S("serviceNetwork").String()))
	d.Set("source_network_name", stripQuotes(cont.S("serviceNetworkName").String()))
	d.Set("source_vrf_name", stripQuotes(cont.S("sourceVRFName").String()))
	d.Set("status", stripQuotes(cont.S("status").String()))
	d.Set("status_details", stripQuotes(cont.S("statusDetails").String()))
	d.Set("attach_details", stripQuotes(cont.S("attachDetails").String()))
	d.Set("destination_interfaces", stripQuotes(cont.S("destinationInterfaces").String()))
	d.Set("source_interfaces", stripQuotes(cont.S("sourceInterfaces").String()))
	d.Set("protocol", stripQuotes(cont.S("nvPair","PROTOCOL").String()))
	d.Set("src_port", stripQuotes(cont.S("nvPair","SRC_PORT").String()))
	d.Set("dest_port", stripQuotes(cont.S("nvPair","DEST_PORT").String()))
	d.Set("route_map_action", stripQuotes(cont.S("nvPair","ROUTE_MAP_ACTION").String()))
	d.Set("next_hop_action", stripQuotes(cont.S("nvPair","NEXT_HOP_ACTION").String()))
	d.Set("reverse", stripQuotes(cont.S("nvPair","REVERSE").String()))
	d.Set("reverse_next_hop_ip", stripQuotes(cont.S("nvPair","REVERSE_NEXT_HOP_IP").String()))
	d.Set("fwd_direction", stripQuotes(cont.S("nvPair","FWD_DIRECTION").String()))

	d.SetId(stripQuotes(cont.S("name").String()))
	return d

}

func resourceDCNMServicePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	attachedFabricName := d.Get("attached_fabric_name").(string)
	serviceNodeName:=     d.Get("service_node_name").(string)

	servicePolicy := models.ServicePolicy{
		PolicyName:                         d.Get("policy_name").(string),
		FabricName:                         d.Get("fabric_name").(string),
		AttachedFabricName:                 attachedFabricName,
		DestinationNetwork:                 d.Get("destination_network").(string),
		DestinationNetworkName:             d.Get("destination_network_name").(string),
		DestinationVRFName: 				d.Get("destination_vrf_name").(string),
		NextHopIP:                          d.Get("next_hop_ip").(string),
		PeeringName:                        d.Get("peering_name").(string),
		PolicyTemplateName: 				d.Get("policy_template_name").(string),
		ReverseEnabled:                     d.Get("reverse_enabled").(string),
		ReverseNextHopIP:                   d.Get("reverse_next_hop_ip").(string),
		ServiceNodeName: 				    serviceNodeName,
		ServiceNodeType:                    d.Get("service_node_type").(string),
		SourceNetwork:                      d.Get("source_network").(string),
		SourceNetworkName: 				    d.Get("source_network_name").(string),
		SourceVRFName:                      d.Get("source_vrf_name").(string),
		
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		servicePolicy.Enabled = enabled.(string)
	}

	if lastUpdate, ok := d.GetOk("last_update"); ok {
		servicePolicy.LastUpdate = lastUpdate.(string)
	}

	if routeMapName, ok := d.GetOk("route_map_name"); ok {
		servicePolicy.RouteMapName = routeMapName.(string)
	}

	if status, ok := d.GetOk("status"); ok {
		servicePolicy.Status = status.(string)
	}

	if statusDetails, ok := d.GetOk("status_details"); ok {
		servicePolicy.StatusDetails = statusDetails.(string)
	}

	if attachDetails, ok := d.GetOk("attach_details"); ok {
		servicePolicy.AttachDetails = attachDetails.(string)
	}

	if destinationInterfaces, ok := d.GetOk("destination_interfaces"); ok {
		servicePolicy.DestinationInterfaces = destinationInterfaces.(string)
	}

	if sourceInterfaces, ok := d.GetOk("source_interfaces"); ok {
		servicePolicy.SourceInterfaces = sourceInterfaces.(string)
	}

	nvPairMap := make(map[string]interface{})

	if protocol, ok := d.GetOk("protocol"); ok {
		nvPairMap["PROTOCOL"] = protocol.(string)
	} else {
		nvPairMap["PROTOCOL"] = ""
	}

	if srcPort, ok := d.GetOk("src_port"); ok {
		nvPairMap["SRC_PORT"] = srcPort.(string)
	} else {
		nvPairMap["SRC_PORT"] = ""
	}

	if destPort, ok := d.GetOk("dest_port"); ok {
		nvPairMap["DEST_PORT"] = destPort.(string)
	} else {
		nvPairMap["DEST_PORT"] = ""
	}

	if routeMapAction, ok := d.GetOk("route_map_action"); ok {
		nvPairMap["ROUTE_MAP_ACTION"] = routeMapAction.(string)
	} else {
		nvPairMap["ROUTE_MAP_ACTION"] = ""
	}

	if nextHopAction, ok := d.GetOk("next_hop_action"); ok {
		nvPairMap["NEXT_HOP_ACTION"] = nextHopAction.(string)
	} else {
		nvPairMap["NEXT_HOP_ACTION"] = ""
	}

	
	if reverse, ok := d.GetOk("reverse"); ok {
		nvPairMap["REVERSE"] = reverse.(string)
	}
	if reverseNextHopIP, ok := d.GetOk("reverse_next_hop_ip"); ok {
		nvPairMap["REVERSE_NEXT_HOP_IP"] = reverseNextHopIP.(string)
	}
	if fwdDirection, ok := d.GetOk("fwd_direction"); ok {
		nvPairMap["FWD_DIRECTION"] = fwdDirection.(string)
	}
	
	if nvPairMap != nil {
		servicePolicy.NVPairs = nvPairMap
	}

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("												", 						)
	} else {
		durl = fmt.Sprintf("/fabrics​/%s/service-nodes​/%s/policies", servicePolicy.FabricName, servicePolicy.ServiceNodeName)
	}

	_, err := dcnmClient.Save(durl, &servicePolicy)
	if err != nil {
		return err
	}

	d.SetId(servicePolicy.PolicyName)
	log.Println("[DEBUG] End of Create ", d.Id())
	return resourceDCNMServicePolicyRead(d, m)
}

func resourceDCNMServicePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ")

	dcnmClient := m.(*client.Client)

	attachedFabricName := d.Get("attached_fabric_name").(string)
	serviceNodeName:=     d.Get("service_node_name").(string)

	servicePolicy := models.ServicePolicy{
		PolicyName:                         d.Get("policy_name").(string),
		FabricName:                         d.Get("fabric_name").(string),
		AttachedFabricName:                 attachedFabricName,
		DestinationNetwork:                 d.Get("destination_network").(string),
		DestinationNetworkName:             d.Get("destination_network_name").(string),
		DestinationVRFName: 				d.Get("destination_vrf_name").(string),
		NextHopIP:                          d.Get("next_hop_ip").(string),
		PeeringName:                        d.Get("peering_name").(string),
		PolicyTemplateName: 				d.Get("policy_template_name").(string),
		ReverseEnabled:                     d.Get("reverse_enabled").(string),
		ReverseNextHopIP:                   d.Get("reverse_next_hop_ip").(string),
		ServiceNodeName: 				    serviceNodeName,
		ServiceNodeType:                    d.Get("service_node_type").(string),
		SourceNetwork:                      d.Get("source_network").(string),
		SourceNetworkName: 				    d.Get("source_network_name").(string),
		SourceVRFName:                      d.Get("source_vrf_name").(string),
		
	}

	if enabled, ok := d.GetOk("enabled"); ok {
		servicePolicy.Enabled = enabled.(string)
	}

	if lastUpdate, ok := d.GetOk("last_update"); ok {
		servicePolicy.LastUpdate = lastUpdate.(string)
	}

	if routeMapName, ok := d.GetOk("route_map_name"); ok {
		servicePolicy.RouteMapName = routeMapName.(string)
	}

	if status, ok := d.GetOk("status"); ok {
		servicePolicy.Status = status.(string)
	}

	if statusDetails, ok := d.GetOk("status_details"); ok {
		servicePolicy.StatusDetails = statusDetails.(string)
	}

	if attachDetails, ok := d.GetOk("attach_details"); ok {
		servicePolicy.AttachDetails = attachDetails.(string)
	}

	if destinationInterfaces, ok := d.GetOk("destination_interfaces"); ok {
		servicePolicy.DestinationInterfaces = destinationInterfaces.(string)
	}

	if sourceInterfaces, ok := d.GetOk("source_interfaces"); ok {
		servicePolicy.SourceInterfaces = sourceInterfaces.(string)
	}

	nvPairMap := make(map[string]interface{})

	if protocol, ok := d.GetOk("protocol"); ok {
		nvPairMap["PROTOCOL"] = protocol.(string)
	} else {
		nvPairMap["PROTOCOL"] = ""
	}

	if srcPort, ok := d.GetOk("src_port"); ok {
		nvPairMap["SRC_PORT"] = srcPort.(string)
	} else {
		nvPairMap["SRC_PORT"] = ""
	}

	if destPort, ok := d.GetOk("dest_port"); ok {
		nvPairMap["DEST_PORT"] = destPort.(string)
	} else {
		nvPairMap["DEST_PORT"] = ""
	}

	if routeMapAction, ok := d.GetOk("route_map_action"); ok {
		nvPairMap["ROUTE_MAP_ACTION"] = routeMapAction.(string)
	} else {
		nvPairMap["ROUTE_MAP_ACTION"] = ""
	}

	if nextHopAction, ok := d.GetOk("next_hop_action"); ok {
		nvPairMap["NEXT_HOP_ACTION"] = nextHopAction.(string)
	} else {
		nvPairMap["NEXT_HOP_ACTION"] = ""
	}

	
	if reverse, ok := d.GetOk("reverse"); ok {
		nvPairMap["REVERSE"] = reverse.(string)
	}
	if reverseNextHopIP, ok := d.GetOk("reverse_next_hop_ip"); ok {
		nvPairMap["REVERSE_NEXT_HOP_IP"] = reverseNextHopIP.(string)
	}
	if fwdDirection, ok := d.GetOk("fwd_direction"); ok {
		nvPairMap["FWD_DIRECTION"] = fwdDirection.(string)
	}
	
	if nvPairMap != nil {
		servicePolicy.NVPairs = nvPairMap
	}

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("												", 						)
	} else {
		durl = fmt.Sprintf("/fabrics​/%s/service-nodes​/%s/policies/%s/%s", servicePolicy.FabricName, servicePolicy.ServiceNodeName, servicePolicy.AttachedFabricName, servicePolicy.PolicyName)
	}

	_, err := dcnmClient.Update(durl, &servicePolicy)
	if err != nil {
		return err
	}

	d.SetId(servicePolicy.PolicyName)
	log.Println("[DEBUG] End of Create ", d.Id())
	return resourceDCNMServicePolicyRead(d, m)
}


func resourceDCNMServicePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	nodeName := d.Get("service_node_name").(string)
	fabricName := d.Get("fabric_name").(string)

	var durl string
	if dcnmClient.GetPlatform() == "nd" {
		durl = fmt.Sprintf("																",  )
	} else {
		durl = fmt.Sprintf("/fabrics​/%s/service-nodes​/%s/policies/", fabricName, nodeName)
	}

	cont, err := dcnmClient.GetviaURL(durl)
	if err != nil {
		d.SetId("")
		return nil
	}

	setServicePolicyAttributes(d, cont)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMServicePolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)
	nodeName := d.Get("service_node_name").(string)
	fabricName := d.Get("fabric_name").(string)

	durl := fmt.Sprintf("/fabrics​/%s/service-nodes​/%s/policies/", fabricName, nodeName)
	_, err := dcnmClient.Delete(durl)
	if err != nil {
		return err
	}
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}