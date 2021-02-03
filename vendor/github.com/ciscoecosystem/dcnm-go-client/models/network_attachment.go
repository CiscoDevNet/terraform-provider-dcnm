package models

type NetworkAttach struct {
	Name       string      `json:",omitempty"`
	AttachList interface{} `json:",omitempty"`
}

func NewNetworkAttachment(networkName string, ianAttach []map[string]interface{}) *NetworkAttach {
	networkAttach := NetworkAttach{}

	networkAttach.Name = networkName
	attachList := make([]interface{}, 0, 1)
	for _, val := range ianAttach {
		attachList = append(attachList, val)
	}

	networkAttach.AttachList = attachList

	return &networkAttach
}

func (networkAttach *NetworkAttach) ToMap() (map[string]interface{}, error) {
	networkAttachMap := make(map[string]interface{})

	A(networkAttachMap, "networkName", networkAttach.Name)

	A(networkAttachMap, "lanAttachList", networkAttach.AttachList)

	return networkAttachMap, nil
}
