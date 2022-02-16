package main

/* Bootstrap the plugin for Pass */

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/mecodia/terraform-provider-pass/pass"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pass.Provider,
	})
}
