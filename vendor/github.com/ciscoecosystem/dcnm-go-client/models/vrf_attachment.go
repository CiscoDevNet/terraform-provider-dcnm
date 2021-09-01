package models

type VRFAttach struct {
	Name       string      `json:",omitempty"`
	AttachList interface{} `json:",omitempty"`
}

type VRFInstance struct {
	LookbackID   int    `json:"loopbackId,omitempty"`
	LoopbackIpv4 string `json:"loopbackIpAddress,omitempty"`
	LoopbackIpv6 string `json:"loopbackIpV6Address,omitempty"`
}

type VRFDot1qID struct {
	ScopeType    string `json:"scopeType,omitempty"`
	UsageType    string `json:"usageType,omitempty"`
	AllocatedTo  string `json:"allocatedTo,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	IfName       string `json:"ifName,omitempty"`
}

type VRFDeploy struct {
	Name string `json:",omitempty"`
}

func NewVRFAttachment(vrfName string, ianAttach []map[string]interface{}) *VRFAttach {
	vrfAttach := VRFAttach{}

	vrfAttach.Name = vrfName

	attachList := make([]interface{}, 0, 1)
	for _, val := range ianAttach {
		attachList = append(attachList, val)
	}

	vrfAttach.AttachList = attachList

	return &vrfAttach
}

func (vrfAttach *VRFAttach) ToMap() (map[string]interface{}, error) {
	vrfAttachMap := make(map[string]interface{})

	A(vrfAttachMap, "vrfName", vrfAttach.Name)

	A(vrfAttachMap, "lanAttachList", vrfAttach.AttachList)

	return vrfAttachMap, nil
}

func (vrfDeploy *VRFDeploy) ToMap() (map[string]interface{}, error) {
	vrfDeployMap := make(map[string]interface{})

	A(vrfDeployMap, "vrfNames", vrfDeploy.Name)

	return vrfDeployMap, nil
}

func (vrfDot1qID *VRFDot1qID) ToMap() (map[string]interface{}, error) {
	vrfDot1qIDMap := make(map[string]interface{})

	A(vrfDot1qIDMap, "scopeType", vrfDot1qID.ScopeType)
	A(vrfDot1qIDMap, "allocatedTo", vrfDot1qID.AllocatedTo)
	A(vrfDot1qIDMap, "ifName", vrfDot1qID.IfName)
	A(vrfDot1qIDMap, "serialNumber", vrfDot1qID.SerialNumber)
	A(vrfDot1qIDMap, "usageType", vrfDot1qID.UsageType)

	return vrfDot1qIDMap, nil
}
