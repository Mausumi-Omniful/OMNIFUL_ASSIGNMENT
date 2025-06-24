package csv

import "strings"

func addTrailingSlash(str string) (result string) {
	str = strings.TrimSpace(str)
	// Check if the string is blank
	if len(str) == 0 {
		return result
	}

	// Check if the string ends with "/"
	if str[len(str)-1:] != "/" {
		// If not, add "/" to the end of the string
		str += "/"
	}
	return str
}
