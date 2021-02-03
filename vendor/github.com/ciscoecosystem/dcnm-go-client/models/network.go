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
	NetworkName     string `json:"networkName,omitempty"`
	VRFName         string `json:"vrfName,omitempty"`
	SegmentID       string `json:"segmentId,omitempty"`
	Vlan            int    `json:"vlanId,omitempty"`
	MTU             int    `json:"mtu,omitempty"`
	GatewayIpv4     string `json:"gatewayIpAddress,omitempty"`
	GatewayIPv6     string `json:"gatewayIpV6Address,omitempty"`
	VlanName        string `json:"vlanName,omitempty"`
	Description     string `json:"intfDescription,omitempty"`
	SecondaryGate1  string `json:"secondaryGW1,omitempty"`
	SecondaryGate2  string `json:"secondaryGW2,omitempty"`
	ARPSuppFlag     bool   `json:"suppressArp,omitempty"`
	IRFlag          bool   `json:"enableIR,omitempty"`
	McastGroup      string `json:"mcastGroup,omitempty"`
	DHCPServer1     string `json:"dhcpServerAddr1,omitempty"`
	DHCPServer2     string `json:"dhcpServerAddr2,omitempty"`
	DHCPServerVRF   string `json:"vrfDhcp,omitempty"`
	LookbackID      int    `json:"loopbackId,omitempty"`
	Tag             string `json:"tag,omitempty"`
	TRMEnable       bool   `json:"trmEnabled,omitempty"`
	RTBothFlag      bool   `json:"rtBothAuto,omitempty"`
	L3GatewayEnable bool   `json:"enableL3OnBorder,omitempty"`
	L2OnlyFlag      bool   `json:"isLayer2Only,omitempty"`
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
