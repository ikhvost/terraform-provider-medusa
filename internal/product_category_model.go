package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
)

// productCategoryResourceModel maps the resource schema data.
type productCategoryResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Handle           types.String `tfsdk:"handle"`
	IsInternal       types.Bool   `tfsdk:"is_internal"`
	IsActive         types.Bool   `tfsdk:"is_active"`
	ParentCategoryId types.String `tfsdk:"parent_category_id"`
}

func (m *productCategoryResourceModel) toCreateInput() medusa.AdminPostProductCategoriesReq {
	return medusa.AdminPostProductCategoriesReq{
		Name:             m.Name.ValueString(),
		Description:      m.Description.ValueStringPointer(),
		Handle:           m.Handle.ValueStringPointer(),
		IsInternal:       m.IsInternal.ValueBoolPointer(),
		IsActive:         m.IsActive.ValueBoolPointer(),
		ParentCategoryId: m.ParentCategoryId.ValueStringPointer(),
	}
}

func (m *productCategoryResourceModel) toUpdateInput() medusa.AdminPostProductCategoriesCategoryReq {
	return medusa.AdminPostProductCategoriesCategoryReq{
		Name:             m.Name.ValueStringPointer(),
		Description:      m.Description.ValueStringPointer(),
		Handle:           m.Handle.ValueStringPointer(),
		IsInternal:       m.IsInternal.ValueBoolPointer(),
		IsActive:         m.IsActive.ValueBoolPointer(),
		ParentCategoryId: m.ParentCategoryId.ValueStringPointer(),
	}
}

func (m *productCategoryResourceModel) fromRemote(c *medusa.AdminProductCategoriesCategoryRes) error {
	if c == nil {
		return fmt.Errorf("product_category is nil")
	}

	m.ID = types.StringValue(c.ProductCategory.Id)
	m.Name = types.StringValue(c.ProductCategory.Name)
	m.Description = types.StringPointerValue(c.ProductCategory.Description)
	m.Handle = types.StringPointerValue(&c.ProductCategory.Handle)
	m.IsInternal = types.BoolValue(c.ProductCategory.IsInternal)
	m.IsActive = types.BoolValue(c.ProductCategory.IsActive)
	m.ParentCategoryId = types.StringPointerValue(c.ProductCategory.ParentCategoryId)

	return nil
}
