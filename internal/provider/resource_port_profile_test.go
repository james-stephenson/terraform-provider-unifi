package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPortProfile_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortProfileConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_profile.test", "poe_mode", "off"),
				),
			},
			importStep("unifi_port_profile.test"),
		},
	})
}

func TestAccPortProfile_speed(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortProfileConfig_speedAuto,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_profile.test", "link_speed", "auto"),
				),
			},
			importStep("unifi_port_profile.test"),
		},
	})
}

const testAccPortProfileConfig = `
resource "unifi_port_profile" "test" {
	name = "provider created"

	poe_mode	  = "off"
	link_speed 		  = 1000
	stp_enabled = false
}
`

const testAccPortProfileConfig_speedAuto = `
resource "unifi_port_profile" "test" {
	name = "provider created"

	poe_mode	  = "off"
	link_speed 		  = "auto"
	stp_enabled = false
	stormctrl_bcast_level = 20
}
`
