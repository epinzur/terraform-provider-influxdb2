Terraform Provider InfluxDB2
==================

Available in the [Terraform Registry](https://registry.terraform.io/providers/rltvty/influxdb2/latest).

Note that the provider currently only supports the following resources & data sources:
* Organizations

Expect additional resources to be supported very soon.

Development has been done using InfluxDB OSS version 2.0.4.  The provider should also work on InfluxDB Cloud, but this has not been tested.  Proceed at your own risk.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.16.3

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using `make dev`. This will place the provider onto your system in a [Terraform 0.13-compliant](https://www.terraform.io/upgrade-guides/0-13.html#in-house-providers) manner.

You'll need to ensure that your Terraform file contains the information necessary to find the plugin when running `terraform init`. `make dev` will use a version number of 0.0.1, so the following block will work:

```hcl
terraform {
        required_providers {
                influxdb2 = {
                        source = "localhost/providers/rltvty/influxdb2"
                        version = "0.0.1"
                }
        }
}
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above). You will also need `docker` & `docker-compose`.

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of Acceptance tests:
* First boot the test InfluxDB server via `docker-compose up`
* In another window, run the tests with `make testacc`

## Generating Docs

From the root of the repo run `make generate`

## TODO:

* Add Warnings on Plan if organization or bucket `name` changes, similiar to the UI.
* Add additional resources & data sources for:
  * Buckets
  * Users
  * Authorizations
  * Labels
  * Tasks

## Using the provider

Please see the detailed docs for individual resource usage. Below is an
example using the InfluxDB provider to configure all resource types available:

```hcl
provider "influxdb2" {
  host     = "http://localhost:8086"
  token    = "super-secret-admin-token"
}

resource "influxdb2_organization" "org" {
  name = "test-org"
  description = "Organization for test users"
}
```
