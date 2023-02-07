package dcnm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDCNMFabric() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMFabricCreate,
		Read:   resourceDCNMFabricRead,
		Update: resourceDCNMFabricUpdate,
		Delete: resourceDCNMFabricDelete,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fabric_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Easy_Fabric",
			},

			"asn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"underlay_interface_numbering": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "p2p",
				ValidateFunc: validation.StringInSlice([]string{
					"p2p",
					"unnumbered",
				}, false),
			},

			"underlay_subnet_mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "30",
				//TODO: Range validation
			},

			"underlay_routing_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ospf",
				ValidateFunc: validation.StringInSlice([]string{
					"ospf",
					"is-is",
				}, false),
			},

			"route_reflectors_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2",
				ValidateFunc: validation.StringInSlice([]string{
					"2",
					"4",
				}, false),
			},

			"anycast_mac": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2020.0000.00aa",
				//TODO: Validate MAC
			},

			"replication_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Multicast",
				ValidateFunc: validation.StringInSlice([]string{
					"Multicast",
					"Ingress",
				}, false),
			},

			"multicast_group_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "239.1.1.0/25",
				//TODO: Validate range
			},

			"rendevous_point_count": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2",
				ValidateFunc: validation.StringInSlice([]string{
					"2",
					"4",
				}, false),
			},

			"rendevous_point_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "asm",
				ValidateFunc: validation.StringInSlice([]string{
					"asm",
					"bidir",
				}, false),
			},

			"rendevous_loopback_id": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      254,
				ValidateFunc: validation.IntBetween(1, 254),
			},

			"vpc_peer_link_vlan": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validation.IntBetween(2, 4094),
			},

			"vpc_peer_keep_alive_option": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "management",
				ValidateFunc: validation.StringInSlice([]string{
					"management",
					"loopback",
				}, false),
			},

			"vpc_auto_recovery_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      360,
				ValidateFunc: validation.IntBetween(240, 3600),
			},

			"vpc_delay_restore_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      150,
				ValidateFunc: validation.IntBetween(1, 3600),
			},

			"underlay_routing_loopback_id": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1023),
			},

			"underlay_vtep_loopback_id": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1023),
			},

			"underlay_routing_protocol_tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UNDERLAY",
			},

			"ospf_area_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0",
			},

			"vrf_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_VRF_Universal",
			},

			"network_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_Network_Universal",
			},

			"vrf_extension_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_VRF_Extension_Universal",
			},

			"network_extension_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_Network_Extension_Universal",
			},

			"intra_fabric_interface_mtu": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      9216,
				ValidateFunc: validation.IntBetween(576, 9216),
			},

			"layer_2_host_interface_mtu": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      9216,
				ValidateFunc: validation.IntBetween(576, 9216),
			},

			"power_supply_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ps-redundant",
				ValidateFunc: validation.StringInSlice([]string{
					"pd-redundant",
					"combined",
					"insrc-redundant",
				}, false),
			},

			"copp_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "strict",
				ValidateFunc: validation.StringInSlice([]string{
					"strict",
					"dense",
					"lenient",
					"moderate",
					"manual",
				}, false),
			},

			"underlay_routing_loopback_ip_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10.2.0.0/22",
			},

			"underlay_vtep_loopback_ip_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10.3.0.0/22",
			},

			"underlay_rp_loopback_ip_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10.254.254.0/24",
			},

			"underlay_subnet_ip_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10.4.0.0/16",
			},

			"layer_2_vxlan_vni_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "30000-49000",
			},

			"layer_3_vxlan_vni_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "50000-59000",
			},

			"network_vlan_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2300-2999",
			},

			"vrf_vlan_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2000-2299",
			},

			"subinterface_dot1q_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "2-511",
			},

			"vrf_lite_deployment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Manual",
				ValidateFunc: validation.StringInSlice([]string{
					"Manual",
					"Back2BackOnly",
					"ToExternalOnly",
					"Back2Back&ToExternal",
				}, false),
			},

			"vrf_lite_subnet_ip_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10.33.0.0/16",
			},

			"vrf_lite_subnet_mask": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntBetween(8, 31),
			},

			"service_network_vlan_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "3000-3199",
			},

			"route_map_sequence_number_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1-65534",
			},
		},
	}
}

func getRemoteFabric(client *client.Client, name string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/control/fabrics/%s", name)

	cont, err := client.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	return cont, nil
}

func setFabricAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {

	d.Set("name", stripQuotes(cont.S("fabricName").String()))
	d.Set("fabric_id", stripQuotes(cont.S("fabricId").String()))
	d.Set("template", stripQuotes(cont.S("templateName").String()))

	cont, err := cleanJsonString(stripQuotes(cont.S("nvPairs").String()))
	if err == nil {
		if cont.Exists("BGP_AS") {
			d.Set("asn", stripQuotes(cont.S("BGP_AS").String()))
		}
		if cont.Exists("FABRIC_INTERFACE_TYPE") {
			d.Set("underlay_interface_numbering", stripQuotes(cont.S("FABRIC_INTERFACE_TYPE").String()))
		}
		if cont.Exists("SUBNET_TARGET_MASK") {
			d.Set("underlay_subnet_mask", stripQuotes(cont.S("SUBNET_TARGET_MASK").String()))
		}
		if cont.Exists("LINK_STATE_ROUTING") {
			d.Set("underlay_routing_protocol", stripQuotes(cont.S("LINK_STATE_ROUTING").String()))
		}
		if cont.Exists("RR_COUNT") {
			d.Set("route_reflectors_count", stripQuotes(cont.S("RR_COUNT").String()))
		}
		if cont.Exists("ANYCAST_GW_MAC") {
			d.Set("anycast_mac", stripQuotes(cont.S("ANYCAST_GW_MAC").String()))
		}
		if cont.Exists("REPLICATION_MODE") {
			d.Set("replication_mode", stripQuotes(cont.S("REPLICATION_MODE").String()))
		}
		if cont.Exists("MULTICAST_GROUP_SUBNET") {
			d.Set("multicast_group_subnet", stripQuotes(cont.S("MULTICAST_GROUP_SUBNET").String()))
		}
		if cont.Exists("RP_COUNT") {
			d.Set("rendevous_point_count", stripQuotes(cont.S("RP_COUNT").String()))
		}
		if cont.Exists("RP_MODE") {
			d.Set("rendevous_point_mode", stripQuotes(cont.S("RP_MODE").String()))
		}
		if cont.Exists("RP_LB_ID") && stripQuotes(cont.S("RP_LB_ID").String()) != "" {
			if rpLbId := stripQuotes(cont.S("RP_LB_ID").String()); err == nil {
				rpLbIdInt, _ := strconv.Atoi(rpLbId)
				d.Set("rendevous_loopback_id", rpLbIdInt)
			}
		}
		if cont.Exists("VPC_PEER_LINK_VLAN") && stripQuotes(cont.S("VPC_PEER_LINK_VLAN").String()) != "" {
			if vpcPeerVlan := stripQuotes(cont.S("VPC_PEER_LINK_VLAN").String()); err == nil {
				vpcPeerVlanInt, _ := strconv.Atoi(vpcPeerVlan)
				d.Set("vpc_peer_link_vlan", vpcPeerVlanInt)
			}
		}
		if cont.Exists("VPC_PEER_KEEP_ALIVE_OPTION") {
			d.Set("vpc_peer_keep_alive_option", stripQuotes(cont.S("VPC_PEER_KEEP_ALIVE_OPTION").String()))
		}
		if cont.Exists("VPC_AUTO_RECOVERY_TIME") && stripQuotes(cont.S("VPC_AUTO_RECOVERY_TIME").String()) != "" {
			if vpcRecoveryTime := stripQuotes(cont.S("VPC_AUTO_RECOVERY_TIME").String()); err == nil {
				vpcRecoveryTimeInt, _ := strconv.Atoi(vpcRecoveryTime)
				d.Set("vpc_auto_recovery_time", vpcRecoveryTimeInt)
			}
		}
		if cont.Exists("VPC_DELAY_RESTORE") && stripQuotes(cont.S("VPC_DELAY_RESTORE").String()) != "" {
			if vpcRestoreTime := stripQuotes(cont.S("VPC_DELAY_RESTORE").String()); err == nil {
				vpcRestoreTimeInt, _ := strconv.Atoi(vpcRestoreTime)
				d.Set("vpc_delay_restore_time", vpcRestoreTimeInt)
			}
		}
		if cont.Exists("BGP_LB_ID") && stripQuotes(cont.S("BGP_LB_ID").String()) != "" {
			if bgpLbId := stripQuotes(cont.S("BGP_LB_ID").String()); err == nil {
				bgpLbIdInt, _ := strconv.Atoi(bgpLbId)
				d.Set("underlay_routing_loopback_id", bgpLbIdInt)
			}
		}
		if cont.Exists("NVE_LB_ID") && stripQuotes(cont.S("NVE_LB_ID").String()) != "" {
			if nveLbId := stripQuotes(cont.S("NVE_LB_ID").String()); err == nil {
				nveLbIdInt, _ := strconv.Atoi(nveLbId)
				d.Set("underlay_vtep_loopback_id", nveLbIdInt)
			}
		}
		if cont.Exists("LINK_STATE_ROUTING_TAG") {
			d.Set("underlay_routing_protocol_tag", stripQuotes(cont.S("LINK_STATE_ROUTING_TAG").String()))
		}
		if cont.Exists("OSPF_AREA_ID") {
			d.Set("ospf_area_id", stripQuotes(cont.S("OSPF_AREA_ID").String()))
		}
		if cont.Exists("default_vrf") {
			d.Set("vrf_template", stripQuotes(cont.S("default_vrf").String()))
		}
		if cont.Exists("default_network") {
			d.Set("network_template", stripQuotes(cont.S("default_network").String()))
		}
		if cont.Exists("vrf_extension_template") {
			d.Set("vrf_extension_template", stripQuotes(cont.S("vrf_extension_template").String()))
		}
		if cont.Exists("network_extension_template") {
			d.Set("network_extension_template", stripQuotes(cont.S("network_extension_template").String()))
		}
		if cont.Exists("FABRIC_MTU") && stripQuotes(cont.S("FABRIC_MTU").String()) != "" {
			if fabricMtu := stripQuotes(cont.S("FABRIC_MTU").String()); err == nil {
				fabricMtuInt, _ := strconv.Atoi(fabricMtu)
				d.Set("intra_fabric_interface_mtu", fabricMtuInt)
			}
		}
		if cont.Exists("L2_HOST_INTF_MTU") && stripQuotes(cont.S("L2_HOST_INTF_MTU").String()) != "" {
			if hostMtu := stripQuotes(cont.S("L2_HOST_INTF_MTU").String()); err == nil {
				hostMtuInt, _ := strconv.Atoi(hostMtu)
				d.Set("layer_2_host_interface_mtu", hostMtuInt)
			}
		}
		if cont.Exists("POWER_REDUNDANCY_MODE") {
			d.Set("power_supply_mode", stripQuotes(cont.S("POWER_REDUNDANCY_MODE").String()))
		}
		if cont.Exists("COPP_POLICY") {
			d.Set("copp_profile", stripQuotes(cont.S("COPP_POLICY").String()))
		}
		if cont.Exists("LOOPBACK0_IP_RANGE") {
			d.Set("underlay_routing_loopback_ip_range", stripQuotes(cont.S("LOOPBACK0_IP_RANGE").String()))
		}
		if cont.Exists("LOOPBACK1_IP_RANGE") {
			d.Set("underlay_vtep_loopback_ip_range", stripQuotes(cont.S("LOOPBACK1_IP_RANGE").String()))
		}
		if cont.Exists("RP_MODANYCAST_RP_IP_RANGEE") {
			d.Set("underlay_rp_loopback_ip_range", stripQuotes(cont.S("ANYCAST_RP_IP_RANGE").String()))
		}
		if cont.Exists("SUBNET_RANGE") {
			d.Set("underlay_subnet_ip_range", stripQuotes(cont.S("SUBNET_RANGE").String()))
		}
		if cont.Exists("L2_SEGMENT_ID_RANGE") {
			d.Set("layer_2_vxlan_vni_range", stripQuotes(cont.S("L2_SEGMENT_ID_RANGE").String()))
		}
		if cont.Exists("L3_PARTITION_ID_RANGE") {
			d.Set("layer_3_vxlan_vni_range", stripQuotes(cont.S("L3_PARTITION_ID_RANGE").String()))
		}
		if cont.Exists("NETWORK_VLAN_RANGE") {
			d.Set("network_vlan_range", stripQuotes(cont.S("NETWORK_VLAN_RANGE").String()))
		}
		if cont.Exists("VRF_VLAN_RANGE") {
			d.Set("vrf_vlan_range", stripQuotes(cont.S("VRF_VLAN_RANGE").String()))
		}
		if cont.Exists("SUBINTERFACE_RANGE") {
			d.Set("subinterface_dot1q_range", stripQuotes(cont.S("SUBINTERFACE_RANGE").String()))
		}
		if cont.Exists("VRF_LITE_AUTOCONFIG") {
			d.Set("vrf_lite_deployment", stripQuotes(cont.S("VRF_LITE_AUTOCONFIG").String()))
		}
		if cont.Exists("DCI_SUBNET_RANGE") {
			d.Set("vrf_lite_subnet_ip_range", stripQuotes(cont.S("DCI_SUBNET_RANGE").String()))
		}
		if cont.Exists("DCI_SUBNET_TARGET_MASK") && stripQuotes(cont.S("DCI_SUBNET_TARGET_MASK").String()) != "" {
			if vrfLiteSubnetMask := stripQuotes(cont.S("DCI_SUBNET_TARGET_MASK").String()); err == nil {
				vrfLiteSubnetMaskInt, _ := strconv.Atoi(vrfLiteSubnetMask)
				d.Set("vrf_lite_subnet_mask", vrfLiteSubnetMaskInt)
			}
		}
		if cont.Exists("SERVICE_NETWORK_VLAN_RANGE") {
			d.Set("service_network_vlan_range", stripQuotes(cont.S("SERVICE_NETWORK_VLAN_RANGE").String()))
		}
		if cont.Exists("ROUTE_MAP_SEQUENCE_NUMBER_RANGE") {
			d.Set("route_map_sequence_number_range", stripQuotes(cont.S("ROUTE_MAP_SEQUENCE_NUMBER_RANGE").String()))
		}
	}

	d.Set("id", stripQuotes(cont.S("id").String()))
	return d
}

func resourceDCNMFabricCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	fabric := models.Fabric{}
	fabric.Name = d.Get("name").(string)
	fabric.Template = d.Get("template").(string)

	configMap := models.FabricConfig{}

	if asn, ok := d.GetOk("asn"); ok {
		configMap.Asn = asn.(string)
	}
	if underlayRoutingNumbering, ok := d.GetOk("underlay_interface_numbering"); ok {
		configMap.UnderlayInterfaceNumbering = underlayRoutingNumbering.(string)
	}
	if underlaySubnetMask, ok := d.GetOk("underlay_subnet_mask"); ok {
		configMap.UnderlaySubnetMask = underlaySubnetMask.(string)
	}
	if uderlayRoutingProtocol, ok := d.GetOk("underlay_routing_protocol"); ok {
		configMap.UnderlayRoutingProcotol = uderlayRoutingProtocol.(string)
	}
	if rrCount, ok := d.GetOk("route_reflectors_count"); ok {
		configMap.RouteReflectorCount = rrCount.(string)
	}
	if anycastMac, ok := d.GetOk("anycast_mac"); ok {
		configMap.AnycastMac = anycastMac.(string)
	}
	if replicationMode, ok := d.GetOk("replication_mode"); ok {
		configMap.ReplicationMode = replicationMode.(string)
	}
	if multicastGroupSubnet, ok := d.GetOk("multicast_group_subnet"); ok {
		configMap.MulticastGroupSubnet = multicastGroupSubnet.(string)
	}
	if rpCount, ok := d.GetOk("rendevous_point_count"); ok {
		configMap.RendevouzPointCount = rpCount.(string)
	}
	if rpMode, ok := d.GetOk("rendevous_point_mode"); ok {
		configMap.RendevouzPointMode = rpMode.(string)
	}
	if rpId, ok := d.GetOk("rendevous_loopback_id"); ok {
		configMap.RendevouzPointLoopbackId = strconv.Itoa(rpId.(int))
	}
	if vpcPlVlan, ok := d.GetOk("vpc_peer_link_vlan"); ok {
		configMap.VpcPeerLinkVlan = strconv.Itoa(vpcPlVlan.(int))
	}
	if vpcPkaOption, ok := d.GetOk("vpc_peer_keep_alive_option"); ok {
		configMap.VpcPeerKeepAliveOption = vpcPkaOption.(string)
	}
	if vpcAutoRectime, ok := d.GetOk("vpc_auto_recovery_time"); ok {
		configMap.VpcAutoRecoveryTime = strconv.Itoa(vpcAutoRectime.(int))
	}
	if vpcDelayResTime, ok := d.GetOk("vpc_delay_restore_time"); ok {
		configMap.VpcDelayRestore = strconv.Itoa(vpcDelayResTime.(int))
	}
	if rtLooId, ok := d.GetOk("underlay_routing_loopback_id"); ok {
		configMap.UnderlayRoutingLoopbackId = strconv.Itoa(rtLooId.(int))
	} else {
		configMap.UnderlayRoutingLoopbackId = "0"
	}
	if vtepLooId, ok := d.GetOk("underlay_vtep_loopback_id"); ok {
		configMap.UnderlayVtepLoopbackId = strconv.Itoa(vtepLooId.(int))
	}
	if rtProtoTag, ok := d.GetOk("underlay_routing_protocol_tag"); ok {
		configMap.UnderlayRoutingProtocolTag = rtProtoTag.(string)
	}
	if ospfAreaId, ok := d.GetOk("ospf_area_id"); ok {
		configMap.OspfAreaId = ospfAreaId.(string)
	}
	if vrfTemp, ok := d.GetOk("vrf_template"); ok {
		configMap.VrfTemplate = vrfTemp.(string)
	}
	if netTemp, ok := d.GetOk("network_template"); ok {
		configMap.NetworkTemplate = netTemp.(string)
	}
	if vrfExtTemp, ok := d.GetOk("vrf_extension_template"); ok {
		configMap.VrfExtensionTemplate = vrfExtTemp.(string)
	}
	if netExtTemp, ok := d.GetOk("network_extension_template"); ok {
		configMap.NetworkExtensionTemplate = netExtTemp.(string)
	}
	if fabMtu, ok := d.GetOk("intra_fabric_interface_mtu"); ok {
		configMap.IntraFabricInterfaceMtu = strconv.Itoa(fabMtu.(int))
	}
	if hostMtu, ok := d.GetOk("layer_2_host_interface_mtu"); ok {
		configMap.Layer2HostInterfaceMtu = strconv.Itoa(hostMtu.(int))
	}
	if psMode, ok := d.GetOk("power_supply_mode"); ok {
		configMap.PowerSupplyMode = psMode.(string)
	}
	if coppProf, ok := d.GetOk("copp_profile"); ok {
		configMap.CoppProfile = coppProf.(string)
	}
	if rtLooIpRange, ok := d.GetOk("underlay_routing_loopback_ip_range"); ok {
		configMap.UnderlayRoutingLoopbackIpRange = rtLooIpRange.(string)
	}
	if vtepLooIpRange, ok := d.GetOk("underlay_vtep_loopback_ip_range"); ok {
		configMap.UnderlayVtepLoopbackIpRange = vtepLooIpRange.(string)
	}
	if rpLooIpRange, ok := d.GetOk("underlay_rp_loopback_ip_range"); ok {
		configMap.UnderlayRpLoopbackIpRange = rpLooIpRange.(string)
	}
	if subIpRange, ok := d.GetOk("underlay_subnet_ip_range"); ok {
		configMap.UnderlaySubnetIpRange = subIpRange.(string)
	}
	if l2VniRange, ok := d.GetOk("layer_2_vxlan_vni_range"); ok {
		configMap.Layer2VxlanVniRange = l2VniRange.(string)
	}
	if l3VniRange, ok := d.GetOk("layer_3_vxlan_vni_range"); ok {
		configMap.Layer3VxlanVniRange = l3VniRange.(string)
	}
	if netVlanRange, ok := d.GetOk("network_vlan_range"); ok {
		configMap.NetworkVlanRange = netVlanRange.(string)
	}
	if vrfVlanRange, ok := d.GetOk("vrf_vlan_range"); ok {
		configMap.VrfVlanRange = vrfVlanRange.(string)
	}
	if vrfLiteDeployment, ok := d.GetOk("vrf_lite_deployment"); ok {
		configMap.VrfLiteDeployment = vrfLiteDeployment.(string)
	}
	if subIfDot1qRange, ok := d.GetOk("subinterface_dot1q_range"); ok {
		configMap.SubinterfaceDot1qRange = subIfDot1qRange.(string)
	}
	if vrfLifeSubnetIpRange, ok := d.GetOk("vrf_lite_subnet_ip_range"); ok {
		configMap.VrfLiteSubnetIpRange = vrfLifeSubnetIpRange.(string)
	}
	if vrfLiteSubnetMask, ok := d.GetOk("vrf_lite_subnet_mask"); ok {
		configMap.VrfLiteSubnetMask = strconv.Itoa(vrfLiteSubnetMask.(int))
	}
	if svcNetVlanRange, ok := d.GetOk("service_network_vlan_range"); ok {
		configMap.ServiceNetworkVlanRange = svcNetVlanRange.(string)
	}
	if rmSeqRange, ok := d.GetOk("route_map_sequence_number_range"); ok {
		configMap.RouteMapSequenceNumberRange = rmSeqRange.(string)
	}

	configMap.FabricTemplate = fabric.Template
	configMap.FabricName = fabric.Name
	fabric.Config = configMap
	fabric.SetConfigDefaults()

	durl := "/rest/control/fabrics"
	cont, err := dcnmClient.Save(durl, &fabric)
	if err != nil {
		return err
	}

	d.SetId(stripQuotes(cont.S("id").String()))

	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMFabricRead(d, m)
}

func resourceDCNMFabricUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ")

	dcnmClient := m.(*client.Client)

	fabric := models.Fabric{}
	fabric.Name = d.Get("name").(string)
	fabric.Template = d.Get("template").(string)
	fabric.FabricId = d.Get("fabric_id").(string)

	configMap := models.FabricConfig{}

	if asn, ok := d.GetOk("asn"); ok {
		configMap.Asn = asn.(string)
	}
	if underlayRoutingNumbering, ok := d.GetOk("underlay_interface_numbering"); ok {
		configMap.UnderlayInterfaceNumbering = underlayRoutingNumbering.(string)
	}
	if underlaySubnetMask, ok := d.GetOk("underlay_subnet_mask"); ok {
		configMap.UnderlaySubnetMask = underlaySubnetMask.(string)
	}
	if uderlayRoutingProtocol, ok := d.GetOk("underlay_routing_protocol"); ok {
		configMap.UnderlayRoutingProcotol = uderlayRoutingProtocol.(string)
	}
	if rrCount, ok := d.GetOk("route_reflectors_count"); ok {
		configMap.RouteReflectorCount = rrCount.(string)
	}
	if anycastMac, ok := d.GetOk("anycast_mac"); ok {
		configMap.AnycastMac = anycastMac.(string)
	}
	if replicationMode, ok := d.GetOk("replication_mode"); ok {
		configMap.ReplicationMode = replicationMode.(string)
	}
	if multicastGroupSubnet, ok := d.GetOk("multicast_group_subnet"); ok {
		configMap.MulticastGroupSubnet = multicastGroupSubnet.(string)
	}
	if rpCount, ok := d.GetOk("rendevous_point_count"); ok {
		configMap.RendevouzPointCount = rpCount.(string)
	}
	if rpMode, ok := d.GetOk("rendevous_point_mode"); ok {
		configMap.RendevouzPointMode = rpMode.(string)
	}
	if rpId, ok := d.GetOk("rendevous_loopback_id"); ok {
		configMap.RendevouzPointLoopbackId = strconv.Itoa(rpId.(int))
	}
	if vpcPlVlan, ok := d.GetOk("vpc_peer_link_vlan"); ok {
		configMap.VpcPeerLinkVlan = strconv.Itoa(vpcPlVlan.(int))
	}
	if vpcPkaOption, ok := d.GetOk("vpc_peer_keep_alive_option"); ok {
		configMap.VpcPeerKeepAliveOption = vpcPkaOption.(string)
	}
	if vpcAutoRectime, ok := d.GetOk("vpc_auto_recovery_time"); ok {
		configMap.VpcAutoRecoveryTime = strconv.Itoa(vpcAutoRectime.(int))
	}
	if vpcDelayResTime, ok := d.GetOk("vpc_delay_restore_time"); ok {
		configMap.VpcDelayRestore = strconv.Itoa(vpcDelayResTime.(int))
	}
	if rtLooId, ok := d.GetOk("underlay_routing_loopback_id"); ok {
		configMap.UnderlayRoutingLoopbackId = strconv.Itoa(rtLooId.(int))
	} else {
		configMap.UnderlayRoutingLoopbackId = "0"
	}
	if vtepLooId, ok := d.GetOk("underlay_vtep_loopback_id"); ok {
		configMap.UnderlayVtepLoopbackId = strconv.Itoa(vtepLooId.(int))
	}
	if rtProtoTag, ok := d.GetOk("underlay_routing_protocol_tag"); ok {
		configMap.UnderlayRoutingProtocolTag = rtProtoTag.(string)
	}
	if ospfAreaId, ok := d.GetOk("ospf_area_id"); ok {
		configMap.OspfAreaId = ospfAreaId.(string)
	}
	if vrfTemp, ok := d.GetOk("vrf_template"); ok {
		configMap.VrfTemplate = vrfTemp.(string)
	}
	if netTemp, ok := d.GetOk("network_template"); ok {
		configMap.NetworkTemplate = netTemp.(string)
	}
	if vrfExtTemp, ok := d.GetOk("vrf_extension_template"); ok {
		configMap.VrfExtensionTemplate = vrfExtTemp.(string)
	}
	if netExtTemp, ok := d.GetOk("network_extension_template"); ok {
		configMap.NetworkExtensionTemplate = netExtTemp.(string)
	}
	if fabMtu, ok := d.GetOk("intra_fabric_interface_mtu"); ok {
		configMap.IntraFabricInterfaceMtu = strconv.Itoa(fabMtu.(int))
	}
	if hostMtu, ok := d.GetOk("layer_2_host_interface_mtu"); ok {
		configMap.Layer2HostInterfaceMtu = strconv.Itoa(hostMtu.(int))
	}
	if psMode, ok := d.GetOk("power_supply_mode"); ok {
		configMap.PowerSupplyMode = psMode.(string)
	}
	if coppProf, ok := d.GetOk("copp_profile"); ok {
		configMap.CoppProfile = coppProf.(string)
	}
	if rtLooIpRange, ok := d.GetOk("underlay_routing_loopback_ip_range"); ok {
		configMap.UnderlayRoutingLoopbackIpRange = rtLooIpRange.(string)
	}
	if vtepLooIpRange, ok := d.GetOk("underlay_vtep_loopback_ip_range"); ok {
		configMap.UnderlayVtepLoopbackIpRange = vtepLooIpRange.(string)
	}
	if rpLooIpRange, ok := d.GetOk("underlay_rp_loopback_ip_range"); ok {
		configMap.UnderlayRpLoopbackIpRange = rpLooIpRange.(string)
	}
	if subIpRange, ok := d.GetOk("underlay_subnet_ip_range"); ok {
		configMap.UnderlaySubnetIpRange = subIpRange.(string)
	}
	if l2VniRange, ok := d.GetOk("layer_2_vxlan_vni_range"); ok {
		configMap.Layer2VxlanVniRange = l2VniRange.(string)
	}
	if l3VniRange, ok := d.GetOk("layer_3_vxlan_vni_range"); ok {
		configMap.Layer3VxlanVniRange = l3VniRange.(string)
	}
	if netVlanRange, ok := d.GetOk("network_vlan_range"); ok {
		configMap.NetworkVlanRange = netVlanRange.(string)
	}
	if vrfVlanRange, ok := d.GetOk("vrf_vlan_range"); ok {
		configMap.VrfVlanRange = vrfVlanRange.(string)
	}
	if subIfDot1qRange, ok := d.GetOk("subinterface_dot1q_range"); ok {
		configMap.SubinterfaceDot1qRange = subIfDot1qRange.(string)
	}
	if vrfLiteDeployment, ok := d.GetOk("vrf_lite_deployment"); ok {
		configMap.VrfLiteDeployment = vrfLiteDeployment.(string)
	}
	if vrfLifeSubnetIpRange, ok := d.GetOk("vrf_lite_subnet_ip_range"); ok {
		configMap.VrfLiteSubnetIpRange = vrfLifeSubnetIpRange.(string)
	}
	if vrfLiteSubnetMask, ok := d.GetOk("vrf_lite_subnet_mask"); ok {
		configMap.VrfLiteSubnetMask = strconv.Itoa(vrfLiteSubnetMask.(int))
	}
	if svcNetVlanRange, ok := d.GetOk("service_network_vlan_range"); ok {
		configMap.ServiceNetworkVlanRange = svcNetVlanRange.(string)
	}
	if rmSeqRange, ok := d.GetOk("route_map_sequence_number_range"); ok {
		configMap.RouteMapSequenceNumberRange = rmSeqRange.(string)
	}

	configMap.FabricTemplate = fabric.Template
	configMap.FabricName = fabric.Name
	fabric.Config = configMap
	idInt, _ := strconv.Atoi(d.Id())
	fabric.Id = idInt
	fabric.SetConfigDefaults()

	dn := fabric.Name
	durl := fmt.Sprintf("/rest/control/fabrics/%s", dn)
	cont, err := dcnmClient.Update(durl, &fabric)
	if err != nil {
		return err
	}

	d.SetId(stripQuotes(cont.S("id").String()))
	log.Println("[DEBUG] End of Update method ", d.Id())
	return resourceDCNMFabricRead(d, m)
}

func resourceDCNMFabricRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Get("name").(string)

	cont, err := getRemoteFabric(dcnmClient, dn)
	if err != nil {
		return err
	}

	setFabricAttributes(d, cont)

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMFabricDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)

	dn := d.Get("name").(string)

	durl := fmt.Sprintf("/rest/control/fabrics/%s", dn)
	_, err := dcnmClient.Delete(durl)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}
