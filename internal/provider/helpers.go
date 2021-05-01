package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createdUpdatedSchema(itemType string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_at": {
			Description: fmt.Sprintf("The string time that the %s was created.", itemType),
			Type:        schema.TypeString,
			Computed:    true,
		},
		"updated_at": {
			Description: fmt.Sprintf("The string time that the %s was last updated.", itemType),
			Type:        schema.TypeString,
			Computed:    true,
		},
		"created_timestamp": {
			Description: fmt.Sprintf("The timestamp that the %s was created.", itemType),
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"updated_timestamp": {
			Description: fmt.Sprintf("The timestamp that the %s was last updated.", itemType),
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
}

func mergeSchemas(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	res := map[string]*schema.Schema{}
	for _, s := range schemas {
		for k, v := range s {
			res[k] = v
		}
	}
	return res
}
