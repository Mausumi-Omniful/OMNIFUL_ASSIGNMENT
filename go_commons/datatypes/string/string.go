package string

import (
	"net/mail"
	"strings"
	"unicode"
)

type String string

func (str String) String() string {
	return string(str)
}

func (str String) Sanitize(includeChars ...String) String {
	var sb strings.Builder
	includeSet := make(map[rune]bool)

	for _, chars := range includeChars {
		for _, ch := range chars {
			includeSet[ch] = true
		}
	}

	for _, ch := range str {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || includeSet[ch] {
			sb.WriteRune(ch)
		}
	}

	return String(sb.String())
}

func (str String) IsValidEmailFormat() bool {
	_, err := mail.ParseAddress(str.String())
	return err == nil
}
