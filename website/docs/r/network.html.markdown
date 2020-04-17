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

- **dhcp_enabled** - (Boolean, Optional)
- **dhcp_lease** - (Number, Optional)
- **dhcp_start** - (String, Optional)
- **dhcp_stop** - (String, Optional)
- **domain_name** - (String, Optional)
- **id** - (String, Optional)
- **igmp_snooping** - (Boolean, Optional)
- **name** - (String, Required)
- **network_group** - (String, Optional)
- **purpose** - (String, Required)
- **subnet** - (String, Optional)
- **vlan_id** - (Number, Optional)


