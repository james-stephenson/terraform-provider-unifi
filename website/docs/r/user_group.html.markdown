---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_user_group"
description: |-
  unifi_user_group manages a user group (called "client group" in the UI), which can be used
to limit bandwidth for groups of users.
---

# Resource: `unifi_user_group`

unifi_user_group manages a user group (called "client group" in the UI), which can be used
to limit bandwidth for groups of users.

## Example Usage

```terraform
resource "unifi_user_group" "wifi" {
  name = "wifi"

  qos_rate_max_down = 2000 # 2mbps
  qos_rate_max_up   = 10   # 10kbps
}
```

## Schema

- **id** - (String, Optional)
- **name** - (String, Required)
- **qos_rate_max_down** - (Number, Optional)
- **qos_rate_max_up** - (Number, Optional)


