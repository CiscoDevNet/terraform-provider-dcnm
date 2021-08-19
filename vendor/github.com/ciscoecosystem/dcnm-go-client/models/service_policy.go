package models

type ServicePolicy struct {
	PolicyName         string      `json:",omitempty"`
	FabricName         string      `json:",omitempty"`
	AttachedFabricName string      `json:",omitempty"`
	DestinationNetwork string      `json:",omitempty"`
	DestinationVrfName string      `json:",omitempty"`
	Enabled            bool        `json:",omitempty"`
	NextHopIp          string      `json:",omitempty"`
	PeeringName        string      `json:",omitempty"`
	PolicyTemplateName string      `json:",omitempty"`
	ReverseEnabled     bool        `json:",omitempty"`
	ReverseNextHopIp   string      `json:",omitempty"`
	ServiceNodeName    string      `json:",omitempty"`
	ServiceNodeType    string      `json:",omitempty"`
	SourceNetwork      string      `json:",omitempty"`
	SourceVrfName      string      `json:",omitempty"`
	Status             string      `json:",omitempty"`
	NvPairs            interface{} `json:",omitempty"`
}

func (servicepolicy *ServicePolicy) ToMap() (map[string]interface{}, error) {
	servicepolicyAttributeMap := make(map[string]interface{})

	A(servicepolicyAttributeMap, "policyName", servicepolicy.PolicyName)
	A(servicepolicyAttributeMap, "fabricName", servicepolicy.FabricName)
	A(servicepolicyAttributeMap, "attachedFabricName", servicepolicy.AttachedFabricName)
	A(servicepolicyAttributeMap, "destinationNetwork", servicepolicy.DestinationNetwork)
	A(servicepolicyAttributeMap, "destinationVrfName", servicepolicy.DestinationVrfName)
	A(servicepolicyAttributeMap, "enabled", servicepolicy.Enabled)
	A(servicepolicyAttributeMap, "nextHopIp", servicepolicy.NextHopIp)
	A(servicepolicyAttributeMap, "peeringName", servicepolicy.PeeringName)
	A(servicepolicyAttributeMap, "policyTemplateName", servicepolicy.PolicyTemplateName)
	A(servicepolicyAttributeMap, "reverseEnabled", servicepolicy.ReverseEnabled)
	A(servicepolicyAttributeMap, "reverseNextHopIp", servicepolicy.ReverseNextHopIp)
	A(servicepolicyAttributeMap, "serviceNodeName", servicepolicy.ServiceNodeName)
	A(servicepolicyAttributeMap, "serviceNodeType", servicepolicy.ServiceNodeType)
	A(servicepolicyAttributeMap, "sourceNetwork", servicepolicy.SourceNetwork)
	A(servicepolicyAttributeMap, "sourceVrfName", servicepolicy.SourceVrfName)
	A(servicepolicyAttributeMap, "status", servicepolicy.Status)

	if servicepolicy.NvPairs != nil {
		A(servicepolicyAttributeMap, "nvPairs", servicepolicy.NvPairs)
	}

	return servicepolicyAttributeMap, nil
}
