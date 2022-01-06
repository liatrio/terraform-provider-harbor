package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func New() *schema.Provider {
	provider := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"harbor_project":       resourceProject(),
			"harbor_robot_account": resourceRobotAccount(),
			"harbor_webhook":       resourceWebhook(),
			"harbor_label":         resourceLabel(),
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_URL", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_PASSWORD", nil),
			},
			"tls_insecure_skip_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		url := data.Get("url").(string)
		username := data.Get("username").(string)
		password := data.Get("password").(string)
		tlsInsecureSkipVerify := data.Get("tls_insecure_skip_verify").(bool)

		userAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", provider.TerraformVersion, meta.SDKVersionString())
		client := harbor.NewClient(url, username, password, tlsInsecureSkipVerify, userAgent)

		return client, diag.Diagnostics{}
	}

	return provider
}
