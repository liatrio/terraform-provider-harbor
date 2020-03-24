package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func New() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"harbor_project":       resourceScaffolding(),
			"harbor_robot_account": resourceScaffolding(),
		},
		Schema: map[string]*schema.Schema{
			"harbor_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"harbor_admin_secret": {},
		},
		ConfigureFunc: func(data *schema.ResourceData) (i interface{}, err error) {

		},
	}
}
