package util

import (
	"fmt"
	"math"

	"golang.org/x/exp/constraints"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Number is a constraint interface for number types
type Number interface {
	constraints.Integer | constraints.Float
}

// FormatNumber formats a number with commas according to the specified language tag.
// This function uses golang.org/x/text/message.Printer for locale-specific number formatting,
// which provides standard internationalization support.
//
// Parameters:
//   - value: The number to format (can be any integer or float type)
//   - lang: The language tag for locale-specific formatting
//     (e.g., language.English for "1,234.56" or language.MustParse("en-IN") for Indian format)
//   - decimalPlaces: The number of decimal places to display for floating point numbers
//
// Examples:
//
//	FormatNumber(1234567, language.English, 0) -> "1,234,567"
//	FormatNumber(1234.56, language.English, 2) -> "1,234.56"
//	FormatNumber(1234567.89, language.MustParse("en-IN"), 2) -> "12,34,567.89"
//
// Note: The actual formatting depends on the language tag and may vary
// between different locales for the same numeric value.
func FormatNumber[T Number](value T, decimalPlaces int, lang language.Tag) string {
	// Create a printer for the specified language
	p := message.NewPrinter(lang)

	// Convert to float64 for checking if it has a decimal part
	var floatValue float64

	// Extract the numeric value properly
	switch v := any(value).(type) {
	case int:
		floatValue = float64(v)
	case int8:
		floatValue = float64(v)
	case int16:
		floatValue = float64(v)
	case int32:
		floatValue = float64(v)
	case int64:
		floatValue = float64(v)
	case uint:
		floatValue = float64(v)
	case uint8:
		floatValue = float64(v)
	case uint16:
		floatValue = float64(v)
	case uint32:
		floatValue = float64(v)
	case uint64:
		floatValue = float64(v)
	case float32:
		floatValue = float64(v)
	case float64:
		floatValue = v
	default:
		// This should never happen due to the Number constraint
		return fmt.Sprintf("%v", value)
	}

	// Determine if the value has a non-zero decimal part
	_, frac := math.Modf(floatValue)
	isWholeNumber := math.Abs(frac) < 1e-10 // Small epsilon to account for floating-point imprecision

	// For whole numbers (integers or floats with zero decimal part)
	if isWholeNumber {
		return p.Sprintf("%d", int64(floatValue))
	}

	// For numbers with decimal parts
	formatStr := "%." + fmt.Sprintf("%d", decimalPlaces) + "f"
	return p.Sprintf(formatStr, floatValue)
}
