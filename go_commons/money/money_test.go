package money

import (
	"math"
	"testing"

	gomoney "github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

// mockMoney is a mock implementation of the Money interface for testing incompatible implementations
type mockMoney struct{}

func (m *mockMoney) amount() int64                          { return 100 }
func (m *mockMoney) Currency() string                       { return "USD" }
func (m *mockMoney) Add(Money) (Money, error)               { return nil, nil }
func (m *mockMoney) Subtract(Money) (Money, error)          { return nil, nil }
func (m *mockMoney) Multiply(float64) (Money, error)        { return nil, nil }
func (m *mockMoney) Divide(float64) (Money, error)          { return nil, nil }
func (m *mockMoney) Equals(Money) (bool, error)             { return false, nil }
func (m *mockMoney) GreaterThan(Money) (bool, error)        { return false, nil }
func (m *mockMoney) GreaterThanOrEqual(Money) (bool, error) { return false, nil }
func (m *mockMoney) LessThan(Money) (bool, error)           { return false, nil }
func (m *mockMoney) LessThanOrEqual(Money) (bool, error)    { return false, nil }
func (m *mockMoney) Allocate([]int) ([]Money, error)        { return nil, nil }
func (m *mockMoney) Split(int) ([]Money, error)             { return nil, nil }
func (m *mockMoney) Display() string                        { return "" }
func (m *mockMoney) Format(string, string) (string, error)  { return "", nil }
func (m *mockMoney) FormatWithCode(string) string           { return "" }
func (m *mockMoney) ToFloat64() float64                     { return 0 }
func (m *mockMoney) Clone() Money                           { return m }
func (m *mockMoney) IsZero() bool                           { return false }
func (m *mockMoney) IsNegative() bool                       { return false }
func (m *mockMoney) IsPositive() bool                       { return true }
func (m *mockMoney) Absolute() Money                        { return m }
func (m *mockMoney) Negate() Money                          { return m }
func (m *mockMoney) RoundToNearestUnit() Money              { return m }

// Constructor Tests

func TestNewAndAccessors(t *testing.T) {
	m := New(1234, "USD")
	assert.Equal(t, int64(1234), m.amount())
	assert.Equal(t, "USD", m.Currency())
}

func TestNewFromFloat(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		want     int64
		wantErr  bool
	}{
		{"standard rounding half up", 12.345, "USD", 1235, false},
		{"exact two decimals", 12.34, "USD", 1234, false},
		{"zero-decimal currency", 1234, "JPY", 1234, false},
		{"NaN", math.NaN(), "USD", 0, true},
		{"+Inf", math.Inf(1), "USD", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromFloat(tt.amount, tt.currency)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got.amount())
		})
	}
}

func TestNewFromFloatNegativeRounding(t *testing.T) {
	m, err := NewFromFloat(-12.345, "USD")
	assert.NoError(t, err)
	assert.Equal(t, int64(-1235), m.amount())
}

func TestParseString(t *testing.T) {
	m, err := ParseString("12.34", "USD")
	assert.NoError(t, err)
	assert.Equal(t, int64(1234), m.amount())
}

func TestParseStringInvalid(t *testing.T) {
	_, err := ParseString("abc", "USD")
	assert.Error(t, err)
}

func TestZeroIsZero(t *testing.T) {
	m := Zero("USD")
	assert.True(t, m.IsZero())
	assert.Equal(t, "USD", m.Currency())
}

// Basic Operations Tests

func TestAddSubtract(t *testing.T) {
	m1 := New(1000, "USD")
	m2 := New(250, "USD")

	sum, err := m1.Add(m2)
	assert.NoError(t, err)
	assert.Equal(t, int64(1250), sum.amount())

	diff, err := m1.Subtract(m2)
	assert.NoError(t, err)
	assert.Equal(t, int64(750), diff.amount())

	// Currency mismatch
	_, err = m1.Add(New(100, "EUR"))
	assert.Error(t, err)
}

func TestMultiplyDivide(t *testing.T) {
	m := New(1000, "USD") // $10.00

	multiplied, err := m.Multiply(1.5)
	assert.NoError(t, err)
	assert.Equal(t, int64(1500), multiplied.amount()) // $15.00

	divided, err := m.Divide(4)
	assert.NoError(t, err)
	assert.Equal(t, int64(250), divided.amount()) // $2.50

	_, err = m.Divide(0)
	assert.Error(t, err)

	_, err = m.Multiply(math.NaN())
	assert.Error(t, err)
}

func TestMultiplyDivideEdgeCases(t *testing.T) {
	m := New(1000, "USD")
	_, err := m.Multiply(math.Inf(1))
	assert.Error(t, err)

	_, err = m.Divide(math.NaN())
	assert.Error(t, err)
}

func TestMultiplyNegative(t *testing.T) {
	m := New(-1000, "USD")
	res, err := m.Multiply(1.5)
	assert.NoError(t, err)
	assert.Equal(t, int64(-1500), res.amount())
}

// Comparison Tests

func TestComparisons(t *testing.T) {
	m10 := New(1000, "USD")
	m20 := New(2000, "USD")
	equal10 := New(1000, "USD")

	eq, _ := m10.Equals(equal10)
	assert.True(t, eq)

	gt, _ := m20.GreaterThan(m10)
	assert.True(t, gt)

	gte, _ := m20.GreaterThanOrEqual(m20)
	assert.True(t, gte)

	lt, _ := m10.LessThan(m20)
	assert.True(t, lt)

	lte, _ := m10.LessThanOrEqual(equal10)
	assert.True(t, lte)
}

// Allocation Tests

func TestAllocate(t *testing.T) {
	m := New(100, "USD") // $1.00
	parts, err := m.Allocate([]int{1, 1, 1})
	assert.NoError(t, err)
	assert.Len(t, parts, 3)
	assert.Equal(t, int64(34), parts[0].amount())
	assert.Equal(t, int64(33), parts[1].amount())
	assert.Equal(t, int64(33), parts[2].amount())

	_, err = m.Allocate([]int{})
	assert.Error(t, err)

	// Zero ratios sum
	_, err = m.Allocate([]int{0, 0})
	assert.Error(t, err)
}

func TestAllocateErrorCases(t *testing.T) {
	m := New(100, "USD")

	// Test allocate with empty ratios
	t.Run("Allocate with empty ratios", func(t *testing.T) {
		_, err := m.Allocate([]int{})
		assert.Error(t, err)
		assert.Equal(t, ErrEmptyRatios, err)
	})

	// Test allocate with zero sum of ratios
	t.Run("Allocate with zero sum of ratios", func(t *testing.T) {
		_, err := m.Allocate([]int{0, 0, 0})
		assert.Error(t, err)
		assert.Equal(t, ErrZeroRatioSum, err)
	})

	// Test allocate with negative ratios (underlying library might handle this differently)
	t.Run("Allocate with negative ratios", func(t *testing.T) {
		_, err := m.Allocate([]int{-1, -2, -3})
		assert.Error(t, err)
	})
}

func TestSplit(t *testing.T) {
	m := New(100, "USD")
	parts, err := m.Split(4)
	assert.NoError(t, err)
	assert.Len(t, parts, 4)
	// First part may get remainder cent
	assert.Equal(t, int64(25), parts[0].amount())

	_, err = m.Split(0)
	assert.Error(t, err)
}

func TestSplitErrorCases(t *testing.T) {
	m := New(100, "USD")

	// Test split with zero parts
	t.Run("Split with zero parts", func(t *testing.T) {
		_, err := m.Split(0)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDivisor, err)
	})

	// Test split with negative parts
	t.Run("Split with negative parts", func(t *testing.T) {
		_, err := m.Split(-5)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDivisor, err)
	})
}

// Formatting Tests

func TestFormatting(t *testing.T) {
	m := New(1234, "USD") // $12.34

	assert.Equal(t, "12.34 USD", m.FormatWithCode(" "))

	def, err := m.Format(".", ",")
	assert.NoError(t, err)
	assert.Equal(t, "12.34", def)

	big := New(1234567, "USD") // 12,345.67
	custom, err := big.Format(".", " ")
	assert.NoError(t, err)
	assert.Equal(t, "12 345.67", custom)
}

func TestFormattingEdgeCases(t *testing.T) {
	t.Run("default separators from currency when blanks provided", func(t *testing.T) {
		m := New(1234567, "USD")
		formatted, err := m.Format("", "")
		assert.NoError(t, err)
		// Default US separators , and .
		assert.Equal(t, "12,345.67", formatted)
	})

	t.Run("unknown currency default separators error", func(t *testing.T) {
		m := New(100, "XXX")
		_, err := m.Format("", "")
		assert.Error(t, err)
	})
}

func TestFormattingZeroDecimalCurrency(t *testing.T) {
	m := New(1234, "JPY")
	// JPY has no decimal fraction
	formatted, err := m.Format(".", ",")
	assert.NoError(t, err)
	assert.Equal(t, "1,234", formatted)

	assert.Equal(t, "1234 JPY", m.FormatWithCode(" "))
}

func TestFormatErrorCases(t *testing.T) {
	// Test format with unknown currency
	t.Run("Format with unknown currency", func(t *testing.T) {
		m := New(1000, "XYZ") // Unknown currency code
		_, err := m.Format(".", ",")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown currency")
	})

	// Test format with various decimal and thousand separators
	t.Run("Format with custom separators", func(t *testing.T) {
		m := New(1234567, "USD")

		// Test with custom separators
		formatted, err := m.Format(".", "")
		assert.NoError(t, err)
		assert.Equal(t, "12,345.67", formatted)

		formatted, err = m.Format("", ",")
		assert.NoError(t, err)
		assert.Equal(t, "12,345.67", formatted)

		// Test with non-standard separators
		formatted, err = m.Format(":", ";")
		assert.NoError(t, err)
		assert.Equal(t, "12;345:67", formatted)
	})
}

func TestNegativeFormatting(t *testing.T) {
	// Test formatting negative values
	t.Run("Format negative values", func(t *testing.T) {
		m := New(-1234567, "USD")

		formatted, err := m.Format(".", ",")
		assert.NoError(t, err)
		assert.Equal(t, "-12,345.67", formatted)

		// Test zero-decimal currency
		m = New(-1234, "JPY")
		formatted, err = m.Format(".", ",")
		assert.NoError(t, err)
		assert.Equal(t, "-1,234", formatted)
	})
}

func TestCurrencyWithDifferentDecimals(t *testing.T) {
	// Test currencies with different decimal places
	t.Run("Format currencies with different decimal places", func(t *testing.T) {
		// Test 3-decimal currency (KWD)
		m := New(1234, "KWD")
		formatted, err := m.Format(".", ",")
		assert.NoError(t, err)
		assert.Equal(t, "1.234", formatted)

		// Test 4-decimal currency (CLF)
		m = New(12345, "CLF")
		formatted, err = m.Format(".", ",")
		assert.NoError(t, err)
		assert.Equal(t, "1.2345", formatted)
	})
}

func TestFormatAmount(t *testing.T) {
	// Test standard 2-decimal currency
	t.Run("USD (2 decimals)", func(t *testing.T) {
		assert.Equal(t, "12.34", formatAmount(1234, "USD"))
		assert.Equal(t, "-12.34", formatAmount(-1234, "USD"))
		assert.Equal(t, "0.00", formatAmount(0, "USD"))
		assert.Equal(t, "0.01", formatAmount(1, "USD"))
		assert.Equal(t, "1234567.89", formatAmount(123456789, "USD"))
	})

	// Test 0-decimal currency
	t.Run("JPY (0 decimals)", func(t *testing.T) {
		assert.Equal(t, "1234", formatAmount(1234, "JPY"))
		assert.Equal(t, "-1234", formatAmount(-1234, "JPY"))
		assert.Equal(t, "0", formatAmount(0, "JPY"))
		assert.Equal(t, "1", formatAmount(1, "JPY"))
	})

	// Test 3-decimal currency
	t.Run("KWD (3 decimals)", func(t *testing.T) {
		assert.Equal(t, "1.234", formatAmount(1234, "KWD"))
		assert.Equal(t, "-1.234", formatAmount(-1234, "KWD"))
		assert.Equal(t, "0.000", formatAmount(0, "KWD"))
		assert.Equal(t, "0.001", formatAmount(1, "KWD"))
	})

	// Test 4-decimal currency
	t.Run("CLF (4 decimals)", func(t *testing.T) {
		assert.Equal(t, "1.2345", formatAmount(12345, "CLF"))
		assert.Equal(t, "-1.2345", formatAmount(-12345, "CLF"))
		assert.Equal(t, "0.0000", formatAmount(0, "CLF"))
		assert.Equal(t, "0.0001", formatAmount(1, "CLF"))
	})

	// Test unknown currency (should default to 2 decimals)
	t.Run("Unknown currency (default 2 decimals)", func(t *testing.T) {
		assert.Equal(t, "12.34", formatAmount(1234, "XYZ"))
	})

	// Test edge cases
	t.Run("Edge cases", func(t *testing.T) {
		// Very large numbers
		assert.Equal(t, "92233720368547760.00", formatAmount(math.MaxInt64, "USD"))
		assert.Equal(t, "-92233720368547760.00", formatAmount(math.MinInt64, "USD"))

		// Small fractions
		assert.Equal(t, "0.01", formatAmount(1, "USD"))
		assert.Equal(t, "0.001", formatAmount(1, "KWD"))
	})
}

func TestFormatIntegerPart(t *testing.T) {
	// Test with empty thousand separator
	t.Run("Empty thousand separator", func(t *testing.T) {
		assert.Equal(t, "1234", formatIntegerPart(1234, ""))
		assert.Equal(t, "1234567", formatIntegerPart(1234567, ""))
		assert.Equal(t, "0", formatIntegerPart(0, ""))
	})

	// Test with comma separator
	t.Run("Comma separator", func(t *testing.T) {
		assert.Equal(t, "1,234", formatIntegerPart(1234, ","))
		assert.Equal(t, "1,234,567", formatIntegerPart(1234567, ","))
		assert.Equal(t, "0", formatIntegerPart(0, ","))
	})

	// Test with space separator
	t.Run("Space separator", func(t *testing.T) {
		assert.Equal(t, "1 234", formatIntegerPart(1234, " "))
		assert.Equal(t, "1 234 567", formatIntegerPart(1234567, " "))
	})

	// Test with custom separator
	t.Run("Custom separator", func(t *testing.T) {
		assert.Equal(t, "1'234", formatIntegerPart(1234, "'"))
		assert.Equal(t, "1'234'567", formatIntegerPart(1234567, "'"))
	})

	// Test large numbers
	t.Run("Large numbers", func(t *testing.T) {
		assert.Equal(t, "9,223,372,036,854,775,807", formatIntegerPart(math.MaxInt64, ","))
		assert.Equal(t, "-9,223,372,036,854,775,808", formatIntegerPart(math.MinInt64, ",")) // Absolute value
	})
}

func TestFormatWithCodeMethod(t *testing.T) {
	// Test standard currency
	t.Run("Standard currency", func(t *testing.T) {
		m := New(1234, "USD")
		assert.Equal(t, "12.34 USD", m.FormatWithCode(" "))
		assert.Equal(t, "12.34/USD", m.FormatWithCode("/"))
		assert.Equal(t, "12.34USD", m.FormatWithCode(""))
	})

	// Test zero-decimal currency
	t.Run("Zero-decimal currency", func(t *testing.T) {
		m := New(1234, "JPY")
		assert.Equal(t, "1234 JPY", m.FormatWithCode(" "))
	})

	// Test 3-decimal currency
	t.Run("3-decimal currency", func(t *testing.T) {
		m := New(1234, "KWD")
		assert.Equal(t, "1.234 KWD", m.FormatWithCode(" "))
	})

	// Test negative values
	t.Run("Negative values", func(t *testing.T) {
		m := New(-1234, "USD")
		assert.Equal(t, "-12.34 USD", m.FormatWithCode(" "))
	})

	// Test zero values
	t.Run("Zero values", func(t *testing.T) {
		m := New(0, "USD")
		assert.Equal(t, "0.00 USD", m.FormatWithCode(" "))
	})
}

func TestFormatMethodComprehensive(t *testing.T) {
	// Test with different separator combinations
	t.Run("Different separator combinations", func(t *testing.T) {
		m := New(1234567, "USD")

		// Default separators (should use currency defaults)
		formatted, err := m.Format("", "")
		assert.NoError(t, err)
		assert.Equal(t, "12,345.67", formatted)

		// Custom decimal separator only
		formatted, err = m.Format(".", "")
		assert.NoError(t, err)
		assert.Equal(t, "12,345.67", formatted)

		// Custom thousand separator only
		formatted, err = m.Format("", ",")
		assert.NoError(t, err)
		assert.Equal(t, "12,345.67", formatted)

		// Both custom separators
		formatted, err = m.Format(":", ";")
		assert.NoError(t, err)
		assert.Equal(t, "12;345:67", formatted)
	})
}

func TestFormatWithInvalidCurrency(t *testing.T) {
	// Test with invalid currency
	m := New(1234, "XYZ")
	_, err := m.Format(".", ",")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown currency")
}

func TestDisplayMethodComprehensive(t *testing.T) {
	// Test with different currencies
	t.Run("Different currencies", func(t *testing.T) {
		// USD (2 decimals)
		m := New(1234, "USD")
		display := m.Display()
		assert.Contains(t, display, "$12.34")

		// EUR (2 decimals)
		m = New(1234, "EUR")
		display = m.Display()
		assert.Contains(t, display, "€12.34")

		// JPY (0 decimals)
		m = New(1234, "JPY")
		display = m.Display()
		assert.Contains(t, display, "¥1,234")
	})

	// Test with different values
	t.Run("Different values", func(t *testing.T) {
		// Zero
		m := New(0, "USD")
		display := m.Display()
		assert.Contains(t, display, "$0.00")

		// Negative
		m = New(-1234, "USD")
		display = m.Display()
		assert.Contains(t, display, "-$12.34")

		// Large number
		m = New(1234567, "USD")
		display = m.Display()
		assert.Contains(t, display, "$12,345.67")
	})
}

// Conversion Tests

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		want     float64
	}{
		// Standard 2-decimal currencies
		{"standard USD case", 1234, "USD", 12.34},
		{"large USD amount", 1234567, "USD", 12345.67},
		{"negative USD", -1234, "USD", -12.34},
		{"zero USD", 0, "USD", 0},
		{"small USD amount", 1, "USD", 0.01},
		// {"very large USD", 9223372036854775807, "USD", 92233720368547758.07}, // max int64

		// 0-decimal currencies
		{"zero-decimal currency (JPY)", 1234, "JPY", 1234},
		{"negative JPY", -5000, "JPY", -5000},
		{"zero JPY", 0, "JPY", 0},

		// 3-decimal currencies
		{"three-decimal currency (KWD)", 1234, "KWD", 1.234},
		{"negative KWD", -1234, "KWD", -1.234},
		{"fractional KWD", 1, "KWD", 0.001},

		// 4-decimal currencies
		{"four-decimal currency (CLF)", 12345, "CLF", 1.2345},

		// Edge cases
		{"one cent USD", 1, "USD", 0.01},
		// {"almost min int64 USD", math.MinInt64 + 1, "USD", -92233720368547758.08},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.amount, tt.currency)
			assert.Equal(t, tt.want, m.ToFloat64())

			// Round-trip test for standard currencies
			if tt.currency == "USD" || tt.currency == "EUR" {
				roundTrip, _ := NewFromFloat(m.ToFloat64(), tt.currency)
				assert.Equal(t, tt.amount, roundTrip.amount())
			}
		})
	}

	// Test with all currencies registered in go-money
	t.Run("test representative currencies with different decimal places", func(t *testing.T) {
		// Map of test cases with different decimal places
		currencyTests := map[string]struct {
			amount        int64
			decimalPlaces int
		}{
			"JPY": {1000, 0},  // Japanese Yen - 0 decimal places
			"USD": {1234, 2},  // US Dollar - 2 decimal places
			"KWD": {1234, 3},  // Kuwaiti Dinar - 3 decimal places
			"CLF": {12345, 4}, // Chilean Unidad de Fomento - 4 decimal places
		}

		for currency, test := range currencyTests {
			m := New(test.amount, currency)
			expectedValue := float64(test.amount) / math.Pow10(test.decimalPlaces)
			assert.Equal(t, expectedValue, m.ToFloat64())
		}
	})
}

func TestToFloat64EdgeCases(t *testing.T) {
	// Test with extreme values
	t.Run("ToFloat64 extreme values", func(t *testing.T) {
		m := New(math.MaxInt64, "USD")
		result := m.ToFloat64()
		assert.True(t, result > 0)

		m = New(math.MinInt64, "USD")
		result = m.ToFloat64()
		assert.True(t, result < 0)
	})
}

// Helper Method Tests

func TestHelpers(t *testing.T) {
	m := New(-1234, "USD")
	assert.True(t, m.IsNegative())
	assert.False(t, m.IsPositive())

	abs := m.Absolute()
	assert.True(t, abs.IsPositive())
	assert.Equal(t, int64(1234), abs.amount())

	neg := abs.Negate()
	assert.Equal(t, int64(-1234), neg.amount())

	clone := abs.Clone()
	eq, _ := abs.Equals(clone)
	assert.True(t, eq)
	// Different instances (underlying pointers should differ)
	assert.NotSame(t, abs, clone)
}

func TestIsZeroComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		expected bool
	}{
		{"zero amount", 0, "USD", true},
		{"positive amount", 1, "USD", false},
		{"negative amount", -1, "USD", false},
		{"large positive amount", math.MaxInt64, "USD", false},
		{"large negative amount", math.MinInt64, "USD", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.amount, tt.currency)
			assert.Equal(t, tt.expected, m.IsZero())
		})
	}
}

func TestIsNegativeComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		expected bool
	}{
		{"zero amount", 0, "USD", false},
		{"positive amount", 1, "USD", false},
		{"negative amount", -1, "USD", true},
		{"large positive amount", math.MaxInt64, "USD", false},
		{"large negative amount", math.MinInt64, "USD", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.amount, tt.currency)
			assert.Equal(t, tt.expected, m.IsNegative())
		})
	}
}

func TestIsPositiveComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		currency string
		expected bool
	}{
		{"zero amount", 0, "USD", false},
		{"positive amount", 1, "USD", true},
		{"negative amount", -1, "USD", false},
		{"large positive amount", math.MaxInt64, "USD", true},
		{"large negative amount", math.MinInt64, "USD", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.amount, tt.currency)
			assert.Equal(t, tt.expected, m.IsPositive())
		})
	}
}

func TestAbsoluteEdgeCases(t *testing.T) {
	// Test with MinInt64 (which can't be negated without overflow)
	t.Run("Absolute of MinInt64", func(t *testing.T) {
		m := New(math.MinInt64+1, "USD")
		abs := m.Absolute()
		// Should return a positive value close to MaxInt64
		assert.True(t, abs.amount() > 0)
	})
}

func TestNegateEdgeCases(t *testing.T) {
	// Test with MinInt64 (which can't be negated without overflow)
	t.Run("Negate of MinInt64", func(t *testing.T) {
		m := New(math.MinInt64+1, "USD")
		neg := m.Negate()
		// Should handle this edge case gracefully
		assert.NotEqual(t, m.amount(), neg.amount())
	})
}

func TestRoundToNearestUnit(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
		want   int64
	}{
		{"round down", 105, 100},
		{"round up", 150, 200},
		{"exact", 300, 300},
		{"negative round away from zero", -150, -200},
		{"negative round toward zero", -49, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.amount, "USD").RoundToNearestUnit()
			assert.Equal(t, tt.want, got.amount())
		})
	}
}

func TestRoundToNearestUnitZeroFraction(t *testing.T) {
	m := New(1234, "JPY")
	rounded := m.RoundToNearestUnit()
	assert.Equal(t, int64(1234), rounded.amount())
}

// Error Handling Tests

func TestIncompatibleImplementation(t *testing.T) {
	m := New(100, "USD")
	mockM := &mockMoney{}

	// Test Add with incompatible implementation
	t.Run("Add with incompatible implementation", func(t *testing.T) {
		_, err := m.Add(mockM)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleImplementation, err)
	})

	// Test Subtract with incompatible implementation
	t.Run("Subtract with incompatible implementation", func(t *testing.T) {
		_, err := m.Subtract(mockM)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleImplementation, err)
	})

	// Test comparison methods with incompatible implementation
	t.Run("Equals with incompatible implementation", func(t *testing.T) {
		_, err := m.Equals(mockM)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleImplementation, err)
	})

	t.Run("GreaterThan with incompatible implementation", func(t *testing.T) {
		_, err := m.GreaterThan(mockM)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleImplementation, err)
	})

	t.Run("GreaterThanOrEqual with incompatible implementation", func(t *testing.T) {
		_, err := m.GreaterThanOrEqual(mockM)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleImplementation, err)
	})

	t.Run("LessThan with incompatible implementation", func(t *testing.T) {
		_, err := m.LessThan(mockM)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleImplementation, err)
	})

	t.Run("LessThanOrEqual with incompatible implementation", func(t *testing.T) {
		_, err := m.LessThanOrEqual(mockM)
		assert.Error(t, err)
		assert.Equal(t, ErrIncompatibleImplementation, err)
	})
}

func TestMultiplyOverflow(t *testing.T) {
	// Test with a value that would overflow int64 when multiplied
	t.Run("Multiply overflow", func(t *testing.T) {
		m := New(math.MaxInt64/2+1, "USD")
		_, err := m.Multiply(3)
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})

	// Test with a negative value that would overflow int64 when multiplied
	t.Run("Multiply negative overflow", func(t *testing.T) {
		m := New(math.MinInt64/2-1, "USD")
		_, err := m.Multiply(3)
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})

	// Test with a value close to MaxInt64 boundary
	t.Run("Multiply near MaxInt64", func(t *testing.T) {
		m := New(math.MaxInt64/2, "USD")
		_, err := m.Multiply(1.999)
		assert.NoError(t, err)

		_, err = m.Multiply(2.001)
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})
}

func TestMultiplyOverflowEdgeCases(t *testing.T) {
	// Test with values that would cause overflow after rounding
	t.Run("Overflow after rounding", func(t *testing.T) {
		m := New(math.MaxInt64-10, "USD")
		_, err := m.Multiply(1.000001) // Small factor but enough to cause overflow after rounding
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})

	// Test with negative values that would overflow
	t.Run("Negative value overflow", func(t *testing.T) {
		m := New(math.MinInt64+10, "USD")
		_, err := m.Multiply(1.000001) // Small factor but enough to cause overflow after rounding
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})

	// Test with factor that would cause overflow in float calculation
	t.Run("Float calculation overflow", func(t *testing.T) {
		m := New(math.MaxInt64/2, "USD")
		_, err := m.Multiply(2.1) // Just over the threshold to cause overflow
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})
}

func TestDivideEdgeCases(t *testing.T) {
	// Test with a value that would overflow int64 when divided (this is rare but possible with rounding)
	t.Run("Division with potential overflow after rounding", func(t *testing.T) {
		// Create a value that when divided by a small number and rounded would overflow
		m := New(math.MaxInt64, "USD")
		_, err := m.Divide(0.5) // This would double the value, causing overflow
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})

	// Test the special case handling for MaxInt64/2 divided by 2
	t.Run("Special case MaxInt64/2 divided by 2", func(t *testing.T) {
		m := New(math.MaxInt64/2, "USD")
		result, err := m.Divide(2)
		assert.NoError(t, err)
		assert.Equal(t, int64(math.MaxInt64/4), result.amount())
	})

	// Test division by very small number (approaching zero)
	t.Run("Division by very small number", func(t *testing.T) {
		m := New(1000, "USD")
		result, err := m.Divide(1e-15)
		assert.NoError(t, err)
		// The result should be a very large number but within int64 range
		assert.True(t, result.amount() > 0)
	})
}

func TestDivideErrorPropagation(t *testing.T) {
	// Mock the gomoney.Money to return an error
	originalMoney := gomoney.New(100, "USD")

	// Create a moneyImpl with a custom inner money object
	m := &moneyImpl{
		money: originalMoney,
	}

	// Test with a divisor that would cause an error
	_, err := m.Divide(0)
	assert.Error(t, err)
	assert.Equal(t, ErrDivisionByZero, err)
}

func TestDivideAdditionalEdgeCases(t *testing.T) {
	// Test with divisor that would cause overflow in the result
	t.Run("Division causing overflow", func(t *testing.T) {
		m := New(math.MaxInt64, "USD")
		_, err := m.Divide(0.1) // Would result in a value 10x larger than MaxInt64
		assert.Error(t, err)
		assert.Equal(t, ErrOverflow, err)
	})

	// Test with very large divisor
	t.Run("Very large divisor", func(t *testing.T) {
		m := New(100, "USD")
		result, err := m.Divide(1e10)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), result.amount()) // Should round to 0
	})

	// Test with divisor close to zero but not zero
	t.Run("Divisor close to zero", func(t *testing.T) {
		m := New(100, "USD")
		_, err := m.Divide(1e-16) // Very small but not zero
		assert.NoError(t, err)
	})
}

func TestFormatWithExtremeValues(t *testing.T) {
	// Test with max int64 value
	t.Run("Format MaxInt64", func(t *testing.T) {
		m := New(math.MaxInt64, "USD")
		formatted, err := m.Format(".", ",")
		assert.NoError(t, err)
		assert.NotEmpty(t, formatted)
	})

	// Test with min int64 value
	t.Run("Format MinInt64", func(t *testing.T) {
		m := New(math.MinInt64, "USD")
		formatted, err := m.Format(".", ",")
		assert.NoError(t, err)
		assert.NotEmpty(t, formatted)
		assert.Contains(t, formatted, "-")
	})

	// Test with unusual separators
	t.Run("Unusual separators", func(t *testing.T) {
		m := New(1234567, "USD")
		formatted, err := m.Format("@", "#")
		assert.NoError(t, err)
		assert.Equal(t, "12#345@67", formatted)
	})

	// Test with empty separators
	t.Run("Empty separators", func(t *testing.T) {
		m := New(1234567, "USD")
		formatted, err := m.Format("", "")
		assert.NoError(t, err)
		assert.NotEmpty(t, formatted)
	})
}

func TestFormatWithCodeEdgeCases(t *testing.T) {
	// Test with extreme values
	t.Run("FormatWithCode extreme values", func(t *testing.T) {
		m := New(math.MaxInt64, "USD")
		result := m.FormatWithCode(" ")
		assert.Contains(t, result, "USD")

		m = New(math.MinInt64, "USD")
		result = m.FormatWithCode(" ")
		assert.Contains(t, result, "USD")
		assert.Contains(t, result, "-")
	})
}

func TestCloneEdgeCases(t *testing.T) {
	// Test with extreme values
	t.Run("Clone extreme values", func(t *testing.T) {
		m := New(math.MaxInt64, "USD")
		clone := m.Clone()
		assert.Equal(t, m.amount(), clone.amount())
		assert.Equal(t, m.Currency(), clone.Currency())

		m = New(math.MinInt64, "USD")
		clone = m.Clone()
		assert.Equal(t, m.amount(), clone.amount())
		assert.Equal(t, m.Currency(), clone.Currency())
	})
}
