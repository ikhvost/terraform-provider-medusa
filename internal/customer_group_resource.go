package internal

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ikhvost/terraform-provider-medusa/internal/utils"

	"github.com/ikhvost/medusajs-go-sdk/medusa"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &customerGroupResource{}
	_ resource.ResourceWithConfigure   = &customerGroupResource{}
	_ resource.ResourceWithImportState = &customerGroupResource{}
)

// NewCustomerGroupResource is a helper function to simplify the provider implementation.
func NewCustomerGroupResource() resource.Resource {
	return &customerGroupResource{}
}

// customerGroupResource is the resource implementation.
type customerGroupResource struct {
	client medusa.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *customerGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customer_group"
}

// Schema defines the schema for the data source.
func (r *customerGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Customer Groups can be used to organize customers that share similar data or attributes into dedicated groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the customer group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the sales channel.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *customerGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *customerGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan customerGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()

	content, err := r.client.PostCustomerGroupsWithResponse(ctx, input)
	if d := utils.CheckCreateError("customer_group", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating customer_group",
			"Could not create customer_group, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *customerGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state customerGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value
	content, err := r.client.GetCustomerGroupsGroupWithResponse(ctx, state.ID.ValueString(), nil)
	if d := utils.CheckGetError("customer_group", state.ID.ValueString(), content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200

	// Overwrite items with refreshed state
	if err := state.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error reading Customer Group",
			"Could not read Customer Group "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *customerGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan customerGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()

	content, err := r.client.PostCustomerGroupsGroupWithResponse(ctx, plan.ID.ValueString(), input)
	if d := utils.CheckUpdateError("customer_group", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating customer_group",
			"Could not create customer_group, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *customerGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state customerGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := r.client.DeleteCustomerGroupsCustomerGroupWithResponse(ctx, state.ID.ValueString())
	if d := utils.CheckDeleteError("customer_group", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}
}

func (r *customerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
