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
