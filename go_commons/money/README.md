# Money Package

The **money** package wraps the excellent
[github.com/Rhymond/go-money](https://github.com/Rhymond/go-money) library and
exposes a convenient, opinionated interface (`Money`) that is consistent with
other Omniful commons utilities.

It focuses on:

* Precise, integer–based representation of amounts to avoid floating-point
  errors.
* Currency awareness – every operation validates that the currency of the
  operands matches.
* A rich helper API that covers the majority of day-to-day financial use-cases
  (allocation, splitting, comparisons, formatting, rounding etc).

---

## Installation

```bash
go get github.com/omniful/go_commons/money
```

> The package depends on Go 1.20+ (tested on 1.22) and has no other
> requirements than the transitive dependency on *go-money*.

---

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/omniful/go_commons/money"
)

func main() {
    salary := money.NewFromFloat(1234.56, "USD") // $1,234.56

    // Give a 10% raise
    raised, _ := salary.Multiply(1.10)

    // Pay monthly
    monthly, _ := raised.Divide(12)

    fmt.Println("Yearly:", raised.Display())  // "USD 1,358.02"
    fmt.Println("Monthly:", monthly.Display()) // "USD 113.17"
}
```

---

## API Overview

### Construction

| Function | Description |
|----------|-------------|
| `New(amount int64, currency string)` | Creates a value where *amount* is stored as minor units (cents, yen, etc.). |
| `NewFromFloat(f float64, currency string)` | Convenience constructor that converts the float to the correct minor unit using **round-half-up**. |
| `ParseString("12.34", "USD")` | Like `NewFromFloat` but accepts a string to avoid binary-float surprises. |
| `Zero("USD")` | Returns a zero instance for the given currency. |

### Arithmetic & Comparison (partial)

```go
m1, _ := money.NewFromFloat(10.00, "USD")
m2 := money.New(250, "USD") // $2.50

sum, _  := m1.Add(m2)           // $12.50
diff, _ := m1.Subtract(m2)      // $7.50
prod, _ := m1.Multiply(2.5)     // $25.00
quot, _ := m1.Divide(3)         // $3.33 rounded

sum.GreaterThan(diff)           // true
sum.Equals(m1)                  // false
```

### Allocation & Splitting

```go
bucket := money.New(100, "USD") // $1.00 total
parts, _ := bucket.Allocate([]int{1, 1, 1}) // three shares => [0.34, 0.33, 0.33]

weekly, _ := bucket.Split(7) // seven equal parts
```

### Formatting Helpers

```go
m := money.New(1234567, "USD")

m.Display()                          // "USD 12,345.67"
m.FormatWithCode(" ")               // "12,345.67 USD"
m.FormatWithDefaultSeparators()      // "12,345.67"
m.FormatWithCustomSeparators(".", " ") // "12 345.67"
```

### Utility methods

* `ToFloat64()` – lossless conversion to a float64 value.
* Sign helpers: `IsZero()`, `IsPositive()`, `IsNegative()`.
* `Absolute()` / `Negate()`.
* `RoundToNearestUnit()` – rounds to the nearest _whole_ unit of the currency
  (useful for cash-only contexts).

---

## Error Handling

All operations that could fail (different currency, invalid arguments, overflow
etc.) return an `error` as the *second* return value following Go conventions:

```go
sum, err := m1.Add(m2)
if err != nil {
    // handle
}
```

---

## Contributing & Tests

Extensive unit tests ensure deterministic behaviour – run them with:

```bash
go test ./money -v -race
```

Feel free to open an issue or pull request if you spot a problem or have a
feature in mind.
