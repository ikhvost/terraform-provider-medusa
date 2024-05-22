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
	_ resource.Resource                = &productCategoryResource{}
	_ resource.ResourceWithConfigure   = &productCategoryResource{}
	_ resource.ResourceWithImportState = &productCategoryResource{}
)

// NewProductCategoryResource is a helper function to simplify the provider implementation.
func NewProductCategoryResource() resource.Resource {
	return &productCategoryResource{}
}

// productCategoryResource is the resource implementation.
type productCategoryResource struct {
	client medusa.ClientWithResponsesInterface
}

// Metadata returns the data source type name.
func (r *productCategoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_category"
}

// Schema defines the schema for the data source.
func (r *productCategoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Products can be categoriezed into categories. A product can be added into more than one category.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the product category.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the product category.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the product category.",
				Computed:    true,
				Optional:    true,
			},
			"handle": schema.StringAttribute{
				Description: "The handle of the product category.",
				Computed:    true,
				Optional:    true,
			},
			"is_internal": schema.BoolAttribute{
				Description: "If set to true, the product category will only be available to admins.",
				Computed:    true,
				Optional:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "If set to false, the product category will not be available in the storefront.",
				Computed:    true,
				Optional:    true,
			},
			"parent_category_id": schema.StringAttribute{
				Description: "The id of the parent product category.",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *productCategoryResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = utils.GetClient(req.ProviderData)
}

// Create creates the resource and sets the initial Terraform state.
func (r *productCategoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan productCategoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toCreateInput()

	content, err := r.client.PostProductCategoriesWithResponse(ctx, nil, input)
	if d := utils.CheckCreateError("product_category", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating product_category",
			"Could not create product_category, unexpected error: "+err.Error(),
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
func (r *productCategoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state productCategoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed value
	content, err := r.client.GetProductCategoriesCategoryWithResponse(ctx, state.ID.ValueString(), nil)
	if d := utils.CheckGetError("product_category", state.ID.ValueString(), content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200

	// Overwrite items with refreshed state
	if err := state.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error reading Product Category",
			"Could not read Product Category "+state.ID.ValueString()+": "+err.Error(),
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
func (r *productCategoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan productCategoryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	input := plan.toUpdateInput()

	content, err := r.client.PostProductCategoriesCategoryWithResponse(ctx, plan.ID.ValueString(), nil, input)
	if d := utils.CheckUpdateError("product_category", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}

	resource := content.JSON200
	tflog.Debug(ctx, spew.Sdump(resource))

	// Map response body to schema
	if err := plan.fromRemote(resource); err != nil {
		resp.Diagnostics.AddError(
			"Error creating product_category",
			"Could not create product_category, unexpected error: "+err.Error(),
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
func (r *productCategoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state productCategoryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := r.client.DeleteProductCategoriesCategoryWithResponse(ctx, state.ID.ValueString())
	if d := utils.CheckDeleteError("product_category", content, err); d != nil {
		resp.Diagnostics.Append(d)
		return
	}
}

func (r *productCategoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
