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

resource "dcnm_rest" "first" {
  path    = "/rest/top-down/fabrics/fab2/networks"
  payload = <<EOF
  {
    "displayName": "ashutosh",
    "fabric": "fab2",
    "networkExtensionTemplate": "Default_Network_Extension_Universal",
    "networkId": "30006",
    "networkName": "import",
    "networkTemplate": "Default_Network_Universal",
    "networkTemplateConfig": "{\"networkName\":\"import\",\"segmentId\":\"30006\",\"vlanId\":2300,\"mtu\":1500,\"gatewayIpAddress\":\"192.0.3.2/24\",\"gatewayIpV6Address\":\"2001:db8::1/64\",\"vlanName\":\"vlan2\",\"intfDescription\":\"second network from terraform\",\"secondaryGW1\":\"192.0.3.1/24\",\"secondaryGW2\":\"192.0.3.1/24\",\"suppressArp\":true,\"mcastGroup\":\"239.1.2.2\",\"dhcpServerAddr1\":\"1.2.3.4\",\"dhcpServerAddr2\":\"1.2.3.4\",\"vrfDhcp\":\"VRF1012\",\"loopbackId\":100,\"tag\":\"1400\",\"trmEnabled\":true,\"rtBothAuto\":true,\"enableL3OnBorder\":true}",
    "vrf": "MyVRF"
  }
 EOF 
}

resource "dcnm_rest" "template_validate" {
  path         = "/rest/config/templates/validate"
  method       = "POST"
  payload      = file("payload.txt")
  payload_type = "text"
}
