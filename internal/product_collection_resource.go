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
	_ resource.Resource                = &productCollectionResource{}
	_ resource.ResourceWithConfigure   = &productCollectionResource{}
	_ resource.ResourceWithImportState = &productCollectionResource{}
)

// NewProductCollectionResource is a helper function to simplify the provider implementation.
func NewProductCollectionResource() resource.Resource {
	return &productCollectionResource{}
}

// productCollectionResource is the resource implementation.
type productCollectionResource struct {
	client medusa.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *productCollectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_collection"
}

// Schema defines the schema for the data source.
func (r *productCollectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A product collection is used to organize products for different purposes such as marketing or discount purposes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the product collection.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The title of the product collection.",
				Required:    true,
			},
			"handle": schema.StringAttribute{
				Description: "The handle of the product collection.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *productCollectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *productCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan productCollectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()

	content, err := r.client.PostCollectionsWithResponse(ctx, input)
	if d := utils.CheckCreateError("product_collection", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating product_collection",
			"Could not create product_collection, unexpected error: "+err.Error(),
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
func (r *productCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state productCollectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value
	content, err := r.client.GetCollectionsCollectionWithResponse(ctx, state.ID.ValueString(), nil)
	if d := utils.CheckGetError("product_collection", state.ID.ValueString(), content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200

	// Overwrite items with refreshed state
	if err := state.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error reading Product Collection",
			"Could not read Product Collection "+state.ID.ValueString()+": "+err.Error(),
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
func (r *productCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan productCollectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()

	content, err := r.client.PostCollectionsCollectionWithResponse(ctx, plan.ID.ValueString(), input)
	if d := utils.CheckUpdateError("product_collection", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating product_collection",
			"Could not create product_collection, unexpected error: "+err.Error(),
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
func (r *productCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state productCollectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := r.client.DeleteCollectionsCollectionWithResponse(ctx, state.ID.ValueString())
	if d := utils.CheckDeleteError("product_collection", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}
}

func (r *productCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
