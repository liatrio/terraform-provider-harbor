package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func New() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"harbor_project": resourceProject(),
			//"harbor_robot_account": resourceRobotAccount(),
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
		ConfigureFunc: configureHarborProvider,
	}
}

func configureHarborProvider(data *schema.ResourceData) (interface{}, error) {
	url := data.Get("url").(string)
	username := data.Get("username").(string)
	password := data.Get("password").(string)
	tlsInsecureSkipVerify := data.Get("tls_insecure_skip_verify").(bool)

	return harbor.NewClient(url, username, password, tlsInsecureSkipVerify)
}
