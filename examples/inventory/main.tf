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
  config_timeout  = 10
  switch_config {
    ip   = ""
    role = "leaf"
  }
}
