---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_port_forward"
description: |-
  unifi_port_forward manages a port forwarding rule on the gateway.
---

# Resource: `unifi_port_forward`

unifi_port_forward manages a port forwarding rule on the gateway.



## Schema

### Optional

- **dst_port** (String, Optional)
- **enabled** (Boolean, Optional)
- **fwd_ip** (String, Optional)
- **fwd_port** (String, Optional)
- **id** (String, Optional)
- **log** (Boolean, Optional) Defaults to `false`.
- **name** (String, Optional)
- **port_forward_interface** (String, Optional)
- **protocol** (String, Optional) Defaults to `tcp_udp`.
- **src_ip** (String, Optional) Defaults to `any`.


