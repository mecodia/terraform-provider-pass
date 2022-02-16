package main

/* Bootstrap the plugin for Pass */

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/mecodia/terraform-provider-pass/pass"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pass.Provider,
	})
}
