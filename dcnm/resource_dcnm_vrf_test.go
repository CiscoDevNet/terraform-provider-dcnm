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

var providerfVrf *schema.Provider

func TestAccDCNMVRF_Basic(t *testing.T) {
	var vrf models.VRF
	var vrfProfile models.VRFProfileConfig

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfVrf),
		CheckDestroy:      testAccCheckDCNMVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMVRFConfig_basic("vrf decription check", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMVRFExists("dcnm_vrf.vrf_check", &vrf, &vrfProfile),
					testAccCheckDCNMVRFAttributes("vrf decription check", &vrf, &vrfProfile),
				),
			},
		},
	})
}

func TestAccDCNMVRF_Update(t *testing.T) {
	var vrf models.VRF
	var vrfProfile models.VRFProfileConfig

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfVrf),
		CheckDestroy:      testAccCheckDCNMVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCNMVRFConfig_basic("vrf decription check", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMVRFExists("dcnm_vrf.vrf_check", &vrf, &vrfProfile),
					testAccCheckDCNMVRFAttributes("vrf decription check", &vrf, &vrfProfile),
				),
			},
			{
				Config: testAccCheckDCNMVRFConfig_basic("vrf update check", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMVRFExists("dcnm_vrf.vrf_check", &vrf, &vrfProfile),
					testAccCheckDCNMVRFAttributes("vrf update check", &vrf, &vrfProfile),
				),
			},
		},
	})
}

func testAccCheckDCNMVRFConfig_basic(desc string, deploy string) string {
	return fmt.Sprintf(`
	resource "dcnm_vrf" "vrf_check" {
		fabric_name = "fab2"
		name = "two" 
		vlan_id = 2002
		vlan_name = "check"
		description = "%s"
		intf_description = "vrf"
		deploy = "%s"
		attachments {
			serial_number = "9ZGMF8CBZK5"
			vlan_id       = 2300
			attach        = %s
			loopback_id   = 70
			loopback_ipv4 = "1.2.3.4"
		}
	}
	`, desc, deploy, deploy)
}

func testAccCheckDCNMVRFExists(name string, vrf *models.VRF, profile *models.VRFProfileConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VRF %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VRF dn was set")
		}

		dcnmClient := (*providerfVrf).Meta().(*client.Client)

		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/%s", "fab2", rs.Primary.ID))
		if err != nil {
			return err
		}

		vrfGet := &models.VRF{}
		vrfGet.Fabric = stripQuotes(cont.S("fabric").String())
		vrfGet.Name = stripQuotes(cont.S("vrfName").String())

		vrfProfile := &models.VRFProfileConfig{}
		configCont, err := cleanJsonString(stripQuotes(cont.S("vrfTemplateConfig").String()))
		if err != nil {
			return err
		}
		if configCont.Exists("mtu") {
			if mtu, err := strconv.Atoi(stripQuotes(configCont.S("mtu").String())); err == nil {
				vrfProfile.Mtu = mtu
			}
		}
		if configCont.Exists("vrfVlanId") && stripQuotes(configCont.S("vrfVlanId").String()) != "" {
			if vlan, err := strconv.Atoi(stripQuotes(configCont.S("vrfVlanId").String())); err == nil {
				vrfProfile.Vlan = vlan
			}
		}
		if configCont.Exists("vrfVlanName") {
			vrfProfile.VlanName = stripQuotes(configCont.S("vrfVlanName").String())
		}
		if configCont.Exists("vrfDescription") {
			vrfProfile.Description = stripQuotes(configCont.S("vrfDescription").String())
		}

		*vrf = *vrfGet
		*profile = *vrfProfile
		return nil
	}
}

func testAccCheckDCNMVRFDestroy(s *terraform.State) error {
	dcnmClient := (*providerfVrf).Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_vrf" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/%s", "fab2", rs.Primary.ID))
			if err == nil {
				return fmt.Errorf("VRF still exists")
			}
		}
	}

	return nil
}

func testAccCheckDCNMVRFAttributes(desc string, vrf *models.VRF, profile *models.VRFProfileConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "two" != vrf.Name {
			return fmt.Errorf("Bad VRF Name %s", vrf.Name)
		}

		if "fab2" != vrf.Fabric {
			return fmt.Errorf("Bad VRF fabric name %s", vrf.Fabric)
		}

		if 2002 != profile.Vlan {
			return fmt.Errorf("Bad VRF VLAN %d", profile.Vlan)
		}

		if "check" != profile.VlanName {
			return fmt.Errorf("Bad VRF VLAN name %s", profile.VlanName)
		}

		if desc != profile.Description {
			return fmt.Errorf("Bad VRF description %s", profile.Description)
		}
		return nil
	}
}
