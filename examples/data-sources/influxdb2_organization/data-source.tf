resource "influxdb2_organization" "org" {
  name        = "test-org"
  description = "Organization for test users"
}

data "influxdb2_organization" "by_name" {
  name = influxdb2_organization.org.name
}

data "influxdb2_organization" "by_id" {
  id = influxdb2_organization.org.id
}
