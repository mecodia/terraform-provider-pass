package pass

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func passwordDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: passwordDataSourceRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full path from which a password will be read",
			},

			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Secret password",
			},

			"data": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Secret data",
			},

			"body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Body of the secret",
			},

			"full": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Entire secret contents",
			},
		},
	}
}

func passwordDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	path := d.Get("path").(string)
	pp := meta.(*passProvider)
	pp.mutex.Lock()
	defer pp.mutex.Unlock()
	d.SetId(path)
	err := populateResourceData(d, pp, path, true)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
