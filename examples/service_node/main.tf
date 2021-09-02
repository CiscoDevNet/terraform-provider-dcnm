resource "dcnm_service_node" "example" {
  name                           = "SN-1"
  node_type                      = "Firewall"
  service_fabric                 = "ISN"
  interface_name                 = "node"
  link_template_name             = "service_link_trunk"
  switches                       = ["9O7Q9J652MN"]
  attached_switch_interface_name = "Ethernet1/9"
  attached_fabric                = "terraform"
  form_factor                    = "Virtual"
  bpdu_guard_flag                = "no"
  speed                          = "Auto"
  mtu                            = "jumbo"
  allowed_vlans                  = "none"
  porttype_fast_enabled          = true
  admin_state                    = true
  policy_description             = "from terraform"
}