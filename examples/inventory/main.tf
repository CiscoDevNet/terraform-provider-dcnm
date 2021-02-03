provider "dcnm" {
  username = ""
  password = ""
  url      = ""
  # expiry   = 900000
}

resource "dcnm_inventory" "first" {
  fabric_name   = "fab1"
  username      = "username"
  password      = "password"
  ip            = "172.25.74.93"
  max_hops      = 0
  auth_protocol = 0
  deploy        = true
}
