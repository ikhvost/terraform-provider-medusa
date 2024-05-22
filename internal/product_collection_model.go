package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
)

// productCollectionResourceModel maps the resource schema data.
type productCollectionResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Title  types.String `tfsdk:"title"`
	Handle types.String `tfsdk:"handle"`
}

func (m *productCollectionResourceModel) toCreateInput() medusa.AdminPostCollectionsReq {
	return medusa.AdminPostCollectionsReq{
		Title:  m.Title.ValueString(),
		Handle: m.Handle.ValueStringPointer(),
	}
}

func (m *productCollectionResourceModel) toUpdateInput() medusa.AdminPostCollectionsCollectionReq {
	return medusa.AdminPostCollectionsCollectionReq{
		Title:  m.Title.ValueStringPointer(),
		Handle: m.Handle.ValueStringPointer(),
	}
}

func (m *productCollectionResourceModel) fromRemote(c *medusa.AdminCollectionsRes) error {
	if c == nil {
		return fmt.Errorf("product_collection is nil")
	}

	m.ID = types.StringValue(c.Collection.Id)
	m.Title = types.StringValue(c.Collection.Title)
	m.Handle = types.StringPointerValue(c.Collection.Handle)

	return nil
}
