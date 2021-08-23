---
layout: "dcnm"
page_title: "DCNM: dcnm_route_peering"
sidebar_current: "docs-dcnm-resource-route_peering"
description: |-
  Manages DCNM Route Peering modules
---

# dcnm_route_peering #
Manages DCNM Route Peering modules

## Example Usage ##

Route Peering when deployment mode is "IntraTenantFW" and service node "Firewall":
```hcl

resource "dcnm_route_peering" first{
    name = "RP-2"
    attached_fabric_name = "main_fabric_2"
    deployment_mode = "IntraTenantFW"
    fabric_name = "testService"
    next_hop_ip = "192.168.1.11"
    option = "None"
    service_networks {
        network_name = "net1"
        network_type = "InsideNetworkFW"
        template_name = "Service_Network_Universal"
        vlan_id = 1000
        vrf_name = "Test_VRF_2"
        gateway_ip_address = "192.168.1.1/24"
    }
    service_networks {
        network_name = "net2"
        network_type = "OutsideNetworkFW"
        template_name = "Service_Network_Universal"
        vlan_id = 1000
        vrf_name = "Test_VRF_2"
        gateway_ip_address = "192.168.2.1/24"
    }
    service_node_name = "SN-1"
    service_node_type = "Firewall"
    deploy = false
}

```
Route Peering when deployment mode is "InterTenantFW" and service node "Firewall":
```hcl
resource "dcnm_route_peering" first{
    name = "tf-5"
    attached_fabric_name = "Test_fabric_1"
    deployment_mode = "InterTenantFW"
    fabric_name = "testService"
    option = "StaticPeering"
    service_networks {
        network_name = "net"
        network_type = "InsideNetworkFW"
        template_name = "Service_Network_Universal"
        vlan_id = 1000
        vrf_name = "Test_VRF_2"
        gateway_ip_address = "12.32.2.1/23"
    }
    service_networks {
        network_name = "NET3"
        network_type = "OutsideNetworkFW"
        template_name = "Service_Network_Universal"
        vlan_id = 1000
        vrf_name = "Test_VRF_1"
        gateway_ip_address = "129.25.36.32/24"
    }
    service_node_name = "SN-1"
    service_node_type = "Firewall"
    routes {
        template_name = "service_static_route"
        vrf_name = "Test_VRF_2"
        route_parmas = {

        }
    }
    routes {
        vrf_name = "Test_VRF_1"
        route_parmas = {
             "VRF_NAME": "Test_VRF_1"
   
        }
    }
}
```
Route Peering when deployment mode is "OneArmADC" and service node "ADC":
```hcl
resource "dcnm_route_peering" first{
    name = "tfsnadc"
    attached_fabric_name = "Test_fabric_1"
    deployment_mode = "OneArmADC"
    fabric_name = "testService"
    option = "EBGPDynamicPeering"
    service_networks {
        network_name = "netadc"
        network_type = "ArmOneADC"
        template_name = "Service_Network_Universal"
        vlan_id = 1000
        vrf_name = "Test_VRF_2"
        gateway_ip_address ="124.168.2.1/24"
    }
    reverse_next_hop_ip = "124.168.2.10"
    service_node_name = "snadc"
    service_node_type = "ADC"
    routes {
        template_name = "service_static_route"
        vrf_name = "Test_VRF_2"
        route_parmas = {
                "VRF_NAME": "Test_VRF_1"
        }
    }
}
```
Route Peering when deployment mode is "TwoArmADC" and service node "ADC":
```hcl
resource "dcnm_route_peering" "adc3"{
   name = "tfsnadc12"
    attached_fabric_name = "Test_fabric_1"
    deployment_mode = "TwoArmADC"
    fabric_name = "testService"
    option = "StaticPeering"
    service_networks {
        network_name = "netadc"
        network_type = "ArmOneADC"
        template_name = "Service_Network_Universal"
        vlan_id = 1009
        vrf_name = "Test_VRF_2"
        gateway_ip_address ="124.168.2.1/24"
    }
    service_networks {
        network_name = "serviceNetwork2"
        network_type = "ArmTwoADC"
        template_name = "Service_Network_Universal"
        vlan_id = 1002
        vrf_name = "Test_VRF_2"
        gateway_ip_address ="123.168.3.1/24"
    }
    reverse_next_hop_ip = "124.168.2.10"
    service_node_name = "snadc"
    service_node_type = "ADC"
    routes {
        template_name = "service_static_route"
        vrf_name = "Test_VRF_2"
        route_parmas = {
                "VRF_NAME": "Test_VRF_1"
        }
    }
}
```
Route Peering when deployment mode is "OneArmVNF" and service node "VNF":
```hcl
resource "dcnm_route_peering" "adc3"{
   name = "tf"
    attached_fabric_name = "Test_fabric_1"
    deployment_mode = "OneArmVNF"
    fabric_name = "testService"
    option = "StaticPeering"
    service_networks {
        network_name = "netadc"
        network_type = "ArmOneVNF"
        template_name = "Service_Network_Universal"
        vlan_id = 1009
        vrf_name = "Test_VRF_2"
        gateway_ip_address ="124.168.2.1/24"
    }
    reverse_next_hop_ip = "124.168.2.10" #required
    service_node_name = "SN-3"
    service_node_type = "VNF"
    routes {
        template_name = "service_static_route"
        vrf_name = "Test_VRF_2"
        route_parmas = {
                "VRF_NAME": "Test_VRF_1"
        }
    }
}
```
## Argument Reference ##

* `name` - (Required) Name of route peering.
* `attached_fabric_name` - (Required) Name of the target fabric for route peering operations.
* `deployment_mode` - (Required) Type of service node.Allowed values are "IntraTenantFW","InterTenantFW","OneArmADC","TwoArmADC","OneArmVNF".
* `fabric_name` - (Required) Name of the target fabric for route peering operations.
* `next_hop_ip` - (Optional) Nexthop IPV4 information.NOTE: This object is applicable only when 'deploy_mode' is 'IntraTenantFW'
* `option` - (Required) Specifies the type of peering.Allowed values are "StaticPeering","EBGPDynamicPeering","None".
* `service_networks` - (Required) List of network under which peering will be created.
* `network_name` - (Required) Network name.
* `network_type` - (Required) Type of network.Allowed values are "InsideNetworkFW","OutsideNetworkFW","ArmOneADC","ArmTwoADC"."ArmOneVNF",
* `reverse_next_hop_ip`- (Optional)  Reverse Nexthop IPV4 information, e.g., 192.169.1.100
NOTE: This object is applicable only when 'deploy_mode' is either 'IntraTenantFW' or 'OneArmADC' or 'TwoArmADC'.
* `template_name` - (Required) Name of template.
* `vrf_name` - (Required) VRF name under which network is created.
* `gateway_ip_address` - (Required) IPV4 gateway information including the mask e.g. 192.168.1.1/24.
* `routes` - (Optional) Routing configuration.
* `deploy` - (Optional) A flag specifying if a route peering is to be deployed on the switches. Default value is "true".
* `service_node_name`- (Required) Name of service node under which route peering is will be created.
* `service_node_type` - (Required) Type of service node.Allowed values are "Firewall","VNF","ADC".

## Attribute Reference

* `status` - Route peering deployment status.

## Importing ##

An existing route peering can be [imported][docs-import] into this resource via its fabric and name, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import dcnm_route_peering.example <peering_name>:<external_fabric>:<service_node>:<attached_fabric_name>
```