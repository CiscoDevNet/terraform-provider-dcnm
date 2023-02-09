package models

type Fabric struct {
	Id       int          `json:",omitempty"`
	Name     string       `json:",omitempty"`
	FabricId string       `json:",omitempty"`
	Template string       `json:",omitempty"`
	Config   FabricConfig `json:",omitempty"`
}

type FabricConfig struct {
	Asn                            string `json:"BGP_AS"`
	UnderlayInterfaceNumbering     string `json:"FABRIC_INTERFACE_TYPE"`
	UnderlaySubnetMask             string `json:"SUBNET_TARGET_MASK"`
	UnderlayRoutingProcotol        string `json:"LINK_STATE_ROUTING"`
	RouteReflectorCount            string `json:"RR_COUNT"`
	AnycastMac                     string `json:"ANYCAST_GW_MAC"`
	ReplicationMode                string `json:"REPLICATION_MODE"`
	MulticastGroupSubnet           string `json:"MULTICAST_GROUP_SUBNET"`
	RendevouzPointCount            string `json:"RP_COUNT"`
	RendevouzPointMode             string `json:"RP_MODE"`
	RendevouzPointLoopbackId       string `json:"RP_LB_ID"`
	VpcPeerLinkVlan                string `json:"VPC_PEER_LINK_VLAN"`
	VpcPeerKeepAliveOption         string `json:"VPC_PEER_KEEP_ALIVE_OPTION"`
	VpcAutoRecoveryTime            string `json:"VPC_AUTO_RECOVERY_TIME"`
	VpcDelayRestore                string `json:"VPC_DELAY_RESTORE"`
	VpcDelayRestoretime            string `json:"VPC_DELAY_RESTORE_TIME"`
	UnderlayRoutingLoopbackId      string `json:"BGP_LB_ID"`
	UnderlayVtepLoopbackId         string `json:"NVE_LB_ID"`
	UnderlayRoutingProtocolTag     string `json:"LINK_STATE_ROUTING_TAG"`
	OspfAreaId                     string `json:"OSPF_AREA_ID"`
	BfdEnable                      string `json:"BFD_ENABLE"`
	BfdOspf                        string `json:"BFD_OSPF_ENABLE"`
	BfdAuthEnable                  string `json:"BFD_AUTH_ENABLE"`
	BfdAuthKeyId                   string `json:"BFD_AUTH_KEY_ID"`
	BfdAuthKey                     string `json:"BFD_AUTH_KEY"`
	IbgpPeerTemplate               string `json:"IBGP_PEER_TEMPLATE"`
	IbgpPeerTemplateLeaf           string `json:"IBGP_PEER_TEMPLATE_LEAF"`
	IbgpOspf                       string `json:"BFD_IBGP_ENABLE"`
	IsisOspf                       string `json:"BFD_ISIS_ENABLE"`
	PimOspf                        string `json:"BFD_PIM_ENABLE"`
	VrfTemplate                    string `json:"default_vrf"`
	NetworkTemplate                string `json:"default_network"`
	VrfExtensionTemplate           string `json:"vrf_extension_template"`
	NetworkExtensionTemplate       string `json:"network_extension_template"`
	OverlayModePrev                string `json:"OVERLAY_MODE_PREV"`
	IntraFabricInterfaceMtu        string `json:"FABRIC_MTU"`
	Layer2HostInterfaceMtu         string `json:"L2_HOST_INTF_MTU"`
	PowerSupplyMode                string `json:"POWER_REDUNDANCY_MODE"`
	CoppProfile                    string `json:"COPP_POLICY"`
	EnableNgoam                    string `json:"ENABLE_NGOAM"`
	EnableNxapi                    string `json:"ENABLE_NXAPI"`
	EnableNxapiHttp                string `json:"ENABLE_NXAPI_HTTP"`
	SnmpServerHostTrap             string `json:"SNMP_SERVER_HOST_TRAP"`
	UnderlayRoutingLoopbackIpRange string `json:"LOOPBACK0_IP_RANGE"`
	UnderlayVtepLoopbackIpRange    string `json:"LOOPBACK1_IP_RANGE"`
	UnderlayRpLoopbackIpRange      string `json:"ANYCAST_RP_IP_RANGE"`
	UnderlaySubnetIpRange          string `json:"SUBNET_RANGE"`
	Layer2VxlanVniRange            string `json:"L2_SEGMENT_ID_RANGE"`
	Layer3VxlanVniRange            string `json:"L3_PARTITION_ID_RANGE"`
	NetworkVlanRange               string `json:"NETWORK_VLAN_RANGE"`
	VrfVlanRange                   string `json:"VRF_VLAN_RANGE"`
	SubinterfaceDot1qRange         string `json:"SUBINTERFACE_RANGE"`
	VrfLiteDeployment              string `json:"VRF_LITE_AUTOCONFIG"`
	VrfLiteSubnetIpRange           string `json:"DCI_SUBNET_RANGE"`
	VrfLiteSubnetMask              string `json:"DCI_SUBNET_TARGET_MASK"`
	ServiceNetworkVlanRange        string `json:"SERVICE_NETWORK_VLAN_RANGE"`
	RouteMapSequenceNumberRange    string `json:"ROUTE_MAP_SEQUENCE_NUMBER_RANGE"`
	// Parameters not exposed in the Terraform resource
	MsoSiteId                      string `json:"MSO_SITE_ID"`
	PhantomRpLbId1                 string `json:"PHANTOM_RP_LB_ID1"`
	PhantomRpLbId2                 string `json:"PHANTOM_RP_LB_ID2"`
	PhantomRpLbId3                 string `json:"PHANTOM_RP_LB_ID3"`
	PhantomRpLbId4                 string `json:"PHANTOM_RP_LB_ID4"`
	AbstractOspf                   string `json:"abstract_ospf"`
	FeaturePtp                     string `json:"FEATURE_PTP"`
	DhcpStartInternal              string `json:"DHCP_START_INTERNAL"`
	SspineCount                    string `json:"SSPINE_COUNT"`
	AdvertisePipBgp                string `json:"ADVERTISE_PIP_BGP"`
	FabricVpcQosPolicyName         string `json:"FABRIC_VPC_QOS_POLICY_NAME"`
	DhcpEnd                        string `json:"DHCP_END"`
	UnderlayIsV6                   string `json:"UNDERLAY_IS_V6"`
	FabricVpcDomainId              string `json:"FABRIC_VPC_DOMAIN_ID"`
	SeedSwitchCoreInterfaces       string `json:"SEED_SWITCH_CORE_INTERFACES"`
	FabricMtuPrev                  string `json:"FABRIC_MTU_PREV"`
	HdTime                         string `json:"HD_TIME"`
	OspfAuthEnable                 string `json:"OSPF_AUTH_ENABLE"`
	Loopback1Ipv6Range             string `json:"LOOPBACK1_IPV6_RANGE"`
	RouterIdRange                  string `json:"ROUTER_ID_RANGE"`
	MsoConnectivityDeployed        string `json:"MSO_CONNECTIVITY_DEPLOYED"`
	EnableMacsec                   string `json:"ENABLE_MACSEC"`
	DeafultQueuingPolicyOther      string `json:"DEAFULT_QUEUING_POLICY_OTHER"`
	UnnumDhcpStartInternal         string `json:"UNNUM_DHCP_START_INTERNAL"`
	MacsecReportTimer              string `json:"MACSEC_REPORT_TIMER"`
	PremsoParentFabric             string `json:"PREMSO_PARENT_FABRIC"`
	UnnumDhcpEndInternal           string `json:"UNNUM_DHCP_END_INTERNAL"`
	PtpDomainId                    string `json:"PTP_DOMAIN_ID"`
	AutoSymmetricVrfLite           bool   `json:"AUTO_SYMMETRIC_VRF_LITE"`
	UseLinkLocal                   bool   `json:"USE_LINK_LOCAL"`
	BgpAsPrev                      string `json:"BGP_AS_PREV"`
	EnablePbr                      string `json:"ENABLE_PBR"`
	VpcPeerLinkPo                  string `json:"VPC_PEER_LINK_PO"`
	IsisAuthEnable                 bool   `json:"ISIS_AUTH_ENABLE"`
	VpcEnableIpv6NdSync            string `json:"VPC_ENABLE_IPv6_ND_SYNC"`
	AbstractIsisInterface          string `json:"abstract_isis_interface"`
	TcamAllocation                 string `json:"TCAM_ALLOCATION"`
	MacsecAlgorithm                string `json:"MACSEC_ALGORITHM"`
	IsisLevel                      string `json:"ISIS_LEVEL"`
	AbstractAnycastRp              string `json:"abstract_anycast_rp"`
	EnableNetflow                  string `json:"ENABLE_NETFLOW"`
	DeafultQueuingPolicyRSeries    string `json:"DEAFULT_QUEUING_POLICY_R_SERIES"`
	TempVpcPeerLink                string `json:"temp_vpc_peer_link"`
	BrownfieldNetworkNameFormat    string `json:"BROWNFIELD_NETWORK_NAME_FORMAT"`
	EnableFabricVpcDomainId        string `json:"ENABLE_FABRIC_VPC_DOMAIN_ID"`
	MgmtGwInternal                 string `json:"MGMT_GW_INTERNAL"`
	GrfieldDebugFlag               string `json:"GRFIELD_DEBUG_FLAG"`
	IsisAuthKeychainName           string `json:"ISIS_AUTH_KEYCHAIN_NAME"`
	AbstractBgpNeighbor            string `json:"abstract_bgp_neighbor"`
	OspfAuthKeyId                  string `json:"OSPF_AUTH_KEY_ID"`
	PimHelloAuthEnable             string `json:"PIM_HELLO_AUTH_ENABLE"`
	AbstractFeatureLeaf            string `json:"abstract_feature_leaf"`
	ExtraConfTor                   string `json:"EXTRA_CONF_TOR"`
	AaaServerConf                  string `json:"AAA_SERVER_CONF"`
	Enablerealtimebackup           string `json:"enableRealTimeBackup"`
	StrictCcMode                   string `json:"STRICT_CC_MODE"`
	V6SubnetTargetMask             string `json:"V6_SUBNET_TARGET_MASK"`
	AbstractTrunkHost              string `json:"abstract_trunk_host"`
	MstInstanceRange               string `json:"MST_INSTANCE_RANGE"`
	BgpAuthEnable                  string `json:"BGP_AUTH_ENABLE"`
	PmEnablePrev                   string `json:"PM_ENABLE_PREV"`
	Enablescheduledbackup          string `json:"enableScheduledBackup"`
	AbstractOspfInterface          string `json:"abstract_ospf_interface"`
	MacsecFallbackAlgorithm        string `json:"MACSEC_FALLBACK_ALGORITHM"`
	UnnumDhcpEnd                   string `json:"UNNUM_DHCP_END"`
	EnableAaa                      bool   `json:"ENABLE_AAA"`
	DeploymentFreeze               string `json:"DEPLOYMENT_FREEZE"`
	L2HostIntfMtuPrev              string `json:"L2_HOST_INTF_MTU_PREV"`
	NetflowMonitorList             string `json:"NETFLOW_MONITOR_LIST"`
	EnableAgent                    string `json:"ENABLE_AGENT"`
	NtpServerIpList                string `json:"NTP_SERVER_IP_LIST"`
	OverlayMode                    string `json:"OVERLAY_MODE"`
	MacsecFallbackKeyString        string `json:"MACSEC_FALLBACK_KEY_STRING"`
	StpRootOption                  string `json:"STP_ROOT_OPTION"`
	IsisOverloadEnable             bool   `json:"ISIS_OVERLOAD_ENABLE"`
	NetflowRecordList              string `json:"NETFLOW_RECORD_LIST"`
	SpineCount                     string `json:"SPINE_COUNT"`
	AbstractExtraConfigBootstrap   string `json:"abstract_extra_config_bootstrap"`
	MplsLoopbackIpRange            string `json:"MPLS_LOOPBACK_IP_RANGE"`
	LinkStateRoutingTagPrev        string `json:"LINK_STATE_ROUTING_TAG_PREV"`
	DhcpEnable                     bool   `json:"DHCP_ENABLE"`
	MsoSiteGroupName               string `json:"MSO_SITE_GROUP_NAME"`
	MgmtPrefixInternal             string `json:"MGMT_PREFIX_INTERNAL"`
	DhcpIpv6EnableInternal         string `json:"DHCP_IPV6_ENABLE_INTERNAL"`
	BgpAuthKeyType                 string `json:"BGP_AUTH_KEY_TYPE"`
	TempAnycastGateway             string `json:"temp_anycast_gateway"`
	BrfieldDebugFlag               string `json:"BRFIELD_DEBUG_FLAG"`
	BootstrapMultisubnet           string `json:"BOOTSTRAP_MULTISUBNET"`
	IsisP2PEnable                  bool   `json:"ISIS_P2P_ENABLE"`
	CdpEnable                      string `json:"CDP_ENABLE"`
	PtpLbId                        string `json:"PTP_LB_ID"`
	DhcpIpv6Enable                 string `json:"DHCP_IPV6_ENABLE"`
	MacsecKeyString                string `json:"MACSEC_KEY_STRING"`
	OspfAuthKey                    string `json:"OSPF_AUTH_KEY"`
	EnableFabricVpcDomainIdPrev    string `json:"ENABLE_FABRIC_VPC_DOMAIN_ID_PREV"`
	ExtraConfLeaf                  string `json:"EXTRA_CONF_LEAF"`
	DhcpStart                      string `json:"DHCP_START"`
	EnableTrm                      string `json:"ENABLE_TRM"`
	FeaturePtpInternal             string `json:"FEATURE_PTP_INTERNAL"`
	AbstractIsis                   string `json:"abstract_isis"`
	MplsLbId                       string `json:"MPLS_LB_ID"`
	FabricVpcDomainIdPrev          string `json:"FABRIC_VPC_DOMAIN_ID_PREV"`
	StaticUnderlayIpAlloc          string `json:"STATIC_UNDERLAY_IP_ALLOC"`
	MgmtV6PrefixInternal           string `json:"MGMT_V6PREFIX_INTERNAL"`
	MplsHandoff                    string `json:"MPLS_HANDOFF"`
	StpBridgePriority              string `json:"STP_BRIDGE_PRIORITY"`
	Scheduledtime                  string `json:"scheduledTime"`
	MacsecCipherSuite              string `json:"MACSEC_CIPHER_SUITE"`
	StpVlanRange                   string `json:"STP_VLAN_RANGE"`
	AnycastLbId                    string `json:"ANYCAST_LB_ID"`
	MsoControlerId                 string `json:"MSO_CONTROLER_ID"`
	AbstractExtraConfigLeaf        string `json:"abstract_extra_config_leaf"`
	AbstractDhcp                   string `json:"abstract_dhcp"`
	ExtraConfSpine                 string `json:"EXTRA_CONF_SPINE"`
	NtpServerVrf                   string `json:"NTP_SERVER_VRF"`
	SpineSwitchCoreInterfaces      string `json:"SPINE_SWITCH_CORE_INTERFACES"`
	IsisOverloadElapseTime         string `json:"ISIS_OVERLOAD_ELAPSE_TIME"`
	BootstrapConf                  string `json:"BOOTSTRAP_CONF"`
	IsisAuthKey                    string `json:"ISIS_AUTH_KEY"`
	DnsServerIpList                string `json:"DNS_SERVER_IP_LIST"`
	DnsServerVrf                   string `json:"DNS_SERVER_VRF"`
	EnableEvpn                     string `json:"ENABLE_EVPN"`
	AbstractMulticast              string `json:"abstract_multicast"`
	AgentIntf                      string `json:"AGENT_INTF"`
	L3VniMcastGroup                string `json:"L3VNI_MCAST_GROUP"`
	UnnumBootstrapLbId             string `json:"UNNUM_BOOTSTRAP_LB_ID"`
	VpcDomainIdRange               string `json:"VPC_DOMAIN_ID_RANGE"`
	HostIntfAdminState             string `json:"HOST_INTF_ADMIN_STATE"`
	SyslogSev                      string `json:"SYSLOG_SEV"`
	AbstractLoopbackInterface      string `json:"abstract_loopback_interface"`
	SyslogServerVrf                string `json:"SYSLOG_SERVER_VRF"`
	ExtraConfIntraLinks            string `json:"EXTRA_CONF_INTRA_LINKS"`
	AbstractExtraConfigSpine       string `json:"abstract_extra_config_spine"`
	PimHelloAuthKey                string `json:"PIM_HELLO_AUTH_KEY"`
	TempVpcDomainMgmt              string `json:"temp_vpc_domain_mgmt"`
	V6SubnetRange                  string `json:"V6_SUBNET_RANGE"`
	AbstractRoutedHost             string `json:"abstract_routed_host"`
	BgpAuthKey                     string `json:"BGP_AUTH_KEY"`
	InbandDhcpServers              string `json:"INBAND_DHCP_SERVERS"`
	IsisAuthKeychainKeyId          string `json:"ISIS_AUTH_KEYCHAIN_KEY_ID"`
	MgmtV6Prefix                   string `json:"MGMT_V6PREFIX"`
	AbstractFeatureSpine           string `json:"abstract_feature_spine"`
	EnableDefaultQueuingPolicy     string `json:"ENABLE_DEFAULT_QUEUING_POLICY"`
	AnycastBgwAdvertisePip         string `json:"ANYCAST_BGW_ADVERTISE_PIP"`
	NetflowExporterList            string `json:"NETFLOW_EXPORTER_LIST"`
	AbstractVlanInterface          string `json:"abstract_vlan_interface"`
	FabricName                     string `json:"FABRIC_NAME"`
	AbstractPimInterface           string `json:"abstract_pim_interface"`
	PmEnable                       string `json:"PM_ENABLE"`
	Loopback0Ipv6Range             string `json:"LOOPBACK0_IPV6_RANGE"`
	EnableVpcPeerLinkNativeVlan    string `json:"ENABLE_VPC_PEER_LINK_NATIVE_VLAN"`
	AbstractRouteMap               string `json:"abstract_route_map"`
	InbandMgmtPrev                 string `json:"INBAND_MGMT_PREV"`
	AbstractVpcDomain              string `json:"abstract_vpc_domain"`
	DhcpEndInternal                string `json:"DHCP_END_INTERNAL"`
	BootstrapEnable                string `json:"BOOTSTRAP_ENABLE"`
	AbstractExtraConfigTor         string `json:"abstract_extra_config_tor"`
	SyslogServerIpList             string `json:"SYSLOG_SERVER_IP_LIST"`
	BootstrapEnablePrev            string `json:"BOOTSTRAP_ENABLE_PREV"`
	EnableTenantDhcp               string `json:"ENABLE_TENANT_DHCP"`
	AnycastRpIpRangeInternal       string `json:"ANYCAST_RP_IP_RANGE_INTERNAL"`
	BootstrapMultisubnetInternal   string `json:"BOOTSTRAP_MULTISUBNET_INTERNAL"`
	MgmtGw                         string `json:"MGMT_GW"`
	UnnumDhcpStart                 string `json:"UNNUM_DHCP_START"`
	MgmtPrefix                     string `json:"MGMT_PREFIX"`
	AbstractBgpRr                  string `json:"abstract_bgp_rr"`
	InbandMgmt                     string `json:"INBAND_MGMT"`
	AbstractBgp                    string `json:"abstract_bgp"`
	EnableNetflowPrev              string `json:"ENABLE_NETFLOW_PREV"`
	DeafultQueuingPolicyCloudscale string `json:"DEAFULT_QUEUING_POLICY_CLOUDSCALE"`
	FabricVpcQos                   string `json:"FABRIC_VPC_QOS"`
	AaaRemoteIpEnabled             string `json:"AAA_REMOTE_IP_ENABLED"`
	FabricTemplate                 string `json:"FF"`
	FabricType                     string `json:"FABRIC_TYPE"`
	SpineAddDelBedugFlag           string `json:"SSPINE_ADD_DEL_DEBUG_FLAG"`
	ActiveMigration                string `json:"ACTIVE_MIGRATION"`
	SiteId                         string `json:"SITE_ID"`
}

func (fabric *Fabric) ToMap() (map[string]interface{}, error) {
	fabricAttributeMap := make(map[string]interface{})

	if fabric.Id != 0 {
		A(fabricAttributeMap, "id", fabric.Id)
	}
	A(fabricAttributeMap, "fabricName", fabric.Name)
	A(fabricAttributeMap, "fabricId", fabric.FabricId)
	A(fabricAttributeMap, "templateName", fabric.Template)
	A(fabricAttributeMap, "nvPairs", fabric.Config)
	return fabricAttributeMap, nil
}

func (config *FabricConfig) SetConfigDefaults() {
	config.MsoSiteId = ""
	config.PhantomRpLbId1 = ""
	config.PhantomRpLbId2 = ""
	config.PhantomRpLbId3 = ""
	config.IbgpPeerTemplate = ""
	config.PhantomRpLbId4 = ""
	config.AbstractOspf = "base_ospf"
	config.FeaturePtp = "false"
	config.DhcpStartInternal = ""
	config.SspineCount = "0"
	config.AdvertisePipBgp = "false"
	config.FabricVpcQosPolicyName = ""
	config.DhcpEnd = ""
	config.UnderlayIsV6 = "false"
	config.FabricVpcDomainId = ""
	config.SeedSwitchCoreInterfaces = ""
	config.FabricMtuPrev = "9216"
	config.HdTime = "180"
	config.OspfAuthEnable = "false"
	config.Loopback1Ipv6Range = ""
	config.RouterIdRange = ""
	config.MsoConnectivityDeployed = ""
	config.EnableMacsec = "false"
	config.DeafultQueuingPolicyOther = ""
	config.UnnumDhcpStartInternal = ""
	config.MacsecReportTimer = ""
	config.PremsoParentFabric = ""
	config.UnnumDhcpEndInternal = ""
	config.PtpDomainId = ""
	config.AutoSymmetricVrfLite = false
	config.UseLinkLocal = false
	config.BgpAsPrev = ""
	config.EnablePbr = "false"
	config.VpcPeerLinkPo = "500"
	config.VpcDelayRestoretime = "60"
	config.IsisAuthEnable = false
	config.VpcEnableIpv6NdSync = "true"
	config.AbstractIsisInterface = "isis_interface"
	config.TcamAllocation = "true"
	config.MacsecAlgorithm = ""
	config.IsisLevel = ""
	config.AbstractAnycastRp = "anycast_rp"
	config.EnableNetflow = "false"
	config.DeafultQueuingPolicyRSeries = ""
	config.TempVpcPeerLink = "int_vpc_peer_link_po"
	config.BrownfieldNetworkNameFormat = "Auto_Net_VNI$$VNI$$_VLAN$$VLAN_ID$$"
	config.EnableFabricVpcDomainId = "false"
	config.IbgpPeerTemplateLeaf = ""
	config.MgmtGwInternal = ""
	config.EnableNxapi = "true"
	config.GrfieldDebugFlag = "Disable"
	config.IsisAuthKeychainName = ""
	config.AbstractBgpNeighbor = "evpn_bgp_rr_neighbor"
	config.OspfAuthKeyId = ""
	config.PimHelloAuthEnable = "false"
	config.AbstractFeatureLeaf = "base_feature_leaf_upg"
	config.BfdAuthEnable = "false"
	config.ExtraConfTor = ""
	config.AaaServerConf = ""
	config.Enablerealtimebackup = ""
	config.StrictCcMode = "false"
	config.V6SubnetTargetMask = ""
	config.AbstractTrunkHost = "int_trunk_host"
	config.MstInstanceRange = ""
	config.BgpAuthEnable = "false"
	config.PmEnablePrev = "false"
	config.Enablescheduledbackup = ""
	config.AbstractOspfInterface = "ospf_interface_11_1"
	config.MacsecFallbackAlgorithm = ""
	config.UnnumDhcpEnd = ""
	config.EnableAaa = false
	config.DeploymentFreeze = "false"
	config.L2HostIntfMtuPrev = "9216"
	config.NetflowMonitorList = ""
	config.EnableAgent = "false"
	config.NtpServerIpList = ""
	config.OverlayMode = "config-profile"
	config.MacsecFallbackKeyString = ""
	config.StpRootOption = "unmanaged"
	config.FabricType = "Switch_Fabric"
	config.IsisOverloadEnable = false
	config.NetflowRecordList = ""
	config.SpineCount = "0"
	config.AbstractExtraConfigBootstrap = "extra_config_bootstrap_11_1"
	config.MplsLoopbackIpRange = ""
	config.LinkStateRoutingTagPrev = ""
	config.DhcpEnable = false
	config.MsoSiteGroupName = ""
	config.MgmtPrefixInternal = ""
	config.DhcpIpv6EnableInternal = ""
	config.BgpAuthKeyType = ""
	config.SiteId = ""
	config.TempAnycastGateway = "anycast_gateway"
	config.BrfieldDebugFlag = "Disable"
	config.BootstrapMultisubnet = ""
	config.IsisP2PEnable = false
	config.EnableNgoam = "true"
	config.CdpEnable = "false"
	config.PtpLbId = ""
	config.DhcpIpv6Enable = ""
	config.MacsecKeyString = ""
	config.OspfAuthKey = ""
	config.EnableFabricVpcDomainIdPrev = ""
	config.ExtraConfLeaf = ""
	config.DhcpStart = ""
	config.EnableTrm = "false"
	config.FeaturePtpInternal = "false"
	config.EnableNxapiHttp = "true"
	config.AbstractIsis = "base_isis_level2"
	config.MplsLbId = ""
	config.FabricVpcDomainIdPrev = ""
	config.StaticUnderlayIpAlloc = "false"
	config.MgmtV6PrefixInternal = ""
	config.MplsHandoff = "false"
	config.StpBridgePriority = ""
	config.Scheduledtime = ""
	config.MacsecCipherSuite = ""
	config.StpVlanRange = ""
	config.AnycastLbId = ""
	config.MsoControlerId = ""
	config.AbstractExtraConfigLeaf = "extra_config_leaf"
	config.AbstractDhcp = "base_dhcp"
	config.ExtraConfSpine = ""
	config.NtpServerVrf = ""
	config.SpineSwitchCoreInterfaces = ""
	config.IsisOverloadElapseTime = ""
	config.BootstrapConf = ""
	config.IsisAuthKey = ""
	config.DnsServerIpList = ""
	config.DnsServerVrf = ""
	config.EnableEvpn = "true"
	config.AbstractMulticast = "base_multicast_11_1"
	config.AgentIntf = "eth0"
	config.L3VniMcastGroup = ""
	config.UnnumBootstrapLbId = ""
	config.VpcDomainIdRange = "1-1000"
	config.HostIntfAdminState = "true"
	config.SyslogSev = ""
	config.AbstractLoopbackInterface = "int_fabric_loopback_11_1"
	config.SyslogServerVrf = ""
	config.ExtraConfIntraLinks = ""
	config.SnmpServerHostTrap = "true"
	config.AbstractExtraConfigSpine = "extra_config_spine"
	config.PimHelloAuthKey = ""
	config.TempVpcDomainMgmt = "vpc_domain_mgmt"
	config.V6SubnetRange = ""
	config.AbstractRoutedHost = "int_routed_host"
	config.BgpAuthKey = ""
	config.InbandDhcpServers = ""
	config.IsisAuthKeychainKeyId = ""
	config.MgmtV6Prefix = "64"
	config.AbstractFeatureSpine = "base_feature_spine_upg"
	config.EnableDefaultQueuingPolicy = "false"
	config.AnycastBgwAdvertisePip = "false"
	config.NetflowExporterList = ""
	config.AbstractVlanInterface = "int_fabric_vlan_11_1"
	config.AbstractPimInterface = "pim_interface"
	config.PmEnable = "false"
	config.Loopback0Ipv6Range = ""
	config.OverlayModePrev = ""
	config.EnableVpcPeerLinkNativeVlan = "false"
	config.AbstractRouteMap = "route_map"
	config.InbandMgmtPrev = "false"
	config.AbstractVpcDomain = "base_vpc_domain_11_1"
	config.ActiveMigration = "false"
	config.DhcpEndInternal = ""
	config.BootstrapEnable = "false"
	config.AbstractExtraConfigTor = "extra_config_tor"
	config.SyslogServerIpList = ""
	config.BootstrapEnablePrev = "false"
	config.EnableTenantDhcp = "true"
	config.AnycastRpIpRangeInternal = ""
	config.BootstrapMultisubnetInternal = ""
	config.MgmtGw = ""
	config.UnnumDhcpStart = ""
	config.MgmtPrefix = ""
	config.AbstractBgpRr = "evpn_bgp_rr"
	config.InbandMgmt = "false"
	config.AbstractBgp = "base_bgp"
	config.EnableNetflowPrev = ""
	config.DeafultQueuingPolicyCloudscale = ""
	config.FabricVpcQos = "false"
	config.AaaRemoteIpEnabled = "false"
	config.FabricType = "Switch_Fabric"
	config.SpineAddDelBedugFlag = "Disable"
	config.ActiveMigration = "false"
	config.BfdEnable = "false"

}
