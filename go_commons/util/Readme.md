# Package util

## Overview
The util package contains a collection of helper functions used across the codebase. These utilities, ranging from string manipulation to data conversion, reduce code duplication and improve maintainability.

## Key Components
- Utility Functions: Reusable helper methods.
- Test Suites: Comprehensive tests to ensure reliability.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/util"
)

func main() {
	output := util.SomeUtilityFunction("input value") // replace with actual function
	fmt.Println("Utility function output:", output)
}
~~~

## Function Documentation

### FormatNumber

`FormatNumber` formats numeric values with locale-specific formatting including appropriate digit grouping (thousands separators).

~~~go
func FormatNumber[T Number](value T, decimalPlaces int, lang language.Tag) string
~~~

#### Parameters

- `value`: The number to format (supports all integer and float types)
- `decimalPlaces`: The number of decimal places to display (for floating point numbers)
- `lang`: The language tag for locale-specific formatting (e.g., `language.English`, `language.MustParse("en-IN")`)

#### Examples

##### English Format (1,000,000.00)
~~~go
// Format an integer with English (Western) grouping
fmt.Println(util.FormatNumber(1234567, 0, language.English))
// Output: 1,234,567

// Format a float with 2 decimal places
fmt.Println(util.FormatNumber(1234567.89, 2, language.English))
// Output: 1,234,567.89
~~~

##### Indian Format (10,00,000.00)
~~~go
// Format an integer with Indian grouping
fmt.Println(util.FormatNumber(1234567, 0, language.MustParse("en-IN")))
// Output: 12,34,567

// Format a float with 2 decimal places
fmt.Println(util.FormatNumber(1234567.89, 2, language.MustParse("en-IN")))
// Output: 12,34,567.89
~~~

#### Notes
- Trailing zeros after the decimal point are automatically removed
- For integer values, decimal places are ignored
- The function supports negative numbers and various numeric types (int, int64, float64, etc.)

## Notes
- Acts as a foundational package for common operations.
