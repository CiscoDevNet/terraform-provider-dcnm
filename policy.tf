terraform {
  required_providers {
    dcnm = {
      source = "hashicorp.com/edu/dcnm"
      version = "0.1"
    }
  }
}
provider "dcnm" {
  username = "admin"
  password = "ins3965!"
  url      = "https://172.25.74.91/"
#   platform = "nd"
  # expiry   = 900000
}
# data "dcnm_policy" p1{
#     policy_id = "POLICY-1490550"
# }
# resource "dcnm_policy" "second" {

#     serial_number   =   "9GKKMREPA58" 
#     template_name   =   "aaa_radius_deadtime"
#     template_props  =   {
#                             "DTIME" : "3",
#                              "AAA_GROUP" : "management"
#                         }
#     priority        =   1201
#     source          =   "policy"
#     entity_name     =   "policy"
#     entity_type     =   "policy"
#     description     =   "This is demo policy81."
#     template_content_type   =   "TEMPLATE_CLI"
#     deploy =false
# }
# resource "dcnm_service_node" "terraform_service_node" {
#   name                           = "terraform_service_node82"
#   node_type                      = "Firewall"
#   service_fabric                 = "testService"
#   interface_name                 = "vPC"
#   link_template_name             = "service_link_trunk"
#   switches                       = ["9GKKMREPA58","9OCBQ3Z52ZM"]
#   attached_switch_interface_name = "vPC1"
#   attached_fabric                = "main_fabric_2"
#  /*  form_factor                    = "Physical"
#   bpdu_guard_flag                = "no"
#   speed                          = "10Mb"
#   mtu                            = "jumbo"
#   allowed_vlans                  = "1000"
#   porttype_fast_enabled          = false
#   admin_state                    = false
#   policy_description             = "Node for Terraform" */
# }