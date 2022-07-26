package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/getbreathelife/terraform-provider-onespansign/internal/helpers"
	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDataManagementPolicy(t *testing.T) {
	tr := generateTransactionRetention()

	// Default config
	tr2 := ossign.TransactionRetention{
		Draft:     json.Number("0"),
		Sent:      json.Number("0"),
		Completed: json.Number("0"),
		Archived:  json.Number("0"),
		Declined:  json.Number("0"),
		OptedOut:  json.Number("0"),
		Expired:   json.Number("0"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "onespansign_data_management_policy" "foo" {
					transaction_retention {
						draft = %s
						sent = %s
						completed = %s
						archived = %s
						declined = %s
						opted_out = %s
						expired = %s
					}
				}
				`, tr.Draft.String(), tr.Sent.String(), tr.Completed.String(),
					tr.Archived.String(), tr.Declined.String(), tr.OptedOut.String(),
					tr.Expired.String(),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_data_management_policy.foo", "transaction_retention.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"onespansign_data_management_policy.foo", "transaction_retention.*", map[string]string{
							"draft":     tr.Draft.String(),
							"sent":      tr.Sent.String(),
							"completed": tr.Completed.String(),
							"archived":  tr.Archived.String(),
							"declined":  tr.Declined.String(),
							"opted_out": tr.OptedOut.String(),
							"expired":   tr.Expired.String(),
						}),
					testAccCheckDataManagementPolicyResourceMatches(ossign.DataManagementPolicy{
						TransactionRetention: tr,
					}),
				),
			},
			{
				ResourceName:      "onespansign_data_management_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(`
				resource "onespansign_data_management_policy" "foo" {
					transaction_retention {
						draft = %s
						sent = %s
						completed = %s
						archived = %s
						declined = %s
						opted_out = %s
						expired = %s
					}
				}
				`, tr2.Draft.String(), tr2.Sent.String(), tr2.Completed.String(),
					tr2.Archived.String(), tr2.Declined.String(), tr2.OptedOut.String(),
					tr2.Expired.String(),
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_data_management_policy.foo", "transaction_retention.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"onespansign_data_management_policy.foo", "transaction_retention.*", map[string]string{
							"draft":     tr2.Draft.String(),
							"sent":      tr2.Sent.String(),
							"completed": tr2.Completed.String(),
							"archived":  tr2.Archived.String(),
							"declined":  tr2.Declined.String(),
							"opted_out": tr2.OptedOut.String(),
							"expired":   tr2.Expired.String(),
						}),
					testAccCheckDataManagementPolicyResourceMatches(ossign.DataManagementPolicy{
						TransactionRetention: tr2,
					}),
				),
			},
		},
	})
}

func testAccCheckDataManagementPolicyResourceMatches(m ossign.DataManagementPolicy) resource.TestCheckFunc {
	return func(*terraform.State) error {
		c := getTestApiClient()

		p, err := c.GetDataManagementPolicy()
		if err != nil {
			return err.GetError()
		}

		if p.TransactionRetention.Draft.String() != m.TransactionRetention.Draft.String() ||
			p.TransactionRetention.Sent.String() != m.TransactionRetention.Sent.String() ||
			p.TransactionRetention.Completed.String() != m.TransactionRetention.Completed.String() ||
			p.TransactionRetention.Archived.String() != m.TransactionRetention.Archived.String() ||
			p.TransactionRetention.Declined.String() != m.TransactionRetention.Declined.String() ||
			p.TransactionRetention.OptedOut.String() != m.TransactionRetention.OptedOut.String() ||
			p.TransactionRetention.Expired.String() != m.TransactionRetention.Expired.String() {
			fmt.Printf("Obtained remote value: %v\n", p)
			fmt.Printf("Obtained local value: %v\n", m)
			return errors.New("Data management policy resource does not match expectation")
		}

		return nil
	}
}

func generateTransactionRetention() ossign.TransactionRetention {
	return ossign.TransactionRetention{
		Draft:     helpers.RandJsonNumber(30, 120),
		Sent:      helpers.RandJsonNumber(30, 120),
		Completed: helpers.RandJsonNumber(30, 120),
		Archived:  helpers.RandJsonNumber(30, 120),
		Declined:  helpers.RandJsonNumber(30, 120),
		OptedOut:  helpers.RandJsonNumber(30, 120),
		Expired:   helpers.RandJsonNumber(30, 120),
	}
}
