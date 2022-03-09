terraform {
  required_providers {
    dcnm = {
      source = "CiscoDevNet/dcnm"
    }
  }
}

provider "dcnm" {
  username = ""
  password = ""
  url      = ""
  # expiry   = 900000
}

resource "dcnm_network" "first" {
  fabric_name     = "fab2"
  name            = "first"
  network_id      = "1234"
  description     = "first network from terraform"
  vrf_name        = "VRF1012"
  l2_only_flag    = false
  vlan_id         = 2300
  vlan_name       = "vlan1"
  ipv4_gateway    = "192.0.3.1/24"
  ipv6_gateway    = "2001:db8::1/64"
  mtu             = 1500
  secondary_gw_1  = "192.0.3.1/24"
  secondary_gw_2  = "192.0.3.1/24"
  arp_supp_flag   = true
  ir_enable_flag  = false
  mcast_group     = "239.1.2.2"
  dhcp_1          = "1.2.3.4"
  dhcp_2          = "1.2.3.5"
  dhcp_vrf        = "VRF1012"
  loopback_id     = 100
  tag             = "1400"
  rt_both_flag    = true
  trm_enable_flag = true
  l3_gateway_flag = true

  deploy = true
  attachments {
    serial_number = dcnm_inventory.example1.serial_number
    vlan_id       = 2300
    attach        = true
    switch_ports = [
      "Ethernet1/5",
      "Ethernet1/6"
    ]
  }
  attachments {
    serial_number = dcnm_inventory.example2.serial_number
    vlan_id       = 0
    attach        = false
  }
  attachments {
    serial_number = dcnm_inventory.example3.serial_number
    vlan_id       = 2300
    attach        = true
    switch_ports = ["Ethernet1/1",
    "Ethernet1/2"]
  }
}
