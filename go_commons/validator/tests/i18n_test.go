package tests

import (
	"context"
	"testing"

	"github.com/omniful/go_commons/validator"
)

// PaymentRequest is a sample struct for i18n testing
type PaymentRequest struct {
	Amount   float64 `json:"amount" validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"required,oneof=USD EUR SAR AED"`
	Method   string  `json:"method" validate:"required,oneof=card bank wallet"`
}

func TestI18nValidation(t *testing.T) {
	// Create an invalid payment request
	invalidPayment := PaymentRequest{
		Amount:   -10.0,
		Currency: "invalid",
		Method:   "crypto", // not a valid option
	}

	ctx := context.Background()

	// Basic validation with context
	t.Log("--- Basic Validation ---")
	customErr := validator.ValidateStruct(ctx, invalidPayment)
	if customErr.Exists() {
		// Get the error map for field-specific errors
		errorsMap := customErr.ErrorMap()
		if len(errorsMap) > 0 {
			t.Logf("Found %d validation errors", len(errorsMap))
			for field, message := range errorsMap {
				t.Logf("Field: %s, Message: %s", field, message)
			}
		}

		t.Logf("Error message: %s", customErr.ErrorMessage())
		t.Logf("Error code: %s", customErr.ErrorCode())
	}

	// Test FormatValidationErrors directly
	t.Log("\n--- Manual Error Formatting ---")
	err := validator.RawValidateStruct(ctx, invalidPayment)
	if err != nil {
		customErr := validator.FormatValidationErrors(ctx, err)
		if customErr.Exists() {
			errorsMap := customErr.ErrorMap()
			t.Logf("Found %d formatted errors", len(errorsMap))

			// Check for specific field error
			if currencyErr, ok := errorsMap["currency"]; ok {
				t.Logf("Currency error: %s", currencyErr)
			}

			t.Logf("Error message: %s", customErr.ErrorMessage())
		}
	}
}
