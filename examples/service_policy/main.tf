resource "dcnm_service_policy" "example" {
  policy_name              = "SP-2"
  service_fabric              = "edge"
  attached_fabric     = "terraform"
  dest_network             = "n1"
  dest_vrf_name            = "check1"
  next_hop_ip              = "10.10.10.2"
  peering_name             = "p1"
  service_node_name        = "SN-2"
  source_network           = "n2"
  source_vrf_name          = "check1"
}