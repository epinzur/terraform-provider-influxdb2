package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// How to run the acceptance tests for this provider:
//
// - Install docker & docker-compose on your machine
//
// - Run the following to start the Influx server in test mode:
//       docker-compose up
//
// - Run the tests `make testacc`

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
func providerFactories(p **schema.Provider) map[string]func() (*schema.Provider, error) {
	*p = New("dev")()
	return map[string]func() (*schema.Provider, error){
		"influxdb2": func() (*schema.Provider, error) {
			return *p, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testConfig(res ...string) string {
	provider := `
		provider "influxdb2" {
			host     = "http://localhost:8086"
			token    = "oops_this_is_committed_to_source_control"
		}
	`

	c := []string{provider}
	c = append(c, res...)
	return strings.Join(c, "\n")
}

func importStep(name string, ignore ...string) resource.TestStep {
	step := resource.TestStep{
		ResourceName:      name,
		ImportState:       true,
		ImportStateVerify: true,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}
