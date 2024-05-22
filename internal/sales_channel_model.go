package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
)

// salesChannelResourceModel maps the resource schema data.
type salesChannelResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsDisabled  types.Bool   `tfsdk:"is_disabled"`
}

func (m *salesChannelResourceModel) toCreateInput() medusa.AdminPostSalesChannelsReq {
	return medusa.AdminPostSalesChannelsReq{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueStringPointer(),
		IsDisabled:  m.IsDisabled.ValueBoolPointer(),
	}
}

func (m *salesChannelResourceModel) toUpdateInput() medusa.AdminPostSalesChannelsSalesChannelReq {
	return medusa.AdminPostSalesChannelsSalesChannelReq{
		Name:        m.Name.ValueStringPointer(),
		Description: m.Description.ValueStringPointer(),
		IsDisabled:  m.IsDisabled.ValueBoolPointer(),
	}
}

func (m *salesChannelResourceModel) fromRemote(c *medusa.AdminSalesChannelsRes) error {
	if c == nil {
		return fmt.Errorf("sales_channel is nil")
	}

	m.ID = types.StringValue(c.SalesChannel.Id)
	m.Name = types.StringValue(c.SalesChannel.Name)
	m.Description = types.StringPointerValue(c.SalesChannel.Description)
	m.IsDisabled = types.BoolValue(c.SalesChannel.IsDisabled)

	return nil
}
