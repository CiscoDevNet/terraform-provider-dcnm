---
layout: "dcnm"
page_title: "DCNM: dcnm_policy"
sidebar_current: "docs-dcnm-data-source-policy"
description: |-
  Data source for DCNM Policy
---

# dcnm_policy #
Data source for DCNM Policy

## Example Usage ##

```hcl

data "dcnm_policy" "example" {
  policy_id   = "POLICY-1197060"
}

```


## Argument Reference ##

* `policy_id` - (Required) A unique ID identifying a policy.
   NOTE: User can specify only empty string value.


## Attribute Reference

* `serial_number` - Serial number of switch under which policy will be created.
* `template_name` -  A unique name identifying the template. Please note that a template name can be used by multiple policies and hence a template name does not identify a policy uniquely.
* `template_props` - Properties of the templates related to template name.
* `priority` - Priority of the policy.Default value is 500.
* `source` - The source of the policy.
* `description`- Description of the policy. The description may include the details regarding the policy.Default value is "".
* `entity_name`- Name of the entity.i.e."SWITCH".
* `entity_type`- Type of the entity.i.e."SWITCH".