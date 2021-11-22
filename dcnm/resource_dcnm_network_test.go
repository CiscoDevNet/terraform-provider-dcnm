package dcnm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerNetwork *schema.Provider

func TestAccDCNMNetwork_Basic(t *testing.T) {
	var network models.Network
	var networkProfile models.NetworkProfileConfig

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerNetwork),
		CheckDestroy:      testAccCheckDCNMNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMNetworkConfig_basic("network decription check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMNetworkExists("dcnm_network.test", &network, &networkProfile),
					testAccCheckDCNMNetworkAttributes("network decription check", &network, &networkProfile),
				),
			},
		},
	})
}

func TestAccDCNMNetwork_Update(t *testing.T) {
	var network models.Network
	var networkProfile models.NetworkProfileConfig

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerNetwork),
		CheckDestroy:      testAccCheckDCNMNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMNetworkConfig_basic("network decription check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMNetworkExists("dcnm_network.test", &network, &networkProfile),
					testAccCheckDCNMNetworkAttributes("network decription check", &network, &networkProfile),
				),
			},

			{
				Config: testAccCheckDCNMNetworkConfig_basic("network update check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMNetworkExists("dcnm_network.test", &network, &networkProfile),
					testAccCheckDCNMNetworkAttributes("network update check", &network, &networkProfile),
				),
			},
		},
	})
}

func testAccCheckDCNMNetworkConfig_basic(desc, deploy string) string {
	return fmt.Sprintf(`
	resource "dcnm_network" "test" {
		fabric_name     = "fab2"
		name            = "import"
		display_name    = "check"
		description     = "%s"
		vrf_name        = "Test-vrf"
		vlan_id         = 2301
		vlan_name       = "vlan1"
		deploy = %s
		attachments {
			serial_number = "9EQ00OGQYV6"
			vlan_id       = 2400
			attach        = %s
		}
	}
	`, desc, deploy, deploy)
}

func testAccCheckDCNMNetworkExists(name string, network *models.Network, profile *models.NetworkProfileConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Network %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network dn was set")
		}

		dcnmClient := (*providerNetwork).Meta().(*client.Client)

		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s", "fab2", rs.Primary.ID))
		if err != nil {
			return err
		}

		networkGet := &models.Network{}
		networkGet.Fabric = stripQuotes(cont.S("fabric").String())
		networkGet.Name = stripQuotes(cont.S("networkName").String())
		networkGet.VRF = stripQuotes(cont.S("vrf").String())

		netProfile := &models.NetworkProfileConfig{}
		configCont, err := cleanJsonString(stripQuotes(cont.S("networkTemplateConfig").String()))
		if err != nil {
			return err
		}
		if configCont.Exists("vlanId") && stripQuotes(configCont.S("vlanId").String()) != "" {
			if vlan, err := strconv.Atoi(stripQuotes(configCont.S("vlanId").String())); err == nil {
				netProfile.Vlan = strconv.Itoa(vlan)
			}
		}
		if configCont.Exists("vlanName") {
			netProfile.VlanName = stripQuotes(configCont.S("vlanName").String())
		}
		if configCont.Exists("intfDescription") {
			netProfile.Description = stripQuotes(configCont.S("intfDescription").String())
		}

		*network = *networkGet
		*profile = *netProfile
		return nil
	}
}

func testAccCheckDCNMNetworkDestroy(s *terraform.State) error {
	dcnmClient := (*providerNetwork).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_network" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s", "fab2", rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Network still exists")
			}
		}
	}

	return nil
}

func testAccCheckDCNMNetworkAttributes(desc string, network *models.Network, profile *models.NetworkProfileConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "import" != network.Name {
			return fmt.Errorf("Bad Network Name %s", network.Name)
		}

		if "fab2" != network.Fabric {
			return fmt.Errorf("Bad Network fabric name %s", network.Fabric)
		}

		if "MyVRF" != network.VRF {
			return fmt.Errorf("Bad Network VRF name %s", network.VRF)
		}

		if strconv.Itoa(2301) != (profile.Vlan) {
			return fmt.Errorf("Bad Network VLAN %d", profile.Vlan)
		}

		if "vlan1" != profile.VlanName {
			return fmt.Errorf("Bad Network VLAN name %s", profile.VlanName)
		}

		if desc != profile.Description {
			return fmt.Errorf("Bad Network description %s", profile.Description)
		}
		return nil
	}
}
