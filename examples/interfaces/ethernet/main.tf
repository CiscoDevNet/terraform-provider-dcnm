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

resource "dcnm_interface" "second" {
  policy        = "int_trunk_host_11_1"
  type          = "ethernet"
  name          = "Ethernet1/1"
  switch_name_1 = "leaf1"
  fabric_name   = "fab2"

  ethernet_speed  = "Auto"
  bpdu_guard_flag = "no"
  allowed_vlans   = "none"
  mtu             = "jumbo"
  port_fast_flag  = true
}