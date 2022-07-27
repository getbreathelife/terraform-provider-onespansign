package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/getbreathelife/terraform-provider-onespansign/internal/helpers"
	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceExpiryTimeConfig(t *testing.T) {
	etc := generateExpiryTimeConfig()

	// Default config
	etc2 := ossign.ExpiryTimeConfiguration{
		Default: json.Number("0"),
		Maximum: json.Number("0"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "onespansign_expiry_time_config" "foo" {
					default = %s
					maximum = %s
				}
				`, etc.Default.String(), etc.Maximum.String()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_expiry_time_config.foo", "default", etc.Default.String()),
					resource.TestCheckResourceAttr("onespansign_expiry_time_config.foo", "maximum", etc.Maximum.String()),
					testAccCheckExpiryTimeConfigResourceMatches(etc),
				),
			},
			{
				ResourceName:      "onespansign_expiry_time_config.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(`
				resource "onespansign_expiry_time_config" "foo" {
					default = %s
					maximum = %s
				}
				`, etc2.Default.String(), etc2.Maximum.String()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_expiry_time_config.foo", "default", etc2.Default.String()),
					resource.TestCheckResourceAttr("onespansign_expiry_time_config.foo", "maximum", etc2.Maximum.String()),
					testAccCheckExpiryTimeConfigResourceMatches(etc2),
				),
			},
		},
	})
}

func testAccCheckExpiryTimeConfigResourceMatches(m ossign.ExpiryTimeConfiguration) resource.TestCheckFunc {
	return func(*terraform.State) error {
		c := getTestApiClient()

		p, err := c.GetExpiryTimeConfiguration()
		if err != nil {
			return err.GetError()
		}

		if !cmp.Equal(*p, m) {
			fmt.Printf("Obtained remote value: %v\n", *p)
			fmt.Printf("Obtained local value: %v\n", m)
			return errors.New("Expiry time configuration resource does not match expectation")
		}

		return nil
	}
}

func generateExpiryTimeConfig() ossign.ExpiryTimeConfiguration {
	return ossign.ExpiryTimeConfiguration{
		Default: helpers.RandJsonNumber(30, 60),
		Maximum: helpers.RandJsonNumber(60, 120),
	}
}
