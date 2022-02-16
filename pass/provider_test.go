package pass

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var testProviderFactory = map[string]func() (*schema.Provider, error){
	"pass": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}
