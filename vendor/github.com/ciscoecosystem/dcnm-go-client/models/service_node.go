package models

type ServiceNode struct {
	Name                        string      `json:",omitempty"`
	Type                        string      `json:",omitempty"`
	FormFactor                  string      `json:",omitempty"`
	FabricName                  string      `json:",omitempty"`
	InterfaceName               string      `json:",omitempty"`
	LinkTemplateName            string      `json:",omitempty"`
	AttachedSwitchSn            string      `json:",omitempty"`
	AttachedSwitchInterfaceName string      `json:",omitempty"`
	AttachedFabricName          string      `json:",omitempty"`
	NVPairs                     interface{} `json:",omitempty"`
}

func NewServiceNode(serviceNode *ServiceNode, nvPairs map[string]interface{}) *ServiceNode {
	if nvPairs != nil {
		serviceNode.NVPairs = nvPairs
	}
	return serviceNode
}

func (servicenode *ServiceNode) ToMap() (map[string]interface{}, error) {
	servicenodeAttributeMap := make(map[string]interface{})
	A(servicenodeAttributeMap, "name", servicenode.Name)
	A(servicenodeAttributeMap, "type", servicenode.Type)
	A(servicenodeAttributeMap, "formFactor", servicenode.FormFactor)
	A(servicenodeAttributeMap, "fabricName", servicenode.FabricName)
	A(servicenodeAttributeMap, "interfaceName", servicenode.InterfaceName)
	A(servicenodeAttributeMap, "linkTemplateName", servicenode.LinkTemplateName)
	A(servicenodeAttributeMap, "attachedSwitchSn", servicenode.AttachedSwitchSn)
	A(servicenodeAttributeMap, "attachedSwitchInterfaceName", servicenode.AttachedSwitchInterfaceName)
	A(servicenodeAttributeMap, "attachedFabricName", servicenode.AttachedFabricName)
	if servicenode.NVPairs != nil {
		A(servicenodeAttributeMap, "nvPairs", servicenode.NVPairs)
	}
	return servicenodeAttributeMap, nil
}
