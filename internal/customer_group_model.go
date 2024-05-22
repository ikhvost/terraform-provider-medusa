package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
)

// customerGroupResourceModel maps the resource schema data.
type customerGroupResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (m *customerGroupResourceModel) toCreateInput() medusa.AdminPostCustomerGroupsReq {
	return medusa.AdminPostCustomerGroupsReq{
		Name: m.Name.ValueString(),
	}
}

func (m *customerGroupResourceModel) toUpdateInput() medusa.AdminPostCustomerGroupsGroupReq {
	return medusa.AdminPostCustomerGroupsGroupReq{
		Name: m.Name.ValueStringPointer(),
	}
}

func (m *customerGroupResourceModel) fromRemote(c *medusa.AdminCustomerGroupsRes) error {
	if c == nil {
		return fmt.Errorf("customer_group is nil")
	}

	m.ID = types.StringValue(c.CustomerGroup.Id)
	m.Name = types.StringValue(c.CustomerGroup.Name)

	return nil
}
