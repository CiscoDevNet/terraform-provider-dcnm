---
layout: "dcnm"
page_title: "DCNM: dcnm_service_node"
sidebar_current: "docs-dcnm-resource-service_node"
description: |-
  Manages DCNM Service Node
---

# dcnm_service_node

Manages DCNM Service Node

## Example Usage

```hcl

resource "dcnm_service_node" "example" {
  name                           = "SN-1"
  node_type                      = "Firewall"
  service_fabric                 = "ISN"
  interface_name                 = "node"
  link_template_name             = "service_link_trunk"
  switches                       = ["9O7Q9J652MN"]
  attached_switch_interface_name = "Ethernet1/9"
  attached_fabric                = "terraform"
  form_factor                    = "Virtual"
  bpdu_guard_flag                = "no"
  speed                          = "Auto"
  mtu                            = "jumbo"
  allowed_vlans                  = "none"
  porttype_fast_enabled          = true
  admin_state                    = true
  policy_description             = "from terraform"
}

```

## Argument Reference

- `name` - (Required) Name of Object Service Node.
- `node_type` - (Required) Name of the service node type. Aloowed values are "Firewall", "ADC" and "VNF".
- `service_fabric` - (Required) Name of external fabric where the service node is located.
- `attached_fabric` - (Required) Name of attached easy fabric to which service node is attached.
- `attached_switch_interface_name` - (Required) Switch interfaces where the service node will be attached.
- `interface_name` - (Required) Name of the service interface.
- `link_template_name` - (Optional) Link template name of the service node.
- `switches` - (Required) List of serial Numbers of the switch where service node will be added.
- `admin_state` - (Optional) Admin state for the Service Node. Allowed values are true and false. Default value is true.
- `allowed_vlans` - (Optional) Allowed vlan names of the Service. Default value is "none".
- `bpdu_guard_flag` - (Optional) BPDU flag for the service node. Allowed values are "yes" and "no". Default value is "no".
- `form_factor` - (Optional) Form factor of the service node. Allowed values are "Physical" and "Virtual". Default value is "Virtual".
- `mtu` - (Optional) MTU of the service node. Default value is "jumbo".
- `policy_description` - (Optional) Description of the attached policy.
- `porttype_fast_enabled` - (Optional) Port-type-fast flag of the service node. Allowed values are true and false. Default value is true.
- `speed` - (Optional) bandwidth of the service node. Default value is "Auto".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Service Node.
