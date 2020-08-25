package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func handleNotFoundError(err error, data *schema.ResourceData) error {
	if harbor.ErrorIs404(err) {
		log.Printf("[WARN] Removing resource with id %s from state as it no longer exists", data.Id())
		data.SetId("")

		return nil
	}

	return err
}
