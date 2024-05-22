package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
)

// shippingProfileResourceModel maps the resource schema data.
type shippingProfileResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func (m *shippingProfileResourceModel) toCreateInput() medusa.AdminPostShippingProfilesReq {
	return medusa.AdminPostShippingProfilesReq{
		Name: m.Name.ValueString(),
		Type: m.Type.ValueString(),
	}
}

func (m *shippingProfileResourceModel) toUpdateInput() medusa.AdminPostShippingProfilesProfileReq {
	return medusa.AdminPostShippingProfilesProfileReq{
		Name: m.Name.ValueStringPointer(),
		Type: m.Type.ValueStringPointer(),
	}
}

func (m *shippingProfileResourceModel) fromRemote(c *medusa.AdminShippingProfilesRes) error {
	if c == nil {
		return fmt.Errorf("shipping_profile is nil")
	}

	m.ID = types.StringValue(c.ShippingProfile.Id)
	m.Name = types.StringValue(c.ShippingProfile.Name)
	m.Type = types.StringValue(c.ShippingProfile.Type)

	return nil
}
