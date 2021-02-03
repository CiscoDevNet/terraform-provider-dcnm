---
layout: "dcnm"
page_title: "DCNM: dcnm_inventory"
sidebar_current: "docs-dcnm-data-source-inventory"
description: |-
  Data source for DCNM inventory module
---

# dcnm_inventory #
Data source for DCNM inventory module

## Example Usage ##

```hcl

data "dcnm_inventory" "check" {
  fabric_name = "fab1"
  switch_name = "${dcnm_inventory.first.switch_name}" 
}

```


## Argument Reference ##

* `fabric_name` - (Required) fabric name under which inventory should be created.
* `switch_name` - (Required) name of switch.


## Attribute Reference

* `id` - Dn for the switch inventory.
* `ip` - Ip address of the switch.
* `role` - Role of the switch.
* `switch_db_id` - Db id for the switch.
* `serial_number` - Serial number of the switch.
* `model` - Model name of the switch.
* `mode` - Mode of the switch.
* `deploy` - (Optional) deploy flag for the switch.