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

var providerfPeering *schema.Provider

func TestAccDCNMPeering_Basic(t *testing.T) {
	var peering models.RoutePeering
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfPeering),
		// CheckDestroy:      testAccCheckDCNMPeeringDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDCNMPeeringConfig_basic("RP-1", 1000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMPeeringExists("dcnm_route_peering.first", &peering),
					testAccCheckDCNMPeeringAttributes("RP-1", &peering),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})

}
func TestAccDCNMPeering_Update(t *testing.T) {
	var peering models.RoutePeering
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerfPeering),
		// CheckDestroy:      testAccCheckDCNMPolicyDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccCheckDCNMPeeringConfig_basic("RP-1", 1000),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCNMPeeringExists("dcnm_route_peering.first", &peering),
					testAccCheckDCNMPeeringAttributes("RP-1", &peering),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})

}
func testAccCheckDCNMPeeringDestroy(s *terraform.State) error {
	dcnmClient := (*providerfPeering).Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "dcnm_route_peering" {
			_, err := dcnmClient.GetviaURL(fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/testService/service-nodes/SN-3/peerings/Test_fabric_1/%s", "RP-1"))
			if err == nil {
				return fmt.Errorf("Peering still exists!!")
			}
		}
	}

	return nil
}
func testAccCheckDCNMPeeringConfig_basic(name string, vlan int) string {
	return fmt.Sprintf(`
	resource "dcnm_route_peering" first{
		name = "%s"
		attached_fabric_name = "Test_fabric_1"
		deployment_mode = "OneArmADC"
		fabric_name = "testService"
		option = "EBGPDynamicPeering"
		service_networks {
			network_name = "netadc"
			network_type = "ArmOneADC"
			template_name = "Service_Network_Universal"
			vlan_id = 1000
			vrf_name = "Test_VRF_2"
			gateway_ip_address ="124.168.2.1/24"
		}
		reverse_next_hop_ip = "124.168.2.10"
		service_node_name = "snadc"
		service_node_type = "ADC"
		deploy = false
		routes {
			template_name = "service_static_route"
			vrf_name = "Test_VRF_2"
			route_parmas = {
					"VRF_NAME": "Test_VRF_1"
			}
		}
	}`, name)
}

func testAccCheckDCNMPeeringExists(name string, peering *models.RoutePeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Peering %s not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Peering dn was set")
		}
		dcnmClient := (*providerfPeering).Meta().(*client.Client)
		cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/appcenter/Cisco/elasticservice/elasticservice-api/fabrics/testService/service-nodes/SN-3/peerings/Test_fabric_1/%s", "RP-1"))
		log.Printf("[DEBUG] before err %s", cont)
		if err != nil {
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
		}
		log.Printf("[DEBUG] after err %s", cont)
		peeringTest := &models.RoutePeering{}
		peeringTest.Name = stripQuotes(cont.S("peeringName").String())
		peeringTest.AttachedFabricName = stripQuotes(cont.S("attachedFabricName").String())
		peeringTest.FabricName = stripQuotes(cont.S("fabricName").String())
		peeringTest.ServiceNodeName = stripQuotes(cont.S("serviceNodeName").String())
		*peering = *peeringTest
		return nil
	}
}

func testAccCheckDCNMPeeringAttributes(name string, peering *models.RoutePeering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "RP-1" != peering.Name {
			return fmt.Errorf("Bad Peering name %s", peering.Name)
		}
		if "Test_fabric_1" != peering.AttachedFabricName {
			return fmt.Errorf("Bad attached fabric name  %s", peering.AttachedFabricName)
		}
		if "SN-3" != peering.ServiceNodeName {
			return fmt.Errorf("Bad service node name:  %s", peering.ServiceNodeName)

		}
		return nil
	}
}
