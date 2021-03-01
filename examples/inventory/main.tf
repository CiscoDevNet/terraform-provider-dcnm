provider "dcnm" {
  username = ""
  password = ""
  url      = ""
  # expiry   = 900000
}

resource "dcnm_inventory" "first" {
  fabric_name   = "fab2"
  switch_config {
    username      = ""
    password      = ""
    ip            = ""
    preserve_config = "false"
    config_timeout = 10
    role = "leaf"
  }
}
