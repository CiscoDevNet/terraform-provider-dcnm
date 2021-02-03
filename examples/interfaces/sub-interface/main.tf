provider "dcnm" {
  username = ""
  password = ""
  url      = ""
  # expiry   = 900000
}

resource "dcnm_interface" "second" {
  policy        = "int_subif_11_1"
  type          = "sub-interface"
  name          = "Ethernet1/41.8"
  fabric_name   = "fab2"
  switch_name_1 = "leaf1"

  vrf               = "MyVRF"
  subinterface_vlan = "8"
  ipv4              = "1.2.3.4"
  ipv6              = "2001::0"
  ipv4_prefix       = "24"
  ipv6_prefix       = "65"
  subinterface_mtu  = "9216"
  description       = "creation from terraform"

  deploy = false
}