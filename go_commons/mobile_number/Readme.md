# Package mobile_number

## Overview
The mobile_number package provides utilities for parsing, validating, and formatting mobile phone numbers. It ensures conformity to international standards and supports various regional formats.

## Key Components
- Validation Functions: Check the correctness of phone numbers.
- Formatting Utilities: Standardize phone number presentation.
- Parsing Helpers: Extract country codes and local numbers.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/mobile_number"
)

func main() {
	formatted, err := mobile_number.Format("+1-800-123-4567")
	if err != nil {
		fmt.Println("Invalid mobile number:", err)
	} else {
		fmt.Println("Formatted mobile number:", formatted)
	}
}
~~~

## Notes
- Supports international number formats robustly.
