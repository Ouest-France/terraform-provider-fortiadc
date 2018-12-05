package main

import (
	"github.com/Ouest-France/terraform-provider-fortiadc/fortiadc"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return fortiadc.Provider()
		},
	})
}
