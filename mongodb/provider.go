package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jlpedrosa/terraform-provider-mongodb/mongodb/internal/provider"
)

// ProviderFunc Hook point to terraform plugin module entry point
func ProviderFunc() terraform.ResourceProvider {
	return provider.Provider()
}
