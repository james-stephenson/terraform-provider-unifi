---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_wlan"
description: |-
  unifi_wlan manages a WiFi network / SSID.
---

# Resource: `unifi_wlan`

unifi_wlan manages a WiFi network / SSID.

## Example Usage

```terraform
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "wifi" {
  name          = "myssid"
  vlan_id       = 10
  passphrase    = "12345678"
  wlan_group_id = data.unifi_wlan_group.default.id
  user_group_id = data.unifi_user_group.default.id
  security      = "wpapsk"
}
```

## Schema

- **hide_ssid** - (Boolean, Optional) Indicates whether or not to hide the SSID from broadcast.
- **id** - (String, Optional)
- **is_guest** - (Boolean, Optional) Indicates that this is a guest WLAN and should use guest behaviors.
- **mac_filter_enabled** - (Boolean, Optional) Indicates whether or not the MAC filter is turned of for the network.
- **mac_filter_list** - (Set of String, Optional) List of MAC addresses to filter (only valid if `mac_filter_enabled` is `true`).
- **mac_filter_policy** - (String, Optional) MAC address filter policy (only valid if `mac_filter_enabled` is `true`).
- **multicast_enhance** - (Boolean, Optional) Indicates whether or not Multicast Enhance is turned of for the network.
- **name** - (String, Required) The SSID of the network.
- **passphrase** - (String, Optional) The passphrase for the network, this is only required if `security` is not set to `open`.
- **security** - (String, Required) The type of WiFi security for this network. Valid values are: `wpapsk`, `wpaeap`, and `open`.
- **user_group_id** - (String, Required) ID of the user group to use for this network.
- **vlan_id** - (Number, Optional) VLAN ID for the network, defaults to `1`.
- **wlan_group_id** - (String, Required) ID of the WLAN group to use for this network.


