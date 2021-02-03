---
layout: "dcnm"
page_title: "DCNM: dcnm_network"
sidebar_current: "docs-dcnm-data-source-network"
description: |-
  Data source for DCNM network
---

# dcnm_network #
Data source for DCNM network

## Example Usage ##

```hcl

data "dcnm_vrf" "check" {
  fabric_name = "fab2"
  name        = "check"
}

```


## Argument Reference ##

* `name` - (Required) name of network object.
* `fabric_name` - (Required) fabric name under which network exists.


## Attribute Reference

* `id` - Attribute id set to the Dn of the network.
* `display_name` -  display name for the network object.
* `description` -  description for the network.
* `vrf_name` -  name of the vrf which should be associated with the network.
* `l2_only_flag` -  layer 2 only flag for the network. 
* `vlan_id` -  vlan number for the network.
* `vlan_name` -  vlan name for the network.
* `ipv4_gateway` -  ipv4 address of gateway for the network.
* `ipv6_gateway` -  ipv6 address of gateway for the network.
* `mtu` -  mtu value for the network.
* `tag` -  tag for the Network.
* `secondary_gw_1` -  ipv4 secondary gateway 1 for the network.
* `secondary_gw_2` -  ipv4 secondary gateway 2 for the network.
* `arp_supp_flag` -  arp suppression flag for the network.
* `ir_enable_flag` -  ingress replication flag for the network.
* `mcast_group` -  multicast group address for the network.
* `dhcp_1` -  ipv4 address of DHCP server 1 for the network.
* `dhcp_2` -  ipv4 address of DHCP server 2 for the network.
* `dhcp_vrf` -  vrf name of DHCP server for the network.
* `loopback_id` -  loopback id for the network.
* `rt_both_flag` -  l2 VNI route-target both enable flag for the network.
* `trm_enable_flag` -  TRM enable flag for the network.
* `l3_gateway_flag` -  enable L3 gateway on border flag for the network. 
* `template` -  template name for the network. Default is "Default_VRF_Universal".
* `extension_template` -  extension Template name for the network. Default is "Default_Network_Extension_Universal".
* `service_template` -  service template name for the network.
* `source` -  source for the network.

* `deploy` - deploy flag, used to deploy the network.

* `attachments` - attachment block, have information regarding the switches which should be attached or detached to/from network.
* `attachments.serial_number` - serial number of the switch.
* `attachments.vlan_id` - vlan ID for the switch associated with network.
* `attachments.attach` - attach flag for switch.
* `attachments.switch_ports` - list of port name(i.e. interface names) for switch attachment.
* `attachments.untagged` -  untagged flag for switch attachment.
* `attachments.free_form_config` -  free form configuration for the switch attachment.
* `attachments.extension_values` -  extension values for switch attachment.
* `attachments.instance_values` -  instance values for switch attachment.
* `attachments.dot1_qvlan` -  dot1 qvlan for switch attachment.