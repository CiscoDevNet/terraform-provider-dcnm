---
layout: "dcnm"
page_title: "DCNM: dcnm_policy"
sidebar_current: "docs-dcnm-resource-policy"
description: |-
  Manages DCNM policy modules
---

# dcnm_policy #
Manages DCNM policy modules

## Example Usage ##

```hcl
resource "dcnm_policy" "second" {
    serial_number   =   "9BH270169LJ" 
    template_name   =   "aaa_radius_deadtime"
    template_props  =   {
                            "DTIME" : "3"
                            "AAA_GROUP" : "management"
                        }
    priority        =   500
    source          =   "Ethernet1/3_FABRIC"
    entity_name     =   "Ethernet1/3"
    entity_type     =   "INTERFACE"
    description     =   "This is demo policy."
    template_content_type   =   "TEMPLATE_CLI"

}
```

## Common Argument Reference ##

* `serial_number` - (Required) Serial number of switch under which policy will be created.
* `template_name` - (Required)  A unique name identifying the template. Please note that a template name can be used by multiple policies and hence a template name does not identify a policy uniquely.
* `template_props` - (Required) Properties of the templates related to template name.
* `template_content_type` - (Optional) Content type of the specified template.
* `priority` - (Optional) Priority of the policy.Default value is 500.
* `source` - (Optional) The source of the policy.
* `description`- (Optional) Description of the policy. The description may include the details regarding the policy.Default value is "".
* `entity_name`- (Optional) Name of the entity.i.e."SWITCH".
* `entity_type`- (Optional) Type of the entity.i.e."SWITCH".
* `template_content_type`- (Optional) Template content type of the policy.

#### `Note`: Destroying Policy will re-deploy the switch.

## Attribute Reference

*  `policy_id` - (Optional) A unique ID identifying a policy.
    NOTE: User can specify only empty string value.

## Importing ##

An existing policy can be [imported][docs-import] into this resource via its policy id using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import dcnm_policy.example <policyId>
```
