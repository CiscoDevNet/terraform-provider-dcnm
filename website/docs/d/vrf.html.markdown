---
layout: "dcnm"
page_title: "DCNM: dcnm_vrf"
sidebar_current: "docs-dcnm-data-source-vrf"
description: |-
  Data source for DCNM VRF
---

# dcnm_vrf #
Data source for DCNM VRF

## Example Usage ##

```hcl

data "dcnm_vrf" "check" {
  fabric_name   = "fab1"
  name          = "two" 
}

```


## Argument Reference ##

* `name` - (Required) Name of Object VRF.
* `fabric_name` - (Required) Fabric name under which VRF exists.


## Attribute Reference

* `id` - Attribute id set to the Dn of the VRF.
* `vlan` - Vlan Id for the VRF.
* `vlan_name` - Vlan name for the VRF.
* `description` - Description for the VRF.
* `intf_description` - Intf desscription for the VRF.
* `tag` - Tag for the VRF.
* `max_bgp_path` - Maximum BGP path value for the VRF.
* `max_ibgp_path` - Maximum iBGP path value for the VRF.
* `trm_enable` - Trm enable flag for the VRF. Allowed values are "true" and "false".
* `rp_external_flag` - Rp external flag for the VRF. Allowed values are "true" and "false".
* `rp_address` - Rp address for the VRF.
* `loopback_id` - Loopback ip address for the VRF.
* `mutlicast_group` - Multicast group address for the VRF.
* `mutlicast_address` - Multicast address for the VRF.
* `ipv6_link_local_flag` - Ipv6 link local enable flag for the VRF. Allowed values are "true" and "false".
* `trm_bgw_msite_flag` - Trm bgw multisite enable flag for the VRF. Allowed values are "true" and "false".
* `advertise_host_route` - Advertise host route enable flag for the VRF. Allowed values are "true" and "false".
* `advertise_default_route` - Advertise default route enable flag for the VRF. Allowed values are "true" and "false".
* `static_default_route` - Configure static default route enable flag for the VRF. Allowed values are "true" and "false".
* `template` - Template name for the VRF. Values allowed "Default_VRF_Universal". Default is "Default_VRF_Universal".
* `mtu` - Mtu value for the VRF. Ranginf from 68 to 9216.
* `extension_template` - Extension Template name for the VRF. Values allowed are "Default_VRF_Extension_Universal". Default is "Default_VRF_Extension_Universal".
* `service_template` - Service template name for the VRF.
* `source` - Source for the VRF.
* `deploy` - Deploy flag, used to deploy the VRF. Default value is "true".

* `attachments` - Attachment block, have information regarding the switches which should be attached or detached to/from VRF.
* `attachments.serial_number` - Serial number of the switch.
* `attachments.vlan_id` - Vlan ID for the switch associated with VRF.
* `attachments.attach` - Attach flag for switch.