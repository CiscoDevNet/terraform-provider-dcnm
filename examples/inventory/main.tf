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

resource "dcnm_inventory" "first" {
  fabric_name     = "fab2"
  username        = ""
  password        = ""
  preserve_config = "false"
  config_timeout  = 300
  switch_config {
    ip   = ""
    role = "leaf"
  }
}
