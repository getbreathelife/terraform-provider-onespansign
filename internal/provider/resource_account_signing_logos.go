package provider

import (
	"context"
	"fmt"
	"strings"

	ossign "github.com/getbreathelife/terraform-provider-onespan-sign/pkg/onespan-sign"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vincent-petithory/dataurl"
)

func resourceAccountSigningLogos() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: `OneSpan Sign account's customized logos used during the Signing Ceremony.
		
		Please note that this resource is a singleton, which means that only one instance of this resource should
		exist. Having multiple instances may produce unexpected result.`,

		CreateContext: resourceAccountSigningLogosCreate,
		ReadContext:   resourceAccountSigningLogosRead,
		UpdateContext: resourceAccountSigningLogosUpdate,
		DeleteContext: resourceAccountSigningLogosDelete,

		Schema: map[string]*schema.Schema{
			"logo": {
				// This description is used by the documentation generator and the language server.
				Description: "Customized logo used during the Signing Ceremony.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"language": {
							Description:      "The language of the Signing Ceremony where the image will be used.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"en", "fr", "it", "ru", "es", "pt", "de", "nl", "da", "el", "zh-CN", "zh-TW", "ja", "ko"}, false)),
						},
						"image": {
							Description:      "Base 64 decoded image (Data URI).",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: isValidImageData,
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

func isValidImageData(v interface{}, p cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	value := v.(string)

	d, err := dataurl.DecodeString(value)

	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to parse the data URI",
			Detail:   err.Error(),
		})
	}

	var supportedContentTypes = []string{"image/jpeg", "image/png", "image/gif", "image/bmp", "image/svg+xml"}
	contentType := d.ContentType()

	isValidType := false

	for _, v := range supportedContentTypes {
		if contentType == v {
			isValidType = true
			break
		}
	}

	if !isValidType {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("invalid or unsupported content type: %s", contentType),
			Detail:   fmt.Sprintf("supported content types are: %s", strings.Join(supportedContentTypes, ", ")),
		})
	}

	size := len(d.Data)

	if size > 1_000_000 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "content is too large",
			Detail:   fmt.Sprintf("maximum content size is 1MB, got: %d", size),
		})
	}

	return diags
}

func flattenAccountSigningLogos(logos []ossign.SigningLogo) []interface{} {
	ls := make([]interface{}, len(logos))

	for i, v := range logos {
		e := make(map[string]interface{}, 2)

		e["language"] = v.Language
		e["image"] = v.Image

		ls[i] = e
	}

	return ls
}

func resourceAccountSigningLogosCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "updating (replacing) resource instead of creating",
		Detail:   "This resource is a singleton. It only supports retrieval or replacement operations.",
	})

	diags = append(diags, resourceAccountSigningLogosUpdate(ctx, d, meta)...)

	return diags
}

func resourceAccountSigningLogosRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	logos, err := c.GetSigningLogos()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	if err := d.Set("logo", flattenAccountSigningLogos(logos)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAccountSigningLogosUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics
	var b []ossign.SigningLogo

	logos := d.Get("logo").(*schema.Set).List()
	for _, item := range logos {
		i := item.(map[string]interface{})

		b = append(b, ossign.SigningLogo{
			Language: i["language"].(string),
			Image:    i["image"].(string),
		})
	}

	if err := c.UpdateSigningLogos(b); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "updated the account signing logos resource")

	return resourceAccountSigningLogosRead(ctx, d, meta)
}

func resourceAccountSigningLogosDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "no deletion will take place",
			Detail: `This resource is a singleton. It only supports retrieval or replacement operations.
			To remove all account signing logos, declare the resource without any "logo" attribute instead.`,
		},
	}
}
