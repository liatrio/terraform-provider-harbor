package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/liatrio/terraform-provider-harbor/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: provider.New})
}
