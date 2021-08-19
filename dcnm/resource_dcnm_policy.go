package dcnm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDCNMPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMPolicyCreate,
		Read:   resourceDCNMPolicyRead,
		Update: resourceDCNMPolicyUpdate,
		Delete: resourceDCNMPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDCNMPolicyImporter,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: schema.SchemaValidateFunc(IsEmpty()),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
			"serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"entity_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_content_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_props": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"deploy": &schema.Schema{
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
func resourceDCNMPolicyCreate(d *schema.ResourceData, m interface{}) error {
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
			return err
		}
		Id := stripQuotes(cont.S("id").String())
		policy.PolicyId = "POLICY-" + Id
		d.SetId(Id)

	} else {

		cont, err := dcnmClient.Save("/rest/control/policies/bulk-create", &policy)
		if err != nil {
			if cont != nil {
				return fmt.Errorf(cont.String())
			}
			return err
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
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		log.Println("[DEBUG] Begining Deployment ", d.Id())

		_, err := dcnmClient.SaveDeploy("/rest/control/policies/deploy", policy.PolicyId)
		if err != nil {
			d.Set("deploy", false)
			return fmt.Errorf("policy is created but failed to deploy with error : %s", err)
		}
		log.Println("[DEBUG] End of Deployment ", d.Id())
	}
	return resourceDCNMPolicyRead(d, m)

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
	d.Set("source", stripQuotes(cont.S("source").String()))
	d.Set("description", stripQuotes(cont.S("description").String()))
	d.Set("entity_type", stripQuotes(cont.S("entityType").String()))
	d.Set("entity_name", stripQuotes(cont.S("entityName").String()))
	d.Set("template_name", stripQuotes(cont.S("templateName").String()))
	d.Set("template_content_type", stripQuotes(cont.S("templateContentType").String()))
	d.Set("priority", stripQuotes(cont.S("priority").String()))
	var strByte []byte
	strJson := stripQuotes(cont.S("nvPairs").String())
	strByte = []byte(strJson)
	var nvPair map[string]interface{}
	json.Unmarshal(strByte, &nvPair)
	props := d.Get("template_props").(map[string]interface{})
	map2 := make(map[string]interface{})
	for k, _ := range props {
		map2[k] = nvPair[k]

	}
	d.Set("template_props", map2)

	return d
}
func resourceDCNMPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read Method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	var policyId string

	policyId = "POLICY-" + dn

	cont, err := getAllPolicy(dcnmClient, policyId)

	if err != nil {
		// d.SetId("")
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}
	setPolicyAttributes(d, cont)
	d.SetId(stripQuotes(cont.S("id").String()))
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil

}
func resourceDCNMPolicyUpdate(d *schema.ResourceData, m interface{}) error {
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
			return fmt.Errorf(cont.String())
		}
		return err
	}
	// Deploy the policy
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		log.Println("[DEBUG] Begining Deployment ", d.Id())

		_, err := dcnmClient.SaveDeploy("/rest/control/policies/deploy", policy.PolicyId)
		if err != nil {
			d.Set("deploy", false)
			return fmt.Errorf("policy is created but failed to deploy with error : %s", err)
		}
		log.Println("[DEBUG] End of Deployment ", d.Id())
	}
	d.SetId(stripQuotes(cont.S("id").String()))
	return resourceDCNMPolicyRead(d, m)

}
func resourceDCNMPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)

	policyId := d.Get("policy_id").(string)
	durl := fmt.Sprintf("/rest/control/policies/%s", policyId)
	cont, err := dcnmClient.Delete(durl)
	if err != nil {
		if cont != nil {
			return fmt.Errorf(cont.String())
		}
		return err
	}

	d.SetId("")

	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil

}
