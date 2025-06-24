package csv

import (
	"strings"
)

type sanitizerFunc func([]string) []string

func Sanitize(strSlice []string, sanitizers []sanitizerFunc) []string {
	sanitized := strSlice
	for _, s := range sanitizers {
		sanitized = s(sanitized)
	}
	return sanitized
}

func SanitizeAsterisks(strSlice []string) []string {
	var sanitized []string
	for _, str := range strSlice {
		// Remove asterisk prefix
		for strings.HasPrefix(str, "*") {
			str = str[1:]
		}
		// Remove asterisk suffix
		for strings.HasSuffix(str, "*") {
			str = str[:len(str)-1]
		}
		sanitized = append(sanitized, str)
	}
	return sanitized
}

func SanitizeToLower(strSlice []string) []string {
	var sanitized []string
	for _, str := range strSlice {
		sanitized = append(sanitized, strings.ToLower(str))
	}

	return sanitized
}

func SanitizeSpace(strSlice []string) []string {
	var sanitized []string
	for _, str := range strSlice {
		sanitized = append(sanitized, strings.TrimSpace(str))
	}

	return sanitized
}

func SanitizeUpperCase(strSlice []string) []string {
	var sanitized []string

	for _, str := range strSlice {
		sanitized = append(sanitized, strings.ToLower(strings.ReplaceAll(str, " ", "_")))
	}

	return sanitized
}

func sanitizeBOM(strSlice []string) []string {
	var sanitized []string
	for _, str := range strSlice {
		if containsBOM(str) {
			str = removeBOM(str)
		}

		sanitized = append(sanitized, str)
	}

	return sanitized
}
