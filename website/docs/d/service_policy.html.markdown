---
layout: "dcnm"
page_title: "DCNM: dcnm_service_policy"
sidebar_current: "docs-dcnm-data-source-service_policy"
description: |-
  Data source for DCNM Service Policy
---

# dcnm_vrf #
Data source for DCNM Service Policy

## Example Usage ##

```hcl

data "dcnm_service_policy" "example" {
  policy_name              = "SP-2"  
  fabric_name              = "external"
}

```

## Argument Reference ##

* `policy_name` - (Required) Name of Object Service Policy.
* `fabric_name` - (Required) Fabric name under which Service Policy should be created.


## Attribute Reference
* `attached_fabric_name` - Attached Fabric name of the Service Policy. 
* `dest_network` - Destination network of the Service Policy.
* `dest_vrf_name` - Destination VRF name of the Service Policy.
* `next_hop_ip` - Next hop IP of the Service Policy.
* `peering_name` - Peering name of the Service Policy. 
* `policy_template_name` - Policy template name of the Service Policy.
* `reverse_enabled` - Reverse enabled of the Service Policy.
* `reverse_next_hop_ip` - Reverse next hop IP of the Service Policy.
* `service_node_name` - Node name of the Service Policy.
* `source_network` - Source network of the Service Policy. 
* `source_vrf_name` - Source VRF name of the Service policy.
* `protocol` - Protocol of the Service Policy.
* `src_port` - Source port of the Service Policy. 
* `dest_port` - Destination Port of the Service Policy.
* `route_map_action` - Route map action of the Service Policy.
* `next_hop_action` - Next hop Action of the Service Policy.
* `fwd_direction` - Forward Direction of the Service Policy.
* `deploy` - Deploy of the Service Policy.
* `deploy_timeout` - Deploy timeout  of the Service Policy.