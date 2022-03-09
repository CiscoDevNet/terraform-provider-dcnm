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

resource "dcnm_vrf" "first" {
  fabric_name             = "fab2"
  name                    = "check"
  vlan_id                 = 2002
  segment_id              = "50016"
  vlan_name               = "check"
  description             = "vrf creation"
  intf_description        = "vrf"
  tag                     = "1250"
  max_bgp_path            = 2
  max_ibgp_path           = 4
  trm_enable              = false
  rp_external_flag        = false
  rp_address              = "1.2.3.4"
  loopback_id             = 15
  mutlicast_address       = "10.0.0.2"
  mutlicast_group         = "224.0.0.1/4"
  ipv6_link_local_flag    = "true"
  trm_bgw_msite_flag      = false
  advertise_host_route    = false
  advertise_default_route = "true"
  static_default_route    = false
  deploy                  = true
  attachments {
    serial_number = "9ZGMF8CBZK5"
    vlan_id       = 2300
    attach        = true
    loopback_id   = 70
    loopback_ipv4 = "1.2.3.4"
  }
}

resource "dcnm_vrf" "second" {
  fabric_name             = "fab2"
  name                    = "check2"
  vlan_id                 = 2003
  segment_id              = "50017"
  vlan_name               = "check2"
  description             = "vrf creation"
  intf_description        = "vrf"
  tag                     = "1250"
  max_bgp_path            = 2
  max_ibgp_path           = 4
  trm_enable              = false
  rp_external_flag        = false
  rp_address              = "1.2.3.4"
  loopback_id             = 15
  mutlicast_address       = "10.0.0.2"
  mutlicast_group         = "224.0.0.1/4"
  ipv6_link_local_flag    = "true"
  trm_bgw_msite_flag      = false
  advertise_host_route    = false
  advertise_default_route = "true"
  static_default_route    = false
  deploy                  = true
  attachments {
    serial_number = "9ZGMF8CBZK5"
    vlan_id       = 2300
    attach        = true
    loopback_id   = 70
    loopback_ipv4 = "1.2.3.4"
  }
}