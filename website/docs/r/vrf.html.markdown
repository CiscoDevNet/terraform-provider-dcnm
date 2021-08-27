---
layout: "dcnm"
page_title: "DCNM: dcnm_vrf"
sidebar_current: "docs-dcnm-resource-vrf"
description: |-
  Manages DCNM VRF
---

# dcnm_vrf

Manages DCNM VRF

## Example Usage

```hcl

resource "dcnm_vrf" "first" {
  fabric_name             = "fab2"
  name                    = "check"
  vlan_id                 = 2002
  segment_id              = "50016"
  vlan_name               = "check"
  description             = "vrf creation"
  intf_description        = "vrf"
  tag                     = "1250"
  max_bgp_path            = 2
  max_ibgp_path           = 4
  trm_enable              = false
  rp_external_flag        = true
  rp_address              = "1.2.3.4"
  loopback_id             = 15
  mutlicast_address       = "10.0.0.2"
  mutlicast_group         = "224.0.0.1/4"
  ipv6_link_local_flag    = "true"
  trm_bgw_msite_flag      = true
  advertise_host_route    = true
  advertise_default_route = "true"
  static_default_route    = false
  deploy                  = true
  attachments {
    serial_number = "9EQ00OGQYV6"
    vlan_id       = 2300
    attach        = false
    loopback_id   = 70
    loopback_ipv4 = "1.2.3.4"
    vrf_lite {
      auto_vrf_lite_flag = false
      dot1q_id = 2
      ip_mask = ""
      ipv6_mask = ""
      ipv6_neighbor = ""
      neighbor_asn = "500"
      neighbor_ip = "10.1.1.1"
      peer_vrf_name = "vrf_lite"
             }
    }
  }
}

```

## Argument Reference

- `name` - (Required) Name of Object VRF.
- `fabric_name` - (Required) Fabric name under which VRF should be created.
- `segment_id` - (Optional) VRF-Segment id. This field is auto-calculated if not provided. However while creating multiple VRFs in the same plan use this field to reserve the VRF id to avoid any conflicts due to concurrent execution.

<strong>Note: </strong> For auto-generation of segment-id while creating multiple VRFs in the same plan, Use the depends on functionality of terraform to avoid any segment-id conflicts.

- `vlan` - (Optional) Vlan Id for the VRF.
- `vlan_name` - (Optional) Vlan name for the VRF.
- `description` - (Optional) Description for the VRF.
- `intf_description` - (Optional) Intf description for the VRF.
- `tag` - (Optional) Tag for the VRF. Ranging from 0 to 4294967295.
- `max_bgp_path` - (Optional) Maximum BGP path value for the VRF. Ranging from 1 to 64.
- `max_ibgp_path` - (Optional) Maximum iBGP path value for the VRF. Ranging from 1 to 64.
- `trm_enable` - (Optional) Trm enable flag for the VRF. Allowed values are "true" and "false".
- `rp_external_flag` - (Optional) Rp external flag for the VRF. Allowed values are "true" and "false".
- `rp_address` - (Optional) Rp address for the VRF.
- `loopback_id` - (Optional) Loopback ip address for the VRF. Ranging from 0 to 1023.
- `mutlicast_group` - (Optional) Multicast group address for the VRF. Ranging from 224.0.0.0/4 to 239.255.255.255/4.
- `mutlicast_address` - (Optional) Multicast address for the VRF.
- `ipv6_link_local_flag` - (Optional) Ipv6 link local enable flag for the VRF. Allowed values are "true" and "false".
- `trm_bgw_msite_flag` - (Optional) Trm bgw multisite enable flag for the VRF. Allowed values are "true" and "false".
- `advertise_host_route` - (Optional) Advertise host route enable flag for the VRF. Allowed values are "true" and "false".
- `advertise_default_route` - (Optional) Advertise default route enable flag for the VRF. Allowed values are "true" and "false".
- `static_default_route` - (Optional) Configure static default route enable flag for the VRF. Allowed values are "true" and "false".
- `template` - (Optional) Template name for the VRF. Values allowed "Default_VRF_Universal". Default is "Default_VRF_Universal".
- `mtu` - (Optional) Mtu value for the VRF. Ranging from 68 to 9216.
- `extension_template` - (Optional) Extension Template name for the VRF. Values allowed are "Default_VRF_Extension_Universal". Default is "Default_VRF_Extension_Universal".
- `service_template` - (Optional) Service template name for the VRF.
- `source` - (Optional) Source for the VRF.

- `deploy` - (Optional) Deploy flag, used to deploy the VRF. Default value is "true".
- `deploy_timeout` - (Optional) Deployment timeout, used as the limiter for the deployment status check for VRF resource. It is in the unit of seconds and default value is "300".

- `attachments` - (Optional) Attachment Block, have information regarding the switches which should be attached or detached to/from VRF. If `deploy` is "true", then atleast one attachment must be configured.
- `attachments.serial_number` - (Required) Serial number of the switch.
- `attachments.vlan_id` - (Optional) Vlan ID for the switch associated with VRF. If not mentioned then VRF's default vlan id will be used for attachment.
- `attachments.attach` - (Optional) Attach flag for switch. Default value is "true".
- `attachments.free_form_config` - (Optional) Free form configuration for the switch attachment.
- `attachments.extension_values` - (Optional) Extension values for switch attachment.
- `attachments.loopback_id` - (Optional) Loopback id for the switch attachment.
- `attachments.loopback_ipv4` - (Optional) Loopback ipv4 address for the switch attachment.
- `attachments.loopback_ipv6` - (Optional) Loopback ipv6 address for the switch attachment.
- `attachments.vrf_lite` - (Optional) Vrf lite for the switch attachment.
- `attachments.vrf_lite.peer_vrf_name` - (Required) Name of vrf lite  for the switch attachment.
- `attachments.vrf_lite.dotq_id` - (Optional) Dotq id of  vrf lite for the switch attachment.
- `attachments.vrf_lite.ip_mask` - (Optional) Ip mask of vrf lite for the switch attachment.
- `attachments.vrf_lite.neighbor_ip` - (Optional) Neighbor ip of vrf lite for the switch attachment.
- `attachments.vrf_lite.neighbor_asn` - (Optional) Neighbor asn of vrf lite for the switch attachment.
- `attachments.vrf_lite.ipv6_mask` - (Optional) Ipv6 mask of vrf lite for the switch attachment.
- `attachments.vrf_lite.ipv6_neighbor` - (Optional) Ipv6 neighbor of vrf lite for the switch attachment.
- `attachments.vrf_lite.auto_vrf_lite_flag` - (Optional) Auto vrf lite flag of vrf lite for the switch attachment.




## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VRF.

## Importing

An existing VRF can be [imported][docs-import] into this resource via its fabric and name, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import dcnm_vrf.example <fabric_name>:<vrf_name>
```

