provider "dcnm" {
  username = ""
  password = ""
  url      = ""
  # expiry   = 900000
}

resource "dcnm_interface" "second" {
  policy        = "int_vpc_trunk_host_11_1"
  type          = "vpc"
  name          = "vPC1"
  fabric_name   = "fab2"
  switch_name_1 = "leaf1"

  switch_name_2           = "leaf2"
  vpc_peer1_id            = "501"
  vpc_peer2_id            = "502"
  mode                    = "active"
  bpdu_guard_flag         = "true"
  mtu                     = "jumbo"
  vpc_peer1_allowed_vlans = "none"
  vpc_peer2_allowed_vlans = "none"
  vpc_peer1_access_vlans  = "10"
  vpc_peer2_access_vlans  = "20"
  vpc_peer1_interface     = ["e1/5", "eth1/7"]
  vpc_peer2_interface     = ["e1/5", "eth1/7"]

  deploy = false
}
