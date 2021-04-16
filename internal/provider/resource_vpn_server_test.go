package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPNServer_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNServerConfig,
			},
			importStep("unifi_vpn_server.test"),
		},
	})
}

const testAccVPNServerConfig = `
data "unifi_radius_profile" "default" {}

resource "unifi_radius_server" "default" {
  enabled = true
  secret  = "secret"

  authentication_port         = 1812
  accounting_port             = 1813
  accounting_interim_interval = 3600

  enable_tunneled_reply   = true
  configure_whole_network = true
}

resource "unifi_vpn_server" "test" {
  name     = "Test VPN Server"
  vpn_type = "l2tp-server"

  subnet = "172.16.1.1/24"

  dhcp_start = "172.16.1.1"
  dhcp_stop  = "172.16.1.254"
  dhcp_dns   = ["1.1.1.1", "1.0.0.1"]

  radius_profile_id = data.unifi_radius_profile.default.id
  require_mschapv2  = true
  pre_shared_key    = "sharedkey"

	depends_on = [
		unifi_radius_server.default,
	]
}
`
