# Package currencyconverter

## Overview
The currencyconverter package offers utilities for converting monetary values between different currencies. It leverages exchange rates (possibly from the exchange_rate package) to perform accurate conversions.

## Key Components
- Conversion Functions: Methods to convert amounts between currencies.
- Data Structures: Represent currency information and conversion rates.
- Error Handling: Manages conversion errors gracefully.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/currencyconverter"
)

func main() {
	converted, err := currencyconverter.Convert(100, "USD", "EUR")
	if err != nil {
		fmt.Println("Conversion error:", err)
	} else {
		fmt.Println("Converted amount:", converted)
	}
}
~~~

## Notes
- Ideal for financial applications.
