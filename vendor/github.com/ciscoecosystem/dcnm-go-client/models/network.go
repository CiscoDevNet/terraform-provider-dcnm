package models

type Network struct {
	Fabric                 string `json:",omitempty"`
	Name                   string `json:",omitempty"`
	DisplayName            string `json:",omitempty"`
	NetworkId              string `json:",omitempty"`
	Template               string `json:",omitempty"`
	Config                 string `json:",omitempty"`
	ExtensionTemplate      string `json:",omitempty"`
	VRF                    string `json:",omitempty"`
	ServiceNetworkTemplate string `json:",omitempty"`
	Source                 string `json:",omitempty"`
}

type NetworkProfileConfig struct {
	NetworkName        string `json:"networkName"`
	VRFName            string `json:"vrfName"`
	SegmentID          string `json:"segmentId"`
	Vlan               string `json:"vlanId"`
	MTU                string `json:"mtu"`
	GatewayIpv4        string `json:"gatewayIpAddress"`
	GatewayIPv6        string `json:"gatewayIpV6Address"`
	VlanName           string `json:"vlanName"`
	Description        string `json:"intfDescription"`
	SecondaryGate1     string `json:"secondaryGW1"`
	SecondaryGate2     string `json:"secondaryGW2"`
	SecondaryGate3     string `json:"secondaryGW3"`
	SecondaryGate4     string `json:"secondaryGW4"`
	ARPSuppFlag        bool   `json:"suppressArp"`
	IRFlag             bool   `json:"enableIR"`
	McastGroup         string `json:"mcastGroup"`
	DHCPServer1        string `json:"dhcpServerAddr1"`
	DHCPServer2        string `json:"dhcpServerAddr2"`
	DHCPServer3        string `json:"dhcpServerAddr3"`
	DHCPServerVRF      string `json:"vrfDhcp"`
	DHCPServerVRF2     string `json:"vrfDhcp2"`
	DHCPServerVRF3     string `json:"vrfDhcp3"`
	LookbackID         string `json:"loopbackId"`
	Tag                string `json:"tag"`
	TRMEnable          bool   `json:"trmEnabled"`
	RTBothFlag         bool   `json:"rtBothAuto"`
	L3GatewayEnable    bool   `json:"enableL3OnBorder"`
	L2OnlyFlag         bool   `json:"isLayer2Only"`
	EnableNetflow      bool   `json:"ENABLE_NETFLOW"`
	SVINetflowMonitor  string `json:"SVI_NETFLOW_MONITOR"`
	VLANNetflowMonitor string `json:"VLAN_NETFLOW_MONITOR"`
	NVEId              string `json:"nveId"`
}

func (network *Network) ToMap() (map[string]interface{}, error) {
	networkAttrMap := make(map[string]interface{})

	A(networkAttrMap, "fabric", network.Fabric)

	A(networkAttrMap, "networkName", network.Name)

	A(networkAttrMap, "displayName", network.DisplayName)

	A(networkAttrMap, "networkId", network.NetworkId)

	A(networkAttrMap, "networkTemplate", network.Template)

	A(networkAttrMap, "networkExtensionTemplate", network.ExtensionTemplate)

	A(networkAttrMap, "networkTemplateConfig", network.Config)

	A(networkAttrMap, "vrf", network.VRF)

	if network.ServiceNetworkTemplate != "" {
		A(networkAttrMap, "serviceNetworkTemplate", network.ServiceNetworkTemplate)
	}

	if network.Source != "" {
		A(networkAttrMap, "source", network.Source)
	}

	return networkAttrMap, nil
}
