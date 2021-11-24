package acctest

import (
	"os"
	"sync"

	"github.com/CiscoDevNet/terraform-provider-dcnm/dcnm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"testing"
)

var testAccProviders map[string]*schema.Provider

var testAccProviderFactories map[string]func() (*schema.Provider, error)

var testAccProvider *schema.Provider

var testAccProviderConfigure sync.Once

func init() {
	testAccProvider = dcnm.Provider()

	testAccProviders = map[string]*schema.Provider{
		"dcnm": testAccProvider,
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"dcnm": func() (*schema.Provider, error) { return dcnm.Provider(), nil },
	}
}

func testAccProviderFactoriesInit(provider **schema.Provider, providerNames string) map[string]func() (*schema.Provider, error) {
	var factories = make(map[string]func() (*schema.Provider, error), len(providerNames))

	p := dcnm.Provider()

	factories[providerNames] = func() (*schema.Provider, error) {
		return p, nil
	}

	if provider != nil {
		*provider = p
	}

	return factories
}

func testAccProviderFactoriesInternal(provider **schema.Provider) map[string]func() (*schema.Provider, error) {
	return testAccProviderFactoriesInit(provider, "dcnm")
}

func TestProvider(t *testing.T) {
	if err := dcnm.Provider().InternalValidate(); err != nil {
		t.Fatalf("err : %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = dcnm.Provider()
}

func testAccPreCheck(t *testing.T) {
	// We will use this function later on to make sure our test environment is valid.
	// For example, you can make sure here that some environment variables are set.
	if v := os.Getenv("DCNM_USERNAME"); v == "" {
		t.Fatal("DCNM_USERNAME env variable must be set for acceptance tests")
	}
	if v := os.Getenv("DCNM_PASSWORD"); v == "" {
		t.Fatal("DCNM_PASSWORD env variable must be set for acceptance tests")
	}
	if v := os.Getenv("DCNM_URL"); v == "" {
		t.Fatal("DCNM_URL env variable must be set for acceptance tests")
	}
}
