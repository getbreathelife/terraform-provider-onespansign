package provider

import (
	"context"
	"net/url"
	"regexp"

	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"environment_url": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Environment URL for the OneSpan sign account.",
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("^https://(www.)?[a-zA-Z0-9.-]{2,256}.[a-z]{2,4}$"), "Please provide a valid environment URL in the format of <scheme>://<host>")),
				},
				"client_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Client ID of the client app created for this provider.",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "Client secret of the client app created for this provider.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{},
			ResourcesMap: map[string]*schema.Resource{
				"onespansign_account_signing_logos":  resourceAccountSigningLogos(),
				"onespansign_account_signing_themes": resourceAccountSigningThemes(),
				"onespansign_data_management_policy": resourceDataManagementPolicy(),
				"onespansign_expiry_time_config":     resourceExpiryTimeConfig(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		eu := d.Get("environment_url").(string)
		id := d.Get("client_id").(string)
		secret := d.Get("client_secret").(string)

		url, err := url.Parse(eu)
		if err != nil {
			panic(err)
		}

		return ossign.NewClient(ossign.ApiClientConfig{
			BaseUrl:      url,
			ClientId:     id,
			ClientSecret: secret,
			UserAgent:    p.UserAgent("terraform-provider-onespan-sign", version),
		}), nil
	}
}
