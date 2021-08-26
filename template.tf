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
  url      = "https://172.25.74.123/"
  platform = "nd"
  # expiry   = 900000
}

# resource "dcnm_template" "ex2"{
#   # name="avin"
# }
resource "dcnm_template" "ex1"{
    name = "avina"
    # file=file()
    content = <<EOF
##template properties
name=avina;
description = "avinaaa";
##
##template variables
#    Copyright (c) 2019 by Cisco Systems, Inc.
#    All rights reserved.

@(DisplayName="BGP AS #", Description="BGP Akkkutonomous System Number")
string BGP_AS;

@(DisplayName="VRF Name", IsVrfName=true)
string VRF_NAME;

@(DisplayName="Roudte map namSe", Description="Redistribute static route map")
string REDIST_ROUTE_MAP {
    defaultValue = FABRIC-RMAP-REDIST-SUBNET;
};

##
##template content

router bgp $$BGP_AS$$
  vrf $$VRF_NAME$$
    address-family ipv4 unicast
      redistribute static route-map $$REDIST_ROUTE_MAP$$
    address-family ipv6 unicast
      redistribute static route-map $$REDIST_ROUTE_MAP$$



##
EOF

}
data "dcnm_template" "ex"{
  name="${dcnm_template.ex1.id}"
}