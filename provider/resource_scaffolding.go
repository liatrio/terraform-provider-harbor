package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceScaffolding() *schema.Resource {
	return &schema.Resource{
		Create: resourceScaffoldingCreate,
		Read:   resourceScaffoldingRead,
		Update: resourceScaffoldingUpdate,
		Delete: resourceScaffoldingDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"robot_account_jwt": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceScaffoldingRead(d *schema.ResourceData, meta interface{}) error {

	d.Set("sample_attribute", "some string")
	d.SetId("${product_id}/${robot_account_id}")

	return nil
}

func resourceScaffoldingCreate(d *schema.ResourceData, meta interface{}) error {
	d.Get("")

	return resourceScaffoldingRead(d, meta)
}

func resourceScaffoldingUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceScaffoldingDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}
