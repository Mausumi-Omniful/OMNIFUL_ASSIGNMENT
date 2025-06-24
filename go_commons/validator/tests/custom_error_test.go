package tests

import (
	"context"
	"testing"

	customError "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/validator"
)

// Product is a sample struct for validation
type Product struct {
	Name        string  `json:"name" validate:"required,min=3,max=200"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description" validate:"required,min=10"`
	Category    string  `json:"category" validate:"required,oneof=electronics clothing food books other"`
}

func TestCustomErrorIntegration(t *testing.T) {
	// Create an invalid product
	invalidProduct := Product{
		Name:        "AB",        // Too short
		Price:       0,           // Must be > 0
		Description: "Too short", // Min 10 chars
		Category:    "invalid",   // Not in allowed list
	}

	ctx := context.Background()

	// Example 1: Using ValidateStruct
	t.Log("--- Example 1: ValidateStructWithError ---")
	customErr := validator.ValidateStruct(ctx, invalidProduct)
	if customErr.Exists() {
		// This is now a CustomError from the error package
		t.Logf("Error code: %s", customErr.ErrorCode())
		t.Logf("Error message: %s", customErr.ErrorMessage())

		// We can add more context to the error
		customErr = customErr.WithParam("additional_info", "This request failed validation")

		// Get the error map to see field-specific errors
		errorsMap := customErr.ErrorMap()
		if len(errorsMap) > 0 {
			t.Logf("Found %d validation errors", len(errorsMap))
			for field, message := range errorsMap {
				t.Logf("- Field: %s, Error: %s", field, message)
			}
		}
	}

	// Example 2: Using error conversion with raw validation
	t.Log("\n--- Example 2: Manual error conversion ---")
	err := validator.RawValidateStruct(ctx, invalidProduct)
	if err != nil {
		// Convert to custom error
		customErr := validator.FormatValidationErrors(ctx, err)

		// Check error map
		errorsMap := customErr.ErrorMap()
		if nameErr, ok := errorsMap["name"]; ok {
			t.Logf("Name error: %s", nameErr)
		}

		t.Logf("Error code: %s", customErr.ErrorCode())
		t.Logf("Error message: %s", customErr.ErrorMessage())
	}

	// Example 3: Service function pattern
	result, err := createProductService(ctx, invalidProduct)
	if err != nil {
		if customErr, ok := err.(customError.CustomError); ok {
			t.Logf("\n--- Example 3: Service function pattern ---")
			t.Logf("Service returned CustomError with code: %s", customErr.ErrorCode())
		}
	} else {
		t.Logf("Product created: %v", result)
	}
}

// Example service function that uses validation
func createProductService(ctx context.Context, product Product) (*Product, error) {
	// Validate the product
	customErr := validator.ValidateStruct(ctx, product)
	if customErr.Exists() {
		// Return the validation error directly
		return nil, customErr
	}

	// Simulate creating the product
	return &product, nil
}
