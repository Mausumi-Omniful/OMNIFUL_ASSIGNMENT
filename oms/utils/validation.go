package utils

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/log"
)

// CSVRowValidator handles validation of CSV rows against IMS data
type CSVRowValidator struct {
	imsClient *IMSClient
}

// ValidationResult represents the result of CSV row validation
type ValidationResult struct {
	IsValid  bool   `json:"is_valid"`
	Reason   string `json:"reason"`
	SKUValid bool   `json:"sku_valid"`
	HubValid bool   `json:"hub_valid"`
}

// NewCSVRowValidator creates a new CSV row validator
func NewCSVRowValidator(imsClient *IMSClient) *CSVRowValidator {
	return &CSVRowValidator{
		imsClient: imsClient,
	}
}

// ValidateCSVRow validates a single CSV row against IMS data
func (v *CSVRowValidator) ValidateCSVRow(ctx context.Context, row CSVRow) ValidationResult {
	result := ValidationResult{
		IsValid:  true,
		SKUValid: false,
		HubValid: false,
	}

	// Validate SKU
	skuValid, err := v.imsClient.ValidateSKU(row.SKU, row.TenantID, row.SellerID)
	if err != nil {
		log.WithError(err).Warnf("Failed to validate SKU %s for tenant %s, seller %s", row.SKU, row.TenantID, row.SellerID)
		result.IsValid = false
		result.Reason = fmt.Sprintf("SKU validation failed: %v", err)
		return result
	}
	result.SKUValid = skuValid

	// Validate Hub
	hubValid, err := v.imsClient.ValidateHub(row.Location, row.TenantID, row.SellerID)
	if err != nil {
		log.WithError(err).Warnf("Failed to validate Hub %s for tenant %s, seller %s", row.Location, row.TenantID, row.SellerID)
		result.IsValid = false
		result.Reason = fmt.Sprintf("Hub validation failed: %v", err)
		return result
	}
	result.HubValid = hubValid

	// Check if both SKU and Hub are valid
	if !skuValid || !hubValid {
		result.IsValid = false
		if !skuValid && !hubValid {
			result.Reason = fmt.Sprintf("Invalid SKU '%s' and Hub '%s' for tenant %s, seller %s", row.SKU, row.Location, row.TenantID, row.SellerID)
		} else if !skuValid {
			result.Reason = fmt.Sprintf("Invalid SKU '%s' for tenant %s, seller %s", row.SKU, row.TenantID, row.SellerID)
		} else {
			result.Reason = fmt.Sprintf("Invalid Hub '%s' for tenant %s, seller %s", row.Location, row.TenantID, row.SellerID)
		}
		return result
	}

	log.Infof("Row %d validation passed: SKU=%s, Hub=%s, Tenant=%s, Seller=%s",
		row.RowNumber, row.SKU, row.Location, row.TenantID, row.SellerID)

	return result
}
