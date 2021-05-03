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
	createBucketDesc = "test bucket"
	updateBucketDesc = "test bucket update"
)

func influxBucket(orgName string, orgDesc string) string {
	return fmt.Sprintf(`
		resource "influxdb2_bucket" "org" {
			name = "%s"
			description = "%s"
		}
`, orgName, orgDesc)
}

func TestAccResourceBucket(t *testing.T) {
	org := acctest.RandomWithPrefix("test-org")

	var provider *schema.Provider

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories(&provider),
		CheckDestroy:      testAccCheckResourceBucketDestroy(t, provider),
		Steps: []resource.TestStep{
			{
				//create
				Config: testConfig(influxBucket(org, createBucketDesc)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb2_bucket.org", "name", org),
					resource.TestCheckResourceAttr("influxdb2_bucket.org", "description", createBucketDesc),
					testAccResourceBucketExists(provider, "influxdb2_bucket.org"),
				),
			},
			importStep("influxdb2_bucket.org"),
			{
				//update
				Config: testConfig(influxOrg(org, updateBucketDesc)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("influxdb2_bucket.org", "name", org),
					resource.TestCheckResourceAttr("influxdb2_bucket.org", "description", updateBucketDesc),
					testAccResourceBucketExists(provider, "influxdb2_bucket.org"),
				),
			},
			importStep("influxdb2_bucket.org"),
		},
	})
}

func testAccResourceBucketExists(testProvider *schema.Provider, name string) resource.TestCheckFunc {
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

		if _, err := client.BucketsAPI().FindBucketByID(context.Background(), id); err != nil {
			return fmt.Errorf("Got an error when reading Bucket %q: %v", id, err)
		}

		return nil
	}
}

func testAccCheckResourceBucketDestroy(t *testing.T, testProvider *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testProvider.Meta() == nil {
			t.Fatal("got nil provider metadata")
		}
		client := testProvider.Meta().(*metaData).client

		for _, rs := range s.RootModule().Resources {
			switch rs.Type {
			case "influxdb2_bucket":
				id := rs.Primary.ID

				_, err := client.BucketsAPI().FindBucketByID(context.Background(), id)
				if !strings.Contains(err.Error(), "not found") {
					//return fmt.Errorf("didn't get a 404 when reading destroyed account %q: %v", id, err)
					return fmt.Errorf("Was able to find destroyed Bucket %q: %v", id, err)
				}

			default:
				continue
			}
		}
		return nil
	}
}
