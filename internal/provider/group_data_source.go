package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	client *APIClient
}

type groupDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
}

// do I really have to do this
type groupSearchResult struct {
	Name        string
	ID          string
	Description string
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		// TODO: merge this with the resource one
		MarkdownDescription: "Group data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Example identifer",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID generated on the server side",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *groupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	tflog.Trace(ctx, "Entered groupDataSource.Configure")
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("I don't know what I expected, but I got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Trace(ctx, "Entered groupDataSource.Read")
	var state groupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	lookupResult, err := lookupGroupByIdentifier(ctx, d.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error looking up group",
			"Error looking up group: "+err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &lookupResult)...)
}
