package dcnm

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDCNMRest() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMRestCreate,
		Update: resourceDCNMRestUpdate,
		Read:   resourceDCNMRestRead,
		Delete: resourceDCNMRestDelete,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GET",
					"PUT",
					"POST",
					"DELETE",
				}, false),
			},

			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"payload_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"json", "text",
				}, false),
				Default: "json",
			},
		},
	}
}

func resourceDCNMRestCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)
	path := d.Get("path").(string)
	payload := d.Get("payload").(string)
	if payload == "" {
		return fmt.Errorf("payload should be given when method is POST")
	}

	var op string

	if method, ok := d.GetOk("method"); ok {
		op = method.(string)
	} else {
		op = "POST"
	}

	if d.Get("payload_type").(string) == "json" {
		_, err := makeAndDoRest(dcnmClient, path, op, payload)
		if err != nil {
			return err
		}
	} else {
		_, err := makeAndDoRestForText(dcnmClient, path, op, payload)
		if err != nil {
			return err
		}
	}

	d.SetId(path)

	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMRestRead(d, m)
}

func resourceDCNMRestUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	dcnmClient := m.(*client.Client)
	path := d.Get("path").(string)
	payload := d.Get("payload").(string)
	if payload == "" {
		return fmt.Errorf("payload should be given when method is PUT")
	}

	var op string

	if method, ok := d.GetOk("method"); ok {
		op = method.(string)
	} else {
		op = "PUT"
	}

	if d.Get("payload_type").(string) == "json" {
		_, err := makeAndDoRest(dcnmClient, path, op, payload)
		if err != nil {
			return err
		}
	} else {
		_, err := makeAndDoRestForText(dcnmClient, path, op, payload)
		if err != nil {
			return err
		}
	}

	d.SetId(path)

	log.Println("[DEBUG] End of Update method ", d.Id())
	return resourceDCNMRestRead(d, m)
}

func resourceDCNMRestRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDCNMRestDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())

	dcnmClient := m.(*client.Client)
	path := d.Get("path").(string)
	payload := d.Get("payload").(string)

	var op string

	if method, ok := d.GetOk("method"); ok {
		op = method.(string)
	} else {
		op = "DELETE"
	}

	if d.Get("payload_type").(string) == "json" {
		_, err := makeAndDoRest(dcnmClient, path, op, payload)
		if err != nil {
			return err
		}
	} else {
		_, err := makeAndDoRestForText(dcnmClient, path, op, payload)
		if err != nil {
			return err
		}
	}

	d.SetId("")

	log.Println("[DEBUG] End of Delete method ")
	return nil
}

func makeAndDoRest(client *client.Client, path, op, payload string) (*container.Container, error) {

	jsonPayload, err := container.ParseJSON([]byte(payload))
	if err != nil {
		return nil, err
	}

	var req *http.Request

	if strings.HasPrefix(path, "/appcenter/cisco/ndfc") {
		req, err = client.MakeRestNDRequest(op, path, jsonPayload, true)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = client.MakeRequest(op, path, jsonPayload, true)
		if err != nil {
			return nil, err
		}
	}

	respCont, resp, err := client.Do(req)
	if err != nil {
		return nil, checkerrorsRest(respCont, resp)
	}

	return respCont, checkerrorsRest(respCont, resp)
}

func makeAndDoRestForText(client *client.Client, path, op, content string) (*container.Container, error) {
	req, err := client.MakeRequestForText(op, path, content, true)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	cont, resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkerrorsRest(cont, resp)
}

func checkerrorsRest(cont *container.Container, resp *http.Response) error {

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if cont != nil {
		return fmt.Errorf("%d Error : %s", resp.StatusCode, cont.S("message").String())
	}

	return fmt.Errorf("%d Error : %s", resp.StatusCode, resp.Status)
}
