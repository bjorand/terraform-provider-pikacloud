package main

import (
	"github.com/bjorand/terraform-provider-pikacloud/pikacloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pikacloud.Provider,
	})
}
