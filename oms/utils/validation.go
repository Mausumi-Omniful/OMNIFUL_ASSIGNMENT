package utils

import (
	"context"
	"fmt"
)

type CSVRowValidator struct {
	imsClient *IMSClient
}

type ValidationResult struct {
	IsValid  bool   `json:"is_valid"`
	Reason   string `json:"reason"`
	SKUValid bool   `json:"sku_valid"`
	HubValid bool   `json:"hub_valid"`
}




func NewCSVRowValidator(imsClient *IMSClient) *CSVRowValidator {
	return &CSVRowValidator{
		imsClient: imsClient,
	}
}




func (v *CSVRowValidator) ValidateCSVRow(ctx context.Context, row CSVRow) ValidationResult {
	result := ValidationResult{
		IsValid:  true,
		SKUValid: false,
		HubValid: false,
	}

	skuValid, err := v.imsClient.ValidateSKU(row.SKU, row.TenantID, row.SellerID)
	if err != nil {
		fmt.Printf("SKU check failed for row %d\n", row.RowNumber)
		result.IsValid = false
		result.Reason = fmt.Sprintf("SKU error: %v", err)
		return result
	}
	result.SKUValid = skuValid

	hubValid, err := v.imsClient.ValidateHub(row.Location, row.TenantID, row.SellerID)
	if err != nil {
		fmt.Printf("Hub check failed for row %d\n", row.RowNumber)
		result.IsValid = false
		result.Reason = fmt.Sprintf("Hub error: %v", err)
		return result
	}
	result.HubValid = hubValid

	if !skuValid || !hubValid {
		result.IsValid = false
		switch {
		case !skuValid && !hubValid:
			result.Reason = "Invalid SKU & Hub"
		case !skuValid:
			result.Reason = "Invalid SKU"
		case !hubValid:
			result.Reason = "Invalid Hub"
		}
		fmt.Printf("Invalid row %d: %s\n", row.RowNumber, result.Reason)
		return result
	}

	fmt.Printf("Valid row %d\n", row.RowNumber)
	return result
}
