---
layout: "dcnm"
page_title: "DCNM: dcnm_interface"
sidebar_current: "docs-dcnm-data-source-interface"
description: |-
  Data source for DCNM interface module
---

# dcnm_interface #
Data source for DCNM interface module

## Example Usage ##

```hcl

data "dcnm_interface" "check" {
  serial_number = "${dcnm_interface.example.serial_number}"
  name          = "Ethernet1/1"
  type          = "ethernet"
  fabric_name   = "fab2"
}

```


## Argument Reference ##

* `serial_number` - (Required) Dn for the interface module.
* `name` - (Required) name of the interface.
* `type` - (Required) type of the interface. Allowed values are "loopback", "port-channel", "vpc", "sub-interface", "ethernet".
**NOTE**: Interface type of "sub-interface" is not supported in NDFC 12.

## Common Attribute Reference ##

* `fabric_name` - fabric name under which interface is created.
* `policy` - policy name for the interface.
* `admin_state` - administrative state for the interface.
* `deploy` - deploy flag for the deployment of interface.
* `switch_name_1` - name of the switch which is associated to the interface.

## Attribute Reference for loopback Interface ##

* `vrf` - vrf name for the loopback interface.
* `ipv4` - ipv4 address for the loopback interface.
* `ipv6` - ipv6 address for the loopback interface.
* `loopback_tag` - tag for the loopback interface.
* `loopback_routing_tag` - routing tag for the loopback interface.
* `loopback_ls_routing` - link state routing protocol for the loopback interface.
* `loopback_router_id` - router id for the loopback interface.
* `loopback_replication_mode` - replication mode for the loopback interface.
* `configuration` - configuration for the loopback interface.
* `description` - description for the loopback interface.

## Attribute Reference for port-channel Interface ##

* `pc_interface` - list of port channel member interface for port-channel interface.
* `access_vlans` - access vlans for the port-channel interface.
* `mode` - mode for the port-channel interface.
* `bpdu_guard_flag` - BPDU flag for the port-channel interface.
* `port_fast_flag` - port type fast flag for the port-channel interface.
* `mtu` - mtu for the port-channel interface.
* `allowed_vlans` - allowed vlans for the port-channel interface.
* `configuration` - configuration for the port-channel interface.
* `description` - description for the port-channel interface.

## Attribute Reference for vPC Interface ##

* `switch_name_2` - name of the second switch with which vpc is associated. 
* `vpc_peer1_id` - peer1 port-channel id for the vPC interface.
* `vpc_peer2_id` - peer2 port-channel id for the vPC interface.
* `vpc_peer1_interface` - list of peer1 member interface for the vPC interface.
* `vpc_peer2_interface` - list of peer2 member interface for the vPC interface.
* `mode` - mode for the vPC interface.
* `bpdu_guard_flag` - BPDU flag for the vPC interface.
* `port_fast_flag` - port type fast flag for the vPC interface.
* `mtu` - mtu for the vPC interface.
* `vpc_peer1_allowed_vlans` - peer1 allowed vlans for the vPC interface.
* `vpc_peer2_allowed_vlans` - peer2 allowed vlans for the vPC interface.
* `vpc_peer1_access_vlans` - peer1 access vlans for the vPC interface.
* `vpc_peer2_access_vlans` - peer2 access vlans for the vPC interface.
* `vpc_peer1_desc` - peer1 description for the vPC interface.
* `vpc_peer2_desc` - peer2 description for the vPC interface.
* `vpc_peer1_conf` - peer1 configuration for the vPC interface.
* `vpc_peer2_conf` - peer2 configuration for the vPC interface.

## Attribute Reference for sub-interface Interface ##

* `subinterface_vlan` - vlan for the sub-interface.
* `vrf` - vrf for the sub-interface.
* `ipv4` - ipv4 address for the sub-interface.
* `ipv6` - ipv6 address for the sub-interface.
* `ipv6_prefix` - ipv6 prefic for the sub-interface.
* `ipv4_prefix` - ipv4 prefix for the sub-interface.
* `subinterface_mtu` - mtu for the sub-interface.
* `configuration` - configuration for the sub-interface.
* `description` - description for the sub-interface.

## Attribute Reference for ethernet Interface ##

* `vrf` - vrf name for the ethernet interface.
* `bpdu_guard_flag` - BPDU flag for the ethernet interface.
* `port_fast_flag` - port type fast flag for the ethernet interface.
* `mtu` - mtu for the ethernet interface. 
* `ethernet_speed` - speed of the ethernet.
* `allowed_vlans` - allowed vlans for the ethernet interface.
* `configuration` - configuration for the ethernet.
* `description` - description for the ethernet.
* `ipv4` - ipv4 address for the ethernet.
* `ipv6` - ipv6 address for the ethernet.
* `ipv6_prefix` - ipv6 prefic for the ethernet.
* `ipv4_prefix` - ipv4 prefix for the ethernet.
* `access_vlans` -  access vlans for the ethernet interface.