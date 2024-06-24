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

resource "dcnm_network" "second" {
  fabric_name        = "fab3"
  name               = "second"
  network_id         = "1235"
  display_name       = "second"
  vrf_name           = "VRF1012"
  template           = "Template_Universal"
  extension_template = "Template_Extension_Universal"
  template_props = {
    "suppressArp" : "true"
    "gatewayIpAddress" : "10.0.4.1/24"
    "enableL3OnBorder" : "false"
    "vlanName" : "first"
    "enableIR" : "false"
    "mtu" : "1500"
    "rtBothAuto" : "false"
    "isLayer2Only" : "false"
    "mcastGroup" : "225.0.0.1"
    "vrfDhcp2" : "VRF1000"
    "dhcpServerAddr1" : "10.1.1.1"
    "dhcpServerAddr2" : "10.1.1.2"
    "vrfDhcp" : "VRF1000"
    "tag" : ""
    "vlanId" : "2303"
    "networkName" : "second"
    "segmentId" : "1235"
    "vrfName" : "VRF1012"
  }
  deploy = true
  attachments {
    serial_number = dcnm_inventory.example3.serial_number
    switch_ports  = ["Ethernet1/22"]
  }
}