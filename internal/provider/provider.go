package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/domain"
)

func init() {
	/// descriptions are written in markdown for docs
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"host": {
					Description: "The host url where influxDB2 lives. Can also be set using the `INFLUX_HOST` environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("INFLUX_HOST", nil),
				},
				"token": {
					Description: "An auth token that has the nesecary permissions to read-from and/or write-to InfluxDB2. Ideally this should be set using the `INFLUX_TOKEN` environment variable, so that the secret is not saved to source control.",
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("INFLUX_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"influxdb2_organization": dataSourceOrganization(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"influxdb2_organization": resourceOrganization(),
			},
		}

		p.ConfigureContextFunc = providerConfigure(version, p)

		return p
	}
}

type metaData struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
	client influxdb2.Client
}

func providerConfigure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		host := d.Get("host").(string)
		token := d.Get("token").(string)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		if host == "" || token == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create InfluxDB2 client",
				Detail:   "Unable to auth authenticated InfluxDB2 client",
			})
			return nil, diags
		}

		options := influxdb2.DefaultOptions()
		//0 error, 1 - warning, 2 - info, 3 - debug
		options.SetLogLevel(3)

		// userAgent := p.UserAgent("terraform-provider-influxdb2", version)
		// TODO: options.UserAgent = userAgent

		client := influxdb2.NewClientWithOptions(host, token, options)

		ok, err := client.Ready(ctx)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			client.Close()
			return nil, diags
		}
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "InfluxDB2 Server is not ready",
				Detail:   "InfluxDB2 Server is not ready",
			})
			client.Close()
			return nil, diags
		}

		check, err := client.Health(ctx)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
			client.Close()
			return nil, diags
		}
		if check.Status != domain.HealthCheckStatusPass {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "InfluxDB2 Health is not passing",
				Detail:   "InfluxDB2 Health is not passing",
			})
			client.Close()
			return nil, diags
		}

		md := &metaData{
			client: client,
		}

		return md, nil
	}
}
