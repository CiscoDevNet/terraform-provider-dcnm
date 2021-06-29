provider "dcnm" {
  username = ""
  password = ""
  url      = ""
  # expiry   = 900000
}

resource "dcnm_interface" "second" {
  policy        = "int_port_channel_access_host_11_1"
  type          = "port-channel"
  name          = "port-channel502"
  fabric_name   = "fab2"
  switch_name_1 = "leaf2"

  mode            = "active"
  bpdu_guard_flag = "true"
  mtu             = "jumbo"
  allowed_vlans   = "none"
  access_vlans    = "10"
  pc_interface    = ["e1/6", "eth1/9"]
}