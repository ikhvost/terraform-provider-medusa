package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
)

func ConvertToStringSlice(slice []types.String) []string {
	if slice == nil {
		return nil
	}

	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = v.ValueString()
	}

	return result
}

func ConvertToPointerStringSlice(slice []types.String) *[]string {
	if slice == nil {
		return nil
	}

	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = v.ValueString()
	}

	return &result
}

func ConvertToTerraformStringSlice(input []string) []types.String {
	result := make([]types.String, len(input))
	for i, v := range input {
		result[i] = types.StringValue(v)
	}
	return result
}

func ConvertToFloat32(n types.Number) float32 {
	if n.IsUnknown() || n.IsNull() {
		return 0.0
	}
	bigFloat := n.ValueBigFloat()
	floatVal, _ := bigFloat.Float32()
	return floatVal
}

func ConvertToPointerFloat32(n types.Number) *float32 {
	if n.IsUnknown() || n.IsNull() {
		return nil
	}
	bigFloat := n.ValueBigFloat()
	floatVal, _ := bigFloat.Float32()
	return &floatVal
}

func ConvertToTerraformNumber(value float32) types.Number {
	return types.NumberValue(big.NewFloat(float64(value)))
}
