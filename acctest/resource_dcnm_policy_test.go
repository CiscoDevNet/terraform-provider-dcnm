package acctest

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerfPolicy *schema.Provider

func TestAccDCNMPolicy_Basic(t *testing.T) {
	var policy_default models.Policy
	var policy_updated models.Policy
	defaultSerialNumber := "9LMU8W6W8VG"
	otherSerialNumber := "9CIWTMB13GP"
	resourceName := "dcnm_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfPolicy),
		CheckDestroy:      testAccCheckDCNMPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testCreatePolicyWithoutTemplate(defaultSerialNumber),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      testCreatePolicyMissingTemplateProps(defaultSerialNumber),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      testCreatePolicyWithoutSerialNumber(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: testCreatePolicyBasic(defaultSerialNumber),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMPolicyExists(resourceName, &policy_default),
					resource.TestCheckResourceAttr(resourceName, "priority", "500"),
					resource.TestCheckResourceAttr(resourceName, "source", ""),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entity_name", ""),
					resource.TestCheckResourceAttr(resourceName, "entity_type", ""),
					resource.TestCheckResourceAttr(resourceName, "template_content_type", "TEMPLATE_CLI"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deploy"},
			},
			{
				Config: testCreatePolicyBasicWithOptionalValues(defaultSerialNumber),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMPolicyExists(resourceName, &policy_updated),
					resource.TestCheckResourceAttr(resourceName, "priority", "500"),
					resource.TestCheckResourceAttr(resourceName, "source", "Ethernet1/3_FABRIC"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is test policy."),
					resource.TestCheckResourceAttr(resourceName, "entity_name", "Ethernet1/3"),
					resource.TestCheckResourceAttr(resourceName, "entity_type", "INTERFACE"),
					resource.TestCheckResourceAttr(resourceName, "template_content_type", "TEMPLATE_CLI"),
					testAccCheckPolicyIdEqual(&policy_default, &policy_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deploy"},
			},
			{
				Config: testCreatePolicyBasic(otherSerialNumber),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMPolicyExists(resourceName, &policy_updated),
					resource.TestCheckResourceAttr(resourceName, "serial_number", otherSerialNumber),
					testAccCheckPolicyIdNotEqual(&policy_default, &policy_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deploy"},
			},
			{
				Config: testCreatePolicyBasic(defaultSerialNumber),
			},
		},
	})
}

// func TestAccDCNMPolicy_Update(t *testing.T) {
// 	var policy models.Policy
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:          func() { testAccPreCheck(t) },
// 		ProviderFactories: testAccProviderFactoriesInternal(&providerfPolicy),
// 		// CheckDestroy:      testAccCheckDCNMPolicyDestroy,
// 		Steps: []resource.TestStep{
// 			{

// 				Config: testAccCheckDCNMPolicyConfig_basic("test-demo-1", "updated-description"),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckDCNMPolicyExists("dcnm_policy.first", &policy),
// 					testAccCheckDCNMPolicyAttributes("test-demo-1", &policy),
// 				),
// 				ExpectNonEmptyPlan: true,
// 			},
// 		},
// 	})
// }

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
		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/control/policies/POLICY-%s", rs.Primary.ID))
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

// Helper functions for tests

func testCreatePolicyWithoutTemplate(serial_number string) string {
	return fmt.Sprintf(`
	resource "dcnm_policy" "test" {
		serial_number = "%s"
	}
	`, serial_number)
}
func testCreatePolicyWithoutSerialNumber() string {
	return `
	resource "dcnm_policy" "test" {
		template_name 	= 	"aaa_radius_deadtime"
		template_props 	= 	{
			"DTIME" : "0"
			"AAA_GROUP" : "management"
		}
	}
	`
}
func testCreatePolicyMissingTemplateProps(serial_number string) string {
	return fmt.Sprintf(`
	resource "dcnm_policy" "test" {
		serial_number = "%s"
		template_name 	= 	"aaa_radius_deadtime"
		template_props 	= 	{
			"DTIME" : "0"
		}
	}
	`, serial_number)
}
func testCreatePolicyBasic(serial_number string) string {
	return fmt.Sprintf(`
	resource "dcnm_policy" "test" {
		serial_number  	= 	"%s"
		template_name  	= 	"aaa_radius_deadtime"
		template_props 	= 	{
				"DTIME" : "0"
				"AAA_GROUP" : "management"
			}
		}
	`, serial_number)
}
func testCreatePolicyBasicWithOptionalValues(serial_number string) string {
	return fmt.Sprintf(`
	resource "dcnm_policy" "test" {
		serial_number  	= 	"%s"
		description    	=	"This is test policy."
		template_name  	= 	"aaa_radius_deadtime"
		template_props 	= {
			"DTIME" : "0"
			"AAA_GROUP" : "management"
		}
		priority        =   500
		source          =   "Ethernet1/3_FABRIC"
		entity_name     =   "Ethernet1/3"
		entity_type     =   "INTERFACE"
		template_content_type   =   "TEMPLATE_CLI"
	}
	`, serial_number)
}
func testAccCheckPolicyIdEqual(pid1, pid2 *models.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pid1.PolicyId != pid2.PolicyId {
			return fmt.Errorf("Poliicy IDs are not equal")
		}
		return nil
	}
}
func testAccCheckPolicyIdNotEqual(pid1, pid2 *models.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pid1.PolicyId == pid2.PolicyId {
			return fmt.Errorf("Poliicy IDs are equal")
		}
		return nil
	}
}

func testAccCheckDCNMPolicyDestroy(s *terraform.State) error {
	dcnmClient := (*providerfPolicy).Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_policy" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/control/policies/POLICY-%s", rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("Policy still exists!!")
			}
		}
	}

	return nil
}
