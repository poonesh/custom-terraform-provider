package main

import (
	"context"
	"strconv"
	"time"

	"custom_terraform_provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFoodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := client.NewClient()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	foodID := d.Get("id").(int)
	food, err := client.GetFood(ctx, foodID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", food.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("origin", food.Origin); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func dataSourceFood() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFoodRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
