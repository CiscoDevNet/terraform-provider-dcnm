package dcnm

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDCNMVRFCustom() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMVRFCustomCreate,
		Read:   resourceDCNMVRFCustomRead,
		Update: resourceDCNMVRFCustomUpdate,
		Delete: resourceDCNMVRFCustomDelete,

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"segment_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_VRF_Universal",
			},
			"extension_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_VRF_Extension_Universal",
			},

			"service_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": &schema.Schema{
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
		},
	}
}

func resourceDCNMVRFCustomRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	fabricName := d.Get("fabric_name").(string)

	cont, err := getRemoteVRF(dcnmClient, fabricName, dn)
	if err != nil {
		return err
	}

	setVRFCustomAttributes(d, cont)

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func setVRFCustomAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	if cont.Exists("fabric") {
		d.Set("fabric_name", stripQuotes(cont.S("fabric").String()))
	}
	d.Set("name", stripQuotes(cont.S("vrfName").String()))
	d.Set("template", stripQuotes(cont.S("vrfTemplate").String()))
	d.Set("extension_template", stripQuotes(cont.S("vrfExtensionTemplate").String()))
	d.Set("segment_id", stripQuotes(cont.S("vrfId").String()))

	if cont.Exists("serviceVrfTemplate") && stripQuotes(cont.S("serviceVrfTemplate").String()) != "null" {
		d.Set("service_template", stripQuotes(cont.S("serviceVrfTemplate").String()))
	}
	if cont.Exists("source") && stripQuotes(cont.S("source").String()) != "null" {
		d.Set("source", stripQuotes(cont.S("source").String()))
	}

	var strByte []byte
	if cont.Exists("vrfTemplateConfig") {
		strJson := models.G(cont, "vrfTemplateConfig")
		strByte = []byte(strJson)
		var vrfTemplateConfig map[string]interface{}
		json.Unmarshal(strByte, &vrfTemplateConfig)
		props, ok := d.GetOk("template_props")

		map2 := make(map[string]interface{})
		for k := range props.(map[string]interface{}) {
			map2[k] = vrfTemplateConfig[k]
		}
		if !ok {
			d.Set("template_props", vrfTemplateConfig)
		} else {

			d.Set("template_props", map2)
		}
	}

	d.SetId(stripQuotes(cont.S("vrfName").String()))
	return d
}

func resourceDCNMVRFCustomCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	vrf := models.VRF{}
	vrf.Name = d.Get("name").(string)
	vrf.Fabric = d.Get("fabric_name").(string)
	vrf.Template = d.Get("template").(string)
	vrf.ExtensionTemplate = d.Get("extension_template").(string)

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); !ok {
			return fmt.Errorf("attachments must be configured if deploy=true")
		}
	}

	if segmentId, ok := d.GetOk("segment_id"); ok {
		vrf.Id = segmentId.(string)
	} else {
		//request to get the next vrf segment id
		if dcnmClient.GetPlatform() == "nd" {
			cont, err := dcnmClient.GetviaURL(fmt.Sprintf("/rest/top-down/fabrics/%s/vrfinfo", vrf.Fabric))
			if err != nil {
				return err
			}
			vrf.Id = cont.S("l3vni").String()
		} else {
			cont, err := dcnmClient.GetSegID(fmt.Sprintf("/rest/managed-pool/fabrics/%s/partitions/ids", vrf.Fabric))
			if err != nil {
				return err
			}
			vrf.Id = cont.S("partitionSegmentId").String()
		}
	}

	if srcTemp, ok := d.GetOk("service_template"); ok {
		vrf.ServiceVRFTemplate = srcTemp.(string)
	}

	if src, ok := d.GetOk("source"); ok {
		vrf.Source = src.(string)
	}
	vrfConfig := d.Get("template_props").(map[string]interface{})

	confStr, err := json.Marshal(vrfConfig)
	if err != nil {
		return err
	}
	vrf.Config = string(confStr)

	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs", vrf.Fabric)
	_, err = dcnmClient.Save(durl, &vrf)
	if err != nil {
		return err
	}

	d.SetId(vrf.Name)
	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMVRFRead(d, m)
}

func resourceDCNMVRFCustomUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	dcnmClient := m.(*client.Client)

	vrf := models.VRF{}
	vrf.Name = d.Get("name").(string)
	vrf.Fabric = d.Get("fabric_name").(string)
	vrf.Template = d.Get("template").(string)
	vrf.ExtensionTemplate = d.Get("extension_template").(string)
	vrf.Id = d.Get("segment_id").(string)

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); !ok {
			return fmt.Errorf("attachments must be configured if deploy=true")
		}
	}

	if srcTemp, ok := d.GetOk("service_template"); ok {
		vrf.ServiceVRFTemplate = srcTemp.(string)
	}

	if src, ok := d.GetOk("source"); ok {
		vrf.Source = src.(string)
	}

	vrfConfig := d.Get("template_props").(map[string]interface{})

	confStr, err := json.Marshal(vrfConfig)
	if err != nil {
		return err
	}
	vrf.Config = string(confStr)

	dn := d.Id()
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/%s", vrf.Fabric, dn)
	_, err = dcnmClient.Update(durl, &vrf)
	if err != nil {
		return err
	}
	d.SetId(vrf.Name)

	//VRF Attachment
	if d.HasChange("deploy") && d.Get("deploy").(bool) == false {
		return fmt.Errorf("Deployed VRF can not be undeployed")
	}

	d.SetId(vrf.Name)
	log.Println("[DEBUG] End of Update method ", d.Id())
	return resourceDCNMVRFRead(d, m)
}

func resourceDCNMVRFCustomDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())
	dcnmClient := m.(*client.Client)

	dn := d.Id()
	fabricName := d.Get("fabric_name").(string)

	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/vrfs/%s", fabricName, dn)
	_, err := dcnmClient.Delete(durl)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}
