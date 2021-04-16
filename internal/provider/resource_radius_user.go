package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/paultyng/go-unifi/unifi"
)

func resourceRadiusUser() *schema.Resource {
	return &schema.Resource{
		Description: `
unifi_radius_user manages user accounts for the built-in RADIUS server.
`,
		Create: resourceRadiusUserCreate,
		Read:   resourceRadiusUserRead,
		Update: resourceRadiusUserUpdate,
		Delete: resourceRadiusUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Description: "Name of the account to create",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "Password to assign to the user. Users will need this paassword and the pre-shared key when connecting to this VPN.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"vlan_id": {
				Description: "When VLAN ID is set, this user will be assigned to that VLAN group and subject to any associated restrictions",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
			"tunnel_type": {
				Description: "The type of tunnel this user will use. Please refer to the Unifi controller documentation to see a mapping between number and type.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
			"tunnel_medium_type": {
				Description: "The type of tunnel medium this user will use. Please refer to the Unifi controller documentation to see a mapping between number and type.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
			},
		},
	}
}

func resourceRadiusUserCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceRadiusUserGetResourceData(d)
	if err != nil {
		return err
	}

	resp, err := c.c.CreateAccount(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceRadiusUserSetResourceData(resp, d)
}

func resourceRadiusUserGetResourceData(d *schema.ResourceData) (*unifi.Account, error) {
	req := &unifi.Account{
		Name:             d.Get("username").(string),
		XPassword:        d.Get("password").(string),
		TunnelType:       d.Get("tunnel_type").(int),
		TunnelMediumType: d.Get("tunnel_medium_type").(int),
		VLAN:             d.Get("vlan_id").(int),
	}

	return req, nil
}

func resourceRadiusUserSetResourceData(resp *unifi.Account, d *schema.ResourceData) error {
	d.Set("username", resp.Name)
	d.Set("password", resp.XPassword)
	d.Set("tunnel_type", resp.TunnelType)
	d.Set("tunnel_medium_type", resp.TunnelMediumType)
	d.Set("vlan_id", resp.VLAN)

	return nil
}

func resourceRadiusUserRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	resp, err := c.c.GetAccount(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceRadiusUserSetResourceData(resp, d)
}

func resourceRadiusUserUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	req, err := resourceRadiusUserGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = c.site

	resp, err := c.c.UpdateAccount(context.TODO(), c.site, req)
	if err != nil {
		return err
	}

	return resourceRadiusUserSetResourceData(resp, d)
}

func resourceRadiusUserDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	err := c.c.DeleteAccount(context.TODO(), c.site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		return nil
	}
	return err
}
