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

* `name` - (Required) name of Object VRF.
* `fabric_name` - (Required) fabric name under which VRF exists.


## Attribute Reference

* `id` - attribute id set to the Dn of the VRF.
* `vlan` - vlan Id for the VRF.
* `vlan_name` - vlan name for the VRF.
* `description` - description for the VRF.
* `intf_description` - intf desscription for the VRF.
* `tag` - tag for the VRF.
* `max_bgp_path` - maximum BGP path value for the VRF.
* `max_ibgp_path` - maximum iBGP path value for the VRF.
* `trm_enable` - trm enable flag for the VRF. Allowed values are "true" and "false".
* `rp_external_flag` - rp external flag for the VRF. Allowed values are "true" and "false".
* `rp_address` - rp address for the VRF.
* `loopback_id` - loopback ip address for the VRF.
* `mutlicast_group` - multicast group address for the VRF.
* `mutlicast_address` - multicast address for the VRF.
* `ipv6_link_local_flag` - ipv6 link local enable flag for the VRF. Allowed values are "true" and "false".
* `trm_bgw_msite_flag` - trm bgw multisite enable flag for the VRF. Allowed values are "true" and "false".
* `advertise_host_route` - advertise host route enable flag for the VRF. Allowed values are "true" and "false".
* `advertise_default_route` - advertise default route enable flag for the VRF. Allowed values are "true" and "false".
* `static_default_route` - configure static default route enable flag for the VRF. Allowed values are "true" and "false".
* `template` - template name for the VRF. Values allowed "Default_VRF_Universal". Default is "Default_VRF_Universal".
* `mtu` - mtu value for the VRF. Ranginf from 68 to 9216.
* `extension_template` - extension Template name for the VRF. Values allowed are "Default_VRF_Extension_Universal". Default is "Default_VRF_Extension_Universal".
* `service_template` - service template name for the VRF.
* `source` - source for the VRF.
* `deploy` - deploy flag, used to deploy the VRF. Default value is "true".

* `attachments` - attachment Block, have information regarding the switches which should be attached or detached to/from VRF.
* `attachments.serial_number` - serial number of the switch.
* `attachments.vlan_id` - vlan ID for the switch associated with VRF.
* `attachments.attach` - attach flag for switch. Default value is "true".