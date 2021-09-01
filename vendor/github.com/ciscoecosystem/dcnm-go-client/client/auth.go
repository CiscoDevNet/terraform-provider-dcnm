package client

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type auth struct {
	token  string
	expiry time.Time
}

func (au *auth) estimateExpiryTime() int64 {
	return time.Now().Unix() + 3
}

func (au *auth) isValid() bool {
	if au.token != "" && au.expiry.Unix() > au.estimateExpiryTime() {
		return true
	}
	return false
}

func (au *auth) calculateExpiry(expiry int64) {
	au.expiry = time.Unix((time.Now().Unix() + expiry/1000), 0)
}

func (client *Client) injectAuthenticationHeader(req *http.Request, path string) (*http.Request, error) {
	log.Println("[DEBUG] Begin Injection")
	if client.authToken == nil || !client.authToken.isValid() {
		err := client.authenticate()
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Content-Type", "application/json")
	if client.platform == "nd" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.authToken.token))
	} else {
		req.Header.Set("dcnm-token", client.authToken.token)
	}
	return req, nil
}
