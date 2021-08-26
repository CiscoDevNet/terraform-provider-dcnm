---
layout: "dcnm"
page_title: "DCNM: dcnm_template"
sidebar_current: "docs-dcnm-data-source-template"
description: |-
  Data source for DCNM Template
---

# dcnm_template #
Data source for DCNM Template

## Example Usage ##

```hcl

data "dcnm_template" "ex"{
  name="test"
}

```


## Argument Reference ##

* `name` - (Required) name of Template.
* `content` - (Optional) File name or file content.


