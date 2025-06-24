package money

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	gomoney "github.com/Rhymond/go-money"
)

// ParseMoneyString parses a string like "10.99 USD" into a Money object
func ParseMoneyString(s string) (Money, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return nil, errors.New("invalid format: expected 'amount currency'")
	}

	amount := parts[0]
	currency := parts[1]

	// Validate currency
	if !IsValidCurrencyCode(currency) {
		return nil, fmt.Errorf("invalid currency code: %s", currency)
	}
	return ParseString(amount, currency)
}

// IsValidCurrencyCode checks if a currency code is valid
func IsValidCurrencyCode(code string) bool {
	return gomoney.GetCurrency(code) != nil
}

// IsValidAmount checks if an amount is valid for a given currency
func IsValidAmount(amount float64, currency string) bool {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return false
	}

	// Additional currency-specific validations could be added here
	return true
}

// sameCurrency checks if all money values have the same currency
func sameCurrency(monies []Money) (string, error) {
	if len(monies) == 0 {
		return "", errors.New("cannot determine currency of empty slice")
	}
	currency := monies[0].Currency()
	for _, m := range monies[1:] {
		if m.Currency() != currency {
			return "", ErrCurrencyMismatch
		}
	}
	return currency, nil
}

// Sum calculates the sum of a slice of Money values
func Sum(monies []Money) (Money, error) {
	if len(monies) == 0 {
		return nil, errors.New("cannot sum empty slice")
	}

	if _, err := sameCurrency(monies); err != nil {
		return nil, err
	}

	result := monies[0].Clone()
	for i := 1; i < len(monies); i++ {
		var err error
		result, err = result.Add(monies[i])
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Min returns the minimum value from a slice of Money values
func Min(monies []Money) (Money, error) {
	if len(monies) == 0 {
		return nil, errors.New("cannot find minimum in empty slice")
	}

	if _, err := sameCurrency(monies); err != nil {
		return nil, err
	}

	result := monies[0]
	for i := 1; i < len(monies); i++ {
		less, err := monies[i].LessThan(result)
		if err != nil {
			return nil, fmt.Errorf("error comparing money values: %w", err)
		}
		if less {
			result = monies[i]
		}
	}

	return result, nil
}

// Max returns the maximum value from a slice of Money values
func Max(monies []Money) (Money, error) {
	if len(monies) == 0 {
		return nil, errors.New("cannot find maximum in empty slice")
	}

	if _, err := sameCurrency(monies); err != nil {
		return nil, err
	}

	result := monies[0]
	for i := 1; i < len(monies); i++ {
		greater, _ := monies[i].GreaterThan(result)
		if greater {
			result = monies[i]
		}
	}

	return result, nil
}

// Sort sorts a slice of Money values in ascending order
func Sort(monies []Money) error {
	if _, err := sameCurrency(monies); err != nil {
		return err
	}

	sort.Slice(monies, func(i, j int) bool {
		less, _ := monies[i].LessThan(monies[j])
		return less
	})
	return nil
}

// GetCurrencySymbol returns the symbol for a currency code
func GetCurrencySymbol(currencyCode string) (string, error) {
	currency := gomoney.GetCurrency(currencyCode)
	if currency == nil {
		return "", fmt.Errorf("unknown currency: %s", currencyCode)
	}
	return currency.Grapheme, nil
}

// ConvertCurrency converts money from one currency to another using an exchange rate
func ConvertCurrency(m Money, targetCurrencyCode string, exchangeRate float64) (Money, error) {
	if exchangeRate <= 0 {
		return nil, errors.New("exchange rate must be positive")
	}

	// Get decimal places for both currencies
	srcCurrencyInfo := gomoney.GetCurrency(m.Currency())
	targetCurrencyInfo := gomoney.GetCurrency(targetCurrencyCode)

	// Default to 2 decimal places if currency info not found
	srcDecimals := 2
	tgtDecimals := 2

	if srcCurrencyInfo != nil {
		srcDecimals = srcCurrencyInfo.Fraction
	}
	if targetCurrencyInfo != nil {
		tgtDecimals = targetCurrencyInfo.Fraction
	}

	// Calculate combined factor: exchange rate and decimal place adjustment
	decimalAdjustment := math.Pow10(tgtDecimals - srcDecimals)
	combinedFactor := exchangeRate * decimalAdjustment

	// Convert in one step
	convertedAmount, err := m.Multiply(combinedFactor)
	if err != nil {
		return nil, err
	}

	return New(convertedAmount.amount(), targetCurrencyCode), nil
}
