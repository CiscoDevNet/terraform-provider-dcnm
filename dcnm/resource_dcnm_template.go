package dcnm

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
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
		"Create":   "/appcenter/cisco/dcnm/api/v1/configtemplate/rest/config/templates/template?templateName=%s",
		"Common":   "/appcenter/cisco/dcnm/api/v1/configtemplate/rest/config/templates/%s",
		"Validate": "/configtemplate/rest/config/templates/validate",
	},
}

func CompareDiffs(old, new string, d *schema.ResourceData) bool {
	var old1 string
	if runtime.GOOS == "windows" {
		old1 = strings.Replace(old, "\\r\\n", "\n", -1)
	} else {
		old1 = strings.Replace(old, "\\n", "\n", -1)
	}

	re, err := regexp.Compile("((##template properties)(.*?)(##))")
	if err != nil {
		return false
	}
	old1 = strings.Replace(old1, "\\\"", "\"", -1)
	old2 := re.FindString(old1)
	old1 = strings.Replace(old1, old2, "", -1)
	old1 = strings.Replace(old1, "\\n", "\n", -1)
	re1 := regexp.MustCompile(`\r?\n`)
	new2, _ := GetStringInBetweenTwoString(new, "##template properties", "##")
	new1 := strings.Replace(new, new2, "", -1)
	new1 = re1.ReplaceAllString(new1, "\n")
	new2 = re1.ReplaceAllString(new2, `\n`)
	serverMap, _ := getTemplateProps(old2)
	localMap, _ := getTemplateProps(new2)
	log.Println("LocalMAPpp", localMap)
	for key, val := range localMap {
		if w, ok := serverMap[key]; ok {

			if val != w {
				return false
			}
		}
	}
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

	cont, err := dcnmClient.ValidateTemplateContent(TemplateURLS[dcnmClient.GetPlatform()]["Validate"], fileContent)
	log.Println("[DEBUG] container", cont)

	if err != nil {
		return err
	}
	if !cont.Exists("status") {
		return fmt.Errorf("Template Content is not valid.")
	}
	temp := models.Template{}
	temp.Name = name
	temp.Content = fileContent
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
	// }
	cont, err := dcnmClient.ValidateTemplateContent(TemplateURLS[dcnmClient.GetPlatform()]["Validate"], fileContent)
	log.Println("[DEBUG] container", cont)

	if err != nil {
		return err
	}
	if !cont.Exists("status") {
		return fmt.Errorf("Template Content is not valid.")
	}

	temp := models.TemplateUpdate{}
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
