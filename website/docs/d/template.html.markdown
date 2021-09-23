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
* `description` - (Optional) Description of template.
* `supported_platforms` - (Optional) Platform supported by the template.
* `template_type` - (Optional) Type of template.
* `template_content_type` - (Optional) Content type of template.
* `tags` - (Optional) Tag of template.
* `template_sub_type` - (Optional) Sub type of template.


