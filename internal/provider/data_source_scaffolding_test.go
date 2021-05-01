package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataSourceScaffolding(t *testing.T) {
	t.Skip("data source not yet implemented, remove this once you add your own code")

	var provider *schema.Provider

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories(&provider),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceScaffolding,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.scaffolding_data_source.foo", "sample_attribute", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccDataSourceScaffolding = `
data "scaffolding_data_source" "foo" {
  sample_attribute = "bar"
}
`
