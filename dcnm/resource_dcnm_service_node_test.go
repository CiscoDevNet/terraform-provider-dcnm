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

var providerServiceNode *schema.Provider

func TestAccDCNMServiceNode_Basic(t *testing.T) {
	var serviceNode models.ServiceNode

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServiceNode),
		CheckDestroy:      testAccCheckDCNMServiceNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMServiceNodeConfig_basic("service node decription check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists("dcnm_service_node.test", &serviceNode),
					testAccCheckDCNMServiceNodeAttributes("service node decription check", &serviceNode),
				),
			},
		},
	})
}

func TestAccDCNMServiceNode_Update(t *testing.T) {
	var serviceNode models.ServiceNode

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerServiceNode),
		CheckDestroy:      testAccCheckDCNMServiceNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMServiceNodeConfig_basic("serviceNode decription check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists("dcnm_serviceNode.test", &serviceNode),
					testAccCheckDCNMServiceNodeAttributes("serviceNode decription check", &serviceNode),
				),
			},

			{
				Config: testAccCheckDCNMServiceNodeConfig_basic("serviceNode update check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMServiceNodeExists("dcnm_serviceNode.test", &serviceNode),
					testAccCheckDCNMServiceNodeAttributes("serviceNode update check", &serviceNode),
				),
			},
		},
	})
}

func testAccCheckDCNMServiceNodeConfig_basic(desc, deploy string) string {
	return fmt.Sprintf(`
	resource "dcnm_serviceNode" "test" {
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

func testAccCheckDCNMServiceNodeExists(name string, serviceNode *models.ServiceNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Network %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Network dn was set")
		}

		dcnmClient := (*providerNetwork).Meta().(*client.Client)

		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/top-down/fabrics/%s/serviceNodes/%s", "fab2", rs.Primary.ID))
		if err != nil {
			return err
		}

		serviceNodeGet := &models.Network{}
		serviceNodeGet.Fabric = stripQuotes(cont.S("fabric").String())
		serviceNodeGet.Name = stripQuotes(cont.S("serviceNodeName").String())
		serviceNodeGet.VRF = stripQuotes(cont.S("vrf").String())

		netProfile := &models.NetworkProfileConfig{}
		configCont, err := cleanJsonString(stripQuotes(cont.S("networkTemplateConfig").String()))
		if err != nil {
			return err
		}
		if configCont.Exists("vlanId") && stripQuotes(configCont.S("vlanId").String()) != "" {
			if vlan, err := strconv.Atoi(stripQuotes(configCont.S("vlanId").String())); err == nil {
				netProfile.Vlan = vlan
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

func testAccCheckDCNMServiceNodeDestroy(s *terraform.State) error {
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

func testAccCheckDCNMServiceNodeAttributes(desc string, serviceNode *models.ServiceNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "import" != serviceNode.Name {
			return fmt.Errorf("Bad Network Name %s", serviceNode.Name)
		}

		if "fab2" != serviceNode.AttachedFabricName {
			return fmt.Errorf("Bad Network fabric name %s", serviceNode.AttachedFabricName)
		}

		if 2301 != profile.Vlan {
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
