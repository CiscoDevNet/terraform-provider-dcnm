package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
)

func (c *Client) GetviaURL(endpoint string) (*container.Container, error) {
	req, err := c.MakeRequest("GET", endpoint, nil, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if cont == nil {
		return nil, errors.New("Empty response body")
	}
	return cont, checkforerrors(cont, resp)
}

func (c *Client) Save(endpoint string, obj models.Model) (*container.Container, error) {
	jsonPayload, err := c.prepareModel(obj)
	if err != nil {
		return nil, err
	}

	req, err := c.MakeRequest("POST", endpoint, jsonPayload, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}
func (c *Client) SaveDeploy(endpoint string, policyIds string) (*container.Container, error) {
	contList := container.New()
	contList.Array()
	contList.ArrayAppend(policyIds)
	req, err := c.MakeRequest("POST", endpoint, contList, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}
func (c *Client) SaveForAttachment(endpoint string, obj models.Model) (*container.Container, error) {
	contList := container.New()
	contList.Array()

	jsonPayload, err := c.prepareModel(obj)
	if err != nil {
		return nil, err
	}
	contList.ArrayAppend(jsonPayload.Data())

	req, err := c.MakeRequest("POST", endpoint, contList, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}

func (c *Client) UpdateCred(endpoint string, body []byte) (*container.Container, error) {
	req, err := c.makeRequestForCred("POST", endpoint, body, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}

func (c *Client) GetSegID(endpoint string) (*container.Container, error) {
	req, err := c.MakeRequest("POST", endpoint, nil, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}

func (c *Client) Update(endpoint string, obj models.Model) (*container.Container, error) {
	jsonPayload, err := c.prepareModel(obj)
	if err != nil {
		return nil, err
	}

	req, err := c.MakeRequest("PUT", endpoint, jsonPayload, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}

func (c *Client) Delete(endpoint string) (*container.Container, error) {
	req, err := c.MakeRequest("DELETE", endpoint, nil, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}

func (c *Client) DeleteWithPayload(endpoint string, obj models.Model) (*container.Container, error) {
	contList := container.New()
	contList.Array()

	jsonPayload, err := c.prepareModel(obj)
	if err != nil {
		return nil, err
	}
	contList.ArrayAppend(jsonPayload.Data())

	req, err := c.MakeRequest("DELETE", endpoint, contList, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return cont, checkforerrors(cont, resp)
}

func (c *Client) SaveAndDeploy(endpoint string) (*container.Container, error) {
	req, err := c.MakeRequest("POST", endpoint, nil, true)
	if err != nil {
		return nil, err
	}

	cont, resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return cont, checkforerrors(cont, resp)
}

func checkforerrors(cont *container.Container, resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("%d Error : %s", resp.StatusCode, cont.S("message").String())
}

func (c *Client) prepareModel(obj models.Model) (*container.Container, error) {
	con, err := obj.ToMap()
	if err != nil {
		return nil, err
	}

	payload := &container.Container{}

	for key, value := range con {
		payload.Set(value, key)
	}
	return payload, nil
}
