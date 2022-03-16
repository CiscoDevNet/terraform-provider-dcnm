package dcnm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDCNMPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDCNMPolicyCreate,
		ReadContext:   resourceDCNMPolicyRead,
		UpdateContext: resourceDCNMPolicyUpdate,
		DeleteContext: resourceDCNMPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDCNMPolicyImporter,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: IsEmpty(),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			"serial_number": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"entity_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_props": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"deploy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func IsEmpty() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if len(v) > 0 {
			es = append(es, fmt.Errorf("expected %s to be empty", k))
			return
		}
		return
	}
}
func resourceDCNMPolicyImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Importer ", d.Id())
	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 1 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	policyId := importInfo[0]
	cont, err := getAllPolicy(dcnmClient, policyId)
	if err != nil {
		return nil, err
	}
	stateImport := setPolicyAttributes(d, cont)
	log.Println("[DEBUG] End of Importer ", d.Id())
	return []*schema.ResourceData{stateImport}, nil

}
func GetID(description string) string {
	policyId := strings.Split(description, " ")[0]
	id := strings.Split(policyId, "-")[1]
	return id
}
func resourceDCNMPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Create method")

	dcnmClient := m.(*client.Client)

	serialNumber := d.Get("serial_number").(string)
	templateName := d.Get("template_name").(string)
	nvPairMap := d.Get("template_props").(map[string]interface{})

	policy := models.Policy{}

	policy.SerialNumber = serialNumber
	policy.TemplateName = templateName
	policy.NVPairs = nvPairMap
	if source, ok := d.GetOk("source"); ok {
		policy.Source = source.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		policy.Description = description.(string)
	}
	if entityType, ok := d.GetOk("entity_type"); ok {
		policy.EntityType = entityType.(string)
	}
	if entityName, ok := d.GetOk("entity_name"); ok {
		policy.EntityName = entityName.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		policy.Priority = priority.(string)
	}
	if templateContentType, ok := d.GetOk("template_content_type"); ok {
		policy.TemplateContentType = templateContentType.(string)
	}
	if dcnmClient.GetPlatform() == "nd" {
		cont, err := dcnmClient.Save("/rest/control/policies", &policy)
		if err != nil {
			return diag.FromErr(err)
		}
		Id := stripQuotes(cont.S("id").String())
		policy.PolicyId = "POLICY-" + Id
		d.SetId(Id)

	} else {

		cont, err := dcnmClient.Save("/rest/control/policies/bulk-create", &policy)
		if err != nil {
			if cont != nil {
				return diag.Errorf(cont.String())
			}
			return diag.FromErr(err)
		}
		// Get the id from resource
		response := stripQuotes(cont.S("successList").String())
		var info []map[string]interface{}
		_ = json.Unmarshal([]byte(response), &info)
		message := info[0]["message"].(string)
		Id := GetID(message)
		policy.PolicyId = "POLICY-" + Id
		d.SetId(Id)
	}

	// Deploy the policy
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) {
		log.Println("[DEBUG] Begining Deployment ", d.Id())

		_, err := dcnmClient.SaveDeploy("/rest/control/policies/deploy", policy.PolicyId)
		if err != nil {
			d.Set("deploy", false)
			return diag.Errorf("policy is created but failed to deploy with error : %s", err)
		}
		log.Println("[DEBUG] End of Deployment ", d.Id())
	}
	return resourceDCNMPolicyRead(ctx, d, m)

}

func getAllPolicy(client *client.Client, policyId string) (*container.Container, error) {
	duro := fmt.Sprintf("/rest/control/policies/%s", policyId)
	cont, err := client.GetviaURL(duro)
	if err != nil {
		return cont, err
	}
	return cont, nil
}
func setPolicyAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	d.Set("policy_id", stripQuotes(cont.S("policyId").String()))
	d.Set("serial_number", stripQuotes(cont.S("serialNumber").String()))

	if cont.Exists("source") {
		d.Set("source", stripQuotes(cont.S("source").String()))
	}
	if cont.Exists("description") {
		d.Set("description", stripQuotes(cont.S("description").String()))
	}
	if cont.Exists("entityType") {
		d.Set("entity_type", stripQuotes(cont.S("entityType").String()))
	}
	if cont.Exists("entityName") {
		d.Set("entity_name", stripQuotes(cont.S("entityName").String()))
	}
	if cont.Exists("templateName") {
		d.Set("template_name", stripQuotes(cont.S("templateName").String()))
	}
	if cont.Exists("templateContentType") {
		d.Set("template_content_type", stripQuotes(cont.S("templateContentType").String()))
	}
	if cont.Exists("priority") {
		d.Set("priority", stripQuotes(cont.S("priority").String()))
	}
	var strByte []byte
	if cont.Exists("nvPairs") {
		strJson := stripQuotes(cont.S("nvPairs").String())
		strByte = []byte(strJson)
		var nvPair map[string]interface{}
		json.Unmarshal(strByte, &nvPair)
		props, ok := d.GetOk("template_props")

		map2 := make(map[string]interface{})
		for k, _ := range props.(map[string]interface{}) {
			map2[k] = nvPair[k]
		}
		if !ok {
			d.Set("template_props", nvPair)
		} else {

			d.Set("template_props", map2)
		}
	}

	return d
}
func resourceDCNMPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Read Method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	policyId := "POLICY-" + dn
	cont, err := getAllPolicy(dcnmClient, policyId)
	if err != nil {
		d.SetId("")
		if cont != nil {
			return diag.Errorf(cont.String())
		}
		return diag.FromErr(err)
	}
	setPolicyAttributes(d, cont)
	d.SetId(stripQuotes(cont.S("id").String()))
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil

}
func resourceDCNMPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Update method")

	dcnmClient := m.(*client.Client)

	serialNumber := d.Get("serial_number").(string)
	templateName := d.Get("template_name").(string)
	nvPairMap := d.Get("template_props").(map[string]interface{})

	policy := models.Policy{}

	policy.SerialNumber = serialNumber
	policy.TemplateName = templateName
	policy.NVPairs = nvPairMap
	if policyId, ok := d.GetOk("policy_id"); ok {
		policy.PolicyId = policyId.(string)

	} else {
		policy.PolicyId = "POLICY-" + d.Id()
	}
	if source, ok := d.GetOk("source"); ok {
		policy.Source = source.(string)
	}
	if description, ok := d.GetOk("description"); ok {
		policy.Description = description.(string)
	}
	if entityType, ok := d.GetOk("entity_type"); ok {
		policy.EntityType = entityType.(string)
	}
	if entityName, ok := d.GetOk("entity_name"); ok {
		policy.EntityName = entityName.(string)
	}
	if priority, ok := d.GetOk("priority"); ok {
		policy.Priority = priority.(string)
	}
	if templateContentType, ok := d.GetOk("template_content_type"); ok {
		policy.TemplateContentType = templateContentType.(string)
	}
	policy.Id = d.Id()
	dUrl := fmt.Sprintf("/rest/control/policies/%s", policy.PolicyId)
	cont, err := dcnmClient.Update(dUrl, &policy)
	if err != nil {
		if cont != nil {
			return diag.Errorf(cont.String())
		}
		return diag.FromErr(err)
	}
	// Deploy the policy
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		log.Println("[DEBUG] Begining Deployment ", d.Id())

		_, err := dcnmClient.SaveDeploy("/rest/control/policies/deploy", policy.PolicyId)
		if err != nil {
			d.Set("deploy", false)
			return diag.Errorf("policy is created but failed to deploy with error : %s", err)
		}
		log.Println("[DEBUG] End of Deployment ", d.Id())
	}
	d.SetId(stripQuotes(cont.S("id").String()))
	return resourceDCNMPolicyRead(ctx, d, m)

}
func resourceDCNMPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)

	policyId := d.Get("policy_id").(string)
	durl := fmt.Sprintf("/rest/control/policies/%s", policyId)
	cont, err := dcnmClient.Delete(durl)
	if err != nil {
		if cont != nil {
			return diag.Errorf(cont.String())
		}
		return diag.FromErr(err)
	}

	err = deploySwitchFabric(dcnmClient, d.Get("serial_number").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}

func deploySwitchFabric(dcnmClient *client.Client, serialNumber string) error {
	// get fabric by switch serial number
	url := fmt.Sprintf("/rest/control/switches/%s/fabric-name", serialNumber)
	cont, err := dcnmClient.GetviaURL(url)
	if err != nil {
		return fmt.Errorf("error deploying fabric after policy deletion: ", err)
	}

	// deploy fabric
	err = deployFabric(dcnmClient, models.G(cont, "fabricName"))
	if err != nil {
		fmt.Errorf("error deploying fabric after policy deletion: ", err)
	}

	return nil
}
