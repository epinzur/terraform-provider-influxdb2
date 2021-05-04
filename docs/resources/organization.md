---
page_title: "influxdb2_organization Resource - terraform-provider-influxdb2"
subcategory: ""
description: |-
  The Organization resource allows you to configure a InfluxDB2 Organization.
---

# Resource `influxdb2_organization`

The Organization resource allows you to configure a InfluxDB2 Organization.

## Example Usage

```terraform
resource "influxdb2_organization" "org" {
  name        = "test-org"
  description = "Organization for test users"
}
```

## Schema

### Required

- **name** (String) Name of the Organization.

### Optional

- **description** (String) The description of the Organization.

### Read-only

- **created_at** (String) The string time that the Organization was created.
- **created_timestamp** (Number) The timestamp that the Organization was created.
- **id** (String) ID of the Organization.
- **updated_at** (String) The string time that the Organization was last updated.
- **updated_timestamp** (Number) The timestamp that the Organization was last updated.

## Import

Import is supported using the following syntax:

```shell
terraform import influxdb2_organization.org <my-id>
```
