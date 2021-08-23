package dcnm

import (
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
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fabric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"attached_fabric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dest_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dest_vrf_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"next_hop_ip": {
				Type:     schema.TypeString,
				Required: true,
			},

			"peering_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policy_template_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "service_pbr",
			},

			"reverse_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_node_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_vrf_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ip",
			},

			"src_port": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
			},

			"dest_port": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
			},

			"route_map_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"permit",
					"deny",
				}, false),
				Default: "permit",
			},

			"next_hop_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"drop-on-fail",
					"drop",
				}, false),
				Default: "none",
			},

			"fwd_direction": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deploy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"deploy_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
		},
	}
}

var servicePolicyURLs = map[string]map[string]string{
	"dcnm": {
		"Create": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/policies",
		"Common": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/policies/%s/%s",
		"Deploy": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/policies/%s/deployments",
		"Attach": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/policies/%s/attachments",
	},
	"nd": {
		"Create": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/policies",
		"Common": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/policies/%s/%s",
		"Deploy": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/policies/%s/deployments",
		"Attach": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/policies/%s/attachments",
	},
}

func getServicePolicy(client *client.Client, attachedFabricName, fabricName, serviceNodeName, name string) (*container.Container, error) {
	dURL := fmt.Sprintf(servicePolicyURLs[client.GetPlatform()]["Common"], fabricName, serviceNodeName, attachedFabricName, name)
	cont, err := client.GetviaURL(dURL)
	return cont, getErrorFromContainer(cont, err)
}

func setServicePolicyAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	policyName := stripQuotes(cont.S("policyName").String())
	fabricName := stripQuotes(cont.S("fabricName").String())
	attachedFabricName := stripQuotes(cont.S("attachedFabricName").String())
	serviceNodeName := stripQuotes(cont.S("serviceNodeName").String())

	d.Set("policy_name", policyName)
	d.Set("fabric_name", fabricName)
	d.Set("attached_fabric_name", attachedFabricName)
	d.Set("dest_network", stripQuotes(cont.S("destinationNetwork").String()))
	d.Set("dest_vrf_name", stripQuotes(cont.S("destinationVrfName").String()))
	d.Set("next_hop_ip", stripQuotes(cont.S("nextHopIp").String()))
	d.Set("peering_name", stripQuotes(cont.S("peeringName").String()))
	d.Set("policy_template_name", stripQuotes(cont.S("policyTemplateName").String()))
	d.Set("service_node_name", serviceNodeName)
	d.Set("source_network", stripQuotes(cont.S("sourceNetwork").String()))
	d.Set("source_vrf_name", stripQuotes(cont.S("sourceVrfName").String()))
	d.Set("status", stripQuotes(cont.S("status").String()))
	d.Set("protocol", stripQuotes(cont.S("nvPairs", "PROTOCOL").String()))
	d.Set("src_port", stripQuotes(cont.S("nvPairs", "SRC_PORT").String()))
	d.Set("dest_port", stripQuotes(cont.S("nvPairs", "DEST_PORT").String()))
	d.Set("next_hop_action", stripQuotes(cont.S("nvPairs", "NEXT_HOP_ACTION").String()))
	if reverseEnabled, err := strconv.ParseBool(stripQuotes(cont.S("reverseEnabled").String())); err == nil {
		d.Set("reverse_enabled", reverseEnabled)
	}
	if attach, err := strconv.ParseBool(stripQuotes(cont.S("enabled").String())); err == nil {
		d.Set("deploy", attach)
	}
	if fwdPassword, err := strconv.ParseBool(stripQuotes(cont.S("nvPairs", "FWD_DIRECTION").String())); err == nil {
		d.Set("fwd_direction", fwdPassword)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s/%s", fabricName, serviceNodeName, attachedFabricName, policyName))
	return d
}

func resourceDCNMServicePolicyImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 2 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	attachedFabricName := importInfo[0]
	fabricName := importInfo[1]
	serviceNodeName := importInfo[2]
	policyName := importInfo[3]

	cont, err := getServicePolicy(dcnmClient, attachedFabricName, fabricName, serviceNodeName, policyName)
	if err != nil {
		return nil, err
	}
	stateImport := setServicePolicyAttributes(d, cont)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return []*schema.ResourceData{stateImport}, nil
}

func resourceDCNMServicePolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")
	dcnmClient := m.(*client.Client)

	policyName := d.Get("policy_name").(string)
	fabricName := d.Get("fabric_name").(string)
	attachedFabricName := d.Get("attached_fabric_name").(string)
	serviceNodeName := d.Get("service_node_name").(string)
	peeringName := d.Get("peering_name").(string)

	peeringCont, err := getRoutePeering(dcnmClient, attachedFabricName, fabricName, serviceNodeName, peeringName)
	if err != nil {
		return err
	}

	servicePolicy := models.ServicePolicy{
		PolicyName:         policyName,
		FabricName:         fabricName,
		AttachedFabricName: attachedFabricName,
		DestinationNetwork: d.Get("dest_network").(string),
		Enabled:            d.Get("deploy").(bool),
		DestinationVrfName: d.Get("dest_vrf_name").(string),
		NextHopIp:          d.Get("next_hop_ip").(string),
		PeeringName:        peeringName,
		PolicyTemplateName: d.Get("policy_template_name").(string),
		ReverseEnabled:     d.Get("reverse_enabled").(bool),
		ReverseNextHopIp:   stripQuotes(peeringCont.S("reverseNextHopIp").String()),
		ServiceNodeName:    serviceNodeName,
		ServiceNodeType:    stripQuotes(peeringCont.S("serviceNodeType").String()),
		SourceNetwork:      d.Get("source_network").(string),
		SourceVrfName:      d.Get("source_vrf_name").(string),
	}

	if attach, ok := d.GetOk("attach"); ok {
		servicePolicy.Enabled = attach.(bool)
	}

	nvPairMap := make(map[string]interface{})

	if protocol, ok := d.GetOk("protocol"); ok {
		nvPairMap["PROTOCOL"] = protocol.(string)
	}

	if srcPort, ok := d.GetOk("src_port"); ok {
		nvPairMap["SRC_PORT"] = srcPort.(string)
	}

	if destPort, ok := d.GetOk("dest_port"); ok {
		nvPairMap["DEST_PORT"] = destPort.(string)
	}

	if routeMapAction, ok := d.GetOk("route_map_action"); ok {
		nvPairMap["ROUTE_MAP_ACTION"] = routeMapAction.(string)
	}

	if nextHopAction, ok := d.GetOk("next_hop_action"); ok {
		nvPairMap["NEXT_HOP_ACTION"] = nextHopAction.(string)
	}

	if fwdDirection, ok := d.GetOk("fwd_direction"); ok {
		nvPairMap["FWD_DIRECTION"] = fwdDirection.(bool)
	}

	nvPairMap["REVERSE"] = servicePolicy.ReverseEnabled
	nvPairMap["REVERSE_NEXT_HOP_IP"] = servicePolicy.ReverseNextHopIp
	nvPairMap["NEXT_HOP_IP"] = servicePolicy.NextHopIp

	if nvPairMap != nil {
		servicePolicy.NvPairs = nvPairMap
	}

	durl := fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Create"], fabricName, serviceNodeName)

	cont, err := dcnmClient.Save(durl, &servicePolicy)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		deployModel := models.ServicePolicyDeploy{
			PolicyNames: []string{policyName},
		}
		log.Println("[DEBUG] Begining of Deploy Method.")

		//attach policy
		dURL := fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Attach"], fabricName, serviceNodeName, attachedFabricName)
		cont, err := dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			d.Set("deploy", false)
			return getErrorFromContainer(cont, err)
		}

		//deploy policy
		dURL = fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Deploy"], fabricName, serviceNodeName, attachedFabricName)

		cont, err = dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			d.Set("deploy", false)
			return getErrorFromContainer(cont, err)
		}
		deployFlag := false
		deployTimeout := d.Get("deploy_timeout").(int)
		for j := 0; j < (deployTimeout / 5); j++ {
			deployStatus, err := getServicePolicyDeploymentStatus(dcnmClient, attachedFabricName, fabricName, serviceNodeName, policyName)
			if err != nil {
				return err
			}
			deployFlag = (deployStatus == "Success" || deployStatus == "In-Sync")
			if !deployFlag {
				time.Sleep(5 * time.Second)
			} else {
				deployFlag = true
				break
			}
		}
		if !deployFlag {
			return fmt.Errorf("Service Policy record is created but not deployed yet. deployment timeout occured")
		}
		log.Println("[DEBUG] End of Deploy Method.")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", fabricName, serviceNodeName, attachedFabricName, policyName))
	log.Println("[DEBUG] End of Create ", d.Id())
	return resourceDCNMServicePolicyRead(d, m)
}

func resourceDCNMServicePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ")
	dcnmClient := m.(*client.Client)

	policyName := d.Get("policy_name").(string)
	fabricName := d.Get("fabric_name").(string)
	attachedFabricName := d.Get("attached_fabric_name").(string)
	serviceNodeName := d.Get("service_node_name").(string)
	peeringName := d.Get("peering_name").(string)

	peeringCont, err := getRoutePeering(dcnmClient, attachedFabricName, fabricName, serviceNodeName, peeringName)
	if err != nil {
		return err
	}

	servicePolicy := models.ServicePolicy{
		PolicyName:         policyName,
		FabricName:         fabricName,
		AttachedFabricName: attachedFabricName,
		DestinationNetwork: d.Get("dest_network").(string),
		Enabled:            d.Get("deploy").(bool),
		DestinationVrfName: d.Get("dest_vrf_name").(string),
		NextHopIp:          d.Get("next_hop_ip").(string),
		PeeringName:        peeringName,
		PolicyTemplateName: d.Get("policy_template_name").(string),
		ReverseEnabled:     d.Get("reverse_enabled").(bool),
		ReverseNextHopIp:   stripQuotes(peeringCont.S("reverseNextHopIp").String()),
		ServiceNodeName:    serviceNodeName,
		ServiceNodeType:    stripQuotes(peeringCont.S("serviceNodeType").String()),
		SourceNetwork:      d.Get("source_network").(string),
		SourceVrfName:      d.Get("source_vrf_name").(string),
	}

	if attach, ok := d.GetOk("attach"); ok {
		servicePolicy.Enabled = attach.(bool)
	}

	nvPairMap := make(map[string]interface{})

	if protocol, ok := d.GetOk("protocol"); ok {
		nvPairMap["PROTOCOL"] = protocol.(string)
	}

	if srcPort, ok := d.GetOk("src_port"); ok {
		nvPairMap["SRC_PORT"] = srcPort.(string)
	}

	if destPort, ok := d.GetOk("dest_port"); ok {
		nvPairMap["DEST_PORT"] = destPort.(string)
	}

	if routeMapAction, ok := d.GetOk("route_map_action"); ok {
		nvPairMap["ROUTE_MAP_ACTION"] = routeMapAction.(string)
	}

	if nextHopAction, ok := d.GetOk("next_hop_action"); ok {
		nvPairMap["NEXT_HOP_ACTION"] = nextHopAction.(string)
	}

	if fwdDirection, ok := d.GetOk("fwd_direction"); ok {
		nvPairMap["FWD_DIRECTION"] = fwdDirection.(bool)
	}

	nvPairMap["REVERSE"] = servicePolicy.ReverseEnabled
	nvPairMap["REVERSE_NEXT_HOP_IP"] = servicePolicy.ReverseNextHopIp
	nvPairMap["NEXT_HOP_IP"] = servicePolicy.NextHopIp

	if nvPairMap != nil {
		servicePolicy.NvPairs = nvPairMap
	}

	dURL := fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Common"], fabricName, serviceNodeName, attachedFabricName, policyName)

	cont, err := dcnmClient.Update(dURL, &servicePolicy)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		deployModel := models.ServicePolicyDeploy{
			PolicyNames: []string{policyName},
		}
		log.Println("[DEBUG] Begining of Deploy Method.")

		//attach policy
		dURL := fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Attach"], fabricName, serviceNodeName, attachedFabricName)
		cont, err := dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			d.Set("deploy", false)
			return getErrorFromContainer(cont, err)
		}

		//deploy policy
		dURL = fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Deploy"], fabricName, serviceNodeName, attachedFabricName)

		cont, err = dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			d.Set("deploy", false)
			return getErrorFromContainer(cont, err)
		}
		deployFlag := false
		deployTimeout := d.Get("deploy_timeout").(int)
		for j := 0; j < (deployTimeout / 5); j++ {
			deployStatus, err := getServicePolicyDeploymentStatus(dcnmClient, attachedFabricName, fabricName, serviceNodeName, policyName)
			if err != nil {
				return err
			}
			deployFlag = (deployStatus == "Success" || deployStatus == "In-Sync")
			if !deployFlag {
				time.Sleep(5 * time.Second)
			} else {
				deployFlag = true
				break
			}
		}
		if !deployFlag {
			return fmt.Errorf("Service Policy record is created but not deployed yet. deployment timeout occured")
		}
		log.Println("[DEBUG] End of Deploy Method.")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", fabricName, serviceNodeName, attachedFabricName, policyName))
	log.Println("[DEBUG] End of Update ", d.Id())
	return resourceDCNMServicePolicyRead(d, m)
}

func resourceDCNMServicePolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)
	policyId := strings.Split(d.Id(), "/")
	fabricName := policyId[0]
	serviceNodeName := policyId[1]
	attachedFabricName := policyId[2]
	policyName := policyId[3]

	cont, err := getServicePolicy(dcnmClient, attachedFabricName, fabricName, serviceNodeName, policyName)
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

	policyName := d.Get("policy_name").(string)
	fabricName := d.Get("fabric_name").(string)
	attachedFabricName := d.Get("attached_fabric_name").(string)
	serviceNodeName := d.Get("service_node_name").(string)

	dURL := fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Attach"]+"?policy-names=%s", fabricName, serviceNodeName, attachedFabricName, policyName)
	cont, err := dcnmClient.Delete(dURL)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}
	attachFlag := false
	deployTimeout := d.Get("deploy_timeout").(int)
	for j := 0; j < (deployTimeout / 2); j++ {
		cont, err := getServicePolicy(dcnmClient, attachedFabricName, fabricName, serviceNodeName, policyName)
		if err != nil {
			return getErrorFromContainer(cont, err)
		}
		attachFlag = stripQuotes(cont.S("enabled").String()) == "false"
		if !attachFlag {
			time.Sleep(2 * time.Second)
		} else {
			attachFlag = true
			break
		}
	}
	dURL = fmt.Sprintf(servicePolicyURLs[dcnmClient.GetPlatform()]["Common"], fabricName, serviceNodeName, attachedFabricName, policyName)
	cont, err = dcnmClient.Delete(dURL)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}

func getServicePolicyDeploymentStatus(dcnmClient *client.Client, attachedFabricName, extFabric, node, name string) (string, error) {
	cont, err := getServicePolicy(dcnmClient, attachedFabricName, extFabric, node, name)
	err = getErrorFromContainer(cont, err)
	status := stripQuotes(cont.S("status").String())
	return status, err
}
