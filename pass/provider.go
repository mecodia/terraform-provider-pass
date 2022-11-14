package pass

import (
	"context"
	"errors"
	"sync"

	"github.com/gopasspw/gopass/pkg/ctxutil"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type passProvider struct {
	store  *api.Gopass
	mutex  *sync.Mutex
	prefix string
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
		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Global key prefix for gopass provider. Must end with a forward slash. Can be used to specify a mount point.",
			},
		},
	}
}

func (pp *passProvider) GetPath(data *schema.ResourceData) string {
	if id := data.Id(); id != "" {
		return id
	}

	return pp.prefix + data.Get("path").(string)
}

func providerConfigureContext(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	store, err := api.New(ctxutil.WithShowParsing(ctx, false))

	if err != nil {
		return nil, diag.FromErr(err)
	}

	// normalize prefix value
	prefix := data.Get("prefix").(string)
	if len(prefix) > 0 && (prefix == "/" || prefix[len(prefix)-1:] != "/") {
		return nil, diag.FromErr(errors.New("the value of prefix must be a string with trailing forward slash"))
	}

	pp := &passProvider{
		store:  store,
		mutex:  &sync.Mutex{},
		prefix: prefix,
	}

	return pp, nil
}
