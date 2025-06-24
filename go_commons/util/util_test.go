package util

import (
	"encoding/json"
	"testing"

	"golang.org/x/text/language"
)

type TestStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		x, y interface{}
		want interface{}
	}{
		{
			name: "int - x greater than y",
			x:    5,
			y:    3,
			want: 5,
		},
		{
			name: "int - y greater than x",
			x:    2,
			y:    4,
			want: 4,
		},
		{
			name: "int - x equal to y",
			x:    7,
			y:    7,
			want: 7,
		},
		{
			name: "float - x greater than y",
			x:    5.5,
			y:    3.3,
			want: 5.5,
		},
		{
			name: "float - y greater than x",
			x:    2.2,
			y:    4.4,
			want: 4.4,
		},
		{
			name: "float - x equal to y",
			x:    7.7,
			y:    7.7,
			want: 7.7,
		},
		{
			name: "string - x greater than y",
			x:    "banana",
			y:    "apple",
			want: "banana",
		},
		{
			name: "string - y greater than x",
			x:    "apple",
			y:    "banana",
			want: "banana",
		},
		{
			name: "string - x equal to y",
			x:    "cherry",
			y:    "cherry",
			want: "cherry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch x := tt.x.(type) {
			case int:
				y := tt.y.(int)
				want := tt.want.(int)
				if got := Max(x, y); got != want {
					t.Errorf("Max() = %v, want %v", got, want)
				}
			case float64:
				y := tt.y.(float64)
				want := tt.want.(float64)
				if got := Max(x, y); got != want {
					t.Errorf("Max() = %v, want %v", got, want)
				}
			case string:
				y := tt.y.(string)
				want := tt.want.(string)
				if got := Max(x, y); got != want {
					t.Errorf("Max() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		x, y interface{}
		want interface{}
	}{
		{
			name: "int - x less than y",
			x:    3,
			y:    5,
			want: 3,
		},
		{
			name: "int - y less than x",
			x:    4,
			y:    2,
			want: 2,
		},
		{
			name: "int - x equal to y",
			x:    7,
			y:    7,
			want: 7,
		},
		{
			name: "float - x less than y",
			x:    3.3,
			y:    5.5,
			want: 3.3,
		},
		{
			name: "float - y less than x",
			x:    4.4,
			y:    2.2,
			want: 2.2,
		},
		{
			name: "float - x equal to y",
			x:    7.7,
			y:    7.7,
			want: 7.7,
		},
		{
			name: "string - x less than y",
			x:    "apple",
			y:    "banana",
			want: "apple",
		},
		{
			name: "string - y less than x",
			x:    "banana",
			y:    "apple",
			want: "apple",
		},
		{
			name: "string - x equal to y",
			x:    "cherry",
			y:    "cherry",
			want: "cherry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch x := tt.x.(type) {
			case int:
				y := tt.y.(int)
				want := tt.want.(int)
				if got := Min(x, y); got != want {
					t.Errorf("Min() = %v, want %v", got, want)
				}
			case float64:
				y := tt.y.(float64)
				want := tt.want.(float64)
				if got := Min(x, y); got != want {
					t.Errorf("Min() = %v, want %v", got, want)
				}
			case string:
				y := tt.y.(string)
				want := tt.want.(string)
				if got := Min(x, y); got != want {
					t.Errorf("Min() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestDeduplicateSlice(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "int slice with duplicates",
			input: []int{1, 2, 2, 3, 4, 4, 5},
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "int slice without duplicates",
			input: []int{1, 2, 3, 4, 5},
			want:  []int{1, 2, 3, 4, 5},
		},
		{
			name:  "string slice with duplicates",
			input: []string{"apple", "banana", "apple", "cherry", "banana"},
			want:  []string{"apple", "banana", "cherry"},
		},
		{
			name:  "string slice without duplicates",
			input: []string{"apple", "banana", "cherry"},
			want:  []string{"apple", "banana", "cherry"},
		},
		{
			name:  "empty slice",
			input: []int{},
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				want := tt.want.([]int)
				if got := DeduplicateSlice(input); !jsonEqual(got, want) {
					t.Errorf("DeduplicateSlice() = %v, want %v", got, want)
				}
			case []string:
				want := tt.want.([]string)
				if got := DeduplicateSlice(input); !jsonEqual(got, want) {
					t.Errorf("DeduplicateSlice() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		val   interface{}
		want  bool
	}{
		{
			name:  "int slice contains value",
			input: []int{1, 2, 3, 4, 5},
			val:   3,
			want:  true,
		},
		{
			name:  "int slice does not contain value",
			input: []int{1, 2, 3, 4, 5},
			val:   6,
			want:  false,
		},
		{
			name:  "string slice contains value",
			input: []string{"apple", "banana", "cherry"},
			val:   "banana",
			want:  true,
		},
		{
			name:  "string slice does not contain value",
			input: []string{"apple", "banana", "cherry"},
			val:   "date",
			want:  false,
		},
		{
			name:  "empty slice",
			input: []int{},
			val:   1,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case []int:
				val := tt.val.(int)
				if got := Contains(input, val); got != tt.want {
					t.Errorf("Contains() = %v, want %v", got, tt.want)
				}
			case []string:
				val := tt.val.(string)
				if got := Contains(input, val); got != tt.want {
					t.Errorf("Contains() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestCommaSeperatedStringFromInt(t *testing.T) {
	tests := []struct {
		name  string
		input []uint64
		want  string
	}{
		{
			name:  "multiple elements",
			input: []uint64{1, 2, 3, 4, 5},
			want:  "1,2,3,4,5",
		},
		{
			name:  "single element",
			input: []uint64{42},
			want:  "42",
		},
		{
			name:  "empty slice",
			input: []uint64{},
			want:  "",
		},
		{
			name:  "large numbers",
			input: []uint64{123456789, 987654321},
			want:  "123456789,987654321",
		},
		{
			name:  "very large numbers, integer overflow",
			input: []uint64{18446744073709551615},
			want:  "18446744073709551615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommaSeperatedStringFromInt(tt.input); got != tt.want {
				t.Errorf("CommaSeperatedStringFromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		x, y int
		want int
	}{
		{
			name: "x greater than y",
			x:    10,
			y:    5,
			want: 5,
		},
		{
			name: "y greater than x",
			x:    5,
			y:    10,
			want: 5,
		},
		{
			name: "x equal to y",
			x:    7,
			y:    7,
			want: 0,
		},
		{
			name: "negative x and positive y",
			x:    -5,
			y:    10,
			want: 15,
		},
		{
			name: "positive x and negative y",
			x:    10,
			y:    -5,
			want: 15,
		},
		{
			name: "both negative x and y",
			x:    -10,
			y:    -5,
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.x, tt.y); got != tt.want {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertUIntSliceToStringSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []uint64
		want  []string
	}{
		{
			name:  "multiple elements",
			input: []uint64{1, 2, 3, 4, 5},
			want:  []string{"1", "2", "3", "4", "5"},
		},
		{
			name:  "single element",
			input: []uint64{42},
			want:  []string{"42"},
		},
		{
			name:  "large numbers",
			input: []uint64{123456789, 987654321},
			want:  []string{"123456789", "987654321"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertUIntSliceToStringSlice(tt.input); !jsonEqual(got, tt.want) {
				t.Errorf("ConvertUIntSliceToStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertStringSliceToUIntSlice(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []uint64
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   []string{"1", "2", "3", "4", "5"},
			want:    []uint64{1, 2, 3, 4, 5},
			wantErr: false,
		},
		{
			name:    "single element",
			input:   []string{"42"},
			want:    []uint64{42},
			wantErr: false,
		},
		{
			name:    "large numbers",
			input:   []string{"123456789", "987654321"},
			want:    []uint64{123456789, 987654321},
			wantErr: false,
		},
		{
			name:    "very large numbers",
			input:   []string{"18446744073709551615"},
			want:    []uint64{18446744073709551615},
			wantErr: false,
		},
		{
			name:    "invalid number format",
			input:   []string{"1", "two", "3"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative number",
			input:   []string{"-1"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringSliceToUIntSlice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStringSliceToUIntSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !jsonEqual(got, tt.want) {
				t.Errorf("ConvertStringSliceToUIntSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	tests := []struct {
		name   string
		input1 interface{}
		input2 interface{}
		want   interface{}
	}{
		{
			name:   "int slices with intersection",
			input1: []int{1, 2, 3, 4, 5},
			input2: []int{3, 4, 5, 6, 7},
			want:   []int{3, 4, 5},
		},
		{
			name:   "string slices with intersection",
			input1: []string{"apple", "banana", "cherry"},
			input2: []string{"banana", "cherry", "date"},
			want:   []string{"banana", "cherry"},
		},
		{
			name:   "duplicate elements in input slices",
			input1: []int{1, 2, 2, 3, 3, 3},
			input2: []int{2, 2, 3, 3, 4},
			want:   []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input1 := tt.input1.(type) {
			case []int:
				input2 := tt.input2.([]int)
				want := tt.want.([]int)
				if got := Intersection(input1, input2); !jsonEqual(got, want) {
					t.Errorf("Intersection() = %v, want %v", got, want)
				}
			case []string:
				input2 := tt.input2.([]string)
				want := tt.want.([]string)
				if got := Intersection(input1, input2); !jsonEqual(got, want) {
					t.Errorf("Intersection() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name   string
		input1 interface{}
		input2 interface{}
		want   interface{}
	}{
		{
			name:   "int slices with difference",
			input1: []int{1, 2, 3, 4, 5},
			input2: []int{3, 4, 5, 6, 7},
			want:   []int{1, 2},
		},
		{
			name:   "string slices with difference",
			input1: []string{"apple", "banana", "cherry"},
			input2: []string{"banana", "cherry", "date"},
			want:   []string{"apple"},
		},
		{
			name:   "one empty slice",
			input1: []int{1, 2, 3},
			input2: []int{},
			want:   []int{1, 2, 3},
		},
		{
			name:   "duplicate elements in input slices",
			input1: []int{1, 2, 2, 3, 3, 3},
			input2: []int{2, 2, 3, 3, 4},
			want:   []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch input1 := tt.input1.(type) {
			case []int:
				input2 := tt.input2.([]int)
				want := tt.want.([]int)
				if got := Difference(input1, input2); !jsonEqual(got, want) {
					t.Errorf("Difference() = %v, want %v", got, want)
				}
			case []string:
				input2 := tt.input2.([]string)
				want := tt.want.([]string)
				if got := Difference(input1, input2); !jsonEqual(got, want) {
					t.Errorf("Difference() = %v, want %v", got, want)
				}
			}
		})
	}
}

func TestAddCommasInInt(t *testing.T) {
	tests := []struct {
		name string
		num  int64
		want string
	}{
		{
			name: "positive number less than 1000",
			num:  123,
			want: "123",
		},
		{
			name: "positive number with thousands",
			num:  1234,
			want: "1,234",
		},
		{
			name: "positive number with millions",
			num:  1234567,
			want: "12,34,567",
		},
		{
			name: "negative number less than 1000",
			num:  -123,
			want: "-123",
		},
		{
			name: "negative number with thousands",
			num:  -1234,
			want: "-1,234",
		},
		{
			name: "negative number with millions",
			num:  -1234567,
			want: "-12,34,567",
		},
		{
			name: "zero",
			num:  0,
			want: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddCommasInInt(tt.num); got != tt.want {
				t.Errorf("AddCommasInInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple string",
			input: "hello",
			want:  "olleh",
		},
		{
			name:  "string with spaces",
			input: "hello world",
			want:  "dlrow olleh",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single character",
			input: "a",
			want:  "a",
		},
		{
			name:  "palindrome",
			input: "madam",
			want:  "madam",
		},
		{
			name:  "string with special characters",
			input: "h@llo!",
			want:  "!oll@h",
		},
		{
			name:  "string with unicode characters",
			input: "こんにちは",
			want:  "はちにんこ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.input); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddCommasInFloat(t *testing.T) {
	tests := []struct {
		name string
		num  float64
		want string
	}{
		{
			name: "positive float less than 1000",
			num:  123.45,
			want: "123.45",
		},
		{
			name: "positive float with thousands",
			num:  1234.56,
			want: "1,234.56",
		},
		{
			name: "positive float with millions",
			num:  1234567.89,
			want: "12,34,567.89",
		},
		{
			name: "negative float less than 1000",
			num:  -123.45,
			want: "-123.45",
		},
		{
			name: "negative float with thousands",
			num:  -1234.56,
			want: "-1,234.56",
		},
		{
			name: "negative float with millions",
			num:  -1234567.89,
			want: "-12,34,567.89",
		},
		{
			name: "zero",
			num:  0,
			want: "0",
		},
		{
			name: "positive float with no fractional part",
			num:  1234.0,
			want: "1,234",
		},
		{
			name: "negative float with no fractional part",
			num:  -1234.0,
			want: "-1,234",
		},
		{
			name: "positive float with fractional part rounding",
			num:  1234.567,
			want: "1,234.57",
		},
		{
			name: "negative float with fractional part rounding",
			num:  -1234.567,
			want: "-1,234.57",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddCommasInFloat(tt.num); got != tt.want {
				t.Errorf("AddCommasInFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFractionZero(t *testing.T) {
	tests := []struct {
		name    string
		fracStr string
		want    bool
	}{
		{
			name:    "all zeros",
			fracStr: "0000",
			want:    true,
		},
		{
			name:    "contains non-zero digit",
			fracStr: "0010",
			want:    false,
		},
		{
			name:    "single zero",
			fracStr: "0",
			want:    true,
		},
		{
			name:    "empty string",
			fracStr: "",
			want:    true,
		},
		{
			name:    "all non-zero digits",
			fracStr: "1234",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFractionZero(tt.fracStr); got != tt.want {
				t.Errorf("isFractionZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsScriptTag(t *testing.T) {
	tests := []struct {
		name       string
		htmlString string
		want       bool
	}{
		{
			name:       "contains script tag",
			htmlString: "<html><body><script>alert('test');</script></body></html>",
			want:       true,
		},
		{
			name:       "does not contain script tag",
			htmlString: "<html><body><p>Hello, world!</p></body></html>",
			want:       false,
		},
		{
			name:       "script tag with attributes",
			htmlString: "<html><body><script type='text/javascript'>alert('test');</script></body></html>",
			want:       true,
		},
		{
			name:       "empty string",
			htmlString: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsScriptTag(tt.htmlString); got != tt.want {
				t.Errorf("ContainsScriptTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddCommasInInternationalFormat(t *testing.T) {
	tests := []struct {
		name string
		num  float64
		want string
	}{
		{
			name: "positive float less than 1000",
			num:  123.45,
			want: "123.45",
		},
		{
			name: "positive float with thousands",
			num:  1234.56,
			want: "1,234.56",
		},
		{
			name: "positive float with millions",
			num:  1234567.89,
			want: "1,234,567.89",
		},
		{
			name: "negative float less than 1000",
			num:  -123.45,
			want: "-123.45",
		},
		{
			name: "negative float with thousands",
			num:  -1234.56,
			want: "-1,234.56",
		},
		{
			name: "negative float with millions",
			num:  -1234567.89,
			want: "-1,234,567.89",
		},
		{
			name: "zero",
			num:  0,
			want: "0",
		},
		{
			name: "positive float with no fractional part",
			num:  1234.0,
			want: "1,234",
		},
		{
			name: "negative float with no fractional part",
			num:  -1234.0,
			want: "-1,234",
		},
		{
			name: "positive float with fractional part rounding",
			num:  1234.567,
			want: "1,234.57",
		},
		{
			name: "negative float with fractional part rounding",
			num:  -1234.567,
			want: "-1,234.57",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddCommasInInternationalFormat(tt.num); got != tt.want {
				t.Errorf("AddCommasInInternationalFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCIEquals(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{
			name: "equal strings with same case",
			a:    "Hello",
			b:    "Hello",
			want: true,
		},
		{
			name: "equal strings with different case",
			a:    "Hello",
			b:    "hello",
			want: true,
		},
		{
			name: "different strings",
			a:    "Hello",
			b:    "World",
			want: false,
		},
		{
			name: "empty strings",
			a:    "",
			b:    "",
			want: true,
		},
		{
			name: "one empty string",
			a:    "Hello",
			b:    "",
			want: false,
		},
		{
			name: "strings with special characters",
			a:    "Hello!",
			b:    "hello!",
			want: true,
		},
		{
			name: "strings with spaces",
			a:    "Hello World",
			b:    "hello world",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CIEquals(tt.a, tt.b); got != tt.want {
				t.Errorf("CIEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name          string
		value         interface{}
		lang          language.Tag
		decimalPlaces int
		want          string
	}{
		// English format test cases
		{
			name:          "English - integer with thousands",
			value:         int64(1234),
			lang:          language.English,
			decimalPlaces: 0,
			want:          "1,234",
		},
		{
			name:          "English - integer with millions",
			value:         int64(1234567),
			lang:          language.English,
			decimalPlaces: 0,
			want:          "1,234,567",
		},
		{
			name:          "English - float with thousands",
			value:         1234.56,
			lang:          language.English,
			decimalPlaces: 2,
			want:          "1,234.56",
		},
		{
			name:          "English - float with millions",
			value:         1234567.89,
			lang:          language.English,
			decimalPlaces: 2,
			want:          "1,234,567.89",
		},
		{
			name:          "English - small number",
			value:         123,
			lang:          language.English,
			decimalPlaces: 0,
			want:          "123",
		},
		{
			name:          "English - negative integer",
			value:         int64(-1234567),
			lang:          language.English,
			decimalPlaces: 0,
			want:          "-1,234,567",
		},
		{
			name:          "English - negative float",
			value:         -1234567.89,
			lang:          language.English,
			decimalPlaces: 2,
			want:          "-1,234,567.89",
		},
		{
			name:          "English - float with no fractional part",
			value:         1234.0,
			lang:          language.English,
			decimalPlaces: 2,
			want:          "1,234",
		},
		{
			name:          "English - rounding with decimal places",
			value:         1234.567,
			lang:          language.English,
			decimalPlaces: 2,
			want:          "1,234.57",
		},
		{
			name:          "English - with 3 decimal places",
			value:         123.456,
			lang:          language.English,
			decimalPlaces: 3,
			want:          "123.456",
		},
		{
			name:          "English - zero value",
			value:         0,
			lang:          language.English,
			decimalPlaces: 0,
			want:          "0",
		},
		{
			name:          "English - int8 value max",
			value:         int8(127), // Max int8
			lang:          language.English,
			decimalPlaces: 0,
			want:          "127",
		},
		{
			name:          "English - int8 value min",
			value:         int8(-128), // Min int8
			lang:          language.English,
			decimalPlaces: 0,
			want:          "-128",
		},
		{
			name:          "English - uint8 value",
			value:         uint8(255), // Max uint8
			lang:          language.English,
			decimalPlaces: 0,
			want:          "255",
		},
		{
			name:          "English - uint16 value",
			value:         uint16(65535), // Max uint16
			lang:          language.English,
			decimalPlaces: 0,
			want:          "65,535",
		},
		{
			name:          "English - very large integer",
			value:         int64(9223372036854775807), // Max int64
			lang:          language.English,
			decimalPlaces: 0,
			want:          "9,223,372,036,854,775,807",
		},
		{
			name:          "English - special pattern billions",
			value:         int64(1234567890123),
			lang:          language.English,
			decimalPlaces: 0,
			want:          "1,234,567,890,123",
		},
		{
			name:          "English - float with trailing zeros",
			value:         123.45000,
			lang:          language.English,
			decimalPlaces: 2,
			want:          "123.45",
		},
		{
			name:          "English - small float",
			value:         0.12345,
			lang:          language.English,
			decimalPlaces: 5,
			want:          "0.12345",
		},
		{
			name:          "English - smallest positive float64",
			value:         4.940656458412465441765687928682213723651e-324, // Smallest positive float64
			lang:          language.English,
			decimalPlaces: 30,
			want:          "0",
		},
		{
			name:          "English - scientific notation threshold",
			value:         1e6,
			lang:          language.English,
			decimalPlaces: 0,
			want:          "1,000,000",
		},
		{
			name:          "English - int32 max",
			value:         int32(2147483647),
			lang:          language.English,
			decimalPlaces: 0,
			want:          "2,147,483,647",
		},
		{
			name:          "English - int32 min",
			value:         int32(-2147483648),
			lang:          language.English,
			decimalPlaces: 0,
			want:          "-2,147,483,648",
		},
		{
			name:          "English - very small decimal",
			value:         0.00000000001,
			lang:          language.English,
			decimalPlaces: 2,
			want:          "0",
		},

		// Indian format test cases
		{
			name:          "Indian - integer with thousands",
			value:         int64(1000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,000",
		},
		{
			name:          "Indian - ten thousands",
			value:         int64(10000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "10,000",
		},
		{
			name:          "Indian - hundred thousands (1 lakh)",
			value:         int64(100000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,00,000",
		},
		{
			name:          "Indian - millions (10 lakhs)",
			value:         int64(1000000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "10,00,000",
		},
		{
			name:          "Indian - ten millions (1 crore)",
			value:         int64(10000000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,00,00,000",
		},
		{
			name:          "Indian - hundred millions (10 crores)",
			value:         int64(100000000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "10,00,00,000",
		},
		{
			name:          "Indian - billions (100 crores)",
			value:         int64(1000000000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,00,00,00,000",
		},
		{
			name:          "Indian - ten billions (1000 crores)",
			value:         int64(10000000000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "10,00,00,00,000",
		},
		{
			name:          "Indian - hundred billions (10000 crores)",
			value:         int64(100000000000),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,00,00,00,00,000",
		},
		{
			name:          "Indian - float with decimal point",
			value:         123456.78,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 2,
			want:          "1,23,456.78",
		},
		{
			name:          "Indian - negative value",
			value:         int64(-123456),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "-1,23,456",
		},
		{
			name:          "Indian - negative float",
			value:         -9876543.21,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 2,
			want:          "-98,76,543.21",
		},
		{
			name:          "Indian - zero decimal places with rounding",
			value:         12345.6789,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "12,346",
		},
		{
			name:          "Indian - very large integer",
			value:         int64(9223372036854775807), // Max int64
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "92,23,37,20,36,85,47,75,807",
		},
		{
			name:          "Indian - single digit value",
			value:         9,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "9",
		},
		{
			name:          "Indian - two digit value",
			value:         99,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "99",
		},
		{
			name:          "Indian - three digit value",
			value:         999,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "999",
		},
		{
			name:          "Indian - four digit value",
			value:         9999,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "9,999",
		},
		{
			name:          "Indian - five digit value",
			value:         99999,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "99,999",
		},
		{
			name:          "Indian - zero value",
			value:         0,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "0",
		},
		{
			name:          "Indian - float with trailing zeros",
			value:         123.45000,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 2,
			want:          "123.45",
		},
		{
			name:          "Indian - small float",
			value:         0.12345,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 5,
			want:          "0.12345",
		},
		{
			name:          "Indian - complex number combination",
			value:         int64(12345678901234),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,23,45,67,89,01,234",
		},
		{
			name:          "Indian - 12-digit number",
			value:         int64(123456789012),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,23,45,67,89,012",
		},
		{
			name:          "Indian - 16-digit number",
			value:         int64(1234567890123456),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "1,23,45,67,89,01,23,456",
		},
		{
			name:          "Indian - 4-decimal places",
			value:         1234.56789,
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 4,
			want:          "1,234.5679",
		},
		{
			name:          "Indian - int32 max",
			value:         int32(2147483647),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "2,14,74,83,647",
		},
		{
			name:          "Indian - int32 min",
			value:         int32(-2147483648),
			lang:          language.MustParse("en-IN"),
			decimalPlaces: 0,
			want:          "-2,14,74,83,648",
		},

		// Polish format test cases
		{
			name:          "Polish - integer with thousands",
			value:         int64(1234),
			lang:          language.Polish,
			decimalPlaces: 0,
			want:          "1\u00A0234",
		},
		{
			name:          "Polish - integer with millions",
			value:         int64(1234567),
			lang:          language.Polish,
			decimalPlaces: 0,
			want:          "1\u00A0234\u00A0567",
		},
		{
			name:          "Polish - float with thousands",
			value:         1234.56,
			lang:          language.Polish,
			decimalPlaces: 2,
			want:          "1\u00A0234,56",
		},
		{
			name:          "Polish - float with millions",
			value:         1234567.89,
			lang:          language.Polish,
			decimalPlaces: 2,
			want:          "1\u00A0234\u00A0567,89",
		},
		{
			name:          "Polish - negative integer",
			value:         int64(-1234567),
			lang:          language.Polish,
			decimalPlaces: 0,
			want:          "-1\u00A0234\u00A0567",
		},
		{
			name:          "Polish - negative float",
			value:         -1234.56,
			lang:          language.Polish,
			decimalPlaces: 2,
			want:          "-1\u00A0234,56",
		},

		// Dutch format test cases
		{
			name:          "Dutch - integer with thousands",
			value:         int64(1234),
			lang:          language.Dutch,
			decimalPlaces: 0,
			want:          "1.234",
		},
		{
			name:          "Dutch - integer with millions",
			value:         int64(1234567),
			lang:          language.Dutch,
			decimalPlaces: 0,
			want:          "1.234.567",
		},
		{
			name:          "Dutch - float with thousands",
			value:         1234.56,
			lang:          language.Dutch,
			decimalPlaces: 2,
			want:          "1.234,56",
		},
		{
			name:          "Dutch - float with millions",
			value:         1234567.89,
			lang:          language.Dutch,
			decimalPlaces: 2,
			want:          "1.234.567,89",
		},
		{
			name:          "Dutch - negative integer",
			value:         int64(-1234567),
			lang:          language.Dutch,
			decimalPlaces: 0,
			want:          "-1.234.567",
		},
		{
			name:          "Dutch - negative float",
			value:         -1234.56,
			lang:          language.Dutch,
			decimalPlaces: 2,
			want:          "-1.234,56",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			switch v := tt.value.(type) {
			case int:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case int8:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case int16:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case int32:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case int64:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case uint:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case uint8:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case uint16:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case uint32:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case uint64:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case float32:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			case float64:
				got = FormatNumber(v, tt.decimalPlaces, tt.lang)
			}

			if got != tt.want {
				t.Errorf("FormatNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func jsonEqual(a, b interface{}) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}
