package main

import (
	"flag"

	"github.com/SAP-cloud-infrastructure/terraform-provider-sci/sci"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

const providerAddr = "registry.terraform.io/SAP-cloud-infrastructure/sci"

func main() {
	// added debugMode to enable debugging for provider per https://www.terraform.io/plugin/sdkv2/debugging
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug:        debugMode,
		ProviderAddr: providerAddr,
		ProviderFunc: sci.Provider,
	})
}
