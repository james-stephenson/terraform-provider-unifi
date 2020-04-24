---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_firewall_rule"
description: |-
  unifi_firewall_rule manages an individual firewall rule on the gateway.
---

# Resource: `unifi_firewall_rule`

unifi_firewall_rule manages an individual firewall rule on the gateway.



## Schema

### Required

- **action** (String, Required)
- **name** (String, Required)
- **protocol** (String, Required)
- **rule_index** (Number, Required)
- **ruleset** (String, Required)

### Optional

- **dst_address** (String, Optional)
- **dst_firewall_group_ids** (Set of String, Optional)
- **dst_network_id** (String, Optional)
- **dst_network_type** (String, Optional) Defaults to `NETv4`.
- **id** (String, Optional)
- **ip_sec** (String, Optional)
- **logging** (Boolean, Optional)
- **src_address** (String, Optional)
- **src_firewall_group_ids** (Set of String, Optional)
- **src_mac** (String, Optional)
- **src_network_id** (String, Optional)
- **src_network_type** (String, Optional) Defaults to `NETv4`.
- **state_established** (Boolean, Optional)
- **state_invalid** (Boolean, Optional)
- **state_new** (Boolean, Optional)
- **state_related** (Boolean, Optional)


