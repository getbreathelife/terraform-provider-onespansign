package provider

import (
	"context"

	ossign "github.com/getbreathelife/terraform-provider-onespan-sign/pkg/onespan-sign"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
					Type:        schema.TypeString,
					Required:    true,
					Description: "Environment URL for the OneSpan sign account. For the list of available URLs, please visit [Environment URLs & IP Addresses (OneSpan Sign)](https://community.onespan.com/documentation/onespan-sign/guides/quick-start-guides/developer/environment-urls-ip-addresses).",
				},
				"client_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Client ID of the [OneSpan Sign Client App](https://community.onespan.com/documentation/onespan-sign/guides/admin-guides/user/integration) created for this provider.",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "Client secret of the [OneSpan Sign Client App](https://community.onespan.com/documentation/onespan-sign/guides/admin-guides/user/integration) created for this provider.",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{},
			ResourcesMap: map[string]*schema.Resource{
				"account_signing_logos_resource": resourceAccountSigningLogos(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		url := d.Get("environment_url").(string)
		id := d.Get("client_id").(string)
		secret := d.Get("client_secret").(string)

		return ossign.NewClient(&ossign.ApiClientConfig{
			BaseUrl:      url,
			ClientId:     id,
			ClientSecret: secret,
			UserAgent:    p.UserAgent("terraform-provider-onespan-sign", version),
		}), nil
	}
}
