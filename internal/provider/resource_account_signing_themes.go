package provider

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccountSigningThemes() *schema.Resource {
	return &schema.Resource{
		Description: `OneSpan Sign account's customized signing themes.
		
Please note that this resource is a singleton, which means that only one instance of this resource should
exist. Having multiple instances may produce unexpected result.`,

		CreateContext: resourceAccountSigningThemesCreate,
		ReadContext:   resourceAccountSigningThemesRead,
		UpdateContext: resourceAccountSigningThemesUpdate,
		DeleteContext: resourceAccountSigningThemesDelete,

		Schema: map[string]*schema.Schema{
			"theme": {
				Description: "Customized signing theme for the account.",
				Type:        schema.TypeSet,
				Required:    true,

				// The current API behaviour is that we are able to create multiple themes for the account.
				// However, only the first theme will be used for the signing ceremony.
				// This constraint is added so that the user does not expect to have multiple themes.
				MaxItems: 1,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the theme.",
							Required:    true,
							Type:        schema.TypeString,
						},
						"primary": {
							Description:      "Primary color hex code.",
							Required:         true,
							Type:             schema.TypeString,
							ValidateDiagFunc: validateColorHex,
						},
						"success": {
							Description:      "Success notification color hex code.",
							Required:         true,
							Type:             schema.TypeString,
							ValidateDiagFunc: validateColorHex,
						},
						"warning": {
							Description:      "Warning notification color hex code.",
							Required:         true,
							Type:             schema.TypeString,
							ValidateDiagFunc: validateColorHex,
						},
						"error": {
							Description:      "Error notification color hex code.",
							Required:         true,
							Type:             schema.TypeString,
							ValidateDiagFunc: validateColorHex,
						},
						"info": {
							Description:      "Info notification color hex code.",
							Required:         true,
							Type:             schema.TypeString,
							ValidateDiagFunc: validateColorHex,
						},
						"signature_button": {
							Description:      "Color hex code for the required signature buttons.",
							Required:         true,
							Type:             schema.TypeString,
							ValidateDiagFunc: validateColorHex,
						},
						"optional_signature_button": {
							Description:      "Color hex code for the optional signature buttons.",
							Required:         true,
							Type:             schema.TypeString,
							ValidateDiagFunc: validateColorHex,
						},
					},
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func validateColorHex(v interface{}, p cty.Path) diag.Diagnostics {
	r := regexp.MustCompile("^#[0-9A-F]{6}$")
	return validation.ToDiagFunc(validation.StringMatch(r, "must be a valid 6-symbol color hex code (e.g. #FFFFFF)"))(v, p)
}

func flattenAccountSigningTheme(name string, t ossign.SigningTheme) interface{} {
	ft := make(map[string]interface{}, 8)

	ft["name"] = name
	ft["primary"] = t.Primary
	ft["success"] = t.Success
	ft["warning"] = t.Warning
	ft["error"] = t.Error
	ft["info"] = t.Info
	ft["signature_button"] = t.SignatureButton
	ft["optional_signature_button"] = t.OptionalSignatureButton

	return ft
}

func buildAccountSigningThemes(d *schema.ResourceData) map[string]ossign.SigningTheme {
	ts := d.Get("theme").(*schema.Set).List()

	r := make(map[string]ossign.SigningTheme, len(ts))

	for _, item := range ts {
		i := item.(map[string]interface{})

		r[i["name"].(string)] = ossign.SigningTheme{
			Primary:                 i["primary"].(string),
			Success:                 i["success"].(string),
			Warning:                 i["warning"].(string),
			Error:                   i["error"].(string),
			Info:                    i["info"].(string),
			SignatureButton:         i["signature_button"].(string),
			OptionalSignatureButton: i["optional_signature_button"].(string),
		}
	}

	return r
}

// getSigningThemeStateChangeConf gets the configuration struct for the `WaitForState` functions.
// c is the OneSpan Sign API client instance, whereas e is the expected map of signing themes state.
func getSigningThemeStateChangeConf(c *ossign.ApiClient, e map[string]ossign.SigningTheme) resource.StateChangeConf {
	return resource.StateChangeConf{
		Delay:                     30 * time.Second,
		Pending:                   []string{"waiting"},
		Target:                    []string{"complete"},
		Timeout:                   5 * time.Minute,
		MinTimeout:                300 * time.Millisecond,
		ContinuousTargetOccurence: 8,
		Refresh: func() (result interface{}, state string, err error) {
			t, apiErr := c.GetAccountSigningThemes()

			if apiErr != nil {
				if apiErr.HttpResponse != nil && apiErr.HttpResponse.StatusCode == http.StatusInternalServerError {
					// This API somtimes return transient 500 errors, we want to continue waiting when that happens
					return t, "waiting", nil
				}
				return nil, "error", apiErr.GetError()
			}

			if len(e) != len(t) {
				return t, "waiting", nil
			}

			eq := true

			for k1, v1 := range e {
				v2 := t[k1]

				if !v1.Equal(v2) {
					eq = false
					break
				}
			}

			if !eq {
				return t, "waiting", nil
			}

			return t, "complete", nil
		},
	}
}

func setResourceData(d *schema.ResourceData, ts map[string]ossign.SigningTheme) diag.Diagnostics {
	var diags diag.Diagnostics

	for k, v := range ts {
		// Only pick the first element that'll be used as the signing theme
		if err := d.Set("theme", []interface{}{flattenAccountSigningTheme(k, v)}); err != nil {
			return diag.FromErr(err)
		}
		break
	}

	return diags
}

func resourceAccountSigningThemesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	b := buildAccountSigningThemes(d)

	if err := c.CreateAccountSigningThemes(b); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	tflog.Trace(ctx, "waiting for the signing theme resource to be created...")

	scc := getSigningThemeStateChangeConf(c, b)
	_, err := scc.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "created the account's signing theme resource")

	d.SetId(c.ClientId)

	return resourceAccountSigningThemesRead(ctx, d, meta)
}

func resourceAccountSigningThemesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	ts, apiErr := c.GetAccountSigningThemes()

	if apiErr != nil {
		return diag.FromErr(apiErr.GetError())
	}

	if len(ts) < 1 {
		d.SetId("")
		return diags
	}

	diags = append(diags, setResourceData(d, ts)...)

	return diags
}

func resourceAccountSigningThemesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	b := buildAccountSigningThemes(d)

	if err := c.UpdateAccountSigningThemes(b); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	tflog.Trace(ctx, "waiting for the signing theme resource to be updated...")

	scc := getSigningThemeStateChangeConf(c, b)
	_, err := scc.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "updated the account's signing theme resource")

	return resourceAccountSigningThemesRead(ctx, d, meta)
}

func resourceAccountSigningThemesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	if err := c.DeleteAccountSigningThemes(); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	tflog.Trace(ctx, "waiting for the signing theme resource to be deleted...")

	scc := getSigningThemeStateChangeConf(c, map[string]ossign.SigningTheme{})
	_, err := scc.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "deleted the account's signing theme resource")

	d.SetId("")

	return diags
}
