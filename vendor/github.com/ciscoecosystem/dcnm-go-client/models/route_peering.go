package models

type RoutePeering struct {
	AttachedFabricName string           `json:"attachedFabricName,omitempty"`
	DeploymentMode     string           `json:"deploymentMode,omitempty"`
	FabricName         string           `json:"fabricName,omitempty"`
	NextHopIP          string           `json:"nextHopIp,omitempty"`
	Name               string           `json:"peeringName,omitempty"`
	Option             string           `json:"peeringOption,omitempty"`
	ReverseNextHopIp   string           `json:"reverseNextHopIp,omitempty"`
	ServiceNetworks    []ServiceNetwork `json:"serviceNetworks,omitempty"`
	Routes             []RouteConfig    `json:"routes,omitempty"`
	ServiceNodeName    string           `json:"serviceNodeName,omitempty"`
	ServiceNodeType    string           `json:"serviceNodeType,omitempty"`
}

type ServiceNetwork struct {
	NetworkName  string      `json:"networkName,omitempty"`
	NetworkType  string      `json:"networkType,omitempty"`
	NVPairs      interface{} `json:"nvPairs,omitempty"`
	TemplateName string      `json:"templateName,omitempty"`
	Vlan         int         `json:"vlanId,omitempty"`
	VrfName      string      `json:"vrfName,omitempty"`
}

type RouteConfig struct {
	TemplateName string      `json:"templateName,omitempty"`
	VrfName      string      `json:"vrfName,omitempty"`
	NVPairs      interface{} `json:"nvPairs,omitempty"`
}
type RoutePeeringDeploy struct {
	PeeringNames []string `json:"peeringNames,omitempty"`
}

func NewNetwork(rp *RoutePeering, sn []*ServiceNetwork) *RoutePeering {
	snList := make([]ServiceNetwork, 0, 1)
	if sn != nil {
		for _, val := range sn {
			snList = append(snList, *val)
		}
	}
	(*rp).ServiceNetworks = snList
	return rp
}
func NewRoute(rp *RoutePeering, route []*RouteConfig) *RoutePeering {
	rList := make([]RouteConfig, 0, 1)
	if route != nil {
		for _, val := range route {
			rList = append(rList, *val)
		}
	}
	(*rp).Routes = rList
	return rp
}
func (serviceNetwork *ServiceNetwork) makeServiceMap() map[string]interface{} {
	serviceMap := make(map[string]interface{})
	A(serviceMap, "networkName", serviceNetwork.NetworkName)
	A(serviceMap, "networkType", serviceNetwork.NetworkType)
	if serviceNetwork.NVPairs != nil {
		A(serviceMap, "nvPairs", serviceNetwork.NVPairs)
	}
	A(serviceMap, "templateName", serviceNetwork.TemplateName)
	A(serviceMap, "vlanId", serviceNetwork.Vlan)
	A(serviceMap, "vrfName", serviceNetwork.VrfName)
	return serviceMap
}

func (route *RouteConfig) makeRouteMap() map[string]interface{} {
	routeMap := make(map[string]interface{})
	A(routeMap, "templateName", route.TemplateName)
	A(routeMap, "vrfName", route.VrfName)
	if route.NVPairs != nil {
		A(routeMap, "nvPairs", route.NVPairs)
	}
	return routeMap
}

func (routePeering *RoutePeering) ToMap() (map[string]interface{}, error) {
	peeringMap := make(map[string]interface{})
	A(peeringMap, "attachedFabricName", routePeering.AttachedFabricName)
	A(peeringMap, "deploymentMode", routePeering.DeploymentMode)
	A(peeringMap, "fabricName", routePeering.FabricName)
	A(peeringMap, "nextHopIp", routePeering.NextHopIP)
	A(peeringMap, "peeringName", routePeering.Name)
	A(peeringMap, "peeringOption", routePeering.Option)
	A(peeringMap, "reverseNextHopIp", routePeering.ReverseNextHopIp)
	if len(routePeering.ServiceNetworks) > 0 {
		netList := make([]interface{}, 0, 1)
		for _, val := range routePeering.ServiceNetworks {
			netMap := val.makeServiceMap()
			netList = append(netList, netMap)
		}
		A(peeringMap, "serviceNetworks", netList)

	}
	if len(routePeering.Routes) > 0 {
		routeList := make([]interface{}, 0, 1)
		for _, val := range routePeering.Routes {
			routeMap := val.makeRouteMap()
			routeList = append(routeList, routeMap)
		}
		A(peeringMap, "routes", routeList)
	}
	A(peeringMap, "serviceNodeName", routePeering.ServiceNodeName)
	A(peeringMap, "serviceNodeType", routePeering.ServiceNodeType)
	return peeringMap, nil

}

func (deploy *RoutePeeringDeploy) ToMap() (map[string]interface{}, error) {
	rDeploy := make(map[string]interface{})
	A(rDeploy, "peeringNames", deploy.PeeringNames)
	return rDeploy, nil
}
