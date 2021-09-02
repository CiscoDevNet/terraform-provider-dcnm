package client

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
)

const authPayload = `{
	"expirationTime": %d
}`

const ndAuthPayload = `{
	"userName": "%s",
  	"userPasswd": "%s",
  	"domain": "local"
}`

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	authToken  *auth
	username   string
	password   string
	insecure   bool
	proxyUrl   string
	expiry     int64
	domain     string
	platform   string
}

var clientImpl *Client

type Option func(*Client)

func Insecure(insecure bool) Option {
	return func(client *Client) {
		client.insecure = insecure
	}
}

func ProxyUrl(pUrl string) Option {
	return func(client *Client) {
		client.proxyUrl = pUrl
	}
}

func Platform(platform string) Option {
	return func(client *Client) {
		client.platform = platform
	}
}

func (c *Client) GetPlatform() string {
	return c.platform
}

func (c *Client) useInsecureHTTPClient(insecure bool) *http.Transport {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       insecure,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS12,
		},
	}

	return transport

}

func (c *Client) configProxy(transport *http.Transport) *http.Transport {
	pUrl, err := url.Parse(c.proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	transport.Proxy = http.ProxyURL(pUrl)
	return transport

}

func initClient(clientURL, username, password string, expiry int64, options ...Option) *Client {
	baseURL, err := url.Parse(clientURL)
	if err != nil {
		log.Fatal(err)
	}

	client := &Client{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		expiry:     expiry,
		insecure:   true,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	transport := client.useInsecureHTTPClient(client.insecure)
	if client.proxyUrl != "" {
		transport = client.configProxy(transport)
	}

	client.httpClient = &http.Client{
		Transport: transport,
	}
	return client
}

func GetClient(clientURL, username, password string, expiry int64, options ...Option) *Client {
	if clientImpl == nil {
		clientImpl = initClient(clientURL, username, password, expiry, options...)
	}
	return clientImpl
}

func (c *Client) MakeRequest(method, path string, body *container.Container, authenticated bool) (*http.Request, error) {

	if c.platform == "nd" && authenticated && !models.IsService(path) && !models.IsTemplate(path) {
		path = fmt.Sprint("/appcenter/cisco/ndfc/api/v1/lan-fabric", path)
	}

	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	reqURL := c.baseURL.ResolveReference(url)
	log.Println("req", reqURL)
	log.Println("req", reqURL.String())

	var req *http.Request
	if body == nil {
		req, err = http.NewRequest(method, reqURL.String(), nil)
	} else {
		req, err = http.NewRequest(method, reqURL.String(), bytes.NewBuffer(body.Bytes()))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if authenticated {
		log.Println("HTTP request ", method, path, req)
	}
	if authenticated {
		req, err = c.injectAuthenticationHeader(req, path)
		if err != nil {
			return req, err
		}
	}
	if authenticated {
		log.Println("HTTP request after injection ", method, path, req)
	}
	return req, nil
}
func (c *Client) makeRequestForText(method, path string, body string, authenticated bool) (*http.Request, error) {

	if c.platform == "nd" && authenticated && !models.IsService(path) {
		path = fmt.Sprint("/appcenter/cisco/ndfc/api/v1", path)
	}

	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	reqURL := c.baseURL.ResolveReference(url)

	var req *http.Request
	if body == "" {
		req, err = http.NewRequest(method, reqURL.String(), nil)
	} else {
		req, err = http.NewRequest(method, reqURL.String(), strings.NewReader(body))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")
	if authenticated {
		log.Println("HTTP request ", method, path, req)
	}
	if authenticated {
		req, err = c.injectAuthenticationHeader(req, path)
		if err != nil {
			return req, err
		}
	}
	if authenticated {
		log.Println("HTTP request after injection ", method, path, req)
	}
	return req, nil
}
func (c *Client) makeRequestForCred(method, path string, body []byte, authenticated bool) (*http.Request, error) {
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	reqURL := c.baseURL.ResolveReference(url)

	var req *http.Request
	req, err = http.NewRequest(method, reqURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	log.Println("HTTP request ", method, path, req)

	if authenticated {
		req, err = c.injectAuthenticationHeader(req, path)
		if err != nil {
			return req, err
		}
	}
	log.Println("HTTP request after injection ", method, path, req)
	return req, nil
}

func (c *Client) authenticate() error {
	method := "POST"

	if c.platform == "nd" {
		path := "/login"

		body, err := container.ParseJSON([]byte(fmt.Sprintf(ndAuthPayload, c.username, c.password)))
		if err != nil {
			return err
		}

		req, err := c.MakeRequest(method, path, body, false)
		if err != nil {
			return err
		}

		obj, resp, err := c.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 401 {
			return fmt.Errorf("Invalid username or password")
		}

		req.Header.Set("Content-Type", "application/json")
		token := models.StripQuotes(obj.S("token").String())

		if c.authToken == nil {
			c.authToken = &auth{}
		}
		c.authToken.token = token
		c.authToken.calculateExpiry(1200)

	} else {
		path := "/rest/logon"

		body, err := container.ParseJSON([]byte(fmt.Sprintf(authPayload, c.expiry)))
		if err != nil {
			return err
		}

		req, err := c.MakeRequest(method, path, body, false)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", getBasicAuth(c.username, c.password)))

		obj, resp, err := c.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 500 {
			return fmt.Errorf("Invalid username or password")
		}

		token := models.StripQuotes(obj.S("Dcnm-Token").String())

		if c.authToken == nil {
			c.authToken = &auth{}
		}
		c.authToken.token = token
		c.authToken.calculateExpiry(c.expiry)
	}
	return nil
}

func (c *Client) Do(req *http.Request) (*container.Container, *http.Response, error) {
	log.Println("[DEBUG] Begining Do method ", req.URL.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	log.Println("[DEBUG] HTTP Request ", req.Method, req.URL.String())
	log.Println("[DEBUG] HTTP Response ", resp.StatusCode, resp)

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	bodystrings := string(bodybytes)
	resp.Body.Close()
	log.Println("[DEBUG] HTTP Response unique string ", req.Method, req.URL.String(), bodystrings)

	obj, err := container.ParseJSON(bodybytes)
	if err != nil && resp.StatusCode != 200 {
		return nil, resp, fmt.Errorf(bodystrings)
	}

	log.Println("[DEBUG] Ending Do method ", req.URL.String())
	return obj, resp, nil
}

func getBasicAuth(username, password string) string {
	authString := fmt.Sprintf("%s:%s", username, password)

	encodedString := b64.StdEncoding.EncodeToString([]byte(authString))

	return encodedString
}
