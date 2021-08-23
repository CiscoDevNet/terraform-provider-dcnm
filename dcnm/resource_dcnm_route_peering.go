package dcnm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var URLS = map[string]map[string]string{
	"DCNMUrl": {
		"Create": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/peerings",
		"Common": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/peerings/%s/%s",
		"Deploy": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/peerings/%s/deployments",
		"Attach": "/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/%s/service-nodes/%s/peerings/%s/attachments",
	},
	"NDUrl": {
		"Create": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/peerings",
		"Common": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/peerings/%s/%s",
		"Deploy": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/peerings/%s/deployments",
		"Attach": "/appcenter/cisco/dcnm/api/v1/elastic-service/fabrics/%s/service-nodes/%s/peerings/%s/attachments",
	},
}

func resourceRoutePeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoutePeeringCreate,
		Update: resourceRoutePeeringUpdate,
		Read:   resourceRoutePeeringRead,
		Delete: resourceRoutePeeringDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRoutePeeringImporter,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attached_fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deployment_mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IntraTenantFW",
					"InterTenantFW",
					"OneArmADC",
					"TwoArmADC",
					"OneArmVNF",
				}, false),
			},
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"option": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
				Default:  nil,
			},
			"service_networks": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"network_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
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
							Required: true,
						},
						"vlan_id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"vrf_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"gateway_ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
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
				Required: true,
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
				Default: nil,
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

func resourceRoutePeeringImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Importer method", d.Id())
	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 4 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	AttachedFabricName := importInfo[3]
	extFabric := importInfo[1]
	node := importInfo[2]
	name := importInfo[0]
	cont, err := getRoutePeering(dcnmClient, AttachedFabricName, extFabric, node, name)
	if err != nil {
		if cont != nil {
			return nil, fmt.Errorf(cont.String())
		}
		return nil, err
	}
	stateImport := setPeeringAttributes(d, cont)
	flag, err := getRoutePeeringDeploymentStatus(dcnmClient, AttachedFabricName, extFabric, node, name)
	if err != nil {
		d.Set("deploy", false)
		return nil, err
	}
	d.Set("deploy", flag)
	log.Println("[DEBUG] End of Importer ", d.Id())
	return []*schema.ResourceData{stateImport}, nil

}
func resourceRoutePeeringCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining of Create Route Peering")
	dcnmClient := m.(*client.Client)

	rp := models.RoutePeering{}

	name := d.Get("name").(string)
	AttachedFabricName := d.Get("attached_fabric_name").(string)
	DeploymentMode := d.Get("deployment_mode").(string)
	FabricName := d.Get("fabric_name").(string)
	NextHopIp := d.Get("next_hop_ip").(string)
	Option := d.Get("option").(string)
	if ReverseNextHopIp, ok := d.GetOk("reverse_next_hop_ip"); ok {
		rp.ReverseNextHopIp = ReverseNextHopIp.(string)
	}
	ServiceNodeName := d.Get("service_node_name").(string)
	ServiceNodeType := d.Get("service_node_type").(string)
	Networks := d.Get("service_networks").(*schema.Set).List()

	rp.AttachedFabricName = AttachedFabricName
	rp.DeploymentMode = DeploymentMode
	rp.FabricName = FabricName
	rp.NextHopIP = NextHopIp
	rp.Name = name
	rp.Option = Option

	rp.ServiceNodeName = ServiceNodeName
	rp.ServiceNodeType = ServiceNodeType

	snObjs := make([]*models.ServiceNetwork, 0, 1)

	// Process the service network list

	for _, val := range Networks {
		netInfo := val.(map[string]interface{})
		nvPairMap := make(map[string]interface{})
		sNet := models.ServiceNetwork{}
		sNet.NetworkName = netInfo["network_name"].(string)
		sNet.NetworkType = netInfo["network_type"].(string)
		sNet.TemplateName = netInfo["template_name"].(string)
		sNet.VrfName = netInfo["vrf_name"].(string)
		sNet.Vlan = netInfo["vlan_id"].(int)
		nvPairMap["gatewayIpAddress"] = netInfo["gateway_ip_address"].(string)
		sNet.NVPairs = nvPairMap
		snObjs = append(snObjs, &sNet)

	}
	rpModel := models.NewNetwork(&rp, snObjs)

	// process the routes

	if r, ok := d.GetOk("routes"); ok {
		routeObjs := make([]*models.RouteConfig, 0, 1)

		routes := r.(*schema.Set).List()
		for _, val := range routes {
			rInfo := val.(map[string]interface{})
			rModel := models.RouteConfig{}
			if rInfo["route_parmas"] != nil {
				nvPairMap := rInfo["route_parmas"].(map[string]interface{})
				rModel.NVPairs = nvPairMap
			}
			if rInfo["template_name"] != nil {

				rModel.TemplateName = rInfo["template_name"].(string)
			}
			if rInfo["vrf_name"] != nil {

				rModel.VrfName = rInfo["vrf_name"].(string)
			}
			routeObjs = append(routeObjs, &rModel)
		}
		rpModel = models.NewRoute(rpModel, routeObjs)
	}
	var dURL string
	if dcnmClient.GetPlatform() == "nd" {
		dURL = fmt.Sprintf(URLS["NDUrl"]["Create"], FabricName, ServiceNodeName)
	} else {

		dURL = fmt.Sprintf(URLS["DCNMUrl"]["Create"], FabricName, ServiceNodeName)
	}
	cont, err := dcnmClient.Save(dURL, rpModel)
	if err != nil {
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}

	// Deploy the route peering
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		deployModel := models.RoutePeeringDeploy{}
		peeringNameList := make([]string, 0, 1)
		peeringNameList = append(peeringNameList, name)
		deployModel.PeeringNames = peeringNameList
		// attach the route peering
		if dcnmClient.GetPlatform() == "nd" {
			dURL = fmt.Sprintf(URLS["NDUrl"]["Attach"], FabricName, ServiceNodeName, AttachedFabricName)
		} else {

			dURL = fmt.Sprintf(URLS["DCNMUrl"]["Attach"], FabricName, ServiceNodeName, AttachedFabricName)
		}

		_, err = dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
		}

		// deploy
		log.Println("[DEBUG] Begining of Deploy Method.")
		if dcnmClient.GetPlatform() == "nd" {
			dURL = fmt.Sprintf(URLS["NDUrl"]["Deploy"], FabricName, ServiceNodeName, AttachedFabricName)
		} else {

			dURL = fmt.Sprintf(URLS["DCNMUrl"]["Deploy"], FabricName, ServiceNodeName, AttachedFabricName)
		}

		cont, err = dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			d.Set("deploy", false)
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
		}

		deployTFlag := false
		deployTimeout := d.Get("deploy_timeout").(int)
		for j := 0; j < (deployTimeout / 5); j++ {
			deployFlag, err := getRoutePeeringDeploymentStatus(dcnmClient, AttachedFabricName, FabricName, ServiceNodeName, name)
			if err != nil {
				return err
			}

			if !deployFlag {
				time.Sleep(5 * time.Second)
			} else {
				deployTFlag = true
				break
			}
		}
		if !deployTFlag {
			return fmt.Errorf("Route Peering record is created but not deployed yet. deployment timeout occured")
		}
		log.Println("[DEBUG] End of Deploy Method.")
	}

	d.SetId(fmt.Sprintf("/fabrics/%s/service-nodes/%s/peerings/%s",
		FabricName, ServiceNodeName, stripQuotes(cont.S("peeringName").String())))
	return resourceRoutePeeringRead(d, m)
}
func getRoutePeeringDeploymentStatus(dcnmClient *client.Client, AttachedFabricName, extFabric, node, name string) (bool, error) {
	cont, err := getRoutePeering(dcnmClient, AttachedFabricName, extFabric, node, name)
	status := stripQuotes(cont.S("status").String())
	if status != "Success" && status != "In-Sync" {
		return false, err
	}
	return true, err
}

func resourceRoutePeeringUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining of Update Route Peering", d.Id())

	dcnmClient := m.(*client.Client)

	rp := models.RoutePeering{}

	name := d.Get("name").(string)
	AttachedFabricName := d.Get("attached_fabric_name").(string)
	DeploymentMode := d.Get("deployment_mode").(string)
	FabricName := d.Get("fabric_name").(string)
	NextHopIp := d.Get("next_hop_ip").(string)
	Option := d.Get("option").(string)
	if ReverseNextHopIp, ok := d.GetOk("reverse_next_hop_ip"); ok {
		rp.ReverseNextHopIp = ReverseNextHopIp.(string)
	}
	ServiceNodeName := d.Get("service_node_name").(string)
	ServiceNodeType := d.Get("service_node_type").(string)
	Networks := d.Get("service_networks").(*schema.Set).List()
	rp.AttachedFabricName = AttachedFabricName
	rp.DeploymentMode = DeploymentMode
	rp.FabricName = FabricName
	rp.NextHopIP = NextHopIp
	rp.Name = name
	rp.Option = Option
	rp.ServiceNodeName = ServiceNodeName
	rp.ServiceNodeType = ServiceNodeType
	snObjs := make([]*models.ServiceNetwork, 0, 1)

	// Process the service network list

	for _, val := range Networks {
		netInfo := val.(map[string]interface{})
		nvPairMap := make(map[string]interface{})
		sNet := models.ServiceNetwork{}
		sNet.NetworkName = netInfo["network_name"].(string)
		sNet.NetworkType = netInfo["network_type"].(string)
		sNet.TemplateName = netInfo["template_name"].(string)
		sNet.VrfName = netInfo["vrf_name"].(string)
		sNet.Vlan = netInfo["vlan_id"].(int)
		nvPairMap["gatewayIpAddress"] = netInfo["gateway_ip_address"].(string)
		sNet.NVPairs = nvPairMap
		snObjs = append(snObjs, &sNet)

	}
	rpModel := models.NewNetwork(&rp, snObjs)

	// process the routes

	if r, ok := d.GetOk("routes"); ok {
		routeObjs := make([]*models.RouteConfig, 0, 1)

		routes := r.(*schema.Set).List()
		for _, val := range routes {
			rInfo := val.(map[string]interface{})
			rModel := models.RouteConfig{}
			if rInfo["route_parmas"] != nil {
				nvPairMap := rInfo["route_parmas"].(map[string]interface{})
				rModel.NVPairs = nvPairMap
			}
			if rInfo["template_name"] != nil {

				rModel.TemplateName = rInfo["template_name"].(string)
			}
			if rInfo["vrf_name"] != nil {

				rModel.VrfName = rInfo["vrf_name"].(string)
			}
			routeObjs = append(routeObjs, &rModel)
		}
		rpModel = models.NewRoute(rpModel, routeObjs)
	}
	var dURL string
	if dcnmClient.GetPlatform() == "nd" {
		dURL = fmt.Sprintf(URLS["NDUrl"]["Common"], FabricName, ServiceNodeName, AttachedFabricName, name)
	} else {
		dURL = fmt.Sprintf(URLS["DCNMUrl"]["Common"], FabricName, ServiceNodeName, AttachedFabricName, name)
	}
	cont, err := dcnmClient.Update(dURL, rpModel)
	if err != nil {
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		deployModel := models.RoutePeeringDeploy{}
		peeringNameList := make([]string, 0, 1)
		peeringNameList = append(peeringNameList, name)
		deployModel.PeeringNames = peeringNameList
		// attach the route peering
		if dcnmClient.GetPlatform() == "nd" {
			dURL = fmt.Sprintf(URLS["NDUrl"]["Attach"], FabricName, ServiceNodeName, AttachedFabricName)
		} else {

			dURL = fmt.Sprintf(URLS["DCNMUrl"]["Attach"], FabricName, ServiceNodeName, AttachedFabricName)
		}

		_, err = dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
		}

		// deploy
		log.Println("[DEBUG] Begining of Deploy Method.")
		if dcnmClient.GetPlatform() == "nd" {
			dURL = fmt.Sprintf(URLS["NDUrl"]["Deploy"], FabricName, ServiceNodeName, AttachedFabricName)
		} else {

			dURL = fmt.Sprintf(URLS["DCNMUrl"]["Deploy"], FabricName, ServiceNodeName, AttachedFabricName)
		}

		cont, err = dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			d.Set("deploy", false)
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
		}

		deployTFlag := false
		deployTimeout := d.Get("deploy_timeout").(int)
		for j := 0; j < (deployTimeout / 5); j++ {
			deployFlag, err := getRoutePeeringDeploymentStatus(dcnmClient, AttachedFabricName, FabricName, ServiceNodeName, name)
			if err != nil {
				return err
			}

			if !deployFlag {
				time.Sleep(5 * time.Second)
			} else {
				deployTFlag = true
				break
			}
		}
		if !deployTFlag {
			return fmt.Errorf("Route Peering  is created but not deployed yet. deployment timeout occured")
		}
		log.Println("[DEBUG] End of Deploy Method.")
	}
	d.SetId(fmt.Sprintf("/fabrics/%s/service-nodes/%s/peerings/%s", FabricName, ServiceNodeName, name))
	log.Println("[DEBUG] End of Update Route Peering", d.Id())
	return resourceRoutePeeringRead(d, m)
}
func resourceRoutePeeringDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining of Delete Route", d.Id())
	dcnmClient := m.(*client.Client)
	AttachedFabricName := d.Get("attached_fabric_name").(string)
	extFabric := d.Get("fabric_name").(string)
	node := d.Get("service_node_name").(string)
	name := d.Get("name").(string)
	var dURL string
	//  Detach the route peering
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {

		if dcnmClient.GetPlatform() == "nd" {
			dURL = fmt.Sprintf(URLS["NDUrl"]["Attach"]+"?peering-names=%s", extFabric, node, AttachedFabricName, name)
		} else {
			dURL = fmt.Sprintf(URLS["DCNMUrl"]["Attach"]+"?peering-names=%s", extFabric, node, AttachedFabricName, name)
		}
		cont, err := dcnmClient.Delete(dURL)

		if err != nil {
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
		}
		deployModel := models.RoutePeeringDeploy{}
		peeringNameList := make([]string, 0, 1)
		peeringNameList = append(peeringNameList, name)
		deployModel.PeeringNames = peeringNameList
		if dcnmClient.GetPlatform() == "nd" {
			dURL = fmt.Sprintf(URLS["NDUrl"]["Deploy"], extFabric, node, AttachedFabricName)
		} else {

			dURL = fmt.Sprintf(URLS["DCNMUrl"]["Deploy"], extFabric, node, AttachedFabricName)
		}

		cont, err = dcnmClient.Save(dURL, &deployModel)
		if err != nil {
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
		}
		deployTFlag := false
		deployTimeout := d.Get("deploy_timeout").(int)
		for j := 0; j < (deployTimeout / 5); j++ {
			cont, err := getRoutePeering(dcnmClient, AttachedFabricName, extFabric, node, name)
			status := stripQuotes(cont.S("status").String())
			if status == "NA" || status == "N/A" {
				deployTFlag = true
				break
			} else {
				time.Sleep(5 * time.Second)
			}
			if err != nil {
				return err
			}
		}
		if !deployTFlag {
			return fmt.Errorf("Route Peering  is created but not deployed yet. deployment timeout occured")
		}
		log.Println("[DEBUG] End of Deploy Method.")
	}
	if dcnmClient.GetPlatform() == "nd" {
		dURL = fmt.Sprintf(URLS["NDUrl"]["Common"], extFabric, node, AttachedFabricName, name)
	} else {
		dURL = fmt.Sprintf(URLS["DCNMUrl"]["Common"], extFabric, node, AttachedFabricName, name)
	}

	cont, err := dcnmClient.Delete(dURL)
	if err != nil {
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}
	return nil
}
func getRoutePeering(client *client.Client, AttachedFabricName, extFabric, node, name string) (*container.Container, error) {
	var dURL string
	if client.GetPlatform() == "nd" {
		dURL = fmt.Sprintf(URLS["NDUrl"]["Common"], extFabric, node, AttachedFabricName, name)
	} else {
		dURL = fmt.Sprintf(URLS["DCNMUrl"]["Common"], extFabric, node, AttachedFabricName, name)
	}
	cont, err := client.GetviaURL(dURL)
	return cont, err
}
func setPeeringAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	var name, FabricName, ServiceNodeName string
	if cont.Exists("peeringName") {
		name = stripQuotes(cont.S("peeringName").String())
		d.Set("name", stripQuotes(cont.S("peeringName").String()))
	}
	if cont.Exists("fabricName") {
		FabricName = stripQuotes(cont.S("fabricName").String())
		d.Set("fabric_name", stripQuotes(cont.S("fabricName").String()))
	}
	if cont.Exists("attachedFabricName") {
		d.Set("attached_fabric_name", stripQuotes(cont.S("attachedFabricName").String()))
	}
	if cont.Exists("deploymentMode") {
		d.Set("deployment_mode", stripQuotes(cont.S("deploymentMode").String()))
	}
	if cont.Exists("nextHopIp") {
		d.Set("next_hop_ip", stripQuotes(cont.S("nextHopIp").String()))
	}
	if cont.Exists("peeringOption") {
		d.Set("option", stripQuotes(cont.S("peeringOption").String()))
	}
	if cont.Exists("reverseNextHopIp") {
		d.Set("reverse_next_hop_ip", stripQuotes(cont.S("reverseNextHopIp").String()))
	}
	if cont.Exists("serviceNodeName") {
		ServiceNodeName = stripQuotes(cont.S("serviceNodeName").String())
		d.Set("service_node_name", stripQuotes(cont.S("serviceNodeName").String()))
	}
	if cont.Exists("serviceNodeType") {
		d.Set("service_node_type", stripQuotes(cont.S("serviceNodeType").String()))
	}

	serviceNetwork := make([]interface{}, 0, 1)
	var network string
	if cont.Exists("serviceNetworks") {
		network = stripQuotes(cont.S("serviceNetworks").String())
	}
	var netinfo []map[string]interface{}
	_ = json.Unmarshal([]byte(network), &netinfo)

	for i := 0; i < len(netinfo); i++ {
		netMap := make(map[string]interface{})
		netMap["network_name"] = netinfo[i]["networkName"].(string)
		netMap["network_type"] = netinfo[i]["networkType"].(string)
		netMap["template_name"] = netinfo[i]["templateName"].(string)
		netMap["vlan_id"] = netinfo[i]["vlanId"].(float64)
		netMap["vrf_name"] = netinfo[i]["vrfName"].(string)
		nvPairs := netinfo[i]["nvPairs"].(map[string]interface{})
		netMap["gateway_ip_address"] = nvPairs["gatewayIpAddress"].(string)
		serviceNetwork = append(serviceNetwork, netMap)
	}
	d.Set("service_networks", serviceNetwork)

	if cont.Exists("routes") {
		routeList := make([]interface{}, 0, 1)
		route := stripQuotes(cont.S("routes").String())
		var rinfo []map[string]interface{}
		_ = json.Unmarshal([]byte(route), &rinfo)
		for i := 0; i < len(rinfo); i++ {
			rMap := make(map[string]interface{})
			if rinfo[i]["templateName"] != nil {
				rMap["template_name"] = rinfo[i]["templateName"].(string)
			}
			if rinfo[i]["vrfName"] != nil {
				rMap["vrf_name"] = rinfo[i]["vrfName"].(string)
			}
			if rinfo[i]["nvPairs"] != nil {
				rMap["route_parmas"] = rinfo[i]["nvPairs"].(map[string]interface{})
			}
			routeList = append(routeList, rMap)

		}
		d.Set("routes", routeList)
	}
	d.SetId(fmt.Sprintf("/fabrics/%s/service-nodes/%s/peerings/%s", FabricName, ServiceNodeName, name))
	return d
}
func resourceRoutePeeringRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method", d.Id())
	dcnmClient := m.(*client.Client)

	AttachedFabricName := d.Get("attached_fabric_name").(string)
	extFabric := d.Get("fabric_name").(string)
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
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}
