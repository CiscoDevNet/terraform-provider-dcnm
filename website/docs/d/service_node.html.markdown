---
layout: "dcnm"
page_title: "DCNM: dcnm_service_node"
sidebar_current: "docs-dcnm-data-source-service_node"
description: |-
  Data source for DCNM Service Node
---

# dcnm_service_node

Data source for DCNM Service Node

## Example Usage

```hcl

data "dcnm_service_node" "example" {
  name           = "SN-1"
  service_fabric = "ISN"
}

```

## Argument Reference 

* `name` - (Required) Name of Object Service Node.
* `service_fabric` - (Required) Name of external fabric where the service node is located.

## Attribute Reference

* `id` - Attribute id is set to the name of the Service Node.
* `admin_state` - Admin state for the Service Node.
* `allowed_vlans` - Allowed VLAN names of the Service.
* `attached_fabric` - Name of attached easy fabric to which service node is attached.
* `attached_switch_interface_name` - Switch interfaces where the service node will be attached.
* `bpdu_guard_flag` - BPDU flag for the service node.
* `dest_fabric_name` - Destination fabric name of the service node.
* `dest_if_name` - Destination interface name of the service node.
* `dest_serial_number` - Destination serial number of the service node.
* `dest_switch_name` - Destination switch name of the service node.
* `form_factor` - Form factor of the service node.
* `interface_name` - Name of the service interface.
* `is_metaswitch` - Meta-switch flag of the service node.
* `link_template_name` - Link template name of the service node.
* `mtu` - MTU of the service node.
* `node_type` - Name of the service node type.
* `policy_description` - Description of the attached policy.
* `policy_id` - ID of the attached policy.
* `porttype_fast_enabled` - Port-type-fast flag of the service node.
* `priority` - Priority of the service node.
* `source_fabric_name` - Source fabric name of the service node.
* `source_if_name` - Source interface name of the service node.
* `source_serial_number` - Source serial number of the service node.
* `source_switch_name` - Source switch name of the service node.
* `speed` - bandwidth of the service node.
* `switches` - Serial Numbers of the switch where service node will be added.
