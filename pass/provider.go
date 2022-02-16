package pass

import (
	"context"
	"sync"

	"github.com/gopasspw/gopass/pkg/ctxutil"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type passProvider struct {
	store *api.Gopass
	mutex *sync.Mutex
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigureContext,
		DataSourcesMap: map[string]*schema.Resource{
			"pass_password": passwordDataSource(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"pass_password": passPasswordResource(),
		},
	}
}

func providerConfigureContext(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	store, err := api.New(ctxutil.WithShowParsing(ctx, false))

	if err != nil {
		return nil, diag.FromErr(err)
	}

	pp := &passProvider{
		store: store,
		mutex: &sync.Mutex{},
	}

	return pp, nil
}
