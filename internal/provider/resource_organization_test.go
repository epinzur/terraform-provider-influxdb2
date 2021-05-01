package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	createDesc = "test org"
	updateDesc = "test org update"
)

func influxOrg(orgName string, orgDesc string) string {
	return fmt.Sprintf(`
		resource "influxdb2_organization" "org" {
			name = "%s"
			description = "%s"
			status = "active"
		}
`, orgName, orgDesc)
}

func TestAccResourceOrganization(t *testing.T) {
	org := acctest.RandomWithPrefix("test-org")

	var provider *schema.Provider

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories(&provider),
		CheckDestroy:      testAccCheckResourceOrganizationDestroy(t, provider),
		Steps: []resource.TestStep{
			{
				//create
				Config: testConfig(influxOrg(org, createDesc)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb2_organization.org", "name", org),
					resource.TestCheckResourceAttr("influxdb2_organization.org", "description", createDesc),
					resource.TestCheckResourceAttr("influxdb2_organization.org", "status", "active"),
					testAccResourceOrganizationExists(provider, "influxdb2_organization.org"),
				),
			},
			importStep("influxdb2_organization.org"),
			{
				//update
				Config: testConfig(influxOrg(org, updateDesc)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb2_organization.org", "name", org),
					resource.TestCheckResourceAttr("influxdb2_organization.org", "description", updateDesc),
					resource.TestCheckResourceAttr("influxdb2_organization.org", "status", "active"),
					testAccResourceOrganizationExists(provider, "influxdb2_organization.org"),
				),
			},
			importStep("influxdb2_organization.org"),
		},
	})
}

func testAccResourceOrganizationExists(testProvider *schema.Provider, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		client := testProvider.Meta().(*metaData).client

		if _, err := client.OrganizationsAPI().FindOrganizationByID(context.Background(), id); err != nil {
			return fmt.Errorf("Got an error when reading organization %q: %v", id, err)
		}

		return nil
	}
}

func testAccCheckResourceOrganizationDestroy(t *testing.T, testProvider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testProvider.Meta() == nil {
			t.Fatal("got nil provider metadata")
		}
		client := testProvider.Meta().(*metaData).client

		for _, rs := range s.RootModule().Resources {
			switch rs.Type {
			case "influxdb2_organization":
				id := rs.Primary.ID

				_, err := client.OrganizationsAPI().FindOrganizationByID(context.Background(), id)
				if !strings.Contains(err.Error(), "not found") {
					//return fmt.Errorf("didn't get a 404 when reading destroyed account %q: %v", id, err)
					return fmt.Errorf("Was able to find destroyed organization %q: %v", id, err)
				}

			default:
				continue
			}
		}
		return nil
	}
}
