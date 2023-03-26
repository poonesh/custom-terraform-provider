package main

import (
	"context"
	"strconv"

	"custom_terraform_provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFood() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFoodCreate,
		ReadContext:   resourceFoodRead,
		UpdateContext: resourceFoodUpdate,
		DeleteContext: resourceFoodDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"origin": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceFoodCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := client.NewClient()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	origin := d.Get("origin").(string)

	food := client.Food{
		Name:   name,
		Origin: origin,
	}

	f, err := c.PostFood(ctx, &food)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(f.Id))
	resourceFoodRead(ctx, d, m)
	return diags
}

func resourceFoodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := client.NewClient()

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	foodID, err := strconv.Atoi(id)
	if err != nil {
		return diag.FromErr(err)
	}

	food, err := c.GetFood(ctx, foodID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", food.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("origin", food.Origin); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFoodUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := client.NewClient()

	id := d.Id()
	foodID, err := strconv.Atoi(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") || d.HasChange("origin") {
		name := d.Get("name").(string)
		origin := d.Get("origin").(string)

		food := client.Food{
			Name:   name,
			Origin: origin,
		}

		_, err = c.UpdateFood(ctx, &food, foodID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFoodRead(ctx, d, m)
}

func resourceFoodDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := client.NewClient()
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()
	foodID, err := strconv.Atoi(id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteFood(ctx, foodID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.

	d.SetId("")
	return diags
}
