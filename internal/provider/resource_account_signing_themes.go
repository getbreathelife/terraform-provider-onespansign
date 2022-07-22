package provider

import (
	"context"
	"regexp"

	ossign "github.com/getbreathelife/terraform-provider-onespan-sign/pkg/onespan-sign"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccountSigningThemes() *schema.Resource {
	return &schema.Resource{
		Description: "OneSpan Sign account's customized signing themes.",

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

	tflog.Trace(ctx, "created the account's signing theme resource")

	return resourceAccountSigningThemesRead(ctx, d, meta)
}

func resourceAccountSigningThemesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	ts, err := c.GetAccountSigningThemes()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	if len(ts) < 1 {
		d.SetId("")
		return diags
	}

	for k, v := range ts {
		// Only pick the first element that'll be used as the signing theme
		if err := d.Set("theme", []interface{}{flattenAccountSigningTheme(k, v)}); err != nil {
			return diag.FromErr(err)
		}
		break
	}

	d.SetId(c.ClientId)

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

	tflog.Trace(ctx, "deleted the account's signing theme resource")

	return resourceAccountSigningThemesRead(ctx, d, meta)
}
