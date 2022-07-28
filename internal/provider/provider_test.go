package provider

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/getbreathelife/terraform-provider-onespansign/pkg/ossign"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joho/godotenv"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"onespansign": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func loadEnvVar() {
	// Load env var from a .env file if it exist
	godotenv.Load("../../.env")

	if v := os.Getenv("ENV_URL"); v == "" {
		panic("ENV_URL must be set for acceptance tests")
	}
	if v := os.Getenv("CLIENT_ID"); v == "" {
		panic("CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("CLIENT_SECRET"); v == "" {
		panic("CLIENT_SECRET must be set for acceptance tests")
	}
}

func getTestApiClient() *ossign.ApiClient {
	loadEnvVar()

	url, err := url.Parse(os.Getenv("ENV_URL"))
	if err != nil {
		panic(err)
	}

	return ossign.NewClient(ossign.ApiClientConfig{
		BaseUrl:      url,
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	})
}

func getTestConfig(c string) string {
	loadEnvVar()

	return fmt.Sprintf(`
	provider onespansign {
		client_id = "%s"
		client_secret = "%s"
		environment_url = "%s"
	}

	%s
	`, os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("ENV_URL"), c)
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
