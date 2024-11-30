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
  platform = ""
}

resource "dcnm_fabric" "easy_default" {
  name = "FAB1"
  asn  = 65003
}

resource "dcnm_fabric" "easy_default_custom" {
  name                               = "FAB2"
  asn                                = 65003
  overlay_mode                       = "cli"
  underlay_routing_loopback_ip_range = "192.168.0.0/20"
  underlay_vtep_loopback_ip_range    = "192.168.1.0/20"
  underlay_subnet_ip_range           = "10.0.0.0/16"
  ospf_bfd                           = true
  pim_bfd                            = true
  enable_vxlan_oam                   = false
  enable_nx_api                      = false
  enable_nx_api_on_http              = false
  enable_ndfc_as_trap_host           = false

}

resource "dcnm_fabric" "custom" {
  name     = "FAB3"
  asn      = 65001
  template = "LAN_Monitor"
  template_props = {
    "FABRIC_NAME" : "FAB3",
    "FABRIC_TECHNOLOGY" : "LANMonitor",
    "FABRIC_TYPE" : "LANMonitor",
    "FF" : "LANMonitor",
    "IS_READ_ONLY" : "true"
  }
}
