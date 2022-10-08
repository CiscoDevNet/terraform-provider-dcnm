---
layout: "dcnm"
page_title: "Provider: DCNM"
sidebar_current: "docs-dcnm-index"
description: |-
  The Cisco DCNM provider is used to interact with the resources provided by Cisco DCNM/NDFC.
  The provider needs to be configured with the proper credentials before it can be used.
---
  

# Overview

The Cisco Nexus Dashboard Fabric Controller (NDFC, formerly DCNM) Terraform Provider is used to manage constructs on the Cisco DCNM/NDFC platform. It lets users represent the infrastructure as code and provides a way to enforce state on the infrastructure managed by Terraform. Customers can use this provider to integrate the Terraform configuration with their DevOps pipeline to manage the DCNM/NDFC fabric policies in a more flexible, consistent and reliable way.

The provider needs to be configured with the proper credentials before it can be used.

To learn more about the DCNM/NDFC, visit the [Cisco Nexus Dashboard Fabric Controller product overview](https://www.cisco.com/c/en/us/products/cloud-systems-management/prime-data-center-network-manager/index.html).

## Example Usage

```hcl
terraform {
  required_providers {
    dcnm = {
      source = "CiscoDevNet/dcnm"
    }
  }
}

# Configure the provider with your Cisco dcnm/ndfc credentials.
provider "dcnm" {
  # cisco-dcnm/ndfc user name
  username = "admin"
  # cisco-dcnm/ndfc password
  password = "password"
  # cisco-dcnm/ndfc url
  url      = "https://my-cisco-dcnm.com"
  insecure = true
  platform = "dcnm"
}

resource "dcnm_vrf" "test-vrf" {
  fabric_name = "fab1"
  name = "MyVRF"
  description = "This VRF is created by Terraform"
}
```

## Argument Reference

Following provider configuration arguments are supported within the `provider "dcnm"` block.

* `username` - (Required) This is the Cisco DCNM/NDFC username, which is required to authenticate with CISCO DCNM/NDFC.
* `password` - (Required) Password of the user mentioned in username argument. It is required when you want to use token-based authentication.
* `url` - (Required) The URL for Cisco DCNM/NDFC.
* `insecure` - (Optional) This determines whether to use insecure HTTP connection or not. Default value is `true`.
* `platform` - (Optional) NDFC/DCNM Platform information (Nexus-Dashboard/DCNM). Allowed values are "nd" or "dcnm". Default value is "dcnm".
