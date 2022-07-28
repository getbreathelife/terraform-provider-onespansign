package provider

import (
	"context"
	"fmt"

	"github.com/getbreathelife/terraform-provider-onespansign/internal/helpers"
	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDataManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Description: `OneSpan Sign account's data management policy.

Please note that this resource is a singleton, which means that only one instance of this resource should
exist. Having multiple instances may produce unexpected result.`,

		CreateContext: resourceDataManagementPolicyCreate,
		ReadContext:   resourceDataManagementPolicyRead,
		UpdateContext: resourceDataManagementPolicyUpdate,
		DeleteContext: resourceDataManagementPolicyDelete,

		Schema: map[string]*schema.Schema{
			"transaction_retention": {
				Description: "Transaction retention settings.",
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"draft": {
							Description:      "Number of days to keep drafts for.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"sent": {
							Description:      "Number of days to keep sent transactions for.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"completed": {
							Description:      "Number of days to keep completed transactions for.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"archived": {
							Description:      "Number of days to keep archived transactions for.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"declined": {
							Description:      "Number of days to keep declined transactions for.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"opted_out": {
							Description:      "Number of days to keep opted-out transactions for.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"expired": {
							Description:      "Number of days to keep expired transactions for.",
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"lifetime_total": {
							Description:      "Number of days to keep the transactions, calculated from the day that the transaction is created.",
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          120,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"lifetime_until_completion": {
							Description:      "Number of days that incomplete transactions will be stored, calculated from the day that the transaction is created.",
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          120,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"include_sent": {
							Description: "Include sent transactions as part of the \"incomplete transactions\".",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
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

func flattenTransactionRetention(tr ossign.TransactionRetention) (interface{}, error) {
	var err error

	trm := make(map[string]interface{}, 7)

	trm["draft"], err = helpers.GetInt(tr.Draft)
	if err != nil {
		return trm, err
	}

	trm["sent"], err = helpers.GetInt(tr.Sent)
	if err != nil {
		return trm, err
	}

	trm["completed"], err = helpers.GetInt(tr.Completed)
	if err != nil {
		return trm, err
	}

	trm["archived"], err = helpers.GetInt(tr.Archived)
	if err != nil {
		return trm, err
	}

	trm["declined"], err = helpers.GetInt(tr.Declined)
	if err != nil {
		return trm, err
	}

	trm["opted_out"], err = helpers.GetInt(tr.OptedOut)
	if err != nil {
		return trm, err
	}

	trm["expired"], err = helpers.GetInt(tr.Expired)
	if err != nil {
		return trm, err
	}

	trm["lifetime_total"], err = helpers.GetInt(tr.LifetimeTotal)
	if err != nil {
		return trm, err
	}

	trm["lifetime_until_completion"], err = helpers.GetInt(tr.LifetimeUntilCompletion)
	if err != nil {
		return trm, err
	}

	trm["include_sent"] = tr.IncludeSent

	return trm, nil
}

func resourceDataManagementPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "updating (replacing) resource instead of creating",
		Detail:   "This resource is a singleton. It only supports retrieval or replacement operations.",
	})

	diags = append(diags, resourceDataManagementPolicyUpdate(ctx, d, meta)...)

	d.SetId(c.ClientId)

	return diags
}

func resourceDataManagementPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics

	dmp, apiErr := c.GetDataManagementPolicy()

	if apiErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  apiErr.Summary,
			Detail:   apiErr.Detail,
		})
		return diags
	}

	tr, err := flattenTransactionRetention(dmp.TransactionRetention)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("transaction_retention", []interface{}{tr}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDataManagementPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	c := meta.(*ossign.ApiClient)

	var diags diag.Diagnostics
	var tr *ossign.TransactionRetention

	trs := d.Get("transaction_retention").(*schema.Set).List()

	if len(trs) > 0 {
		i := trs[0].(map[string]interface{})

		tr = &ossign.TransactionRetention{
			Draft:                   helpers.GetJsonNumber(int64(i["draft"].(int)), 10),
			Sent:                    helpers.GetJsonNumber(int64(i["sent"].(int)), 10),
			Completed:               helpers.GetJsonNumber(int64(i["completed"].(int)), 10),
			Archived:                helpers.GetJsonNumber(int64(i["archived"].(int)), 10),
			Declined:                helpers.GetJsonNumber(int64(i["declined"].(int)), 10),
			OptedOut:                helpers.GetJsonNumber(int64(i["opted_out"].(int)), 10),
			Expired:                 helpers.GetJsonNumber(int64(i["expired"].(int)), 10),
			LifetimeTotal:           helpers.GetJsonNumber(int64(i["lifetime_total"].(int)), 10),
			LifetimeUntilCompletion: helpers.GetJsonNumber(int64(i["lifetime_until_completion"].(int)), 10),
			IncludeSent:             i["include_sent"].(bool),
		}
	}

	if tr != nil {
		b := ossign.DataManagementPolicy{
			TransactionRetention: *tr,
		}

		if apiErr := c.UpdateDataManagementPolicy(b); apiErr != nil {
			// There are undocumented validation errors that occur sometimes on a seemingly valid payload.
			// This special handling is added to easily debug the issue.
			if apiErr.HttpResponse != nil && apiErr.HttpResponse.StatusCode%400 < 100 {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  apiErr.Summary,
					Detail:   fmt.Sprintf("4xx error occurred while updating the data management policy: %v", b),
				})
			}

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  apiErr.Summary,
				Detail:   apiErr.Detail,
			})
			return diags
		}

		tflog.Trace(ctx, "updated the account's data management policy resource")
	}

	return resourceDataManagementPolicyRead(ctx, d, meta)
}

func resourceDataManagementPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "no deletion will take place",
			Detail:   "This resource is a singleton. It only supports retrieval or replacement operations.",
		},
	}
}
