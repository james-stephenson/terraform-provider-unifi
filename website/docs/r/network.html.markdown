---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_network"
description: |-
  unifi_network manages LAN/VLAN networks.
---

# Resource: `unifi_network`

unifi_network manages LAN/VLAN networks.

## Example Usage

```terraform
variable "vlan_id" {
  default = 10
}

resource "unifi_network" "vlan" {
  name    = "wifi-vlan"
  purpose = "corporate"

  subnet       = "10.0.0.1/24"
  vlan_id      = var.vlan_id
  dhcp_start   = "10.0.0.6"
  dhcp_stop    = "10.0.0.254"
  dhcp_enabled = true
}
```

## Schema

### Required

- **name** (String, Required)
- **purpose** (String, Required)

### Optional

- **dhcp_enabled** (Boolean, Optional)
- **dhcp_lease** (Number, Optional) Defaults to `86400`.
- **dhcp_start** (String, Optional)
- **dhcp_stop** (String, Optional)
- **domain_name** (String, Optional)
- **id** (String, Optional)
- **igmp_snooping** (Boolean, Optional)
- **network_group** (String, Optional) Defaults to `LAN`.
- **subnet** (String, Optional)
- **vlan_id** (Number, Optional)


