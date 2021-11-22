# firewall "IntraTenantFW"
resource "dcnm_route_peering" "first" {
  name            = "RP-3"
  attached_fabric = "terraform"
  deployment_mode = "IntraTenantFW"
  service_fabric  = "edge"
  next_hop_ip     = "192.168.1.11"
  option          = "None"
  service_networks {
    network_name       = "net1"
    network_type       = "InsideNetworkFW"
    template_name      = "Service_Network_Universal"
    vlan_id            = 2000
    vrf_name           = "check"
    gateway_ip_address = "192.168.1.1/24"
  }
  service_networks {
    network_name       = "net2"
    network_type       = "OutsideNetworkFW"
    template_name      = "Service_Network_Universal"
    vlan_id            = 1000
    vrf_name           = "check"
    gateway_ip_address = "192.168.2.1/24"
  }
  service_node_name = "SN-1"
  service_node_type = "Firewall"
  deploy            = false
  #     deploy_timeout = 200
}