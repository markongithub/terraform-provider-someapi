package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &groupResource{}
	_ resource.ResourceWithConfigure = &groupResource{}
)

// NewGroupResource is a helper function to simplify the provider implementation.
func NewGroupResource() resource.Resource {
	return &groupResource{}
}

// groupResource is the resource implementation.
type groupResource struct {
	client *APIClient
}

type groupResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Configure adds the provider configured client to the resource.
func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	r.client = client
}

// Metadata returns the resource type name.
func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema defines the schema for the resource.
func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		// TODO: merge this with the datasource one
		MarkdownDescription: "Group resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Example identifer",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

/*
curl -X POST \
  --url 'https://team2.someapi.cloud/api/rest/2.0/groups/create'  \
	  -H 'Accept: application/json' \
		  -H 'Content-Type: application/json' \
			  -H 'Authorization: Bearer bear' \
				  --data-raw '{
						  "name": "name6",
							  "display_name": "display_name6",
								  "type": "LOCAL_GROUP",
									  "visibility": "SHARABLE"
									}'
*/
// Create creates the resource and sets the initial Terraform state.
func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	postData, _ := json.Marshal(map[string]interface{}{
		"name":         plan.Name.ValueString(),
		"display_name": plan.Name.ValueString(),
		"description":  plan.Description.ValueString(),
		"type":         "LOCAL_GROUP",
		"visibility":   "NON_SHARABLE",
	})

	respBytes, err := r.client.Post(ctx, "/groups/create", postData, "200 OK")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			"Could not create group, unexpected error: "+err.Error(),
		)
		return
	}

	var groupReturned groupSearchResult
	err = json.Unmarshal(respBytes, &groupReturned)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing JSON",
			"Could not parse create JSON, unexpected error: "+err.Error(),
		)
		return
	}
	plan.Name = types.StringValue(groupReturned.Name)
	plan.Description = types.StringValue(groupReturned.Description)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed group value from API
	lookupResult, err := lookupGroupByName(ctx, r.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error looking up group",
			"Error looking up group: "+err.Error())
		return
	}
	state.Name = lookupResult.Name
	state.Description = lookupResult.Description

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

/*
curl -X POST \
  --url 'https://team2.someapi.cloud/api/rest/2.0/groups/deletethisplease/update'  \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer no' \
  --data-raw '{
  "operation": "ADD",
  "description": "foo",
  "name": "deletethisplease"
}'
*/
// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	postData, _ := json.Marshal(map[string]interface{}{
		"operation":   "ADD",
		"name":        plan.Name.ValueString(),
		"description": plan.Description.ValueString(),
	})

	// won't work if you want to rename the ord. need to store and use the ID.
	_, err := r.client.Post(ctx, fmt.Sprintf("/groups/%s/update", plan.Name.ValueString()), postData, "204 No Content")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			"Could not create group, unexpected error: "+err.Error(),
		)
		return
	}
	// /update returns no content so we need to do another read here
	lookupResult, err := lookupGroupByName(ctx, r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error looking up group",
			"Error looking up group: "+err.Error())
		return
	}
	plan.Name = lookupResult.Name
	plan.Description = lookupResult.Description
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

/*
curl -X POST \
  --url 'https://team2.someapi.cloud/api/rest/2.0/groups/deletemeplease/delete'  \
  -H 'Authorization: Bearer no'
  returns 204
*/
// Delete deletes the resource and removes the Terraform state on success.
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var emptyPostData []byte

	_, err := r.client.Post(ctx, fmt.Sprintf("/groups/%s/delete", state.Name.ValueString()), emptyPostData, "204 No Content")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting group",
			"Could not delete group, unexpected error: "+err.Error(),
		)
		return
	}
	// TODO: If there is an error in the JSON output, we should surface that.
}
