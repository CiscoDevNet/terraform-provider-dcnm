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

var providerfTemplate *schema.Provider

func TestAccDCNMTemplate_Basic(t *testing.T) {
	var template models.Template

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfTemplate),
		// CheckDestroy:      testAccCheckDCNMTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMTemplateConfig_basic("aa"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMTemplateExists("dcnm_template.ex", &template),
					testAccCheckDCNMTemplateAttributes("aa", &template),
				),
			},
		},
	})
}

func testAccCheckDCNMTemplateConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "dcnm_template" "ex"{
    name = "%s"
    content=file("poc.txt")
}
	`, name)
}

func testAccCheckDCNMTemplateExists(name string, template *models.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Template %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Template dn was set")
		}

		dcnmClient := (*providerfTemplate).Meta().(*client.Client)

		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/config/templates/%s", name))
		if err != nil {
			return err
		}

		TemplateGet := &models.Template{}
		TemplateGet.Name = stripQuotes(cont.S("name").String())
		TemplateGet.Content = stripQuotes(cont.S("content").String())

		*template = *TemplateGet
		return nil
	}
}

func testAccCheckDCNMTemplateDestroy(s *terraform.State) error {
	dcnmClient := (*providerfTemplate).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_template" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/config/templates/%s", "aa"))
			if err == nil {
				return fmt.Errorf("Template still exists")
			}
		}
	}

	return nil
}

func testAccCheckDCNMTemplateAttributes(name string, template *models.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "t1" != template.Name {
			return fmt.Errorf("Bad Template Name %s", template.Name)
		}

		return nil
	}
}
