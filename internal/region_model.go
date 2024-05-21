package internal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ikhvost/medusajs-go-sdk/medusa"
	"github.com/ikhvost/terraform-provider-medusa/internal/utils"
)

// regionResourceModel maps the resource schema data.
type regionResourceModel struct {
	ID                   types.String   `tfsdk:"id"`
	Name                 types.String   `tfsdk:"name"`
	CurrencyCode         types.String   `tfsdk:"currency_code"`
	TaxRate              types.Number   `tfsdk:"tax_rate"`
	PaymentProviders     []types.String `tfsdk:"payment_providers"`
	FulfillmentProviders []types.String `tfsdk:"fulfillment_providers"`
	Countries            []types.String `tfsdk:"countries"`
	TaxCode              types.String   `tfsdk:"tax_code"`
	IncludesTax          types.Bool     `tfsdk:"includes_tax"`
}

func (m *regionResourceModel) toCreateInput() medusa.AdminPostRegionsReq {
	return medusa.AdminPostRegionsReq{
		Name:                 m.Name.ValueString(),
		CurrencyCode:         m.CurrencyCode.ValueString(),
		TaxRate:              utils.ConvertToFloat32(m.TaxRate),
		PaymentProviders:     utils.ConvertToStringSlice(m.PaymentProviders),
		FulfillmentProviders: utils.ConvertToStringSlice(m.FulfillmentProviders),
		Countries:            utils.ConvertToStringSlice(m.Countries),
		TaxCode:              m.TaxCode.ValueStringPointer(),
		IncludesTax:          m.IncludesTax.ValueBoolPointer(),
	}
}
func (m *regionResourceModel) toUpdateInput() medusa.AdminPostRegionsRegionReq {
	return medusa.AdminPostRegionsRegionReq{
		Name:                 m.Name.ValueStringPointer(),
		CurrencyCode:         m.CurrencyCode.ValueStringPointer(),
		TaxRate:              utils.ConvertToPointerFloat32(m.TaxRate),
		PaymentProviders:     utils.ConvertToPointerStringSlice(m.PaymentProviders),
		FulfillmentProviders: utils.ConvertToPointerStringSlice(m.FulfillmentProviders),
		Countries:            utils.ConvertToPointerStringSlice(m.Countries),
		TaxCode:              m.TaxCode.ValueStringPointer(),
		IncludesTax:          m.IncludesTax.ValueBoolPointer(),
	}
}

func (m *regionResourceModel) fromRemote(c *medusa.AdminRegionsRes) error {
	if c == nil {
		return fmt.Errorf("region is nil")
	}

	fulfillmentIDs := utils.ExtractIDs(
		c.Region.FulfillmentProviders,
		func(fp medusa.FulfillmentProvider) string {
			return fp.Id
		},
	)
	paymentIDs := utils.ExtractIDs(
		c.Region.PaymentProviders,
		func(pp medusa.PaymentProvider) string {
			return pp.Id
		},
	)
	countryIDs := utils.ExtractIDs(
		c.Region.Countries,
		func(country medusa.Country) string {
			return country.Iso2
		},
	)

	m.ID = types.StringValue(c.Region.Id)
	m.Name = types.StringValue(c.Region.Name)
	m.CurrencyCode = types.StringValue(c.Region.CurrencyCode)
	m.FulfillmentProviders = utils.ConvertToTerraformStringSlice(fulfillmentIDs)
	m.PaymentProviders = utils.ConvertToTerraformStringSlice(paymentIDs)
	m.Countries = utils.ConvertToTerraformStringSlice(countryIDs)
	m.TaxRate = utils.ConvertToTerraformNumber(c.Region.TaxRate)
	m.TaxCode = types.StringPointerValue(c.Region.TaxCode)
	m.IncludesTax = types.BoolPointerValue(c.Region.IncludesTax)

	return nil
}
