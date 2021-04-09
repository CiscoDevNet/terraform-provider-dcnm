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

var providerfInv *schema.Provider

func TestAccDCNMInventory_Basic(t *testing.T) {
	var inv models.Inventory

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfInv),
		CheckDestroy:      testAccCheckDCNMInventoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMInventoryConfig_basic("172.25.74.93"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMInventoryExists("dcnm_inventory.test", &inv),
					testAccCheckDCNMInventoryAttributes("172.25.74.93", &inv),
				),
			},
		},
	})
}

func TestAccDCNMInventory_Update(t *testing.T) {
	var inv models.Inventory

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfInv),
		CheckDestroy:      testAccCheckDCNMInventoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMInventoryConfig_basic("172.25.74.93"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMInventoryExists("dcnm_inventory.test", &inv),
					testAccCheckDCNMInventoryAttributes("172.25.74.93", &inv),
				),
			},
			{
				Config: testAccCheckDCNMInventoryConfig_basic("172.25.74.93"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMInventoryExists("dcnm_inventory.test", &inv),
					testAccCheckDCNMInventoryAttributes("172.25.74.93", &inv),
				),
			},
		},
	})
}

func testAccCheckDCNMInventoryConfig_basic(ip string) string {
	return fmt.Sprintf(`
	resource "dcnm_inventory" "test" {
		fabric_name   = "fab2"
		username      = "admin"
		password      = "ins3965!"
		max_hops      = 0
		preserve_config = "false"
		auth_protocol = 0
		config_timeout = 10
		switch_config {
		  ip   = "%s"
		  role = "leaf"
		}
	  }
	`, ip)
}

func testAccCheckDCNMInventoryExists(name string, inv *models.Inventory) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Inventory %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Inventory dn was set")
		}

		dcnmClient := (*providerfInv).Meta().(*client.Client)

		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/control/fabrics/%s/inventory", "fab2"))
		if err != nil {
			return err
		}

		sGet := &models.Inventory{}
		var flag bool
		for i := 0; i < len(cont.Data().([]interface{})); i++ {
			switchCont := cont.Index(i)

			ipGet := stripQuotes(switchCont.S("ipAddress").String())
			if ipGet == rs.Primary.ID {
				sGet.SeedIP = stripQuotes(switchCont.S("ipAddress").String())
				flag = true
				break
			}
		}

		if flag != true {
			return fmt.Errorf("Desired switch not found")
		}
		*inv = *sGet
		return nil
	}
}

func testAccCheckDCNMInventoryDestroy(s *terraform.State) error {
	dcnmClient := (*providerfInv).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_inventory" {
			cont, _ := dcnmClient.GetviaURL(fmt.Sprintf("/rest/control/fabrics/%s/inventory", "fab2"))
			var flag bool
			for i := 0; i < len(cont.Data().([]interface{})); i++ {
				switchCont := cont.Index(i)

				ipGet := stripQuotes(switchCont.S("ipAddress").String())
				if ipGet == rs.Primary.ID {
					flag = true
					break
				}
			}
			if flag {
				return fmt.Errorf("Switch inventory still exists")
			}
		}
	}

	return nil
}

func testAccCheckDCNMInventoryAttributes(ip string, inv *models.Inventory) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ip != inv.SeedIP {
			return fmt.Errorf("Bad Switch IP %s", inv.SeedIP)
		}
		return nil
	}
}
