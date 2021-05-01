package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func testDataSourceOrganizationConfig(orgName string) string {
	return fmt.Sprintf(`
		resource "influxdb2_organization" "org" {
			name = "%s"
			description = "test org"
			status = "active"
		}
		data "influxdb2_organization" "by_name" {
			name = influxdb2_organization.org.name
		}
		data "influxdb2_organization" "by_id" {
			id = influxdb2_organization.org.id
		}
`, orgName)
}

func TestAccDataSourceOrganization(t *testing.T) {
	org := acctest.RandomWithPrefix("test-org")

	var provider *schema.Provider

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories(&provider),
		PreCheck:          func() { testAccPreCheck(t) },
		//CheckDestroy:      testAccCheckDataSourceOrganizationDestroy(t, provider),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationConfig(org),
				Check: resource.ComposeTestCheckFunc(
					//testAccDataSourceOrganizationExists(provider, "data.influxdb2_organization.org"),
					resource.TestCheckResourceAttr("data.influxdb2_organization.by_id", "name", org),
					resource.TestCheckResourceAttr("data.influxdb2_organization.by_name", "name", org),
					resource.TestCheckResourceAttr("data.influxdb2_organization.by_id", "description", "test org"),
					resource.TestCheckResourceAttr("data.influxdb2_organization.by_name", "description", "test org"),
					resource.TestCheckResourceAttr("data.influxdb2_organization.by_id", "status", "active"),
					resource.TestCheckResourceAttr("data.influxdb2_organization.by_name", "status", "active"),
				),
			},
		},
	})
}

/*
func testAccDataSourceOrganizationExists(testProvider *schema.Provider, name string) resource.TestCheckFunc {
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

		if _, err := client.OrganizationsAPI().FindOrganizationByID(id); err != nil {
			return fmt.Errorf("Got an error when reading organization %q: %v", id, err)
		}

		return nil
	}
}*/

/*
func testAccCheckDataSourceOrganizationDestroy(t *testing.T, testProvider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testProvider.Meta() == nil {
			t.Fatal("got nil provider metadata")
		}
		client := testProvider.Meta().(*metaData).client

		for _, rs := range s.RootModule().Resources {
			switch rs.Type {
			case "influxdb2oss_organization":
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
}*/
