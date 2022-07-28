package provider

import (
	"context"
	"time"

	"github.com/getbreathelife/terraform-provider-onespansign/internal/helpers"
	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceExpiryTimeConfig() *schema.Resource {
	return &schema.Resource{
		Description: `OneSpan Sign account's expiry configurations.
		
Please note that this resource is a singleton, which means that only one instance of this resource should
exist. Having multiple instances may produce unexpected result.`,

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

func validateFields(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	dft := d.Get("default").(int)
	mxm := d.Get("maximum").(int)

	if dft > mxm {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "validation error",
			Detail:   "the `default` value cannot be larger than the `maximum` value",
		})
	}

	return diags
}

// getExpiryTimeConfigStateChangeConf gets the configuration struct for the `WaitForState` functions.
// c is the OneSpan Sign API client instance, whereas e is the expected expiry time configuration state.
func getExpiryTimeConfigStateChangeConf(c *ossign.ApiClient, e ossign.ExpiryTimeConfiguration) resource.StateChangeConf {
	return resource.StateChangeConf{
		Delay:                     10 * time.Second,
		Pending:                   []string{"waiting"},
		Target:                    []string{"complete"},
		Timeout:                   3 * time.Minute,
		MinTimeout:                300 * time.Millisecond,
		ContinuousTargetOccurence: 5,
		Refresh: func() (result interface{}, state string, err error) {
			t, apiErr := c.GetExpiryTimeConfiguration()

			if apiErr != nil {
				return nil, "error", apiErr.GetError()
			}

			if t.Default.String() != e.Default.String() || t.Maximum.String() != e.Maximum.String() {
				return t, "waiting", nil
			}

			return t, "complete", nil
		},
	}
}

func resourceExpiryTimeConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ossign.ApiClient)

	diags := validateFields(ctx, d, meta)

	if len(diags) > 0 {
		return diags
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "updating (replacing) resource instead of creating",
		Detail:   "This resource is a singleton. It only supports retrieval or replacement operations.",
	})

	diags = append(diags, resourceExpiryTimeConfigUpdate(ctx, d, meta)...)

	d.SetId(c.ClientId)

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

	diags := validateFields(ctx, d, meta)

	if len(diags) > 0 {
		return diags
	}

	b := ossign.ExpiryTimeConfiguration{
		Default: helpers.GetJsonNumber(int64(d.Get("default").(int)), 10),
		Maximum: helpers.GetJsonNumber(int64(d.Get("maximum").(int)), 10),
	}

	if err := c.UpdateExpiryTimeConfiguration(b); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	tflog.Trace(ctx, "waiting for the account's expiry time configuration resource to be updated...")

	scc := getExpiryTimeConfigStateChangeConf(c, b)
	_, err := scc.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
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
