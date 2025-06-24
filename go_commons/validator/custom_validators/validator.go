package custom_validators

import (
	"github.com/go-playground/validator/v10"
)

// customValidator interface defines methods that all custom validators must implement
type customValidator interface {
	// name returns the validation tag name
	name() string

	// validate performs the validation logic
	validate(fl validator.FieldLevel) bool
}

// registry holds all the validators to be registered
var registry []customValidator

// registerValidator adds a validator to the registry
// Returns a boolean to allow usage with global variable initialization
func registerValidator(validator customValidator) bool {
	registry = append(registry, validator)
	return true
}

// Register validators using global variables
var (
	// String validators
	_ = registerValidator(newStrongPasswordValidator())
)

// RegisterAll registers all custom validators with the validator engine
func RegisterAll(v *validator.Validate) error {
	// Register all validators in the registry with empty value check
	for _, validator := range registry {
		err := v.RegisterValidation(validator.name(), validator.validate)
		if err != nil {
			return err
		}
	}
	return nil
}
