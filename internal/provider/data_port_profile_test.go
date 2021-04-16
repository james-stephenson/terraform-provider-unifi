package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataPortProfile_default(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPortProfileConfig_default,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.unifi_port_profile.default", "name", "All"),
					resource.TestCheckResourceAttr("data.unifi_port_profile.default", "forward", "all"),
				),
			},
		},
	})
}

func TestAccDataPortProfile_disabled(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPortProfileConfig_disabled,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.unifi_port_profile.disabled", "name", "Disabled"),
					resource.TestCheckResourceAttr("data.unifi_port_profile.disabled", "forward", "disabled"),
				),
			},
		},
	})
}

func TestAccDataPortProfile_lan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPortProfileConfig_lan,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.unifi_port_profile.lan", "name", "LAN"),
					resource.TestCheckResourceAttr("data.unifi_port_profile.lan", "forward", "native"),
				),
			},
		},
	})
}

func TestAccDataPortProfile_multiple_providers(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"unifi2": func() (*schema.Provider, error) {
				return New("acctest")(), nil
			},
			"unifi3": func() (*schema.Provider, error) {
				return New("acctest")(), nil
			},
		},
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: `
				data "unifi_port_profile" "unifi2" {
					provider = "unifi2"
				}
				data "unifi_port_profile" "unifi3" {
					provider = "unifi3"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

const testAccDataPortProfileConfig_default = `
data "unifi_port_profile" "default" {
}
`

const testAccDataPortProfileConfig_disabled = `
data "unifi_port_profile" "disabled" {
	name = "Disabled"
}
`

const testAccDataPortProfileConfig_lan = `
data "unifi_port_profile" "lan" {
	name = "LAN"
}
`
