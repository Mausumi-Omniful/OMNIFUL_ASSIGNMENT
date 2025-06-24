package mobile_number

import (
	"errors"
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

var ErrCountryCodeNotValid = errors.New("country code is not valid")

// GetCountryCallingCode returns the international calling code for a given ISO 3166-1 alpha-2 country code.
// The country code should be a two-letter code like "IN" for India or "SA" for Saudi Arabia.
//
// Examples:
//   - GetCountryCallingCode("IN") returns ("+91", nil) for India
//   - GetCountryCallingCode("SA") returns ("+966", nil) for Saudi Arabia
//   - GetCountryCallingCode("US") returns ("+1", nil) for United States
//   - GetCountryCallingCode("XX") returns ("", error) for invalid country code
//
// Returns:
//   - The calling code as a string with "+" prefix (e.g. "+91" for India)
//   - An error if the provided country code is invalid
func GetCountryCallingCode(countryCode string) (string, error) {
	callingCode := phonenumbers.GetCountryCodeForRegion(countryCode)
	if callingCode == 0 {
		return "", ErrCountryCodeNotValid
	}
	return fmt.Sprintf("+%d", callingCode), nil
}
