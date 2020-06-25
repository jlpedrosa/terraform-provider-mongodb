package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/jlpedrosa/terraform-provider-mongodb/mongodb"
)

func main() {
	plugin.Serve(&plugin.ServeOpts {
		ProviderFunc: mongodb.ProviderFunc,
	})
}
