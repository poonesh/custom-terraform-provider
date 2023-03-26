package main

import (
	"context"
	"custom_terraform_provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// Provider Configuration such as API keys or username/password goes here
		},
		ResourcesMap: map[string]*schema.Resource{
			// Define resources here
			"food": resourceFood(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// Define data sources here
			"food": dataSourceFood(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Necessary client configuration such as authentication goes here
	c := client.NewClient()
	return c, diags
}
