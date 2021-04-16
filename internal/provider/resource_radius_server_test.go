package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRadiusServer_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRadiusServerConfig,
			},
			importStep("unifi_radius_server.test"),
		},
	})
}

const testAccRadiusServerConfig = `
resource "unifi_radius_server" "test" {
  enabled = false
  secret  = "secret"

  authentication_port         = 1812
  accounting_port             = 1813
  accounting_interim_interval = 3600

  enable_tunneled_reply   = true
  configure_whole_network = true
}
`
