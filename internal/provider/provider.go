package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/eagletmt/terraform-provider-kuroko2/internal/kuroko2"
)

var (
	_ provider.Provider = (*kuroko2Provider)(nil)
)

type kuroko2Provider struct{}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &kuroko2Provider{}
	}
}

func (p *kuroko2Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kuroko2"
}

func (p *kuroko2Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"apikey": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

type kuroko2ProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Apikey   types.String `tfsdk:"apikey"`
}

type kuroko2ProviderData struct {
	client kuroko2.Client
}

func (p *kuroko2Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data kuroko2ProviderModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := kuroko2.NewClient(data.Endpoint.ValueString(), data.Username.ValueString(), data.Apikey.ValueString())
	resp.ResourceData = &kuroko2ProviderData{client: client}
}

func (p *kuroko2Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *kuroko2Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJobDefinitionResource,
	}
}
