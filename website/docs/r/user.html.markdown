---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_user"
description: |-
  unifi_user manages a user (or "client" in the UI) of the network, these are identified
by unique MAC addresses.

Users are "created" in the controller when observed on the network, so the resource defaults to allowing
itself to just take over management of a MAC address, but this can be turned off.
---

# Resource: `unifi_user`

unifi_user manages a user (or "client" in the UI) of the network, these are identified
by unique MAC addresses.

Users are "created" in the controller when observed on the network, so the resource defaults to allowing
itself to just take over management of a MAC address, but this can be turned off.

## Example Usage

```terraform
resource "unifi_user" "test" {
  mac  = "01:23:45:67:89:AB"
  name = "some client"
  note = "my note"

  fixed_ip   = "10.1.10.50"
  network_id = unifi_network.my_vlan.id
}
```

## Schema

- **allow_existing** - (Boolean, Optional)
- **blocked** - (Boolean, Optional)
- **fixed_ip** - (String, Optional)
- **hostname** - (String, Read-only)
- **id** - (String, Optional)
- **ip** - (String, Read-only)
- **mac** - (String, Required)
- **name** - (String, Required)
- **network_id** - (String, Optional)
- **note** - (String, Optional)
- **skip_forget_on_destroy** - (Boolean, Optional)
- **user_group_id** - (String, Optional)


