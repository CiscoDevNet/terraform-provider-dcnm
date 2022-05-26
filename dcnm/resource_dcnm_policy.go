package dcnm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var policyDeployMutexMap = make(map[string]*sync.Mutex, 0)

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
		Id := models.G(cont, "id")
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
		response := models.G(cont, "successList")
		var info []map[string]interface{}
		_ = json.Unmarshal([]byte(response), &info)
		message := info[0]["message"].(string)
		Id := GetID(message)
		policy.PolicyId = "POLICY-" + Id
		d.SetId(Id)
	}

	// Deploy the policy
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) {
		err := deployPolicy(dcnmClient, policy.PolicyId, serialNumber)
		if err != nil {
			d.Set("deploy", false)
			return diag.FromErr(err)
		}
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
	d.Set("policy_id", models.G(cont, "policyId"))
	d.Set("serial_number", models.G(cont, "serialNumber"))

	if cont.Exists("source") {
		d.Set("source", models.G(cont, "source"))
	}
	if cont.Exists("description") {
		d.Set("description", models.G(cont, "description"))
	}
	if cont.Exists("entityType") {
		d.Set("entity_type", models.G(cont, "entityType"))
	}
	if cont.Exists("entityName") {
		d.Set("entity_name", models.G(cont, "entityName"))
	}
	if cont.Exists("templateName") {
		d.Set("template_name", models.G(cont, "templateName"))
	}
	if cont.Exists("templateContentType") {
		d.Set("template_content_type", models.G(cont, "templateContentType"))
	}
	if cont.Exists("priority") {
		d.Set("priority", models.G(cont, "priority"))
	}
	var strByte []byte
	if cont.Exists("nvPairs") {
		strJson := models.G(cont, "nvPairs")
		strByte = []byte(strJson)
		var nvPair map[string]interface{}
		json.Unmarshal(strByte, &nvPair)
		props, ok := d.GetOk("template_props")

		map2 := make(map[string]interface{})
		for k := range props.(map[string]interface{}) {
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
	policyId := dn
	if dcnmClient.GetPlatform() != "nd" {
		policyId = "POLICY-" + dn
	}
	cont, err := getAllPolicy(dcnmClient, policyId)
	if err != nil {
		d.SetId("")
		if cont != nil {
			log.Println(cont.String())
		}
		log.Println(cont.String())
		return nil
	}
	setPolicyAttributes(d, cont)
	d.SetId(models.G(cont, "id"))
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
	policy.Deleted = false
	policy.NVPairs = nvPairMap
	policy.PolicyId = "POLICY-" + d.Id()
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
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) {
		err := deployPolicy(dcnmClient, policy.PolicyId, serialNumber)
		if err != nil {
			d.Set("deploy", false)
			return diag.FromErr(err)
		}
	}
	d.SetId(models.G(cont, "id"))
	return resourceDCNMPolicyRead(ctx, d, m)

}
func resourceDCNMPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)

	policy := models.Policy{
		Id:           d.Id(),
		PolicyId:     "POLICY-" + d.Id(),
		SerialNumber: d.Get("serial_number").(string),
		TemplateName: d.Get("template_name").(string),
		NVPairs:      d.Get("template_props").(map[string]interface{}),
		Deleted:      true,
	}

	policy.NVPairs = d.Get("template_props").(map[string]interface{})

	dUrl := fmt.Sprintf("/rest/control/policies/%s", policy.PolicyId)
	cont, err := dcnmClient.Update(dUrl, &policy)
	if err != nil {
		if cont != nil {
			return diag.Errorf(cont.String())
		}
		return diag.FromErr(err)
	}

	if _, ok := policyDeployMutexMap[policy.SerialNumber]; !ok {
		policyDeployMutexMap[policy.SerialNumber] = &sync.Mutex{}
	}

	policyDeployMutexMap[policy.SerialNumber].Lock()
	_, err = getAllPolicy(dcnmClient, policy.PolicyId)
	if err == nil {
		recurSwitchDeployment(dcnmClient, d.Get("serial_number").(string))
	}
	policyDeployMutexMap[policy.SerialNumber].Unlock()

	d.SetId("")
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}

func deploySwitchFabric(dcnmClient *client.Client, serialNumber string) error {
	// get fabric by switch serial number
	url := fmt.Sprintf("/rest/control/switches/%s/fabric-name", serialNumber)
	cont, err := dcnmClient.GetviaURL(url)
	if err != nil {
		return fmt.Errorf("error deploying fabric after policy deletion: %w", err)
	}

	fabric := models.G(cont, "fabricName")

	// deploy fabric
	err = deployswitch(dcnmClient, fabric, serialNumber, 300)
	if err != nil {
		return fmt.Errorf("error deploying fabric after policy deletion: %w", err)
	}

	return nil
}

func recurSwitchDeployment(dcnmClient *client.Client, serialNumber string) {
	err := deploySwitchFabric(dcnmClient, serialNumber)
	if err != nil {
		recurSwitchDeployment(dcnmClient, serialNumber)
	}
}

func deployPolicy(dcnmClient *client.Client, policyId, serialNumber string) error {
	log.Println("[DEBUG] Begining Deployment ", policyId)

	if _, ok := policyDeployMutexMap[serialNumber]; !ok {
		policyDeployMutexMap[serialNumber] = &sync.Mutex{}
	}

	policyDeployMutexMap[serialNumber].Lock()
	_, err := dcnmClient.SaveDeploy("/rest/control/policies/deploy", policyId)
	if err != nil {
		return fmt.Errorf("policy is created but failed to deploy with error : %s", err)
	}
	policyDeployMutexMap[serialNumber].Unlock()
	log.Println("[DEBUG] End of Deployment ", policyId)
	return nil
}
