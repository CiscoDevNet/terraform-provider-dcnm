terraform {
  required_providers {
    dcnm = {
      source = "CiscoDevNet/dcnm"
    }
  }
}

provider "dcnm" {
  username = "admin"
  password = "c1sco123"
  url      = "https://10.122.18.70"
  insecure = true
  platform = "nd"
}

resource "dcnm_vrf_custom" "my_vrf" {
  fabric_name        = "SIMPL-BROWNFIELD"
  name               = "MyVRF"
  segment_id         = 50016
  template           = "JPMC_VRF_Universal"
  extension_template = "JPMC_VRF_Extension_Universal"
  template_props = {
    "advertiseDefaultRouteFlag" : "true",
    "vrfVlanId" : "123",
    "isRPExternal" : "false",
    "vrfDescription": ""
    "vrfSegmentId": "50016"
    "maxBgpPaths": "6"
    "maxIbgpPaths": "6"
    "borderMaxBgpPaths" : "16"
    "ipv6LinkLocalFlag": "false"
    "vrfRouteMap": "RM_ADVERTISE_CONNECTED_SVI"
    "ENABLE_NETFLOW": "false"
    "bgpPassword": ""
    "mtu": "9000"
    "multicastGroup": ""
    "isRPAbsent": "false"
    "advertiseHostRouteFlag": "true"
    "vrfVlanName": ""
    "trmEnabled": "false"
    "asn": "4201020601"
    "vrfIntfDescription" : ""
    "vrfName": "MyVRF"
  }
}