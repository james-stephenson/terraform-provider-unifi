package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/paultyng/go-unifi/unifi"
)

func dataDevice() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_device` data source retrieves Unifi device information using the given MAC address.",
		Read:        resourceDeviceRead,
		Schema:      deviceSchema(),
	}
}

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Description: "`unifi_device` manages a device of the network.\n\n" +
			"Devices are adopted by the controller, so it is not possible for this resource to be created through " +
			"Terraform, the create operation instead will simply start managing the device specified by MAC address. " +
			"It's safer to start this process with an explicit import of the device.",

		Create:        resourceDeviceCreate,
		Read:          resourceDeviceRead,
		Update:        resourceDeviceUpdate,
		DeleteContext: resourceDeviceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDeviceImport,
		},

		Schema: deviceSchema(),
	}
}

func deviceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the device.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"site": {
			Description: "The name of the site to associate the device with.",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			ForceNew:    true,
		},
		"mac": {
			Description:      "The MAC address of the device. This can be specified so that the provider can take control of a device since devices are created through adoption.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			DiffSuppressFunc: macDiffSuppressFunc,
			ValidateFunc:     validation.StringMatch(macAddressRegexp, "Mac address is invalid"),
		},

		// General settings
		"name": {
			Description: "The name of the device.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		"disabled": {
			Description: "Specifies whether this device should be disabled.",
			Type:        schema.TypeBool,
			Computed:    true,
		},

		"led_override": {
			Description:  "Mode for the device's LED. Must be one of: `default`, `on`, or `off`",
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "default",
			ValidateFunc: validation.StringInSlice([]string{"default", "on", "off"}, false),
		},

		// Network config
		"mgmt_network_id": {
			Description: "The ID for the management network that this device belongs to.",
			Type:        schema.TypeString,
			Optional:    true,
		},

		"network_config": {
			Description: "Network configuration.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        resourceDeviceConfigNetworkResource(),
			MaxItems:    1,
		},

		//-- Network Switches
		"jumboframe_enabled": {
			Description: "Whether to enable jumbo frames.",
			Type:        schema.TypeBool,
			Optional:    true,
		},
		"flowctrl_enabled": {
			Description: "Whether to enable flow control.",
			Type:        schema.TypeBool,
			Optional:    true,
		},

		"stp_version": {
			Description:  "Spanning tree protocol (network switches only). Must be one of: `stp`, `rstp`, `disabled`",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"stp", "rstp", "disabled"}, false),
		},

		"stp_priority": {
			Description:  "Spanning tree protocol priority (network switches only). Must be a multiple of 4096.",
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.Any(validation.IntDivisibleBy(4096), validation.IntAtMost(61444)),
		},

		"ethernet_override": {
			Description: "Defines overrides for ethernet devices.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        resourceDeviceEthernetOverride(),
		},

		"port_override": {
			// TODO: this should really be a map or something when possible in the SDK
			// see https://github.com/hashicorp/terraform-plugin-sdk/issues/62
			Description: "Settings overrides for specific switch ports.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        resourceDevicePortOverride(),
		},

		//-- Access Points
		"radio_2g": {
			Description: "Configure the 2.4GHz radio for this access point.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        resourceDeviceRadioTable(),
			MaxItems:    1,
		},
		"radio_5g": {
			Description: "Configure the 5GHz radio for this access point.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        resourceDeviceRadioTable(),
			MaxItems:    1,
		},
	}
}

func resourceDeviceConfigNetworkResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Description: "Fixed IP for this device.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dns1": {
				Description: "First DNS server to use for this device.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dns2": {
				Description: "Second DNS server to use for this device.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dns_suffix": {
				Description: "DNS search domain for this device.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"gateway": {
				Description: "IP of the gateway for this device.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"netmask": {
				Description: "Subnet netmask for this device's network configuration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  "Type of network configuration to use. Must be one of: `dhcp`, `static`",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"dhcp", "static"}, false),
			},
			"bonding_enabled": {
				Description: "Whether bonding is enabled for this device.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceDeviceEthernetOverride() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"interface": {
				Description: "Ethernet interface",
				Type:        schema.TypeString,
				Required:    true,
			},
			"network_group": {
				Description: "Network group for the interface",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceDevicePortOverride() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"number": {
				Description: "Switch port number. Must be a valid port number on this device.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description: "Name of this port.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"port_profile_id": {
				Description: "ID of the Port Profile used on this port.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"poe_mode": {
				Description:  "Override the PoE mode of this port's profile. Must be one of: `auto`, `pasv24`, `passthrough`, `off`",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"auto", "pasv24", "passthrough", "off"}, false),
			},
			"op_mode": {
				Description:  "Override operation mode of this port. Must be one of: `switch`, `mirror`, `aggregate`. Defaults to `switch`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"switch", "mirror", "aggregate"}, false),
			},
			"link_speed": {
				Description:  "Sets the link speed for ports which support it",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(validation.IntInSlice([]int{10, 100}), validation.StringInSlice([]string{"auto"}, false)),
			},
			"duplex": {
				Description:  "Sets the duplexing for ports which support it, e.g. ports on the Unifi gateway products",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"full", "half"}, false),
			},
		},
	}
}

func resourceDeviceRadioTable() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"channel": {
				Description: "Channel for this radio. Specify a numeric channel appropriate for the radio type, or `auto` to auto select.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"channel_width": {
				Description:  "Set the channel width for this radio. Must be onee  of: `20`, `40`, `80`, `160`, `1080`, `2160`",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{20, 40, 80, 160, 1080, 2160}),
			},
			"transmit_power": {
				Description: "Set the transmission power and mode for this radio. Can either be an integer value between 6 and 26, or a string from this list: `low`, `medium`, `high`, or `auto`",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.Any(
					validation.StringInSlice([]string{"low", "medium", "high", "auto"}, false),
					func(i interface{}, k string) (warnings []string, errors []error) {
						val, err := strconv.Atoi(i.(string))
						if err != nil {
							return warnings, append(errors, err)
						}

						f := validation.IntBetween(6, 26)
						return f(val, k)
					},
				),
			},
		},
	}
}

func resourceDeviceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*client)
	id := d.Id()
	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	if colons := strings.Count(id, ":"); colons == 1 || colons == 6 {
		importParts := strings.SplitN(id, ":", 2)
		site = importParts[0]
		id = importParts[1]
	}

	if macAddressRegexp.MatchString(id) {
		// look up id by mac
		find := cleanMAC(id)

		devices, err := c.c.ListDevice(ctx, site)
		if err != nil {
			return nil, err
		}
		for _, d := range devices {
			if cleanMAC(d.MAC) == find {
				id = d.ID
				break
			}
		}
	}

	if id != "" {
		d.SetId(id)
	}
	if site != "" {
		d.Set("site", site)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	mac := d.Get("mac").(string)
	if mac == "" {
		return fmt.Errorf("no MAC address specified, please import the device using terraform import")
	}

	mac = cleanMAC(mac)
	devices, err := c.c.ListDevice(context.TODO(), site)
	if err != nil {
		return fmt.Errorf("unable to list devices: %w", err)
	}

	var found *unifi.Device
	for _, dev := range devices {
		if cleanMAC(dev.MAC) == mac {
			found = &dev
			break
		}
	}
	if found == nil {
		return fmt.Errorf("device not found using mac %q", mac)
	}

	d.SetId(found.ID)

	return resourceDeviceSetResourceData(found, d, site)
}

func resourceDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	req, err := resourceDeviceGetResourceData(d)
	if err != nil {
		return err
	}

	req.ID = d.Id()
	req.SiteID = site

	resp, err := c.c.UpdateDevice(context.TODO(), site, req)
	if err != nil {
		switch err.(type) {
		case *unifi.NotFoundError:
			return resourceDeviceRead(d, meta)
		default:
			return err
		}
	} else {
		return resourceDeviceSetResourceData(resp, d, site)
	}
}

func resourceDeviceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting a device via Terraform is not supported, the device will just be removed from state.",
		},
	}
}

func resourceDeviceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client)

	id := d.Id()

	site := d.Get("site").(string)
	if site == "" {
		site = c.site
	}

	resp, err := c.c.GetDevice(context.TODO(), site, id)
	if _, ok := err.(*unifi.NotFoundError); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	return resourceDeviceSetResourceData(resp, d, site)
}

func resourceDeviceSetResourceData(resp *unifi.Device, d *schema.ResourceData, site string) error {
	d.Set("site", site)
	d.Set("mac", resp.MAC)
	d.Set("name", resp.Name)
	d.Set("disabled", resp.Disabled)
	d.Set("led_override", resp.LedOverride)
	d.Set("mgmt_network_id", resp.MgmtNetworkID)
	d.Set("jumboframe_enabled", resp.JumboframeEnabled)
	d.Set("flowctrl_enabled", resp.FlowctrlEnabled)

	if resp.StpVersion != "" && resp.StpVersion != "disabled" {
		d.Set("stp_version", resp.StpVersion)
		val, err := strconv.Atoi(resp.StpPriority)
		if err != nil {
			d.Set("err", fmt.Sprintf("%+v", err))
			return err
		}

		d.Set("stp_priority", val)
	}

	portOverrides, err := setFromPortOverrides(resp.PortOverrides)
	if err != nil {
		return err
	}
	d.Set("port_override", portOverrides)

	ethernetOverrides, err := setFromEthernetOverrides(resp.EthernetOverrides)
	if err != nil {
		return err
	}
	d.Set("ethernet_override", ethernetOverrides)

	cfg, err := setFromNetworkConfig(resp.ConfigNetwork)
	if err != nil {
		return err
	}
	d.Set("network_config", cfg)

	var network_config []interface{}
	networkConfig, err := fromNetworkConfig(resp.ConfigNetwork)
	if err != nil {
		return err
	}
	d.Set("network_config", append(network_config, networkConfig))

	for _, i := range resp.RadioTable {
		radioTable, err := setFromRadioTable(i)
		if err != nil {
			return err
		}

		switch i.Radio {
		case "ng":
			d.Set("radio_2g", radioTable)
		case "na":
			d.Set("radio_5g", radioTable)
		}
	}

	return nil
}

func resourceDeviceGetResourceData(d *schema.ResourceData) (*unifi.Device, error) {
	//TODO: pass Disabled once we figure out how to enable the device afterwards
	device := &unifi.Device{
		MAC:           d.Get("mac").(string),
		Name:          d.Get("name").(string),
		LedOverride:   d.Get("led_override").(string),
		MgmtNetworkID: d.Get("mgmt_network_id").(string),

		// Network Switches
		JumboframeEnabled: d.Get("jumboframe_enabled").(bool),
		FlowctrlEnabled:   d.Get("flowctrl_enabled").(bool),
	}

	if stp_version, ok := d.GetOk("stp_version"); ok {
		device.StpVersion = stp_version.(string)
	}

	if stp_priority, ok := d.GetOk("stp_priority"); ok {
		device.StpPriority = strconv.Itoa(stp_priority.(int))
	}

	pos, err := setToPortOverrides(d.Get("port_override").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to process port_override block: %w", err)
	}
	device.PortOverrides = pos

	eos, err := setToEthernetOverrides(d.Get("ethernet_override").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to process ethernet_override block: %w", err)
	}
	device.EthernetOverrides = eos

	cfg, err := setToNetworkConfig(d.Get("network_config").([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("unable to process network_config block: %w", err)
	}
	device.ConfigNetwork = cfg

	if radio_2g, ok := d.GetOk("radio_2g"); ok {
		rt, err := toRadioTable(radio_2g.([]interface{})[0].(map[string]interface{}), "ra0", "ng")
		if err != nil {
			return nil, fmt.Errorf("unable to process radio_2g block: %w", err)
		}

		device.RadioTable = append(device.RadioTable, rt)
	}

	if radio_5g, ok := d.GetOk("radio_5g"); ok {
		rt, err := toRadioTable(radio_5g.([]interface{})[0].(map[string]interface{}), "rai0", "na")
		if err != nil {
			return nil, fmt.Errorf("unable to process radio_5g block: %w", err)
		}

		device.RadioTable = append(device.RadioTable, rt)
	}

	return device, nil
}

// Network configuration data type
func setToNetworkConfig(data []interface{}) (unifi.DeviceConfigNetwork, error) {
	if data != nil && len(data) > 0 {
		return toNetworkConfig(data[0].(map[string]interface{}))
	}

	return unifi.DeviceConfigNetwork{}, nil
}

func toNetworkConfig(cfg map[string]interface{}) (unifi.DeviceConfigNetwork, error) {
	return unifi.DeviceConfigNetwork{
		IP:             cfg["ip"].(string),
		DNS1:           cfg["dns1"].(string),
		DNS2:           cfg["dns2"].(string),
		DNSsuffix:      cfg["dns_suffix"].(string),
		Gateway:        cfg["gateway"].(string),
		Netmask:        cfg["netmask"].(string),
		Type:           cfg["type"].(string),
		BondingEnabled: cfg["bonding_enabled"].(bool),
	}, nil
}

func setFromNetworkConfig(cfg unifi.DeviceConfigNetwork) ([]interface{}, error) {
	var items []interface{}
	networkConfig, err := fromNetworkConfig(cfg)
	if err != nil {
		return nil, err
	}

	return append(items, networkConfig), nil
}

func fromNetworkConfig(cfg unifi.DeviceConfigNetwork) (map[string]interface{}, error) {
	return map[string]interface{}{
		"ip":              cfg.IP,
		"dns1":            cfg.DNS1,
		"dns2":            cfg.DNS2,
		"dns_suffix":      cfg.DNSsuffix,
		"gateway":         cfg.Gateway,
		"netmask":         cfg.Netmask,
		"type":            cfg.Type,
		"bonding_enabled": cfg.BondingEnabled,
	}, nil
}

// Radio tables
func toRadioTable(cfg map[string]interface{}, name string, radio string) (unifi.DeviceRadioTable, error) {
	radio_table := unifi.DeviceRadioTable{
		Name:    name,
		Radio:   radio,
		Channel: cfg["channel"].(string),
		Ht:      strconv.Itoa(cfg["channel_width"].(int)),
	}

	power := cfg["transmit_power"].(string)
	if _, err := strconv.Atoi(power); err == nil {
		radio_table.TxPower = power
		radio_table.TxPowerMode = "custom"
	} else {
		radio_table.TxPowerMode = cfg["transmit_power"].(string)
	}

	return radio_table, nil
}

func setFromRadioTable(cfg unifi.DeviceRadioTable) ([]interface{}, error) {
	var items []interface{}
	rt, err := fromRadioTable(cfg)
	if err != nil {
		return nil, err
	}

	return append(items, rt), nil
}

func fromRadioTable(cfg unifi.DeviceRadioTable) (map[string]interface{}, error) {
	ht, _ := strconv.Atoi(cfg.Ht)
	radio := map[string]interface{}{
		"channel":       cfg.Channel,
		"channel_width": ht,
	}

	if cfg.TxPowerMode == "custom" {
		radio["transmit_power"] = cfg.TxPower
	} else {
		radio["transmit_power"] = cfg.TxPowerMode
	}

	return radio, nil
}

// Ethernet override data type
func setToEthernetOverrides(ethernet_overrides []interface{}) ([]unifi.DeviceEthernetOverrides, error) {
	// use a map here to remove any duplication
	overrideMap := map[string]unifi.DeviceEthernetOverrides{}
	for _, item := range ethernet_overrides {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}

		o, err := toEthernetOverride(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create ethernet override: %w", err)
		}

		overrideMap[o.Ifname] = o
	}

	result := make([]unifi.DeviceEthernetOverrides, 0, len(overrideMap))
	for _, item := range overrideMap {
		result = append(result, item)
	}

	return result, nil
}

func setFromEthernetOverrides(items []unifi.DeviceEthernetOverrides) ([]interface{}, error) {
	list := make([]interface{}, 0, len(items))
	for _, i := range items {
		v, err := fromEthernetOverride(i)
		if err != nil {
			return nil, fmt.Errorf("unable to parse ethernet override: %w", err)
		}
		list = append(list, v)
	}
	return list, nil
}

func toEthernetOverride(data map[string]interface{}) (unifi.DeviceEthernetOverrides, error) {
	ethernet_override := unifi.DeviceEthernetOverrides{
		Ifname:       data["interface"].(string),
		NetworkGroup: data["network_group"].(string),
	}

	return ethernet_override, nil
}

func fromEthernetOverride(input unifi.DeviceEthernetOverrides) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"interface":     input.Ifname,
		"network_group": input.NetworkGroup,
	}
	return data, nil
}

// Port override data type
func setToPortOverrides(port_overrides []interface{}) ([]unifi.DevicePortOverrides, error) {
	// use a map here to remove any duplication
	overrideMap := map[int]unifi.DevicePortOverrides{}
	for _, item := range port_overrides {
		data, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data in block")
		}

		po, err := toPortOverride(data)
		if err != nil {
			return nil, fmt.Errorf("unable to create port override: %w", err)
		}

		overrideMap[po.PortIDX] = po
	}

	pos := make([]unifi.DevicePortOverrides, 0, len(overrideMap))
	for _, item := range overrideMap {
		pos = append(pos, item)
	}

	return pos, nil
}

func setFromPortOverrides(pos []unifi.DevicePortOverrides) ([]interface{}, error) {
	list := make([]interface{}, 0, len(pos))
	for _, po := range pos {
		v, err := fromPortOverride(po)
		if err != nil {
			return nil, fmt.Errorf("unable to parse port override: %w", err)
		}
		list = append(list, v)
	}
	return list, nil
}

func toPortOverride(data map[string]interface{}) (unifi.DevicePortOverrides, error) {
	// TODO: error check these?
	port_override := unifi.DevicePortOverrides{
		PortIDX:       data["number"].(int),
		Name:          data["name"].(string),
		PortProfileID: data["port_profile_id"].(string),
		PoeMode:       data["poe_mode"].(string),
		OpMode:        data["op_mode"].(string),
	}

	if val, ok := data["duplex"]; ok {
		port_override.FullDuplex = (val.(string) == "full")
	}

	if val, ok := data["link_speed"]; ok {
		if val.(string) == "" {
		} else if val.(string) == "auto" {
			port_override.Autoneg = true
		} else {
			i, err := strconv.Atoi(val.(string))
			if err != nil {
				return port_override, fmt.Errorf("link_speed = %v", val)
			}

			port_override.Speed = i
		}
	}

	return port_override, nil
}

func fromPortOverride(po unifi.DevicePortOverrides) (map[string]interface{}, error) {
	port := map[string]interface{}{
		"number":          po.PortIDX,
		"name":            po.Name,
		"port_profile_id": po.PortProfileID,
		"poe_mode":        po.PoeMode,
		"op_mode":         po.OpMode,
	}

	if po.Autoneg {
		port["link_speed"] = "auto"
	} else if po.Speed != 0 {
		port["link_speed"] = strconv.Itoa(po.Speed)
		if po.FullDuplex {
			port["duplex"] = "full"
		} else {
			port["duplex"] = "half"
		}
	}

	return port, nil
}
