# Weight Package

The `weight` package provides utilities for handling weights and unit conversions in Go applications. It supports various units of measurement (UOM) including weight, volume, and countable units, making it ideal for e-commerce and inventory management systems.

## Features

- Unit conversion between different weight and volume measurements
- Support for minimum sellable weight calculations
- Validation of unit compatibility
- Distinction between weighted and non-weighted units
- Error handling for invalid conversions

## Supported Units

### Weight Units
- `kg` - Kilograms
- `g` - Grams
- `lbs` - Pounds
- `oz` - Ounces

### Volume Units
- `l` - Liters
- `ml` - Milliliters

### Countable Units
- `ea` - Each (piece)
- `pack` - Pack

## Usage Examples

### Converting Between Units

```go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/weight"
)

func main() {
	// Convert 1 kg to grams
	w, err := weight.ConvertUOM(weight.Weight{UOM: weight.Kg, Value: 1}, weight.G)
	if err != nil {
		panic(err)
	}
	fmt.Printf("1 kg = %.2f g\n", w.Value) // Output: 1 kg = 1000.00 g

	// Convert 16 oz to pounds
	w, err = weight.ConvertUOM(weight.Weight{UOM: weight.Oz, Value: 16}, weight.Lbs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("16 oz = %.2f lbs\n", w.Value) // Output: 16 oz = 1.00 lbs
}
```

### Working with Minimum Sellable Weights

```go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/weight"
)

func main() {
	// Convert 1.5 kg with quantity 2 to minimum sellable weight (in grams)
	msw, err := weight.ConvertToMinimumSellableWeight(weight.Kg, 1.5, 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Minimum sellable weight: %d %s\n", msw.Value, msw.UOM) // Output: Minimum sellable weight: 3000 g

	// Working with countable units
	msw, err = weight.ConvertToMinimumSellableWeight(weight.EA, 1, 5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Minimum sellable quantity: %d %s\n", msw.Value, msw.UOM) // Output: Minimum sellable quantity: 5 ea
}
```

### Validation and Utility Functions

```go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/weight"
)

func main() {
	// Check if a unit is weighted
	fmt.Printf("Is kg weighted? %v\n", weight.IsUOMWeighted(weight.Kg))     // Output: true
	fmt.Printf("Is ea weighted? %v\n", weight.IsUOMWeighted(weight.EA))     // Output: false

	// Validate UOM compatibility
	fmt.Printf("Is kg valid for g? %v\n", weight.IsValidUOM(weight.Kg, weight.G))   // Output: true
	fmt.Printf("Is kg valid for ea? %v\n", weight.IsValidUOM(weight.Kg, weight.EA)) // Output: false

	// Get minimum sellable UOM
	fmt.Printf("Minimum sellable UOM for kg: %s\n", weight.GetMinimumSellableUOM(weight.Kg)) // Output: g
	fmt.Printf("Minimum sellable UOM for ea: %s\n", weight.GetMinimumSellableUOM(weight.EA)) // Output: ea
}
```

## Error Handling

The package uses custom error types from `github.com/omniful/go_commons/error`. All functions return appropriate error messages for invalid operations:

- Invalid UOM
- Incompatible unit conversions
- Other validation errors

## Conversion Table

| From | To | Multiplier |
|------|-----|------------|
| kg | g | 1000 |
| kg | lbs | 2.20462 |
| kg | oz | 35.27396 |
| lbs | g | 453.59237 |
| lbs | kg | 0.453592 |
| lbs | oz | 16 |
| g | kg | 0.001 |
| g | lbs | 0.00220462 |
| g | oz | 0.03527396 |
| l | ml | 1000 |
| ml | l | 0.001 |
| oz | g | 28.3495 |
| oz | kg | 0.0283495 |
| oz | lbs | 0.0625 |
