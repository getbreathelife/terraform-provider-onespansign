package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDataRetentionSettings() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: `OneSpan Sign account's data retention settings.
		
		Please note that this resource is a singleton, which means that only one instance of this resource should
		exist. Having multiple instances may produce unexpected result.`,

		Schema: map[string]*schema.Schema{
			"data_management_policy": {
				// This description is used by the documentation generator and the language server.
				Description: "Data management policy for the account.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem:        resourceDataManagementPolicy(),
			},
			"expiry_time_config": {
				Description: "Expiry configurations defined for the account.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem:        resourceExpiryTimeConfig(),
			},
		},
	}
}
