---
layout: "dcnm"
page_title: "DCNM: dcnm_route_peering"
sidebar_current: "docs-dcnm-data-source-route_peering"
description: |-
  Data source for DCNM Route Peering
---

# dcnm_route_peering #
Data source for DCNM Route Peering

## Example Usage ##

```hcl

data "dcnm_route_peering" "example"{
  name = "tf"
  attached_fabric = "Test_fabric_1"
  service_fabric   = "testService"
  service_node_name = "snadc"
}

```


## Argument Reference ##

* `name` - (Required) name of route peering.
* `service_fabric` - (Required) Name of the target fabric for route peering operations.
* `attached_fabric` - (Required) Name of the target fabric for route peering operations.
* `service_node_name`- (Required) Name of service node under which route peering is will be created.


## Attribute Reference

* `deployment_mode` - (Required) Type of service node.Allowed values are "IntraTenantFW","InterTenantFW","OneArmADC","TwoArmADC","OneArmVNF".
* `next_hop_ip` - (Optional) Nexthop IPv4 information.NOTE: This object is applicable only when 'deploy_mode' is 'IntraTenantFW'
* `option` - (Required) Specifies the type of peering.Allowed values are "StaticPeering","EBGPDynamicPeering","None".
* `service_networks` - (Required) List of network under which peering will be created.
* `service_networks.network_name` - (Required) Network name.
* `reverse_next_hop_ip`- (Optional)  Reverse Nexthop IPv4 information, e.g., 192.169.1.100
* `service_networks.network_type` - (Required) Type of network.Allowed values are "InsideNetworkFW"(service node = Firewall),"OutsideNetworkFW"(service node = Firewall),"ArmOneADC"(service node = ADC),"ArmTwoADC"(service node = ADC),"ArmOneVNF"(service node= VNF).
* `service_networks.template_name` - (Required) Name of template.
* `service_networks.vrf_name` - (Required) VRF name under which network is created.
* `service_networks.vlan_id` - (Required) VLAN Id of network.
* `service_networks.gateway_ip_address` - (Required) IPv4 gateway information including the mask e.g. 192.168.1.1/24.
* `routes` - (Optional) Routing configuration.
* `routes.template_name` - (Optional) Template name for routing.
* `routes.route_parmas` - (Optional) NVPair map for routing.
* `routes.vrf_name` - (Optional) VRF name for routing.
* `deploy` - (Optional) A flag specifying if a route peering is to be deployed on the switches. Default value is "true".
* `deploy_timeout` - (Optional) Timeout seconds for deployment. Default value is 300s.
* `service_node_type` - (Required) Type of service node.Allowed values are "Firewall","VNF","ADC".