package validator

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	customError "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/i18n"
	customValidators "github.com/omniful/go_commons/validator/custom_validators"
)

var (
	validate *validator.Validate
	once     sync.Once
)


type GinValidatorImpl struct {
	validator *validator.Validate
}

func NewGinValidator() *GinValidatorImpl {
	return &GinValidatorImpl{
		validator: GetValidator(),
	}
}

type SliceValidationError []error

// Error concatenates all error elements in SliceValidationError into a single string separated by \n.
func (err SliceValidationError) Error() string {
	n := len(err)
	switch n {
	case 0:
		return ""
	default:
		var b strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&b, "[%d]: %s", 0, err[0].Error())
		}
		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					b.WriteString("\n")
					fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
				}
			}
		}
		return b.String()
	}
}

func (v *GinValidatorImpl) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		if value.Elem().Kind() != reflect.Struct {
			return v.ValidateStruct(value.Elem().Interface())
		}
		return v.validator.Struct(obj)
	case reflect.Struct:
		return v.validator.Struct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

func (v *GinValidatorImpl) Engine() any {
	return v.validator
}

func initialize() {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())

		// Register custom tag name function to use JSON tags
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		//validate.RegisterCustomTypeFunc(
		//	ValidateValuer,
		//	null.Int{},
		//	null.Float{},
		//	null.String{},
		//	null.Bool{},
		//	null.Time{},
		//)

		// Register all custom validators
		err := customValidators.RegisterAll(validate)
		if err != nil {
			log.Panicf("Failed to register custom validators: %v", err)
		}

		// Initialize i18n
		i18n.Initialize()
	})
}

// GetValidator returns a singleton validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		initialize()
	}
	return validate
}

// Utility functions for validation

// ValidateStruct validates a struct and always returns a CustomError
func ValidateStruct(ctx context.Context, s any) customError.CustomError {
	err := RawValidateStruct(ctx, s)
	if err != nil {
		return FormatValidationErrors(ctx, err)
	}
	return customError.CustomError{}
}

// ValidateVar validates a single variable against a tag and returns a CustomError
func ValidateVar(ctx context.Context, field any, tag string) customError.CustomError {
	err := GetValidator().VarCtx(ctx, field, tag)
	if err != nil {
		return FormatValidationErrors(ctx, err)
	}
	return customError.CustomError{}
}

// RawValidateStruct validates a struct and returns the original validation error
// This is used internally for middleware and other utilities
func RawValidateStruct(ctx context.Context, s any) error {
	return GetValidator().StructCtx(ctx, s)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

// GetValidationErrors extracts validation errors from an error
func GetValidationErrors(err error) validator.ValidationErrors {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return validationErrors
	}
	return nil
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}

	return nil
}
