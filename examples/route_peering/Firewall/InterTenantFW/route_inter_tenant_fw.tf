resource "dcnm_route_peering" first{
    name = "tf5"
    attached_fabric = "terraform"
    deployment_mode = "InterTenantFW"
    service_fabric = "edge"
    option = "StaticPeering"
    service_networks {
        network_name = "net"
        network_type = "InsideNetworkFW"
        template_name = "Service_Network_Universal"
        vlan_id = 1000
        vrf_name = "check"
        gateway_ip_address = "12.32.2.1/23"
    }
    service_networks {
        network_name = "NET3"
        network_type = "OutsideNetworkFW"
        template_name = "Service_Network_Universal"
        vlan_id = 1000
        vrf_name = "check1"
        gateway_ip_address = "129.25.36.32/24"
    }
    service_node_name = "SN-1"
    service_node_type = "Firewall"
    routes {
        template_name = "service_static_route"
        vrf_name = "check"
        route_parmas = {

        }
    }
    routes {
        vrf_name = "check1"
        route_parmas = {
             "VRF_NAME": "check1"
   
        }
    }
    deploy=false
}