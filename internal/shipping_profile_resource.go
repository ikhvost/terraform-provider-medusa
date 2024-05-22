package internal

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ikhvost/terraform-provider-medusa/internal/utils"

	"github.com/ikhvost/medusajs-go-sdk/medusa"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &shippingProfileResource{}
	_ resource.ResourceWithConfigure   = &shippingProfileResource{}
	_ resource.ResourceWithImportState = &shippingProfileResource{}
)

// NewShippingProfileResource is a helper function to simplify the provider implementation.
func NewShippingProfileResource() resource.Resource {
	return &shippingProfileResource{}
}

// shippingProfileResource is the resource implementation.
type shippingProfileResource struct {
	client medusa.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *shippingProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shipping_profile"
}

// Schema defines the schema for the data source.
func (r *shippingProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A shipping profile is used to group products that can be shipped in the same manner.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the shipping profile.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the shipping profile.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the shipping profile.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("default", "gift_card", "custom"),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *shippingProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *shippingProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan shippingProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()

	content, err := r.client.PostShippingProfilesWithResponse(ctx, input)
	if d := utils.CheckCreateError("shipping_profile", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating shipping_profile",
			"Could not create shipping_profile, unexpected error: "+err.Error(),
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
func (r *shippingProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state shippingProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value
	content, err := r.client.GetShippingProfilesProfileWithResponse(ctx, state.ID.ValueString())
	if d := utils.CheckGetError("shipping_profile", state.ID.ValueString(), content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200

	// Overwrite items with refreshed state
	if err := state.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error reading Shipping Profile",
			"Could not read Shipping Profile "+state.ID.ValueString()+": "+err.Error(),
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
func (r *shippingProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan shippingProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()

	content, err := r.client.PostShippingProfilesProfileWithResponse(ctx, plan.ID.ValueString(), input)
	if d := utils.CheckUpdateError("shipping_profile", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating shipping_profile",
			"Could not create shipping_profile, unexpected error: "+err.Error(),
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
func (r *shippingProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state shippingProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := r.client.DeleteShippingProfilesProfileWithResponse(ctx, state.ID.ValueString())
	if d := utils.CheckDeleteError("shipping_profile", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}
}

func (r *shippingProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
