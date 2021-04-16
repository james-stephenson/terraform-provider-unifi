package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

var setting_key string = "radius"

func resourceRadiusServer() *schema.Resource {
	return &schema.Resource{
		Description: `
unifi_radius_server manages the RADIUS server on the gateway.
`,
		Create: resourceRadiusServerCreate,
		Read:   resourceRadiusServerRead,
		Update: resourceRadiusServerUpdate,
		Delete: resourceRadiusServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Whether to enable the RADIUS server on the gateway",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"secret": {
				Description: "The RADIUS secret is used between the RADIUS server and RADIUS clients such as switches and access points",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"authentication_port": {
				Description:  "The port for authentication communications",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1812",
				ValidateFunc: validatePortRange,
			},
			"accounting_port": {
				Description:  "The port for accounting communications",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "1813",
				ValidateFunc: validatePortRange,
			},
			"accounting_interim_interval": {
				Description:  "Statistics will be collected from connected clients at this interval",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validation.IntBetween(60, 86400),
			},
			"enable_tunneled_reply": {
				Description: "Encrypt communication between the server and the client",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"configure_whole_network": {
				Description: "Create a full tunnel",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceRadiusServerCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceRadiusServerGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateSettingRadius(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)
	return resourceRadiusServerSetResourceData(resp, d)
}

func resourceRadiusServerGetResourceData(d *schema.ResourceData) (*unifi.SettingRadius, error) {
	authport, _ := strconv.Atoi(d.Get("authentication_port").(string))
	acctport, _ := strconv.Atoi(d.Get("accounting_port").(string))

	return &unifi.SettingRadius{
		Enabled:               d.Get("enabled").(bool),
		XSecret:               d.Get("secret").(string),
		AuthPort:              authport,
		AcctPort:              acctport,
		InterimUpdateInterval: d.Get("accounting_interim_interval").(int),
		TunneledReply:         d.Get("enable_tunneled_reply").(bool),
		ConfigureWholeNetwork: d.Get("configure_whole_network").(bool),
	}, nil
}

func resourceRadiusServerSetResourceData(resp *unifi.SettingRadius, d *schema.ResourceData) error {
	d.Set("enabled", resp.Enabled)
	d.Set("secret", resp.XSecret)
	d.Set("authentication_port", strconv.Itoa(resp.AuthPort))
	d.Set("accounting_port", strconv.Itoa(resp.AcctPort))
	d.Set("accounting_interim_interval", resp.InterimUpdateInterval)
	d.Set("enable_tunneled_reply", resp.TunneledReply)
	d.Set("configure_whole_network", resp.ConfigureWholeNetwork)

	return nil
}

func resourceRadiusServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	resp, err := c.c.GetSettingRadius(context.TODO(), c.site)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceRadiusServerSetResourceData(resp, d)
}

func resourceRadiusServerUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceRadiusServerGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateSettingRadius(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourceRadiusServerSetResourceData(resp, d)
}

func resourceRadiusServerDelete(d *schema.ResourceData, meta interface{}) error {
	d.Set("enabled", false)

	return resourceRadiusServerUpdate(d, meta)
}
