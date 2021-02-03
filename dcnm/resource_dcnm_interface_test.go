package dcnm

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerIntf *schema.Provider

func TestAccDCNMInterface_Basic(t *testing.T) {
	var intf models.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerIntf),
		CheckDestroy:      testAccCheckDCNMInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMInterfaceConfig_basic("creation from terraform"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMInterfaceExists("dcnm_interface.test", &intf),
					testAccCheckDCNMInterfaceAttributes("creation from terraform", &intf),
				),
			},
		},
	})
}

func TestAccDCNMInterface_Update(t *testing.T) {
	var intf models.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerIntf),
		CheckDestroy:      testAccCheckDCNMInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMInterfaceConfig_basic("creation from terraform"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMInterfaceExists("dcnm_interface.test", &intf),
					testAccCheckDCNMInterfaceAttributes("creation from terraform", &intf),
				),
			},
			{
				Config: testAccCheckDCNMInterfaceConfig_basic("updation from terraform"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMInterfaceExists("dcnm_interface.test", &intf),
					testAccCheckDCNMInterfaceAttributes("updation from terraform", &intf),
				),
			},
		},
	})
}

func testAccCheckDCNMInterfaceConfig_basic(desc string) string {
	return fmt.Sprintf(`
	resource "dcnm_interface" "test" {
		fabric_name = "fab2"
		name        = "loopback5"
		type        = "loopback"
		policy      = "int_loopback_11_1"
	  
		switch_name_1             = "leaf1"
		ipv4                      = "1.2.3.4"
		loopback_tag              = "1234"
		vrf                       = "MyVRF"
		loopback_ls_routing       = "ospf"
		loopback_replication_mode = "Multicast"
		description               = "%s"
		ipv6                      = "2001::0"

		deploy = false
	}
	`, desc)
}

func testAccCheckDCNMInterfaceExists(name string, inv *models.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface dn was set")
		}

		dcnmClient := (*providerIntf).Meta().(*client.Client)

		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/interface?ifName=%s", rs.Primary.ID))
		if err != nil {
			return err
		}
		contGet := cont.Index(0)

		intfGet := &models.Interface{}
		intfAttr := &models.InterfaceConfig{}

		intfGet.Policy = stripQuotes(contGet.S("policy").String())
		interfaces := contGet.S("interfaces").Index(0)
		intfAttr.SerialNumber = stripQuotes(interfaces.S("serialNumber").String())
		intfAttr.InterfaceName = stripQuotes(interfaces.S("nvPairs", "INTF_NAME").String())
		intfAttr.Fabric = stripQuotes(interfaces.S("nvPairs", "FABRIC_NAME").String())

		nvPairs := make(map[string]interface{})
		nvPairs["ipv4"] = stripQuotes(interfaces.S("nvPairs", "IP").String())
		nvPairs["ipv6"] = stripQuotes(interfaces.S("nvPairs", "V6IP").String())
		nvPairs["description"] = stripQuotes(interfaces.S("nvPairs", "DESC").String())
		nvPairs["vrf"] = stripQuotes(interfaces.S("nvPairs", "INTF_VRF").String())
		intfAttr.NVPairs = nvPairs

		intfGet.Interfaces = []models.InterfaceConfig{*intfAttr}

		*inv = *intfGet
		return nil
	}
}

func testAccCheckDCNMInterfaceDestroy(s *terraform.State) error {
	dcnmClient := (*providerIntf).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_interface" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/interface?ifName=%s", rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Interface still exists")
			}
		}
	}

	return nil
}

func testAccCheckDCNMInterfaceAttributes(desc string, intf *models.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		intfAtt := (intf.Interfaces)[0]
		nvPair := (intfAtt.NVPairs).(map[string]interface{})

		if "int_loopback_11_1" != intf.Policy {
			return fmt.Errorf("Bad interface policy %s", intf.Policy)
		}

		if "fab2" != intfAtt.Fabric {
			return fmt.Errorf("Bad interface fabric %s", intfAtt.Fabric)
		}

		if "loopback5" != intfAtt.InterfaceName {
			return fmt.Errorf("Bad interface name %s", intfAtt.InterfaceName)
		}

		if "1.2.3.4" != nvPair["ipv4"] {
			return fmt.Errorf("Bad interface ipv4 address %s", nvPair["ipv4"])
		}

		if "2001::0" != nvPair["ipv6"] {
			return fmt.Errorf("Bad interface ipv6 address %s", nvPair["ipv6"])
		}

		if desc != nvPair["description"] {
			return fmt.Errorf("Bad interface description %s", nvPair["description"])
		}

		if "MyVRF" != nvPair["vrf"] {
			return fmt.Errorf("Bad interface vrf %s", nvPair["vrf"])
		}
		return nil
	}
}
