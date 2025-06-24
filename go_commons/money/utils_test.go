package money

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMoneyString(t *testing.T) {
	// Test valid money string
	t.Run("Valid money string", func(t *testing.T) {
		m, err := ParseMoneyString("12.34 USD")
		assert.NoError(t, err)
		assert.Equal(t, int64(1234), m.amount())
		assert.Equal(t, "USD", m.Currency())
	})

	// Test invalid format
	t.Run("Invalid format", func(t *testing.T) {
		_, err := ParseMoneyString("12.34USD")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid format")
	})

	// Test too many parts
	t.Run("Too many parts", func(t *testing.T) {
		_, err := ParseMoneyString("12.34 USD extra")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid format")
	})

	// Test invalid currency
	t.Run("Invalid currency", func(t *testing.T) {
		_, err := ParseMoneyString("12.34 123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid currency code")
	})

	// Test invalid amount
	t.Run("Invalid amount", func(t *testing.T) {
		_, err := ParseMoneyString("abc USD")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse amount")
	})
}

func TestIsValidCurrencyCode(t *testing.T) {
	// Test valid currency codes
	t.Run("Valid currency codes", func(t *testing.T) {
		assert.True(t, IsValidCurrencyCode("USD"))
		assert.True(t, IsValidCurrencyCode("EUR"))
		assert.True(t, IsValidCurrencyCode("JPY"))
		assert.True(t, IsValidCurrencyCode("GBP"))
	})

	// Test invalid currency codes
	t.Run("Invalid currency codes", func(t *testing.T) {
		assert.False(t, IsValidCurrencyCode("US"))   // too short
		assert.False(t, IsValidCurrencyCode("USDD")) // too long
		assert.False(t, IsValidCurrencyCode("123"))  // not letters
		assert.False(t, IsValidCurrencyCode(""))     // empty
	})
}

func TestIsValidAmount(t *testing.T) {
	// Test valid amounts
	t.Run("Valid amounts", func(t *testing.T) {
		assert.True(t, IsValidAmount(0, "USD"))
		assert.True(t, IsValidAmount(12.34, "USD"))
		assert.True(t, IsValidAmount(-12.34, "USD"))
		assert.True(t, IsValidAmount(1000000, "USD"))
	})

	// Test invalid amounts
	t.Run("Invalid amounts", func(t *testing.T) {
		assert.False(t, IsValidAmount(math.NaN(), "USD"))
		assert.False(t, IsValidAmount(math.Inf(1), "USD"))
		assert.False(t, IsValidAmount(math.Inf(-1), "USD"))
	})
}

func TestSameCurrency(t *testing.T) {
	// Test empty slice
	t.Run("Empty slice", func(t *testing.T) {
		_, err := sameCurrency([]Money{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty slice")
	})

	// Test same currency
	t.Run("Same currency", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(200, "USD"),
			New(300, "USD"),
		}
		currency, err := sameCurrency(monies)
		assert.NoError(t, err)
		assert.Equal(t, "USD", currency)
	})

	// Test different currencies
	t.Run("Different currencies", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(200, "EUR"),
			New(300, "USD"),
		}
		_, err := sameCurrency(monies)
		assert.Error(t, err)
		assert.Equal(t, ErrCurrencyMismatch, err)
	})

	// Test single item
	t.Run("Single item", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
		}
		currency, err := sameCurrency(monies)
		assert.NoError(t, err)
		assert.Equal(t, "USD", currency)
	})
}

func TestSum(t *testing.T) {
	// Test empty slice
	t.Run("Empty slice", func(t *testing.T) {
		_, err := Sum([]Money{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty slice")
	})

	// Test different currencies
	t.Run("Different currencies", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(200, "EUR"),
		}
		_, err := Sum(monies)
		assert.Error(t, err)
		assert.Equal(t, ErrCurrencyMismatch, err)
	})

	// Test valid sum
	t.Run("Valid sum", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(200, "USD"),
			New(300, "USD"),
		}
		sum, err := Sum(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(600), sum.amount())
		assert.Equal(t, "USD", sum.Currency())
	})

	// Test single item
	t.Run("Single item", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
		}
		sum, err := Sum(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), sum.amount())
	})

	// Test with negative values
	t.Run("With negative values", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(-50, "USD"),
			New(25, "USD"),
		}
		sum, err := Sum(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(75), sum.amount())
	})
}

func TestMin(t *testing.T) {
	// Test empty slice
	t.Run("Empty slice", func(t *testing.T) {
		_, err := Min([]Money{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty slice")
	})

	// Test different currencies
	t.Run("Different currencies", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(200, "EUR"),
		}
		_, err := Min(monies)
		assert.Error(t, err)
		assert.Equal(t, ErrCurrencyMismatch, err)
	})

	// Test valid minimum
	t.Run("Valid minimum", func(t *testing.T) {
		monies := []Money{
			New(300, "USD"),
			New(100, "USD"),
			New(200, "USD"),
		}
		min, err := Min(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), min.amount())
	})

	// Test with negative values
	t.Run("With negative values", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(-50, "USD"),
			New(25, "USD"),
		}
		min, err := Min(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(-50), min.amount())
	})

	// Test single item
	t.Run("Single item", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
		}
		min, err := Min(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), min.amount())
	})
}

func TestMax(t *testing.T) {
	// Test empty slice
	t.Run("Empty slice", func(t *testing.T) {
		_, err := Max([]Money{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty slice")
	})

	// Test different currencies
	t.Run("Different currencies", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(200, "EUR"),
		}
		_, err := Max(monies)
		assert.Error(t, err)
		assert.Equal(t, ErrCurrencyMismatch, err)
	})

	// Test valid maximum
	t.Run("Valid maximum", func(t *testing.T) {
		monies := []Money{
			New(300, "USD"),
			New(100, "USD"),
			New(200, "USD"),
		}
		max, err := Max(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(300), max.amount())
	})

	// Test with negative values
	t.Run("With negative values", func(t *testing.T) {
		monies := []Money{
			New(-100, "USD"),
			New(-50, "USD"),
			New(-200, "USD"),
		}
		max, err := Max(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(-50), max.amount())
	})

	// Test single item
	t.Run("Single item", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
		}
		max, err := Max(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), max.amount())
	})
}

func TestSort(t *testing.T) {
	// Test different currencies
	t.Run("Different currencies", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(200, "EUR"),
		}
		err := Sort(monies)
		assert.Error(t, err)
		assert.Equal(t, ErrCurrencyMismatch, err)
	})

	// Test valid sort
	t.Run("Valid sort", func(t *testing.T) {
		monies := []Money{
			New(300, "USD"),
			New(100, "USD"),
			New(200, "USD"),
		}
		err := Sort(monies)
		assert.NoError(t, err)

		// Verify sorted order
		assert.Equal(t, int64(100), monies[0].amount())
		assert.Equal(t, int64(200), monies[1].amount())
		assert.Equal(t, int64(300), monies[2].amount())
	})

	// Test with negative values
	t.Run("With negative values", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
			New(-50, "USD"),
			New(25, "USD"),
		}
		err := Sort(monies)
		assert.NoError(t, err)

		// Verify sorted order
		assert.Equal(t, int64(-50), monies[0].amount())
		assert.Equal(t, int64(25), monies[1].amount())
		assert.Equal(t, int64(100), monies[2].amount())
	})

	// Test empty slice
	t.Run("Empty slice", func(t *testing.T) {
		err := Sort([]Money{})
		assert.Error(t, err)
	})

	// Test single item (should be no-op)
	t.Run("Single item", func(t *testing.T) {
		monies := []Money{
			New(100, "USD"),
		}
		err := Sort(monies)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), monies[0].amount())
	})
}

func TestGetCurrencySymbol(t *testing.T) {
	// Test valid currency codes
	t.Run("Valid currency codes", func(t *testing.T) {
		symbol, err := GetCurrencySymbol("USD")
		assert.NoError(t, err)
		assert.Equal(t, "$", symbol)

		symbol, err = GetCurrencySymbol("EUR")
		assert.NoError(t, err)
		assert.Equal(t, "€", symbol)

		symbol, err = GetCurrencySymbol("GBP")
		assert.NoError(t, err)
		assert.Equal(t, "£", symbol)

		symbol, err = GetCurrencySymbol("JPY")
		assert.NoError(t, err)
		assert.Equal(t, "¥", symbol)
	})

	// Test invalid currency code
	t.Run("Invalid currency code", func(t *testing.T) {
		_, err := GetCurrencySymbol("XYZ")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown currency")
	})

	// Test empty currency code
	t.Run("Empty currency code", func(t *testing.T) {
		_, err := GetCurrencySymbol("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown currency")
	})
}

func TestConvertCurrency(t *testing.T) {
	// Test valid conversion
	t.Run("Valid conversion", func(t *testing.T) {
		m := New(1000, "USD") // $10.00
		converted, err := ConvertCurrency(m, "EUR", 0.85)
		assert.NoError(t, err)
		assert.Equal(t, "EUR", converted.Currency())
		assert.Equal(t, int64(850), converted.amount()) // €8.50
	})

	// Test with zero exchange rate
	t.Run("Zero exchange rate", func(t *testing.T) {
		m := New(1000, "USD")
		_, err := ConvertCurrency(m, "EUR", 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exchange rate must be positive")
	})

	// Test with negative exchange rate
	t.Run("Negative exchange rate", func(t *testing.T) {
		m := New(1000, "USD")
		_, err := ConvertCurrency(m, "EUR", -0.85)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exchange rate must be positive")
	})

	// Test with different decimal places
	t.Run("Different decimal places", func(t *testing.T) {
		// USD (2 decimals) to JPY (0 decimals)
		m := New(1000, "USD") // $10.00
		converted, err := ConvertCurrency(m, "JPY", 110.0)
		assert.NoError(t, err)
		assert.Equal(t, "JPY", converted.Currency())
		assert.Equal(t, int64(1100), converted.amount()) // ¥1100

		// JPY (0 decimals) to USD (2 decimals)
		m = New(1000, "JPY") // ¥1000
		converted, err = ConvertCurrency(m, "USD", 0.009)
		assert.NoError(t, err)
		assert.Equal(t, "USD", converted.Currency())
		assert.Equal(t, int64(900), converted.amount()) // $9.00
	})

	// Test with 3-decimal currency
	t.Run("Three-decimal currency", func(t *testing.T) {
		// USD (2 decimals) to KWD (3 decimals)
		m := New(1000, "USD") // $10.00
		converted, err := ConvertCurrency(m, "KWD", 0.3)
		assert.NoError(t, err)
		assert.Equal(t, "KWD", converted.Currency())
		assert.Equal(t, int64(3000), converted.amount()) // 0.300 KWD
	})

	// Test overflow case
	t.Run("Overflow case", func(t *testing.T) {
		m := New(math.MaxInt64/2, "USD")
		_, err := ConvertCurrency(m, "EUR", 3.0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "overflows")
	})

}

func TestDisplay(t *testing.T) {
	// Test Display method
	t.Run("Display method", func(t *testing.T) {
		m := New(1234, "USD")
		display := m.Display()
		assert.NotEmpty(t, display)
		assert.Contains(t, display, "$")
		assert.Contains(t, display, "12.34")
	})

	// Test Display with different currencies
	t.Run("Display with different currencies", func(t *testing.T) {
		// USD (2 decimals)
		m := New(1234, "USD")
		display := m.Display()
		assert.Contains(t, display, "12.34")

		// JPY (0 decimals)
		m = New(1234, "JPY")
		display = m.Display()
		assert.Contains(t, display, "1,234")

		// KWD (3 decimals)
		m = New(1234, "KWD")
		display = m.Display()
		assert.Contains(t, display, "1.234")
	})

	// Test Display with negative values
	t.Run("Display with negative values", func(t *testing.T) {
		m := New(-1234, "USD")
		display := m.Display()
		assert.Contains(t, display, "-")
	})
}
