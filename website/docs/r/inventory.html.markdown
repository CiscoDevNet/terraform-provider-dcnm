---
layout: "dcnm"
page_title: "DCNM: dcnm_inventory"
sidebar_current: "docs-dcnm-resource-inventory"
description: |-
  Manages DCNM inventory modules
---

# dcnm_inventory #
Manages DCNM inventory modules

## Example Usage ##

```hcl

resource "dcnm_inventory" "first" {
  fabric_name   = "fab2"
  switch_config {
    username      = "username for DCNM switch"
    password      = "password for DCNM switch"
    ip            = "ip for DCNM switch"
    preserve_config = "false"
    config_timeout = 10
    role = "leaf"
  }
}

```


## Argument Reference ##

* `fabric_name` - (Required) fabric name under which inventory should be created.

* `switch_config` - (Required) switch configuration block for inventory resource. It consists of the information regarding switches.
* `switch_config.ip` - (Required) ip Address of switch.
* `switch_config.username` - (Required) username for the the switch.
* `switch_config.password` - (Required) password for the the switch.
* `switch_config.role` - (Optional) role of the switch. Allowed values are "leaf", "spine", "border", "border_spine", "border_gateway", "border_gateway_spine", "super_spine", "border_super_spine", "border_gateway_super_spine".
* `switch_config.max_hops` - (Optional) maximum number hops for switch. Ranging from 0 to 10, default value is 0.
* `switch_config.auth_protocol` - (Optional) authentication protocol for switch. Mapping is as `0 : "MD5", 1: "SHA", 2 : "MD5_DES", 3 : "MD5_AES", 4 : "SHA_DES", 5 : "SHA_AES"`
* `switch_config.preserve_config` - (Optional) flag to preserve the configuration of switch. Default value is "false".
* `switch_config.platform` - (Optional) platform name for the switch.
* `switch_config.second_timeout` - (Optional) second timeout value for switch.
* `switch_config.config_timeout` - (Optional) configuration timeout value in minutes. Default value is "5".

* `deploy` - (Optional) deploy flag for the switch. Default value is "true".

## Attribute Reference

* `id` - Dn for the switch inventory.
* `switch_config` - Switch configuration block for inventory.
* `switch_config.switch_name` - Name of the switch.
* `switch_config.switch_db_id` - DB ID for the switch.
* `switch_config.serial_number` - Serial number of the switch.
* `switch_config.model` - Model name of the switch.
* `switch_config.mode` - Mode of the switch.

## Importing ##

An existing switch inventory can be [imported][docs-import] into this resource via its fabric and name, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import dcnm_inventory.example <fabric_name>:<switch_name>
```