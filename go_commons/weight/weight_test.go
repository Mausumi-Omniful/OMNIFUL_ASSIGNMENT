package weight

import (
	"fmt"
	"testing"

	error2 "github.com/omniful/go_commons/error"
	"github.com/stretchr/testify/assert"
)

func TestConvertToMinimumSellableWeight(t *testing.T) {
	tests := []struct {
		name     string
		uom      UOM
		weight   float64
		quantity uint64
		expected MinimumSellableWeight
		err      error2.CustomError
	}{
		{
			name:     "Kg to G",
			uom:      Kg,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: G, Value: 2000},
			err:      error2.CustomError{},
		},
		{
			name:     "Lbs to G",
			uom:      Lbs,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: G, Value: 907},
			err:      error2.CustomError{},
		},
		{
			name:     "G to G",
			uom:      G,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: G, Value: 2},
			err:      error2.CustomError{},
		},
		{
			name:     "L to Ml",
			uom:      L,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: Ml, Value: 2000},
			err:      error2.CustomError{},
		},
		{
			name:     "Ml to Ml",
			uom:      Ml,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: Ml, Value: 2},
			err:      error2.CustomError{},
		},
		{
			name:     "EA to EA",
			uom:      EA,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: EA, Value: 2},
			err:      error2.CustomError{},
		},
		{
			name:     "Pack to Pack",
			uom:      Pack,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: Pack, Value: 2},
			err:      error2.CustomError{},
		},
		{
			name:     "Oz to G",
			uom:      Oz,
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{UOM: G, Value: 56},
			err:      error2.CustomError{},
		},
		{
			name:     "Invalid UOM",
			uom:      "invalid",
			weight:   1,
			quantity: 2,
			expected: MinimumSellableWeight{},
			err:      error2.NewCustomError(error2.BadRequestError, "INVALID UOM"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToMinimumSellableWeight(tt.uom, tt.weight, tt.quantity)
			if !assert.Equal(t, tt.expected, result) {
				t.Log(result)
			}
			if !assert.Equal(t, tt.err.ErrorCode(), err.ErrorCode()) {
				fmt.Println("Error: ", err)
			}
		})
	}
}
func TestConvertUOM(t *testing.T) {
	tests := []struct {
		name          string
		weight        Weight
		conversionUOM UOM
		expected      Weight
		err           error2.CustomError
	}{
		{
			name:          "Kg to G",
			weight:        Weight{UOM: Kg, Value: 1},
			conversionUOM: G,
			expected:      Weight{UOM: G, Value: 1000},
			err:           error2.CustomError{},
		},
		{
			name:          "Kg to Lbs",
			weight:        Weight{UOM: Kg, Value: 1},
			conversionUOM: Lbs,
			expected:      Weight{UOM: Lbs, Value: 2.20462},
			err:           error2.CustomError{},
		},
		{
			name:          "Kg to Oz",
			weight:        Weight{UOM: Kg, Value: 1},
			conversionUOM: Oz,
			expected:      Weight{UOM: Oz, Value: 35.27396},
			err:           error2.CustomError{},
		},
		{
			name:          "Lbs to Kg",
			weight:        Weight{UOM: Lbs, Value: 1},
			conversionUOM: Kg,
			expected:      Weight{UOM: Kg, Value: 0.453592},
			err:           error2.CustomError{},
		},
		{
			name:          "Lbs to G",
			weight:        Weight{UOM: Lbs, Value: 1},
			conversionUOM: G,
			expected:      Weight{UOM: G, Value: 453.59237},
			err:           error2.CustomError{},
		},
		{
			name:          "Lbs to Oz",
			weight:        Weight{UOM: Lbs, Value: 1},
			conversionUOM: Oz,
			expected:      Weight{UOM: Oz, Value: 16},
			err:           error2.CustomError{},
		},
		{
			name:          "G to Kg",
			weight:        Weight{UOM: G, Value: 1000},
			conversionUOM: Kg,
			expected:      Weight{UOM: Kg, Value: 1},
			err:           error2.CustomError{},
		},
		{
			name:          "G to Lbs",
			weight:        Weight{UOM: G, Value: 1000},
			conversionUOM: Lbs,
			expected:      Weight{UOM: Lbs, Value: 2.20462},
			err:           error2.CustomError{},
		},
		{
			name:          "G to Oz",
			weight:        Weight{UOM: G, Value: 1000},
			conversionUOM: Oz,
			expected:      Weight{UOM: Oz, Value: 35.27396},
			err:           error2.CustomError{},
		},
		{
			name:          "L to Ml",
			weight:        Weight{UOM: L, Value: 1},
			conversionUOM: Ml,
			expected:      Weight{UOM: Ml, Value: 1000},
			err:           error2.CustomError{},
		},
		{
			name:          "Ml to L",
			weight:        Weight{UOM: Ml, Value: 1000},
			conversionUOM: L,
			expected:      Weight{UOM: L, Value: 1},
			err:           error2.CustomError{},
		},
		{
			name:          "EA to Pack",
			weight:        Weight{UOM: EA, Value: 1},
			conversionUOM: Pack,
			expected:      Weight{UOM: Pack, Value: 1},
			err:           error2.CustomError{},
		},
		{
			name:          "Pack to EA",
			weight:        Weight{UOM: Pack, Value: 1},
			conversionUOM: EA,
			expected:      Weight{UOM: EA, Value: 1},
			err:           error2.CustomError{},
		},
		{
			name:          "Oz to G",
			weight:        Weight{UOM: Oz, Value: 1},
			conversionUOM: G,
			expected:      Weight{UOM: G, Value: 28.3495},
			err:           error2.CustomError{},
		},
		{
			name:          "Oz to Lbs",
			weight:        Weight{UOM: Oz, Value: 1},
			conversionUOM: Lbs,
			expected:      Weight{UOM: Lbs, Value: 0.0625},
			err:           error2.CustomError{},
		},
		{
			name:          "Oz to Kg",
			weight:        Weight{UOM: Oz, Value: 1},
			conversionUOM: Kg,
			expected:      Weight{UOM: Kg, Value: 0.0283495},
			err:           error2.CustomError{},
		},
		{
			name:   "Weight UOM is same as conversion UOM",
			weight: Weight{UOM: Kg, Value: 1},
			expected: Weight{
				UOM:   Kg,
				Value: 1,
			},
			conversionUOM: Kg,
			err:           error2.CustomError{},
		},
		{
			name:          "Invalid UOM",
			weight:        Weight{UOM: "invalid", Value: 1},
			conversionUOM: Kg,
			expected:      Weight{},
			err:           error2.NewCustomError(error2.BadRequestError, "INVALID UOM"),
		},
		{
			name:          "Invalid Conversion UOM",
			weight:        Weight{UOM: Kg, Value: 1},
			conversionUOM: "invalid",
			expected:      Weight{},
			err:           error2.NewCustomError(error2.BadRequestError, "INVALID UOM"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertUOM(tt.weight, tt.conversionUOM)
			if !assert.Equal(t, tt.expected, result) {
				t.Log(result)
			}
			if !assert.Equal(t, tt.err.ErrorCode(), err.ErrorCode()) {
				fmt.Println("Error: ", err)
			}
		})
	}
}
func TestGetMinimumSellableUOM(t *testing.T) {
	tests := []struct {
		name     string
		input    UOM
		expected UOM
	}{
		{
			name:     "Kg to G",
			input:    Kg,
			expected: G,
		},
		{
			name:     "Lbs to G",
			input:    Lbs,
			expected: G,
		},
		{
			name:     "G to G",
			input:    G,
			expected: G,
		},
		{
			name:     "L to Ml",
			input:    L,
			expected: Ml,
		},
		{
			name:     "Ml to Ml",
			input:    Ml,
			expected: Ml,
		},
		{
			name:     "EA to EA",
			input:    EA,
			expected: EA,
		},
		{
			name:     "Pack to Pack",
			input:    Pack,
			expected: Pack,
		},
		{
			name:     "Oz to G",
			input:    Oz,
			expected: G,
		},
		{
			name:     "Invalid UOM",
			input:    "invalid",
			expected: EA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMinimumSellableUOM(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestIsValidUOM(t *testing.T) {
	tests := []struct {
		name     string
		inputUOM UOM
		skuUOM   UOM
		expected bool
	}{
		{
			name:     "Kg to Kg",
			inputUOM: Kg,
			skuUOM:   Kg,
			expected: true,
		},
		{
			name:     "Kg to Lbs",
			inputUOM: Kg,
			skuUOM:   Lbs,
			expected: true,
		},
		{
			name:     "Kg to G",
			inputUOM: Kg,
			skuUOM:   G,
			expected: true,
		},
		{
			name:     "Kg to Oz",
			inputUOM: Kg,
			skuUOM:   Oz,
			expected: true,
		},
		{
			name:     "Kg to Ml",
			inputUOM: Kg,
			skuUOM:   Ml,
			expected: false,
		},
		{
			name:     "G to Kg",
			inputUOM: G,
			skuUOM:   Kg,
			expected: true,
		},
		{
			name:     "G to Lbs",
			inputUOM: G,
			skuUOM:   Lbs,
			expected: true,
		},
		{
			name:     "G to G",
			inputUOM: G,
			skuUOM:   G,
			expected: true,
		},
		{
			name:     "G to Oz",
			inputUOM: G,
			skuUOM:   Oz,
			expected: true,
		},
		{
			name:     "L to Ml",
			inputUOM: L,
			skuUOM:   Ml,
			expected: true,
		},
		{
			name:     "Ml to L",
			inputUOM: Ml,
			skuUOM:   L,
			expected: true,
		},
		{
			name:     "EA to Pack",
			inputUOM: EA,
			skuUOM:   Pack,
			expected: true,
		},
		{
			name:     "Pack to EA",
			inputUOM: Pack,
			skuUOM:   EA,
			expected: true,
		},
		{
			name:     "Oz to Kg",
			inputUOM: Oz,
			skuUOM:   Kg,
			expected: true,
		},
		{
			name:     "Invalid UOM",
			inputUOM: "invalid",
			skuUOM:   Kg,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUOM(tt.inputUOM, tt.skuUOM)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestIsUOMWeighted(t *testing.T) {
	tests := []struct {
		name     string
		input    UOM
		expected bool
	}{
		{
			name:     "Kg is weighted",
			input:    Kg,
			expected: true,
		},
		{
			name:     "Lbs is weighted",
			input:    Lbs,
			expected: true,
		},
		{
			name:     "G is weighted",
			input:    G,
			expected: true,
		},
		{
			name:     "Ml is weighted",
			input:    Ml,
			expected: true,
		},
		{
			name:     "L is weighted",
			input:    L,
			expected: true,
		},
		{
			name:     "Oz is weighted",
			input:    Oz,
			expected: true,
		},
		{
			name:     "EA is not weighted",
			input:    EA,
			expected: false,
		},
		{
			name:     "Pack is not weighted",
			input:    Pack,
			expected: false,
		},
		{
			name:     "Invalid UOM is not weighted",
			input:    "invalid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUOMWeighted(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
