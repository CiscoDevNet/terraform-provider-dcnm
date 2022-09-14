package dcnm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const POLICY_PREFIX string = "POLICY-"
const MAX_RETRY_CREATE int = 10
const MAX_RETRY_DEL int = 4

var switchDeployMutexMap = make(map[string]*sync.Mutex, 0)
var policyURLs = map[string]string{
	"Create":        "/rest/control/policies",
	"PolicyDeploy":  "/rest/control/policies/deploy",
	"MarkDelete":    "/rest/control/policies/%s/mark-delete",
	"IntentConfig":  "/rest/control/policies/%s/intent-config",
	"Common":        "/rest/control/policies/%s",
	"GetFabricName": "/rest/control/switches/%s/fabric-name",
}

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
			"deploy_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
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
	log.Println("[DEBUG] Beginning Importer ", d.Id())
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
	log.Println("[DEBUG] Beginning Create method")

	dcnmClient := m.(*client.Client)

	serialNumber := d.Get("serial_number").(string)
	templateName := d.Get("template_name").(string)
	nvPairMap := d.Get("template_props").(map[string]interface{})
	deployTimeout := d.Get("deploy_timeout").(int)
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

	cont, err := dcnmClient.Save(policyURLs["Create"], &policy)
	if err != nil {
		return diag.FromErr(err)
	}
	Id := models.G(cont, "id")
	policy.PolicyId = POLICY_PREFIX + Id
	d.SetId(Id)

	// Deploy the policy
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) {
		err := deployPolicyWithTimeout(dcnmClient, policy.PolicyId, serialNumber, deployTimeout)
		if err != nil {
			d.Set("deploy", false)
			return diag.FromErr(err)
		}
	}

	return resourceDCNMPolicyRead(ctx, d, m)
}

func getAllPolicy(client *client.Client, policyId string) (*container.Container, error) {
	duro := fmt.Sprintf(policyURLs["Common"], policyId)
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
	log.Println("[DEBUG] Beginning Read Method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	policyId := dn
	if dcnmClient.GetPlatform() != "nd" {
		policyId = POLICY_PREFIX + dn
	}
	cont, err := getAllPolicy(dcnmClient, policyId)
	if err != nil {
		d.SetId("")
		if cont != nil {
			log.Printf("[DEBUG] error while reading policy(%s): %v", policyId, cont.String())
		}
		log.Printf("[DEBUG] error while reading policy(%v): %v", policyId, err)
		return nil
	}
	setPolicyAttributes(d, cont)
	d.SetId(models.G(cont, "id"))
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil

}

func resourceDCNMPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Beginning Update method")

	dcnmClient := m.(*client.Client)

	serialNumber := d.Get("serial_number").(string)
	templateName := d.Get("template_name").(string)
	nvPairMap := d.Get("template_props").(map[string]interface{})
	deployTimeout := d.Get("deploy_timeout").(int)

	policy := models.Policy{}

	policy.SerialNumber = serialNumber
	policy.TemplateName = templateName
	policy.Deleted = false
	policy.NVPairs = nvPairMap
	policy.PolicyId = POLICY_PREFIX + d.Id()
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
	dUrl := fmt.Sprintf(policyURLs["Common"], policy.PolicyId)
	cont, err := dcnmClient.Update(dUrl, &policy)
	if err != nil {
		if cont != nil {
			return diag.Errorf(cont.String())
		}
		return diag.FromErr(err)
	}
	// Deploy the policy
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) {
		err := deployPolicyWithTimeout(dcnmClient, policy.PolicyId, serialNumber, deployTimeout)
		if err != nil {
			d.Set("deploy", false)
			return diag.FromErr(err)
		}
	}
	d.SetId(models.G(cont, "id"))
	return resourceDCNMPolicyRead(ctx, d, m)

}

func resourceDCNMPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Beginning Delete method ", d.Id())
	dcnmClient := m.(*client.Client)
	serialNumber := d.Get("serial_number").(string)

	url := fmt.Sprintf(policyURLs["GetFabricName"], serialNumber)
	cont, err := dcnmClient.GetviaURL(url)
	if err != nil {
		return diag.Errorf("error deploying fabric after policy deletion: %w", err)
	}
	fabric := models.G(cont, "fabricName")

	deleteFlag := false

	//Mark delete policy
	url = fmt.Sprintf(policyURLs["MarkDelete"], d.Id())
	cont, err = deletePolicy(url, dcnmClient)
	if err != nil {
		if cont != nil {
			return diag.Errorf(cont.String())
		}
		return diag.FromErr(err)
	}

	//Intent-config checking
	url = fmt.Sprintf(policyURLs["IntentConfig"], d.Id())
	cont, err = dcnmClient.GetviaURL(url)
	if err != nil {
		if err.Error() == fmt.Sprintf("Policy %s does not exist", d.Id()) {
			deleteFlag = true
		} else {
			return diag.Errorf("error deletion policy: %s", err)
		}
	}

	markDeleteConfig := models.G(cont, "markDeletedConfig")

	if markDeleteConfig == "No config is available" && !deleteFlag {
		dUrl := fmt.Sprintf(policyURLs["Common"], d.Id())
		cont, err = dcnmClient.Delete(dUrl)
		if err != nil && err.Error() != fmt.Sprintf("Policy %s does not exist", d.Id()) {
			if cont != nil {
				return diag.Errorf(cont.String())
			}
			return diag.Errorf("error while destroying policy %v", err)
		}
	}

	if _, ok := switchDeployMutexMap[serialNumber]; !ok {
		switchDeployMutexMap[serialNumber] = &sync.Mutex{}
	}

	switchDeployMutexMap[serialNumber].Lock()
	defer switchDeployMutexMap[serialNumber].Unlock()
	for count := 1; count <= MAX_RETRY_DEL; count++ {

		isDeployed, err := checkDeploy(dcnmClient, fabric, serialNumber)
		if err != nil {
			return diag.Errorf("error deploying fabric after policy deletion: %w", err)
		}
		if isDeployed {
			break
		}

		err = deployswitch(dcnmClient, fabric, serialNumber)
		if err == nil {
			break
		}
		if count == MAX_RETRY_DEL {
			return diag.Errorf("error deploying fabric after policy deletion: %s", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}

	d.SetId("")
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}

func deployPolicyWithTimeout(dcnmClient *client.Client, policyId, serialNumber string, timeout int) error {
	log.Println("[DEBUG] Beginning Deployment for Create ", policyId)

	for count := 1; count <= MAX_RETRY_CREATE; count++ {
		cont, err := saveDeployWithTimeout(dcnmClient, policyURLs["PolicyDeploy"], policyId, timeout)
		if err != nil {
			return fmt.Errorf("policy is created but failed to deploy with error : %s", err)
		}
		idFailed := models.G(cont.Index(0), "failedPTIList")
		idSuccess := models.G(cont.Index(0), "successPTIList")
		if idFailed == policyId && count == MAX_RETRY_CREATE {
			return fmt.Errorf("policy is created but failed to deploy policy with id: %s", policyId)
		}
		if idSuccess == policyId {
			break
		}
		time.Sleep(time.Second * 5)
	}
	log.Println("[DEBUG] End of Deployment ", policyId)
	return nil
}
func saveDeployWithTimeout(dcnmClient *client.Client, url, policyId string, timeout int) (*container.Container, error) {
	cont := make(chan *container.Container, 1)
	result := make(chan error, 1)
	go func() {
		container, err := dcnmClient.SaveDeploy(url, policyId)
		cont <- container
		result <- err
	}()
	// Wait until timeout occurs or a response is received
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		log.Println("[DEBUG] Retry Deployment timeout :", policyId)
		return nil, nil
	case container := <-cont:
		return container, nil
	case res := <-result:
		return nil, res
	}
}

func deletePolicy(url string, dcnmClient *client.Client) (*container.Container, error) {

	req, err := dcnmClient.MakeRequest("PUT", url, nil, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := dcnmClient.Do(req, false)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return cont, nil
	}

	return cont, fmt.Errorf("%d Error : %s", resp.StatusCode, cont.S("message").String())
}
