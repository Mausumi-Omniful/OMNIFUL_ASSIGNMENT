package validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	customError "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/i18n"
)

// ValidationError represents a custom validation error with formatted message
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors []ValidationError

// FormatValidationErrors converts validator.ValidationErrors to our custom format
func FormatValidationErrors(ctx context.Context, err error) customError.CustomError {
	var errors ValidationErrors

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			// Format the field name using the namespace, which includes full path for nested fields
			// For example: Parent.Child or Array[0].Field
			field := formatFieldName(e)

			errors = append(errors, ValidationError{
				Field:   field,
				Tag:     e.Tag(),
				Value:   fmt.Sprintf("%v", e.Value()),
				Message: strings.TrimSpace(formatErrorMessageWithContext(ctx, field, e)),
			})
		}
	}

	return errors.ToCustomError(ctx)
}

// formatFieldName creates a properly formatted field name using dot notation for nested fields
// and array indices for slice elements
func formatFieldName(e validator.FieldError) string {
	// Get the namespace which includes the full path (Parent.Child)
	namespace := e.Namespace()

	// Skip the first part of the namespace which is the struct name
	parts := strings.Split(namespace, ".")
	if len(parts) <= 1 {
		return e.Field()
	}

	// Remove the struct name (first element)
	parts = parts[1:]
	return strings.Join(parts, ".")
}

// formatErrorMessageWithContext generates a localized error message
func formatErrorMessageWithContext(ctx context.Context, field string, e validator.FieldError) string {
	tag := e.Tag()
	value := e.Value()
	param := e.Param()

	// Try to get localized message
	messageKey := fmt.Sprintf("validation.%s", tag)
	message := i18n.Translate(ctx, messageKey)

	// If no translation found, use a simple format
	if message == messageKey || message == "" {
		// Default simple message with just field and tag information
		message = fmt.Sprintf("%s failed %s validation", field, tag)

		// For common validators, add the parameter
		if value != "" {
			message = fmt.Sprintf("%s %s %s", field, tag, value)
		}
	}

	// Replace placeholders
	message = strings.ReplaceAll(message, "{field}", field)
	message = strings.ReplaceAll(message, "{value}", fmt.Sprint(value))
	message = strings.ReplaceAll(message, "{tag}", tag)
	message = strings.ReplaceAll(message, "{param}", param)

	return message
}

// GetFirstError returns the first validation error from ValidationErrors
func (ve ValidationErrors) GetFirstError() *ValidationError {
	if len(ve) > 0 {
		return &ve[0]
	}
	return nil
}

// GetFieldError returns the validation error for a specific field
func (ve ValidationErrors) GetFieldError(field string) *ValidationError {
	for i, e := range ve {
		if e.Field == field {
			return &ve[i]
		}
	}
	return nil
}

// ByField returns all errors for a specific field
func (ve ValidationErrors) ByField(field string) []ValidationError {
	var fieldErrors []ValidationError
	for _, e := range ve {
		if e.Field == field {
			fieldErrors = append(fieldErrors, e)
		}
	}
	return fieldErrors
}

// ToCustomError converts ValidationErrors to CustomError from internal error package
func (ve ValidationErrors) ToCustomError(ctx context.Context) customError.CustomError {
	// Get validation failed message with i18n support
	message := i18n.Translate(ctx, "validation.failed")
	if message == "validation.failed" || message == "" {
		message = "Request is Invalid"
	}

	// Also add to error data for API responses
	errorData := make(map[string]string, 0)
	for _, e := range ve {
		errorData[e.Field] = e.Message
	}

	return customError.RequestInvalidError(message, customError.WithErrors(errorData))
}
