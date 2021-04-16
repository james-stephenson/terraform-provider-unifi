package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceVPNServer() *schema.Resource {
	return &schema.Resource{
		Description: `
unifi_network manages LAN/VLAN networks.
`,

		Create: resourceVPNServerCreate,
		Read:   resourceVPNServerRead,
		Update: resourceVPNServerUpdate,
		Delete: resourceVPNServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the VPN.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"vpn_type": {
				Description:  "Type of the VPN. Must be one of: `l2tp-server`, `pptp-server`",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"l2tp-server", "pptp-server"}, false),
			},
			"subnet": {
				Description:      "The subnet of the network. Must be a valid CIDR address.",
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: cidrDiffSuppress,
			},
			"dhcp_start": {
				Description:  "The IPv4 address where the DHCP range of addresses starts.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp_stop": {
				Description:  "The IPv4 address where the DHCP range of addresses stops.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp_dns": {
				Description: "Specifies the IPv4 addresses for the DNS server to be returned from the DHCP " +
					"server. Leave blank to disable this feature.",
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 4,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.IsIPv4Address,
						// this doesn't let blank through
						validation.StringLenBetween(1, 50),
					),
				},
			},
			"radius_profile_id": {
				Description: "RADIUS profile ID for this VPN server to use for authentication",
				Type:        schema.TypeString,
				Required:    true,
			},
			"require_mschapv2": {
				Description: "Require MS-CHAP v2",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"pre_shared_key": {
				Description: "A secret key that will be used when connecting to your VPN",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceVPNServerCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceVPNServerGetResourceData(d)
	if err != nil {
		return err
	}

	resp, err := c.c.CreateNetwork(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceVPNServerSetResourceData(resp, d)
}

func resourceVPNServerGetResourceData(d *schema.ResourceData) (*unifi.Network, error) {
	dhcpDNS, err := listToStringSlice(d.Get("dhcp_dns").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to convert dhcp_dns to string slice: %w", err)
	}

	return &unifi.Network{
		Name:               d.Get("name").(string),
		Enabled:            true,
		Purpose:            "remote-user-vpn",
		VPNType:            d.Get("vpn_type").(string),
		IPSubnet:           cidrOneBased(d.Get("subnet").(string)),
		DHCPDStart:         d.Get("dhcp_start").(string),
		DHCPDStop:          d.Get("dhcp_stop").(string),
		RADIUSProfileID:    d.Get("radius_profile_id").(string),
		RequireMschapv2:    d.Get("require_mschapv2").(bool),
		XIPSecPreSharedKey: d.Get("pre_shared_key").(string),
		IsNAT:              true,
		VLANEnabled:        false,

		DHCPDDNSEnabled: len(dhcpDNS) > 0,
		// this is kinda hacky but ¯\_(ツ)_/¯
		DHCPDDNS1: append(dhcpDNS, "")[0],
		DHCPDDNS2: append(dhcpDNS, "", "")[1],
		DHCPDDNS3: append(dhcpDNS, "", "", "")[2],
		DHCPDDNS4: append(dhcpDNS, "", "", "", "")[3],

		IPV6InterfaceType: "none",
		// IPV6InterfaceType string `json:"ipv6_interface_type"` // "none"
		// IPV6PDStart       string `json:"ipv6_pd_start"`       // "::2"
		// IPV6PDStop        string `json:"ipv6_pd_stop"`        // "::7d1"
	}, nil
}

func resourceVPNServerSetResourceData(resp *unifi.Network, d *schema.ResourceData) error {
	dhcpLease := resp.DHCPDLeaseTime
	if resp.DHCPDEnabled && dhcpLease == 0 {
		dhcpLease = 86400
	}

	dhcpDNS := []string{}
	if resp.DHCPDDNSEnabled {
		for _, dns := range []string{
			resp.DHCPDDNS1,
			resp.DHCPDDNS2,
			resp.DHCPDDNS3,
			resp.DHCPDDNS4,
		} {
			if dns == "" {
				continue
			}
			dhcpDNS = append(dhcpDNS, dns)
		}
	}

	d.Set("name", resp.Name)
	d.Set("vpn_type", resp.VPNType)
	d.Set("subnet", cidrZeroBased(resp.IPSubnet))
	d.Set("dhcp_start", resp.DHCPDStart)
	d.Set("dhcp_stop", resp.DHCPDStop)
	d.Set("dhcp_dns", dhcpDNS)
	d.Set("radius_profile_id", resp.RADIUSProfileID)
	d.Set("require_mschapv2", resp.RequireMschapv2)
	d.Set("pre_shared_key", resp.XIPSecPreSharedKey)

	return nil
}

func resourceVPNServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetNetwork(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceVPNServerSetResourceData(resp, d)
}

func resourceVPNServerUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceVPNServerGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateNetwork(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourceVPNServerSetResourceData(resp, d)
}

func resourceVPNServerDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	name := d.Get("name").(string)
	id := d.Id()

	err := c.c.DeleteNetwork(context.TODO(), c.site, id, name)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
