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
