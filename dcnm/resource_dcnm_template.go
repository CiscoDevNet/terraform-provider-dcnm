package dcnm

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDCNMTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMTemplateCreate,
		Read:   resourceDCNMTemplateRead,
		Update: resourceDCNMTemplateUpdate,
		Delete: resourceDCNMTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDCNMTemplateImporter,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return CompareDiffs(old, new, d)
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"supported_platforms": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_sub_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_content_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

var TemplateURLS = map[string]map[string]string{
	"dcnm": {
		"Create":   "/rest/config/templates/template?templateName=%s",
		"Common":   "/rest/config/templates/%s",
		"Validate": "/rest/config/templates/validate",
	},
	"nd": {
		"Create":   "/appcenter/cisco/ndfc/api/v1/configtemplate/rest/config/templates/template?templateName=%s",
		"Common":   "/appcenter/cisco/ndfc/api/v1/configtemplate/rest/config/templates/%s",
		"Validate": "/configtemplate/rest/config/templates/validate",
	},
}

func CompareDiffs(old, new string, d *schema.ResourceData) bool {
	var old1 string
	re, err := regexp.Compile("((##template properties)(.*?)(##))")
	if err != nil {
		return false
	}
	old1 = strings.Replace(old, "\\r\\n", "\n", -1)
	old1 = strings.Replace(old1, "\\\"", "\"", -1)
	old2 := re.FindString(old1)
	old1 = strings.Replace(old1, old2, "", -1)
	old1 = strings.Replace(old1, "\\n", "\n", -1)
	re1 := regexp.MustCompile(`\r?\n`)
	new2, _ := GetStringInBetweenTwoString(new, "##template properties", "##")
	new1 := strings.Replace(new, new2, "", -1)
	new1 = re1.ReplaceAllString(new1, "\n")

	if len(old1) != 0 {
		old1 = old1[1:]
	}
	new1 = strings.Replace(new1, "\n\n", "\n", -1)
	old1 = strings.Replace(old1, "\n\n", "\n", -1)
	new1 = strings.Replace(new1, "\n\n", "\n", -1)
	old1 = strings.Replace(old1, "\n\n", "\n", -1)
	if old1 == new1 {
		return true
	}
	return false
}
func getTemplateProps(template string) (map[string]string, error) {

	templaterPropsStr := template
	templaterPropsStr = strings.Replace(templaterPropsStr, "template properties", "", -1)
	templaterPropsStr = strings.Replace(templaterPropsStr, "\\n", "", -1)
	templaterPropsStr = strings.Replace(templaterPropsStr, "\\", "", -1)
	templaterPropsStr = strings.Replace(templaterPropsStr, "\"", "", -1)
	templaterPropsStr = strings.Replace(templaterPropsStr, "##", "", -1)
	templaterPropsStr = strings.Replace(templaterPropsStr, " ", "", -1)
	templaterPropsList := strings.Split(templaterPropsStr, ";")

	templateProps := make(map[string]string)
	for _, prop := range templaterPropsList {
		l := strings.Split(prop, "=")
		if len(l) == 2 {
			templateProps[l[0]] = l[1]
		}
	}

	return templateProps, nil
}
func resourceDCNMTemplateImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Importer ", d.Id())
	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 1 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	name := importInfo[0]
	cont, err := getTemplate(dcnmClient, name)
	if err != nil {
		return nil, getErrorFromContainer(cont, err)
	}
	stateImport := setTemplateAttribute(d, cont)
	d.SetId(stripQuotes(cont.S("name").String()))
	log.Println("[DEBUG] End of Importer ", d.Id())
	return []*schema.ResourceData{stateImport}, nil

}
func resourceDCNMTemplateCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining of Creating template")
	dcnmClient := m.(*client.Client)
	name := d.Get("name").(string)
	var fileContent string
	fileContent = d.Get("content").(string)

	temp := models.Template{}
	temp.Name = name

	propertyMap := make(map[string]string)
	if description, ok := d.GetOk("description"); ok {
		propertyMap["description"] = description.(string)
	} else {
		propertyMap["description"] = ""
	}
	if tags, ok := d.GetOk("tags"); ok {
		propertyMap["tags"] = tags.(string)
	} else {
		propertyMap["tags"] = ""
	}
	if supported_platforms, ok := d.GetOk("supported_platforms"); ok {
		platformList := make([]string, 0, 1)
		for _, val := range supported_platforms.([]interface{}) {
			platformList = append(platformList, val.(string))
		}
		supported_platforms := strings.Join(platformList, ",")
		propertyMap["supported_platforms"] = supported_platforms
	} else {
		propertyMap["supported_platforms"] = ""
	}
	if template_type, ok := d.GetOk("template_type"); ok {
		propertyMap["template_type"] = template_type.(string)
	} else {
		propertyMap["template_type"] = ""
	}
	if template_sub_type, ok := d.GetOk("template_sub_type"); ok {
		propertyMap["template_sub_type"] = template_sub_type.(string)
	} else {
		propertyMap["template_sub_type"] = ""
	}
	if template_content_type, ok := d.GetOk("template_content_type"); ok {
		propertyMap["template_content_type"] = template_content_type.(string)
	} else {
		propertyMap["template_content_type"] = ""
	}
	content := fmt.Sprintf(`
	##template properties
	name=%s;
	description = %s;
	tags=%s;
	supportedPlatforms=%s;
	contentType=%s;
	templateType=%s;
	templateSubType=%s;
	##`, name, propertyMap["description"], propertyMap["tags"], propertyMap["supported_platforms"], propertyMap["template_content_type"], propertyMap["template_type"], propertyMap["template_sub_type"])
	fileContent = fmt.Sprintf("%s \n %s", content, fileContent)
	if !(strings.Contains(fileContent, "##template variables")) {
		fileContent = fmt.Sprintf("%s \n %s \n %s", fileContent, "##template variables", "##")
	}
	if !(strings.Contains(fileContent, "##template content")) {
		fileContent = fmt.Sprintf("%s \n %s \n %s", fileContent, "##template content", "##")
	}
	temp.Content = fileContent
	cont, err := dcnmClient.ValidateTemplateContent(TemplateURLS[dcnmClient.GetPlatform()]["Validate"], fileContent)

	if err != nil {
		return err
	}
	if !cont.Exists("status") && cont.S("reportItemType").String() == "ERROR" {
		return fmt.Errorf("Template Content is not valid.")
	}
	dURL := fmt.Sprintf(TemplateURLS[dcnmClient.GetPlatform()]["Create"], name)
	cont, err = dcnmClient.Save(dURL, &temp)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}

	d.SetId(name)

	log.Println("[DEBUG] End of Creating template")
	return resourceDCNMTemplateRead(d, m)
}

func resourceDCNMTemplateRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read Method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()

	cont, err := getTemplate(dcnmClient, dn)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}
	setTemplateAttribute(d, cont)
	d.SetId(dn)
	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func setTemplateAttribute(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	if cont.Exists("name") {

		d.Set("name", stripQuotes(cont.S("name").String()))
	}
	if cont.Exists("content") {

		d.Set("content", stripQuotes(cont.S("content").String()))

	}
	return d
}
func getTemplate(dcnmClient *client.Client, name string) (*container.Container, error) {
	cont, err := dcnmClient.GetviaURL(fmt.Sprintf(TemplateURLS[dcnmClient.GetPlatform()]["Common"], name))
	if err != nil {
		return cont, err
	}
	return cont, nil
}
func resourceDCNMTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining of Updating template")
	dcnmClient := m.(*client.Client)
	name := d.Get("name").(string)
	cont, _ := getTemplate(dcnmClient, name)
	var fileContent string

	fileContent = d.Get("content").(string)
	temp := models.TemplateUpdate{}
	propertyMap := make(map[string]string)
	if description, ok := d.GetOk("description"); ok {
		propertyMap["description"] = description.(string)
	} else {
		propertyMap["description"] = ""
	}
	if tags, ok := d.GetOk("tags"); ok {
		propertyMap["tags"] = tags.(string)
	} else {
		propertyMap["tags"] = ""
	}
	if supported_platforms, ok := d.GetOk("supported_platforms"); ok {
		platformList := make([]string, 0, 1)
		for _, val := range supported_platforms.([]interface{}) {
			platformList = append(platformList, val.(string))
		}
		supported_platforms := strings.Join(platformList, ",")
		propertyMap["supported_platforms"] = supported_platforms
	} else {
		propertyMap["supported_platforms"] = ""
	}
	if template_type, ok := d.GetOk("template_type"); ok {
		propertyMap["template_type"] = template_type.(string)
	} else {
		propertyMap["template_type"] = ""
	}
	if template_sub_type, ok := d.GetOk("template_sub_type"); ok {
		propertyMap["template_sub_type"] = template_sub_type.(string)
	} else {
		propertyMap["template_sub_type"] = ""
	}
	if template_content_type, ok := d.GetOk("template_content_type"); ok {
		propertyMap["template_content_type"] = template_content_type.(string)
	} else {
		propertyMap["template_content_type"] = ""
	}
	new2, _ := GetStringInBetweenTwoString(fileContent, "##template properties", "##")
	fileContent = strings.Replace(fileContent, "\\r\\n", "\n", -1)
	content := fmt.Sprintf("\n##template properties\nname=%s;\ndescription = %s;\ntags=%s;\nsupportedPlatforms=%s;\ncontentType=%s;\ntemplateType=%s;\ntemplateSubType=%s;\n##", name, propertyMap["description"], propertyMap["tags"], propertyMap["supported_platforms"], propertyMap["template_content_type"], propertyMap["template_type"], propertyMap["template_sub_type"])
	content = strings.Replace(content, "\n", "\n\t", -1)
	fileContent = strings.Replace(fileContent, new2, "", -1)
	fileContent = fmt.Sprintf("%s \n %s", content, fileContent)
	re1 := regexp.MustCompile(`\n\t`)
	fileContent = re1.ReplaceAllString(fileContent, "\n")
	fileContent = strings.Replace(fileContent, "\n", `\n`, -1)
	fileContent = strings.Replace(fileContent, "\\n\\n", "\r\n", -1)
	fileContent = strings.Replace(fileContent, "\\\"", "\"", -1)
	if !(strings.Contains(fileContent, "##template variables")) {
		fileContent = fmt.Sprintf("%s \n %s \n %s", fileContent, "##template variables", "##")
	}
	if !(strings.Contains(fileContent, "##template content")) {
		fileContent = fmt.Sprintf("%s \n %s \n %s", fileContent, "##template content", "##")
	}
	cont, err := dcnmClient.ValidateTemplateContent(TemplateURLS[dcnmClient.GetPlatform()]["Validate"], fileContent)

	if err != nil {
		return err
	}
	if !cont.Exists("status") && cont.S("reportItemType").String() == "ERROR" {
		return fmt.Errorf("Template Content is not valid.")
	}

	temp.Content = fileContent
	cont, err = dcnmClient.Update(fmt.Sprintf(TemplateURLS[dcnmClient.GetPlatform()]["Common"], name), &temp)
	if err != nil {
		return getErrorFromContainer(cont, err)
	}
	cont, _ = getTemplate(dcnmClient, name)

	d.SetId(name)
	return resourceDCNMTemplateRead(d, m)
}

func resourceDCNMTemplateDelete(d *schema.ResourceData, m interface{}) error {
	dcnmClient := m.(*client.Client)
	idList := strings.Split(d.Id(), "/")
	name := idList[0]
	cont, err := dcnmClient.Delete(fmt.Sprintf(TemplateURLS[dcnmClient.GetPlatform()]["Common"], name))
	if err != nil {
		return getErrorFromContainer(cont, err)
	}
	d.SetId("")
	return nil
}
func GetStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}

	newS := str[s:]
	e := strings.Index(newS[2:], endS)

	if e == -1 {
		return result, false
	}
	result = newS[:e+len(endS)+2]
	return result, true
}

func ValidateSpaces(data string) bool {
	out, _ := GetStringInBetweenTwoString(data, "##template content", "##")
	// Remove blank lines
	re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	out = re.ReplaceAllString(out, "")
	out1 := strings.Split(out, "\n")
	// Count space from left
	NoOfSpace := len(out1[0]) - len(strings.TrimLeft(out1[0], " "))
	for i := 1; i < len(out1); i++ {
		if NoOfSpace > 0 && NoOfSpace%2 == 0 {
			continue
		} else {
			return false
		}
	}
	return true

}
