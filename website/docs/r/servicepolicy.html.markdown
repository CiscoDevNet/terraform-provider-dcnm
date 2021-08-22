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
  fabric_name              = "external"
  destination_network      = "dev_network_two"
  destination_vrf_name     = "dev_vrf_one"
  next_hop_ip              = "10.10.10.2"
  peering_name             = "RP-1"
  service_node_name        = "SN-1"
  service_node_type        = "Firewall"
  source_network           = "dev_network_one"
  source_vrf_name          = "dev_vrf_one"
}

```


## Argument Reference ##

* `policy_name` - (Required) Name of Object Service Policy.
* `fabric_name` - (Required) Fabric name under which Service Policy should be created.
* `attached_fabric_name` - (Optional) Attached Fabric name of the Service Policy. 
* `destination_network` - (Required) Destination network of the Service Policy.
* `destination_vrf_name` - (Required) Destination VRF name of the Service Policy.
* `next_hop_ip` - (Required) Next hop IP of the Service Policy.
* `peering_name` - (Required) Peering name of the Service Policy. 
* `policy_template_name` - (Optional) Policy template name of the Service Policy. Default value is "service_pbr".
* `reverse_enabled` - (Optional) Reverse enabled of the Service Policy. Default value is false.
* `reverse_next_hop_ip` - (Optional) Reverse next hop IP of the Service Policy.
* `service_node_name` - (Required) Node name of the Service Policy.
* `service_node_type` - (Required)Node Type of the Service Policy.
* `source_network` - (Required) Source network of the Service Policy. 
* `source_vrf_name` - (Required) Source VRF name of the Service policy.
* `protocol` - (Optional) Protocol of the Service Policy. Default value is "ip".
* `src_port` - (Optional) Source port of the Service Policy. Default value is "any".
* `dest_port` - (Optional) Destination Port of the Service Policy. Default value is "any".
* `route_map_action` - (Optional) Route map action of the Service Policy. Allowed values are "deny" and "permit". Default value is "permit".
* `next_hop_action` - (Optional) Next hop Action of the Service Policy. Allowed values are "none", "drop-on-fail" and "drop". Default value is "none".
* `fwd_direction` - (Optional) Forward Direction of the Service Policy.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Service Policy.

## Importing ##

An existing Service Policy can be [imported][docs-import] into this resource via its fabric and name, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import dcnm_service_policy.example <fabric_name>:<policy_name>
```