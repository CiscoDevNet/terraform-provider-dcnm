resource "dcnm_route_peering" "adc3" {
  name            = "tf"
  attached_fabric = "Test_fabric_1"
  deployment_mode = "OneArmVNF" # Should not have "None" peering option
  service_fabric  = "testService"
  option          = "StaticPeering"
  service_networks {
    network_name       = "netadc"
    network_type       = "ArmOneVNF"
    template_name      = "Service_Network_Universal"
    vlan_id            = 1009
    vrf_name           = "Test_VRF_2"
    gateway_ip_address = "124.168.2.1/24"
  }
  reverse_next_hop_ip = "124.168.2.10" #required
  service_node_name   = "SN-3"
  service_node_type   = "VNF"
  routes {
    template_name = "service_static_route"
    vrf_name      = "Test_VRF_2"
    route_parmas = {
      "VRF_NAME" : "Test_VRF_1"
    }
  }
}
