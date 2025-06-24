package tests

import (
	"context"
	"testing"

	"github.com/omniful/go_commons/validator"
)

// TestRequestStruct demonstrates a simple request structure
type TestRequestStruct struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	Age       int    `json:"age" validate:"required,gte=18,lte=120"`
}

func TestBasicValidation(t *testing.T) {
	// Create an invalid struct
	req := TestRequestStruct{
		Email:     "invalid-email",
		Password:  "short",
		FirstName: "A",
		Age:       16,
	}

	// Validate with context
	ctx := context.Background()
	customErr := validator.ValidateStruct(ctx, req)
	if customErr.Exists() {
		// Check error map
		errorsMap := customErr.ErrorMap()
		if len(errorsMap) > 0 {
			// Print each error
			t.Log("Validation errors:")
			for field, message := range errorsMap {
				t.Logf("Field: %s, Message: %s", field, message)
			}
		}

		// Log the custom error
		t.Logf("Error code: %s", customErr.ErrorCode())
		t.Logf("Error message: %s", customErr.ErrorMessage())
	}
}
