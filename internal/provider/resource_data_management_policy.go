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

func resourceDataManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "OneSpan Sign account's data management policy.",

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
					},
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flattenTransactionRetention(tr ossign.TransactionRetention) interface{} {
	trm := make(map[string]interface{}, 7)

	trm["draft"] = tr.Draft
	trm["sent"] = tr.Sent
	trm["completed"] = tr.Completed
	trm["archived"] = tr.Archived
	trm["declined"] = tr.Declined
	trm["opted_out"] = tr.OptedOut
	trm["expired"] = tr.Expired

	return trm
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

	dmp, err := c.GetDataManagementPolicy()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Summary,
			Detail:   err.Detail,
		})
		return diags
	}

	tr := flattenTransactionRetention(dmp.TransactionRetention)

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
			Draft:     helpers.GetJsonNumber(int64(i["draft"].(int)), 10),
			Sent:      helpers.GetJsonNumber(int64(i["sent"].(int)), 10),
			Completed: helpers.GetJsonNumber(int64(i["completed"].(int)), 10),
			Archived:  helpers.GetJsonNumber(int64(i["archived"].(int)), 10),
			Declined:  helpers.GetJsonNumber(int64(i["declined"].(int)), 10),
			OptedOut:  helpers.GetJsonNumber(int64(i["opted_out"].(int)), 10),
			Expired:   helpers.GetJsonNumber(int64(i["expired"].(int)), 10),
		}
	}

	if tr != nil {
		b := ossign.DataManagementPolicy{
			TransactionRetention: *tr,
		}

		if err := c.UpdateDataManagementPolicy(b); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Summary,
				Detail:   err.Detail,
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
