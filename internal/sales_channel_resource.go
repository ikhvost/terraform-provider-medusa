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
	_ resource.Resource                = &salesChannelResource{}
	_ resource.ResourceWithConfigure   = &salesChannelResource{}
	_ resource.ResourceWithImportState = &salesChannelResource{}
)

// NewSalesChannelResource is a helper function to simplify the provider implementation.
func NewSalesChannelResource() resource.Resource {
	return &salesChannelResource{}
}

// salesChannelResource is the resource implementation.
type salesChannelResource struct {
	client medusa.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *salesChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sales_channel"
}

// Schema defines the schema for the data source.
func (r *salesChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A sales channel indicates a channel where products can be sold in.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the sales channel.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the sales channel.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the sales channel.",
				Required:    true,
			},
			"is_disabled": schema.BoolAttribute{
				Description: "Whether the sales channel is disabled.",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *salesChannelResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *salesChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan salesChannelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()

	content, err := r.client.PostSalesChannelsWithResponse(ctx, input)
	if d := utils.CheckCreateError("sales_channel", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating sales_channel",
			"Could not create sales_channel, unexpected error: "+err.Error(),
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
func (r *salesChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state salesChannelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value
	content, err := r.client.GetSalesChannelsSalesChannelWithResponse(ctx, state.ID.ValueString())
	if d := utils.CheckGetError("sales_channel", state.ID.ValueString(), content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200

	// Overwrite items with refreshed state
	if err := state.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error reading Sales Channel",
			"Could not read Sales Channel "+state.ID.ValueString()+": "+err.Error(),
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
func (r *salesChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan salesChannelResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()

	content, err := r.client.PostSalesChannelsSalesChannelWithResponse(ctx, plan.ID.ValueString(), input)
	if d := utils.CheckUpdateError("sales_channel", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating sales_channel",
			"Could not create sales_channel, unexpected error: "+err.Error(),
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
func (r *salesChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state salesChannelResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := r.client.DeleteSalesChannelsSalesChannelWithResponse(ctx, state.ID.ValueString())
	if d := utils.CheckDeleteError("sales_channel", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}
}

func (r *salesChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
