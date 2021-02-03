package models

type VRF struct {
	Fabric             string `json:",omitempty"`
	Name               string `json:",omitempty"`
	Id                 string `json:",omitempty"`
	Template           string `json:",omitempty"`
	Config             string `json:",omitempty"`
	ExtensionTemplate  string `json:",omitempty"`
	ServiceVRFTemplate string `json:",omitempty"`
	Source             string `json:",omitempty"`
}

type VRFProfileConfig struct {
	VrfName         string `json:"vrfName"`
	SegmentID       string `json:"vrfSegmentId"`
	Vlan            int    `json:"vrfVlanId,omitempty"`
	Mtu             int    `json:"mtu,omitempty"`
	Tag             string `json:"tag,omitempty"`
	VlanName        string `json:"vrfVlanName,omitempty"`
	Description     string `json:"vrfDescription,omitempty"`
	IntfDescription string `json:"vrfIntfDescription,omitempty"`
	BGP             int    `json:"maxBgpPaths,omitempty"`
	IBGP            int    `json:"maxIbgpPaths,omitempty"`
	TRM             string `json:"trmEnabled,omitempty"`
	RPexternal      string `json:"isRPExternal,omitempty"`
	Lookback        int    `json:"loopbackNumber,omitempty"`
	RPaddress       string `json:"rpAddress,omitempty"`
	Mcastaddr       string `json:"L3VniMcastGroup,omitempty"`
	IPv6Link        string `json:"ipv6LinkLocalFlag,omitempty"`
	Mcastgroup      string `json:"multicastGroup,omitempty"`
	TRMBGW          string `json:"trmBGWMSiteEnabled,omitempty"`
	AdhostRoute     string `json:"advertiseHostRouteFlag,omitempty"`
	AdDefaultRoute  string `json:"advertiseDefaultRouteFlag,omitempty"`
	StaticRoute     string `json:"configureStaticDefaultRouteFlag,omitempty"`
}

func (vrf *VRF) ToMap() (map[string]interface{}, error) {
	vrfAttributeMap := make(map[string]interface{})
	A(vrfAttributeMap, "fabric", vrf.Fabric)
	A(vrfAttributeMap, "vrfName", vrf.Name)
	A(vrfAttributeMap, "vrfId", vrf.Id)
	A(vrfAttributeMap, "vrfTemplate", vrf.Template)
	A(vrfAttributeMap, "vrfTemplateConfig", vrf.Config)
	if vrf.ExtensionTemplate != "" {
		A(vrfAttributeMap, "vrfExtensionTemplate", vrf.ExtensionTemplate)
	}
	if vrf.ServiceVRFTemplate != "" {
		A(vrfAttributeMap, "serviceVrfTemplate", vrf.ServiceVRFTemplate)
	}
	if vrf.Source != "" {
		A(vrfAttributeMap, "source", vrf.Source)
	}
	return vrfAttributeMap, nil
}
