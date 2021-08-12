package dcnm

import (
	"fmt"
	"log"
	"testing"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerfPolicy *schema.Provider

func TestAccDCNMPolicy_Basic(t *testing.T) {
	var policy models.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfPolicy),
		// CheckDestroy:      testAccCheckDCNMPolicyDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDCNMPolicyConfig_basic("test-demo-1", "description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMPolicyExists("dcnm_policy.first", &policy),
					testAccCheckDCNMPolicyAttributes("test-demo-1", &policy),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDCNMPolicy_Update(t *testing.T) {
	var policy models.Policy
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfPolicy),
		// CheckDestroy:      testAccCheckDCNMPolicyDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDCNMPolicyConfig_basic("test-demo-1", "updated-description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMPolicyExists("dcnm_policy.first", &policy),
					testAccCheckDCNMPolicyAttributes("test-demo-1", &policy),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}
func testAccCheckDCNMPolicyExists(name string, policy *models.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Policy %s not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Policy dn was set")
		}

		dcnmClient := (*providerfPolicy).Meta().(*client.Client)
		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/control/policies/%s", "test-demo-1"))
		log.Printf("[DEBUG] before err %s", cont)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] after err %s", cont)
		policyTest := &models.Policy{}
		policyTest.PolicyId = stripQuotes(cont.S("policyId").String())
		policyTest.SerialNumber = stripQuotes(cont.S("serialNumber").String())
		policyTest.TemplateName = stripQuotes(cont.S("templateName").String())
		// policyTest.NVPairs = stripQuotes(cont.S("template_props").String())
		*policy = *policyTest
		return nil
	}

}
func testAccCheckDCNMPolicyConfig_basic(policyId string, description string) string {
	return fmt.Sprintf(`
	resource "dcnm_policy" "first" {
		serial_number = "9BH270169LJ"
		description="%s"
		template_name = "aaa_radius_deadtime"
		template_props = {
        "DTIME" : "0"
        "AAA_GROUP" : "%s"
      }
	}
	`, description, "management")
}

func testAccCheckDCNMPolicyDestroy(s *terraform.State) error {
	dcnmClient := (*providerfPolicy).Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_policy" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/control/policies/%s", "test-demo-1"))
			if err == nil {
				return fmt.Errorf("Policy still exists!!")
			}
		}
	}

	return nil
}
func testAccCheckDCNMPolicyAttributes(name string, policy *models.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test-demo-1" != policy.PolicyId {
			return fmt.Errorf("Bad Policy name %s", policy.PolicyId)
		}
		if "9BH270169LJ" != policy.SerialNumber {
			return fmt.Errorf("Bad serial number %s", policy.SerialNumber)
		}
		if "aaa_radius_deadtime" != policy.TemplateName {
			return fmt.Errorf("Bas template name %s", policy.TemplateName)
		}
		return nil
	}
}
