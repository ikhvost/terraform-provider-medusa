package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
	"github.com/ikhvost/terraform-provider-medusa/internal/utils"
)

// storeResourceModel maps the resource schema data.
type storeResourceModel struct {
	ID                  types.String   `tfsdk:"id"`
	Name                types.String   `tfsdk:"name"`
	DefaultCurrencyCode types.String   `tfsdk:"default_currency_code"`
	Currencies          []types.String `tfsdk:"currencies"`
	SwapLinkTemplate    types.String   `tfsdk:"swap_link_template"`
	PaymentLinkTemplate types.String   `tfsdk:"payment_link_template"`
	InviteLinkTemplate  types.String   `tfsdk:"invite_link_template"`
}

func (m *storeResourceModel) toUpdateInput() medusa.AdminPostStoreReq {
	return medusa.AdminPostStoreReq{
		Name:                m.Name.ValueStringPointer(),
		DefaultCurrencyCode: m.DefaultCurrencyCode.ValueStringPointer(),
		Currencies:          utils.ConvertToPointerStringSlice(m.Currencies),
		SwapLinkTemplate:    m.SwapLinkTemplate.ValueStringPointer(),
		PaymentLinkTemplate: m.PaymentLinkTemplate.ValueStringPointer(),
		InviteLinkTemplate:  m.InviteLinkTemplate.ValueStringPointer(),
	}
}

func (m *storeResourceModel) toDeleteInput() medusa.AdminPostStoreReq {
	return medusa.AdminPostStoreReq{
		SwapLinkTemplate:    nil,
		PaymentLinkTemplate: nil,
		InviteLinkTemplate:  nil,
	}
}

func (m *storeResourceModel) fromGetRemote(c *medusa.AdminExtendedStoresRes) error {
	if c == nil {
		return fmt.Errorf("store is nil")
	}

	currencyCodes := utils.ExtractIDs(
		c.Store.Currencies,
		func(currency medusa.Currency) string {
			return currency.Code
		},
	)

	m.ID = types.StringValue(c.Store.Id)
	m.Name = types.StringValue(c.Store.Name)
	m.DefaultCurrencyCode = types.StringValue(c.Store.DefaultCurrencyCode)
	m.Currencies = utils.ConvertToTerraformStringSlice(currencyCodes)
	m.SwapLinkTemplate = types.StringPointerValue(c.Store.SwapLinkTemplate)
	m.PaymentLinkTemplate = types.StringPointerValue(c.Store.PaymentLinkTemplate)
	m.InviteLinkTemplate = types.StringPointerValue(c.Store.InviteLinkTemplate)

	return nil
}

func (m *storeResourceModel) fromUpdateRemote(c *medusa.AdminStoresRes) error {
	if c == nil {
		return fmt.Errorf("store is nil")
	}

	currencyCodes := utils.ExtractIDs(
		c.Store.Currencies,
		func(currency medusa.Currency) string {
			return currency.Code
		},
	)

	m.ID = types.StringValue(c.Store.Id)
	m.Name = types.StringValue(c.Store.Name)
	m.DefaultCurrencyCode = types.StringValue(c.Store.DefaultCurrencyCode)
	m.Currencies = utils.ConvertToTerraformStringSlice(currencyCodes)
	m.SwapLinkTemplate = types.StringPointerValue(c.Store.SwapLinkTemplate)
	m.PaymentLinkTemplate = types.StringPointerValue(c.Store.PaymentLinkTemplate)
	m.InviteLinkTemplate = types.StringPointerValue(c.Store.InviteLinkTemplate)

	return nil
}
