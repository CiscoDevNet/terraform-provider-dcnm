resource "dcnm_route_peering" "adc"{
   name = "tf6"
    attached_fabric_name = "terraform"
    deployment_mode = "OneArmADC"
    fabric_name = "edge"
    option = "EBGPDynamicPeering"
    service_networks {
        network_name = "netadc"
        network_type = "ArmOneADC"
        template_name = "Service_Network_Universal"
        vlan_id = 1002
        vrf_name = "check"
        gateway_ip_address ="124.168.2.1/24"
    }
    reverse_next_hop_ip = "124.168.2.10"
    service_node_name = "QA-ADC"
    service_node_type = "ADC"
    routes {
        template_name = "service_static_route"
        vrf_name = "check"
        route_parmas = {
                "VRF_NAME": "check"
        }
    }
    deploy = false
    deploy_timeout = 300
}