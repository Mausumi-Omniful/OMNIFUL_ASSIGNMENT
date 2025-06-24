package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/omniful/go_commons/constants"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
)

func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func DeduplicateSlice[T comparable](input []T) []T {
	unique := make([]T, 0, len(input))
	occurrenceMap := make(map[T]bool)

	for _, val := range input {
		if _, ok := occurrenceMap[val]; !ok {
			occurrenceMap[val] = true
			unique = append(unique, val)
		}
	}
	return unique
}

func Contains[T comparable](input []T, val T) bool {
	for _, v := range input {
		if v == val {
			return true
		}
	}
	return false
}

func CommaSeperatedStringFromInt(input []uint64) string {
	b := make([]string, len(input))
	for i, v := range input {
		b[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(b, ",")
}

func Abs[T constraints.Integer](x, y T) T {
	if x < y {
		return y - x
	}
	return x - y
}

func ConvertUIntSliceToStringSlice(input []uint64) (output []string) {
	for _, v := range input {
		output = append(output, strconv.FormatUint(v, 10))
	}

	return
}

// ConvertStringSliceToUIntSlice converts a slice of strings to a slice of uint64
// If any of the strings cannot be converted to uint64, it returns an error
func ConvertStringSliceToUIntSlice(input []string) (output []uint64, err error) {
	for _, v := range input {
		parsedUint, parseErr := strconv.ParseUint(v, 10, 64)
		if parseErr != nil {
			err = parseErr
			return
		}
		output = append(output, parsedUint)
	}

	return
}

func Intersection[T comparable](input1, input2 []T) (output []T) {
	uniqueInput1 := DeduplicateSlice(input1)
	uniqueInput2 := DeduplicateSlice(input2)
	input1Map := make(map[T]bool)
	for _, v := range uniqueInput1 {
		input1Map[v] = true
	}

	for _, v := range uniqueInput2 {
		if _, ok := input1Map[v]; ok {
			output = append(output, v)
		}
	}

	return
}

func Difference[T comparable](input1, input2 []T) (output []T) {
	uniqueInput1 := DeduplicateSlice(input1)
	uniqueInput2 := DeduplicateSlice(input2)
	input2Map := make(map[T]bool)
	for _, v := range uniqueInput2 {
		input2Map[v] = true
	}

	for _, v := range uniqueInput1 {
		if _, ok := input2Map[v]; !ok {
			output = append(output, v)
		}
	}

	return
}

func AddCommasInInt(num int64) string {
	isNeg := num < 0
	if isNeg {
		num *= -1
	}

	strNumber := strconv.FormatInt(num, 10)
	if len(strNumber) < 4 {
		if isNeg {
			strNumber = "-" + strNumber
		}

		return strNumber
	}

	reverseNum := ReverseString(strNumber)
	formattedNumber := reverseNum[:3]
	for i := 3; i < len(reverseNum); i++ {
		if (i-1)%2 == 0 {
			formattedNumber += ","
		}

		formattedNumber += string(reverseNum[i])
	}

	if isNeg {
		formattedNumber += "-"
	}

	return ReverseString(formattedNumber)
}

func AddCommasInFloat(num float64) string {
	isNeg := num < 0
	if isNeg {
		num *= -1
	}

	intPart := int64(num)
	fracPart := num - float64(intPart)
	roundedFracPart := math.Round(fracPart*100) / 100
	fracStr := "0"
	if roundedFracPart > 0 {
		fracStr = strings.TrimPrefix(strconv.FormatFloat(roundedFracPart, 'f', -1, 64), "0.")
	}

	strIntPart := strconv.FormatInt(intPart, 10)
	if len(strIntPart) < 4 {
		if isNeg {
			strIntPart = "-" + strIntPart
		}

		if len(fracStr) > 0 && !isFractionZero(fracStr) {
			strIntPart += "." + fracStr
		}

		return strIntPart
	}

	reverseIntPart := ReverseString(strIntPart)
	formattedNum := reverseIntPart[:3]
	for i := 3; i < len(reverseIntPart); i++ {
		if (i-1)%2 == 0 {
			formattedNum += ","
		}
		formattedNum += string(reverseIntPart[i])
	}

	if isNeg {
		formattedNum += "-"
	}

	formattedNum = ReverseString(formattedNum)
	if len(fracStr) > 0 && !isFractionZero(fracStr) {
		formattedNum += "." + fracStr
	}

	return formattedNum
}

func AddCommasInInternationalFormat(num float64) string {
	isNeg := num < 0
	if isNeg {
		num *= -1
	}

	intPart := int64(num)
	fracPart := num - float64(intPart)
	roundedFracPart := math.Round(fracPart*100) / 100
	fracStr := "0"
	if roundedFracPart > 0 {
		fracStr = strings.TrimPrefix(strconv.FormatFloat(roundedFracPart, 'f', -1, 64), "0.")
	}

	strIntPart := strconv.FormatInt(intPart, 10)
	if len(strIntPart) < 4 {
		if isNeg {
			strIntPart = "-" + strIntPart
		}

		if len(fracStr) > 0 && !isFractionZero(fracStr) {
			strIntPart += "." + fracStr
		}

		return strIntPart
	}

	reverseIntPart := ReverseString(strIntPart)
	formattedNum := reverseIntPart[:2]
	for i := 2; i < len(reverseIntPart); i++ {
		if i%3 == 0 {
			formattedNum += ","
		}
		formattedNum += string(reverseIntPart[i])
	}

	if isNeg {
		formattedNum += "-"
	}

	formattedNum = ReverseString(formattedNum)
	if len(fracStr) > 0 && !isFractionZero(fracStr) {
		formattedNum += "." + fracStr
	}

	return formattedNum
}

func ReverseString(s string) string {
	reversed := []rune(s)
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	return string(reversed)
}

func isFractionZero(fracStr string) bool {
	for _, s := range fracStr {
		if s != '0' {
			return false
		}
	}

	return true
}

func ContainsScriptTag(htmlString string) bool {
	return constants.ScriptTagRegex.MatchString(htmlString)
}

// CIEquals Case Insensitive Equality between two strings
func CIEquals(a, b string) bool {
	return cases.Fold().String(a) == cases.Fold().String(b)
}

func FirstNonEmptyString(v ...string) string {
	for _, s := range v {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

func RemoveEmptyStrings(input []string) []string {
	output := make([]string, 0, len(input))
	for _, s := range input {
		if len(s) > 0 {
			output = append(output, s)
		}
	}
	return output
}
