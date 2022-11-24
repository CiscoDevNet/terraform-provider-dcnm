## 1.2.7 (November 25, 2022)
BUG FIXES:
- Fix idempotency of VRF Attachments with VRF Lite peering enabled (#108)
- Mark Child policies to be deleted when deleting source policy (#107)

## 1.2.6 (November 18, 2022)
BUG FIXES:
- Fix issue when VRF Attachments with VRF Lite peering enabled are not idempotent (#98)

## 1.2.5 (October 27, 2022)
BUG FIXES:
- Fix dcnm_network l3_gateway_flag not set correctly for Multi-Site Domain (MSD) fabrics (#99)

## 1.2.4 (October 12, 2022)
BUG FIXES:
- Fix waiting logic for switch config save in inventory
- Fix dcnm_template has missing parameters required for the HTTP request (#95)
- Fix dcnm_interface admin_state = false cause error (#94)
- Fix config removal triggered by destroy of the policy (#93)
- Improvement to documentation

## 1.2.3 (August 3, 2022)
BUG FIXES:
- Fix free_form_config typo issue (#90)
- Fix dcnm_network removal issue by setting DHCP attributes to non-computed (#88)
- Fix dcnm_policy destroy issue and dcnm_inventory issue (#85)

## 1.2.2 (April 20, 2022)
BUG FIXES:
- Add M1 MacOS support.
- Fix dcnm_policy resource destroy and deployment issue when modifying multiple policies.

## 1.2.1 (April 6, 2022)
BUG FIXES:
- Fix dcnm_policy resource destroy issue and add redeployement of switch to policy destroy workflow.
- Fix dcnm_rest resource to work with ndfc, accept any URL and not require payload when not needed.

## 1.2.0 (March 9, 2022)
IMPROVEMENTS:
- Add capability to post text file with dcnm_rest resource to support template validation by introducing payload_type attribute

BUG FIXES:
- Add provider source to examples
- Fix handling of Multicast Group setting when not provided by user in dcnm_network.

## 1.1.0 (December 17, 2021)
IMPROVEMENTS:
- New resource and data source for dcnm_policy, dcnm_route_peering, dcnm_service_node, dcnm_service_policy, dcnm_template
- Add support for NDFC 12.x
- Add support for secondary_gw_3, secondary_gw_4, dhcp_3, dhcp_vrf_2, dhcp_vrf_3, netflow_flag, svi_netflow_monitor, vlan_netflow_monitor, nve_id in dcnm_network resource and data source
- Add support for vrf_lite attachment in dcnm_vrf resource

BUG FIXES:
- Fix typo in bpdu_guard_flag in dcnm_interface resource

## 1.0.0 (May 28, 2021)

IMPROVEMENTS:
- Improved speed of dcnm_inventory resource
- Support for import on dcnm_interface
- Common timer for mode and configuration
- Role validation update for dcnm_inventory
- Support for VRF assignment on l2 interface

## 0.0.5 (March 11, 2021)

IMPROVEMENTS:
- For dcnm_vrf resource added a way to provide segment_id manually in order to create multiple VRFs in single plan.

## 0.0.4 (March 10, 2021)

IMPROVEMENTS:
- Added network_id to docs and examples.

## 0.0.3 (March 2, 2021)

IMPROVEMENTS:
- Removed computed from description for interface resource.
- Don't delete the ethernet interfaces
- Added handling for the safe destroy for vrf and network resource.

## 0.0.2 (March 1, 2021)

IMPROVEMENTS:
- Changed the inventory resource to have in-line blocks of switches.
- Network-id is user-configurable now.

## 0.0.1 (February 4, 2021)

- Initial Release
