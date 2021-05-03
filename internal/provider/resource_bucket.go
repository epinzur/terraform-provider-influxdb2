package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/influxdata/influxdb-client-go/domain"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "The Bucket resource allows you to configure a InfluxDB2 Bucket.",

		CreateContext: resourceBucketCreate,
		ReadContext:   resourceBucketRead,
		UpdateContext: resourceBucketUpdate,
		DeleteContext: resourceBucketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: mergeSchemas(map[string]*schema.Schema{
			// Required Inputs
			"name": {
				Description: "Name of the Bucket.",
				Type:        schema.TypeString,
				Required:    true,
			},
			// Optional Inputs
			"description": {
				Description: "The description of the Bucket.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			// Computed outputs
			"id": {
				Description: "ID of the Bucket.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		}, createdUpdatedSchema("Bucket")),
	}
}

func resourceBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*metaData).client
	orgsAPI := client.BucketsAPI()

	name := d.Get("name").(string)

	// Check for an existing Bucket
	_, err := orgsAPI.FindBucketByName(ctx, name)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return diag.Errorf("unable to check for presence of an existing Bucket (%s): %v", name, err)
		}
		log.Printf("[INFO] Bucket (%s) not found, proceeding with create", name)
	} else {
		return diag.Errorf("unable to create Bucket (%s) - an Bucket with this name already exists; see resouce documentation for influxdb2_bucket for instructions on how to add an already existing Bucket to the state", name)
	}

	description := d.Get("description").(string)
	org := *&domain.Bucket{
		Name:        name,
		Description: &description,
	}

	log.Printf("[INFO] Creating Bucket (%s)", name)
	returnedOrg, err := orgsAPI.CreateBucket(ctx, &org)
	if err != nil {
		return diag.Errorf("unable to create Bucket (%s): %v", name, err)
	}

	if returnedOrg.Id == nil {
		return diag.Errorf("unable to create Bucket (%s): <unknown error occurred>", name)
	}

	id := *returnedOrg.Id

	d.SetId(id)

	log.Printf("[INFO] Created Bucket (%s) (%s)", name, id)

	// Get the updated Bucket
	updatedOrg, err := orgsAPI.FindBucketByID(ctx, id)
	if err != nil {
		return diag.Errorf("unable to retrieve Bucket (%s) (%s): %v", name, id, err)
	}

	if err := setBucketResourceData(d, updatedOrg); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*metaData).client
	orgsAPI := client.BucketsAPI()

	id := d.Id()

	log.Printf("[INFO] Reading Bucket (%s)", id)

	org, err := orgsAPI.FindBucketByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] Bucket (%s) not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return diag.Errorf("unable to retrieve Bucket (%s): %v", id, err)
	}

	// Bucket found, update resource data
	if err := setBucketResourceData(d, org); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceBucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*metaData).client
	orgsAPI := client.BucketsAPI()

	id := d.Id()

	log.Printf("[INFO] Reading Bucket (%s)", id)

	org, err := orgsAPI.FindBucketByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] Bucket (%s) not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return diag.Errorf("unable to retrieve Bucket (%s): %v", id, err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	org.Name = name
	org.Description = &description

	log.Printf("[INFO] Updating Bucket (%s)", id)
	updatedOrg, err := orgsAPI.UpdateBucket(ctx, org)
	if err != nil {
		return diag.Errorf("unable to update Bucket (%s): %v", id, err)
	}

	log.Printf("[INFO] Updated Bucket (%s)", id)

	if err := setBucketResourceData(d, updatedOrg); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*metaData).client
	orgsAPI := client.BucketsAPI()

	id := d.Id()

	log.Printf("[INFO] Deleting Bucket (%s)", id)

	err := orgsAPI.DeleteBucketWithID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] Bucket (%s) not found, so no action was taken", id)
			return nil
		}
		return diag.Errorf("unable to delete Bucket (%s): %v", id, err)
	}

	log.Printf("[INFO] Deleting (%s) deleted, removing from state", id)

	return nil
}

func setBucketResourceData(d *schema.ResourceData, org *domain.Bucket) error {
	if err := d.Set("id", org.Id); err != nil {
		return err
	}
	if err := d.Set("name", org.Name); err != nil {
		return err
	}
	if err := d.Set("description", org.Description); err != nil {
		return err
	}
	if err := d.Set("created_at", org.CreatedAt.UTC().String()); err != nil {
		return err
	}
	if err := d.Set("updated_at", org.UpdatedAt.UTC().String()); err != nil {
		return err
	}
	if err := d.Set("created_timestamp", org.CreatedAt.Unix()); err != nil {
		return err
	}
	if err := d.Set("updated_timestamp", org.UpdatedAt.Unix()); err != nil {
		return err
	}
	return nil
}

// resourceBucketImport implements the logic necessary to import an un-tracked
// (by Terraform) Bucket resource into Terraform state.
func resourceBucketImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*metaData).client
	orgsAPI := client.BucketsAPI()

	id := d.Id()

	// Get the imported Bucket
	importedOrg, err := orgsAPI.FindBucketByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to import Bucket (%s) : %v", id, err)
	}

	if err := setBucketResourceData(d, importedOrg); err != nil {
		return nil, err
	}

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
