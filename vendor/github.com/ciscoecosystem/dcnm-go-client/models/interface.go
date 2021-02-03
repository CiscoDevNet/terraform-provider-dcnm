package models

type Interface struct {
	Policy            string            `json:",omitempty"`
	Type              string            `json:",omitempty"`
	Interfaces        []InterfaceConfig `json:",omitempty"`
	SkipResourceCheck bool              `json:",omitempty"`
}

type InterfaceConfig struct {
	SerialNumber  string      `json:",omitempty"`
	InterfaceType string      `json:",omitempty"`
	InterfaceName string      `json:",omitempty"`
	Fabric        string      `json:",omitempty"`
	NVPairs       interface{} `json:",omitempty"`
}

type InterfaceDelete struct {
	SerialNumber string `json:",omitempty"`
	Name         string `json:",omitempty"`
}

func NewInterface(intf *Interface, intfConf *InterfaceConfig, nvPairs map[string]interface{}) *Interface {
	intfList := make([]InterfaceConfig, 0, 1)

	if intfConf != nil {
		if nvPairs != nil {
			(*intfConf).NVPairs = nvPairs
		}
		intfList = append(intfList, *intfConf)
	}

	(*intf).Interfaces = intfList

	return intf
}

func (intfConf *InterfaceConfig) makeConfMap() map[string]interface{} {
	interfaceConfMap := make(map[string]interface{})

	A(interfaceConfMap, "serialNumber", intfConf.SerialNumber)

	A(interfaceConfMap, "interfaceType", intfConf.InterfaceType)

	A(interfaceConfMap, "ifName", intfConf.InterfaceName)

	A(interfaceConfMap, "fabricName", intfConf.Fabric)

	if intfConf.NVPairs != nil {
		A(interfaceConfMap, "nvPairs", intfConf.NVPairs)
	}

	return interfaceConfMap
}

func (intf *Interface) ToMap() (map[string]interface{}, error) {
	interfaceMap := make(map[string]interface{})

	A(interfaceMap, "policy", intf.Policy)

	A(interfaceMap, "interfaceType", intf.Type)

	A(interfaceMap, "skipResourceCheck", intf.SkipResourceCheck)

	if len(intf.Interfaces) > 0 {
		intfList := make([]interface{}, 0, 1)
		for _, val := range intf.Interfaces {
			intfMap := val.makeConfMap()
			intfList = append(intfList, intfMap)
		}

		A(interfaceMap, "interfaces", intfList)
	}

	return interfaceMap, nil
}

func (intfDel *InterfaceDelete) ToMap() (map[string]interface{}, error) {
	intfDelMap := make(map[string]interface{})

	A(intfDelMap, "serialNumber", intfDel.SerialNumber)

	A(intfDelMap, "ifName", intfDel.Name)

	return intfDelMap, nil
}
