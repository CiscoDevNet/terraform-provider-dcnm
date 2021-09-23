resource "dcnm_policy" "example" {
    serial_number   =   "9LMU8W6W8VG" 
    template_name   =   "aaa_radius_deadtime"
    template_props  =   {
                            "DTIME" : "3"
                            "AAA_GROUP" : "management"
                        }
    priority        =   500
    description     =   "This is demo policy."

}
data "dcnm_policy" "example"{
  policy_id = "${dcnm_policy.second.policy_id}"
}
