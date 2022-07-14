package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getbreathelife/terraform-provider-onespan-sign/internal/api_client"
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
		Description: "OneSpan Sign account's customized logos used during the Signing Ceremony.",

		CreateContext: resourceAccountSigningLogosCreate,
		ReadContext:   resourceScaffoldingRead,
		UpdateContext: resourceScaffoldingUpdate,
		DeleteContext: resourceScaffoldingDelete,

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

func resourceAccountSigningLogosCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*api_client.ApiClient)

	var diags diag.Diagnostics
	var b []api_client.SigningLogo

	logos := d.Get("logo").(*schema.Set).List()
	for _, v := range logos {
		b = append(b, v.(api_client.SigningLogo))
	}

	res, err := c.UpdateSigningLogos(b)
	if err != nil {
		return diag.FromErr(err)
	}

	if res.StatusCode != http.StatusOK {
		diags = append(diags, api_client.GetApiErrorDiag(res))
		return diags
	}

	d.SetId(c.ClientId)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created the account signing logos resource")

	return diags
}

func resourceScaffoldingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceScaffoldingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}

func resourceScaffoldingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	return diag.Errorf("not implemented")
}
