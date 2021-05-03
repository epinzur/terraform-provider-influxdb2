package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func testDataSourceBucketConfig(orgName string) string {
	return fmt.Sprintf(`
		resource "influxdb2_bucket" "org" {
			name = "%s"
			description = "test org"
		}
		data "influxdb2_bucket" "by_name" {
			name = influxdb2_bucket.org.name
		}
		data "influxdb2_bucket" "by_id" {
			id = influxdb2_bucket.org.id
		}
`, orgName)
}

func TestAccDataSourceBucket(t *testing.T) {
	org := acctest.RandomWithPrefix("test-org")

	var provider *schema.Provider

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories(&provider),
		Steps: []resource.TestStep{
			{
				Config: testConfig(testDataSourceBucketConfig(org)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.influxdb2_bucket.by_id", "name", org),
					resource.TestCheckResourceAttr("data.influxdb2_bucket.by_name", "name", org),
					resource.TestCheckResourceAttr("data.influxdb2_bucket.by_id", "description", "test org"),
					resource.TestCheckResourceAttr("data.influxdb2_bucket.by_name", "description", "test org"),
				),
			},
		},
	})
}
