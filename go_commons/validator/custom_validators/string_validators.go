package custom_validators

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

// strongPasswordValidator validates that a password meets strength requirements
type strongPasswordValidator struct{}

// newStrongPasswordValidator creates a new strong password validator
func newStrongPasswordValidator() customValidator {
	return &strongPasswordValidator{}
}

// name returns the validation tag name
func (v *strongPasswordValidator) name() string {
	return "strong_password"
}

// validate performs the validation logic
func (v *strongPasswordValidator) validate(fl validator.FieldLevel) bool {
	// This would give panic if the type of the key is not string
	// So we should use custom validator only on the particular datatype defined in the custom validator
	password := fl.Field().Interface().(string)

	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSymbol bool
	for _, ch := range password {
		switch {
		case unicode.In(ch, unicode.Lu, unicode.Ll):
			hasUpper = hasUpper || unicode.IsUpper(ch)
			hasLower = hasLower || unicode.IsLower(ch)
		case unicode.In(ch, unicode.Nd):
			hasNumber = true
		case unicode.In(ch, unicode.P, unicode.S):
			hasSymbol = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSymbol
}
