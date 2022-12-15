package pass

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func passPasswordResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: passPasswordResourceWrite,
		UpdateContext: passPasswordResourceWrite,
		DeleteContext: passPasswordResourceDelete,
		ReadContext:   passPasswordResourceRead,

		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Full path where the pass data will be written",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secret password",
				Sensitive:   true,
			},

			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Additional key-value data",
				Sensitive:   true,
			},

			"yaml": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "YAML encoded data",
				Sensitive:   true,
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

func passPasswordResourceWrite(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pp := meta.(*passProvider)
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	path := pp.GetPath(d)
	log.Printf("writing resource secret at %s", path)

	password := d.Get("password").(string)
	data := d.Get("data").(map[string]interface{})
	yaml := d.Get("yaml").(string)

	if len(data) != 0 && yaml != "" {
		return diag.Errorf("can't set data and yaml at the same time")
	}

	value := password

	if yaml != "" {
		value += "\n---\n" + yaml
	} else {

		elems := make([]string, len(data))
		for k, v := range data {
			elems = append(elems, fmt.Sprintf("%s: %s", k, v))
		}

		value += strings.Join(elems, "\n")
	}

	container := secretContainer{value: value}
	err := pp.store.Set(context.Background(), path, container)

	if err != nil {
		return diag.Errorf("failed to write secret at %s: %s", path, err)
	}

	d.SetId(path)

	return nil
}

func passPasswordResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pp := meta.(*passProvider)
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	path := pp.GetPath(d)
	log.Printf("deleting resource secret at %s", path)

	err := pp.store.Remove(context.Background(), path)
	if err != nil {
		return diag.Errorf("failed to delete secret at %s: %s", path, err)
	}

	return nil
}

func passPasswordResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pp := meta.(*passProvider)
	pp.mutex.Lock()
	defer pp.mutex.Unlock()

	path := pp.GetPath(d)
	log.Printf("reading resource secret at %s", path)

	err := populateResourceData(d, pp, path, false)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
