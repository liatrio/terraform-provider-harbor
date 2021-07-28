package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/liatrio/terraform-provider-harbor/harbor"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProvider          *schema.Provider
	testAccProviderFactories map[string]func() (*schema.Provider, error)
	harborClient             *harbor.Client
)

var requiredEnvironmentVariables = []string{
	"HARBOR_USERNAME",
	"HARBOR_PASSWORD",
	"HARBOR_URL",
}

func init() {
	testAccProvider = New()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"harbor": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
	userAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", testAccProvider.TerraformVersion, meta.SDKVersionString())
	harborClient = harbor.NewClient(os.Getenv("HARBOR_URL"), os.Getenv("HARBOR_USERNAME"), os.Getenv("HARBOR_PASSWORD"), true, userAgent)
}

func TestProvider(t *testing.T) {
	if err := testAccProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
		if value := os.Getenv(requiredEnvironmentVariable); value == "" {
			t.Fatalf("%s must be set before running acceptance tests.", requiredEnvironmentVariable)
		}
	}
}
