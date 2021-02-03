package models

type Inventory struct {
	SeedIP         string   `json:",omitempty"`
	V3auth         int      `json:",omitempty"`
	Username       string   `json:",omitempty"`
	Password       string   `json:",omitempty"`
	MaxHops        int      `json:",omitempty"`
	SecondTimeout  int      `json:",omitempty"`
	PreserveConfig string   `json:",omitempty"`
	Switches       []Switch `json:",omitempty"`
	Platform       string   `json:",omitempty"`
}

type Switch struct {
	Reachable   string `json:",omitempty"`
	Auth        string `json:",omitempty"`
	Known       string `json:",omitempty"`
	Valid       string `json:",omitempty"`
	Selectable  string `json:",omitempty"`
	SysName     string `json:",omitempty"`
	IP          string `json:",omitempty"`
	Platform    string `json:",omitempty"`
	Version     string `json:",omitempty"`
	LastChange  string `json:",omitempty"`
	Hops        int    `json:",omitempty"`
	DeviceIndex string `json:",omitempty"`
	StatReason  string `json:",omitempty"`
}

type SwitchRole struct {
	SerialNumber string `json:",omitempty"`
	Role         string `json:",omitempty"`
}

func NewSwitch(inv *Inventory, s *Switch) *Inventory {
	switchList := make([]Switch, 0, 1)

	if s != nil {
		switchList = append(switchList, *s)
	}

	(*inv).Switches = switchList
	return inv
}

func (s *Switch) MakeMap() (map[string]interface{}, error) {
	switchMap := make(map[string]interface{})

	A(switchMap, "reachable", s.Reachable)

	A(switchMap, "auth", s.Auth)

	A(switchMap, "known", s.Known)

	A(switchMap, "valid", s.Valid)

	A(switchMap, "selectable", s.Selectable)

	A(switchMap, "sysName", s.SysName)

	A(switchMap, "ipaddr", s.IP)

	A(switchMap, "platform", s.Platform)

	A(switchMap, "version", s.Version)

	A(switchMap, "lastChange", s.LastChange)

	A(switchMap, "hopCount", s.Hops)

	A(switchMap, "deviceIndex", s.DeviceIndex)

	A(switchMap, "statusReason", s.StatReason)

	return switchMap, nil
}

func (inv *Inventory) ToMap() (map[string]interface{}, error) {
	inventoryMap := make(map[string]interface{})

	A(inventoryMap, "seedIP", inv.SeedIP)

	A(inventoryMap, "snmpV3AuthProtocol", inv.V3auth)

	A(inventoryMap, "username", inv.Username)

	A(inventoryMap, "password", inv.Password)

	A(inventoryMap, "maxHops", inv.MaxHops)

	A(inventoryMap, "cdpSecondTimeout", inv.SecondTimeout)

	A(inventoryMap, "preserveConfig", inv.PreserveConfig)

	A(inventoryMap, "platform", inv.Platform)

	if len(inv.Switches) > 0 {
		switchList := make([]interface{}, 0, 1)
		for _, s := range inv.Switches {
			sMap, _ := s.MakeMap()

			switchList = append(switchList, sMap)
		}
		A(inventoryMap, "switches", switchList)
	}

	return inventoryMap, nil
}

func (sRole *SwitchRole) ToMap() (map[string]interface{}, error) {
	sroleMap := make(map[string]interface{})

	A(sroleMap, "serialNumber", sRole.SerialNumber)

	A(sroleMap, "role", sRole.Role)

	return sroleMap, nil
}
