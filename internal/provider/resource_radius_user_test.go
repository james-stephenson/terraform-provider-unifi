package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRadiusUser_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRadiusUserConfig,
			},
			importStep("unifi_radius_user.test"),
		},
	})
}

const testAccRadiusUserConfig = `
resource "unifi_radius_user" "test" {
  username = "test-user"
  password = "password"

	vlan_id						 = 10
  tunnel_type        = 3
  tunnel_medium_type = 1
}
`
