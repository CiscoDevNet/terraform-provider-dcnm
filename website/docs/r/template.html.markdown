---
layout: "dcnm"
page_title: "DCNM: dcnm_template"
sidebar_current: "docs-dcnm-resource-template"
description: |-
  Manages DCNM Template
---

# dcnm_template #
Manages DCNM Template

## Example Usage ##

```hcl

resource "dcnm_template" "example" {
  name = "test"
  content = file("<<TXT File Name>></TXT>")
}

resource "dcnm_template" "example1" {
      name = "test"
    content = <<EOF
##template properties
name=test;
description = "Created template resource from terraform";
##
##template variables
#    Copyright (c) 2019 by Cisco Systems, Inc.
#    All rights reserved.

@(DisplayName="BGP AS #", Description="BGP Autonomous System Number")
string BGP_AS;

@(DisplayName="VRF Name", IsVrfName=true)
string VRF_NAME;

@(DisplayName="Roudte map namSe", Description="Redistribute static route map")
string REDIST_ROUTE_MAP {
    defaultValue = FABRIC-RMAP-REDIST-SUBNET;
};

##
##template content

router bgp $$BGP_AS$$
vrf $$VRF_NAME$$
    address-family ipv4 unicast
    redistribute static route-map $$REDIST_ROUTE_MAP$$
    address-family ipv6 unicast
    redistribute static route-map $$REDIST_ROUTE_MAP$$



##
EOF
}


```


## Argument Reference ##

* `name` - (Required) Name of Template.
* `content` - (Required) Content of file or file name.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the template.

## Importing ##

An existing Template can be [imported][docs-import] into this resource via template name, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import dcnm_template.example <template_name>
```