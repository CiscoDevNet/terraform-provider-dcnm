provider "dcnm" {
  username = ""
  password = ""
  url      = ""
  # expiry   = 900000
}

resource "dcnm_interface" "second" {
  fabric_name = "fab2"
  name        = "loopback5"
  type        = "loopback"
  policy      = "int_loopback_11_1"

  switch_name_1             = "leaf1"
  ipv4                      = "1.2.3.4"
  loopback_tag              = "1234"
  vrf                       = "MyVRF"
  loopback_ls_routing       = "ospf"
  loopback_routing_tag      = "1234"
  loopback_router_id        = "10"
  loopback_replication_mode = "Multicast"
  description               = "creation from terraform"
  ipv6                      = "2001::0"
}
