package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceSigningThemes(t *testing.T) {
	th := generateSigningTheme()
	th2 := generateSigningTheme()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSigningThemesResourceDestroyed,
		Steps: []resource.TestStep{
			{
				PreConfig: testAccSigningThemesPreTestCleanup,
				Config: fmt.Sprintf(`
				resource "onespansign_account_signing_themes" "foo" {
					theme {
						name = "default"
						primary = "%s"
						success = "%s"
						warning = "%s"
						error = "%s"
						info = "%s"
						signature_button = "%s"
						optional_signature_button = "%s"
					}
				}
				`, th.Primary, th.Success, th.Warning, th.Error, th.Info, th.SignatureButton, th.OptionalSignatureButton),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_account_signing_themes.foo", "theme.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"onespansign_account_signing_themes.foo", "theme.*", map[string]string{
							"name":                      "default",
							"primary":                   th.Primary,
							"success":                   th.Success,
							"warning":                   th.Warning,
							"error":                     th.Error,
							"info":                      th.Info,
							"signature_button":          th.SignatureButton,
							"optional_signature_button": th.OptionalSignatureButton,
						}),
					testAccCheckSigningThemesResourceMatches(map[string]ossign.SigningTheme{
						"default": th,
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
				resource "onespansign_account_signing_themes" "foo" {
					theme {
						name = "default"
						primary = "%s"
						success = "%s"
						warning = "%s"
						error = "%s"
						info = "%s"
						signature_button = "%s"
						optional_signature_button = "%s"
					}
				}
				`, th2.Primary, th2.Success, th2.Warning, th2.Error, th2.Info, th2.SignatureButton, th2.OptionalSignatureButton),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("onespansign_account_signing_themes.foo", "theme.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"onespansign_account_signing_themes.foo", "theme.*", map[string]string{
							"name":                      "default",
							"primary":                   th2.Primary,
							"success":                   th2.Success,
							"warning":                   th2.Warning,
							"error":                     th2.Error,
							"info":                      th2.Info,
							"signature_button":          th2.SignatureButton,
							"optional_signature_button": th2.OptionalSignatureButton,
						}),
					testAccCheckSigningThemesResourceMatches(map[string]ossign.SigningTheme{
						"default": th2,
					}),
				),
			},
			{
				ResourceName:      "onespansign_account_signing_themes.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(`
				resource "onespansign_account_signing_themes" "foo" {
					theme {
						name = "default"
						primary = "%s"
						success = "%s"
						warning = "%s"
						error = "%s"
						info = "%s"
						signature_button = "%s"
						optional_signature_button = "%s"
					}
				}
				`, th2.Primary, th2.Success, th2.Warning, th2.Error, th2.Info, th2.SignatureButton, th2.OptionalSignatureButton),
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSigningThemesResourceDestroyed,
				),
			},
		},
	})
}

func testAccCheckSigningThemesResourceMatches(m map[string]ossign.SigningTheme) resource.TestCheckFunc {
	return func(*terraform.State) error {
		c := getTestApiClient()

		l, err := c.GetAccountSigningThemes()
		if err != nil {
			panic(err.GetError())
		}

		for k1, v1 := range l {
			v2 := m[k1]

			if !cmp.Equal(v1, v2) {
				fmt.Printf("Obtained first value: %v\n", v1)
				fmt.Printf("Obtained second value: %v\n", v2)
				return errors.New("Signing themes resource does not match expectation")
			}
		}

		return nil
	}
}

func testAccCheckSigningThemesResourceDestroyed(*terraform.State) error {
	c := getTestApiClient()

	l, apiErr := c.GetAccountSigningThemes()

	if apiErr != nil {
		return apiErr.GetError()
	}

	if len(l) == 0 {
		return nil
	}

	var d string
	b, err := json.MarshalIndent(l, "", "\t")

	if err != nil {
		d = fmt.Sprintf("Unable to marshal response: %s", err.Error())
	} else {
		d = string(b)
	}

	return fmt.Errorf("Signing theme resource not destroyed. Got:\n%s", d)
}

func testAccSigningThemesPreTestCleanup() {
	c := getTestApiClient()

	l, apiErr := c.GetAccountSigningThemes()

	if len(l) > 0 {
		c.DeleteAccountSigningThemes()
		return
	}

	if apiErr != nil {
		panic(apiErr.GetError())
	}

	scc := resource.StateChangeConf{
		Delay:                     1 * time.Minute,
		Pending:                   []string{"waiting"},
		Target:                    []string{"complete"},
		Timeout:                   2 * time.Minute,
		MinTimeout:                300 * time.Millisecond,
		ContinuousTargetOccurence: 2,
		Refresh: func() (result interface{}, state string, err error) {
			t, apiErr := c.GetAccountSigningThemes()

			if apiErr != nil {
				return nil, "error", apiErr.GetError()
			}

			if len(t) == 0 {
				return t, "complete", nil
			}

			return t, "waiting", nil
		},
	}

	_, err := scc.WaitForState()
	if err != nil {
		panic(err)
	}
}

func generateSigningTheme() ossign.SigningTheme {
	return ossign.SigningTheme{
		Primary:                 randHex(),
		Success:                 randHex(),
		Warning:                 randHex(),
		Error:                   randHex(),
		Info:                    randHex(),
		SignatureButton:         randHex(),
		OptionalSignatureButton: randHex(),
	}
}

func randHex() string {
	return fmt.Sprintf("#%s", acctest.RandStringFromCharSet(6, "0123456789ABCDEF"))
}
