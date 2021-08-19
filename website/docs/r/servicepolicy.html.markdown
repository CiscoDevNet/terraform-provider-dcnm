---
layout: "dcnm"
page_title: "DCNM: dcnm_service_policy"
sidebar_current: "docs-dcnm-resource-service-policy"
description: |-
  Manages DCNM Service Policy
---

# dcnm_service_policy #
Manages DCNM Service Policy

## Example Usage ##

```hcl

resource "dcnm_service_policy" "example" {

  policy_name              = "SP-2"  
  fabric_name              = "fab1"
  attached_fabric_name     = "check"
  destination_network      = "12.1.1.2/32"
  destination_network_name = "destNet1"
  destination_vrf_name     = "vrf1"
  next_hop_ip              = "1.2.3.4"
  peering_name             = "RP-1"
  policy_template_name     = "service_pbr"
  reverse_enabled          = true
  reverse_next_hop_ip      = "2.3.4.5"
  service_node_name        = "sn1"
  service_node_type        = "Firewall"
  source_network           = "11.1.1.1/24"
  source_network_name      = "srcNet1"
  source_vrf_name          = "vrf1"
  
}

```


## Argument Reference ##

* `policy_name` - (Required) Name of Object Service Policy.
* `fabric_name` - (Required) Fabric name under which Service Policy should be created.
* `attached_fabric_name` - (Required) Attached Fabric name of the Service Policy. 
* `destination_network` - (Required) Destination network IP of the Service Policy.
* `destination_network_name` - (Required) Destination network name of the Service Policy.
* `destination_vrf_name` - (Required) Destination VRF name of the Service Policy.
* `next_hop_ip` - (Required) Next hop IP of the Service Policy.
* `peering_name` - (Required) Peering name of the Service Policy. 
* `policy_template_name` - (Required) Policy template name of the Service Policy. 
* `reverse_enabled` - (Required) Reverse enabled of the Service Policy.
* `reverse_next_hop_ip` - (Required) Reverse next hop IP of the Service Policy.
* `service_node_name` - (Required) Node name of the Service Policy.
* `service_node_type` - (Required)Node Type of the Service Policy.
* `source_network` - (Required) Source network of the Service Policy.
* `source_network_name` - (Required) Source network name of the Service Policy. 
* `source_vrf_name` - (Required) Source VRF name of the Service policy. 