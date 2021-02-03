package dcnm

import (
	"fmt"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DCNM_USERNAME", nil),
				Description: "Username for the DCNM account",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DCNM_PASSWORD", nil),
				Description: "Password for the DCNM account",
			},

			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DCNM_URL", nil),
				Description: "URL for the DCNM server",
			},

			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow insecure HTTPS client",
			},

			"proxy_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Proxy server URL for DCNM",
			},

			"expiry": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     900000,
				Description: "Expiration time in miliseconds for DCNM server",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dcnm_vrf":       resourceDCNMVRF(),
			"dcnm_inventory": resourceDCNMInventroy(),
			"dcnm_network":   resourceDCNMNetwork(),
			"dcnm_interface": resourceDCNMInterface(),
			"dcnm_rest":      resourceDCNMRest(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"dcnm_vrf":       datasourceDCNMVRF(),
			"dcnm_inventory": datasourceDCNMInventory(),
			"dcnm_network":   datasourceDCNMNetwork(),
			"dcnm_interface": datasourceDCNMInterface(),
		},

		ConfigureFunc: configClient,
	}
}

func configClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		URL:        d.Get("url").(string),
		IsInsecure: d.Get("insecure").(bool),
		ProxyURL:   d.Get("proxy_url").(string),
		Expiry:     d.Get("expiry").(int),
	}

	if err := config.Valid(); err != nil {
		return nil, err
	}

	return config.getClient(), nil
}

func (c Config) Valid() error {

	if c.Username == "" {
		return fmt.Errorf("Username must be provided for the DCNM provider")
	}

	if c.Password == "" {
		return fmt.Errorf("Password must be provided for the DCNM provider")
	}

	if c.URL == "" {
		return fmt.Errorf("The URL must be provided for the DCNM provider")
	}

	return nil
}

func (c Config) getClient() interface{} {
	return client.GetClient(c.URL, c.Username, c.Password, int64(c.Expiry), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyURL))
}

type Config struct {
	Username   string
	Password   string
	URL        string
	IsInsecure bool
	ProxyURL   string
	Expiry     int
}
