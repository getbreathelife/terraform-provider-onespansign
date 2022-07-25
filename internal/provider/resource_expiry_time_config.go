package provider

import (
	"context"

	"github.com/getbreathelife/terraform-provider-onespansign/internal/helpers"
	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceExpiryTimeConfig() *schema.Resource {
	return &schema.Resource{
		Description: "OneSpan Sign account's expiry configurations.",

		CreateContext: resourceExpiryTimeConfigCreate,
		ReadContext:   resourceExpiryTimeConfigRead,
		UpdateContext: resourceExpiryTimeConfigUpdate,
		DeleteContext: resourceExpiryTimeConfigDelete,

		Schema: map[string]*schema.Schema{
			"default": {
				Description:      "Default expiry time for transactions in days. 0 for no limit.",
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
			"maximum": {
				Description:      "Maximum allowed value for expiry time for transactions in days. 0 for no limit.",
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceExpiryTimeConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "updating (replacing) resource instead of creating",
		Detail:   "This resource is a singleton. It only supports retrieval or replacement operations.",
	})

	diags = append(diags, resourceExpiryTimeConfigUpdate(ctx, d, meta)...)

	return diags
}

func resourceExpiryTimeConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	etc, err := c.GetExpiryTimeConfiguration()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	if err := d.Set("default", etc.Default); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("maximum", etc.Maximum); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceExpiryTimeConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	b := ossign.ExpiryTimeConfiguration{
		Default: helpers.GetJsonNumber(d.Get("default").(int64), 10),
		Maximum: helpers.GetJsonNumber(d.Get("maximum").(int64), 10),
	}

	if err := c.UpdateExpiryTimeConfiguration(b); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	tflog.Trace(ctx, "updated the account's expiry time configuration resource")

	return resourceExpiryTimeConfigRead(ctx, d, meta)
}

func resourceExpiryTimeConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "no deletion will take place",
			Detail:   "This resource is a singleton. It only supports retrieval or replacement operations.",
		},
	}
}
