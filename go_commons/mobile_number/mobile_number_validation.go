package mobile_number

import (
	"errors"
	"fmt"
	"github.com/nyaruka/phonenumbers"
)

func ValidateMobileNumber(callingCode, mobileNumber, countryCode string) (bool, error) {
	parsedNumber, err := phonenumbers.Parse(callingCode+mobileNumber, countryCode)
	if err != nil {
		return false, err
	}
	isValid := phonenumbers.IsValidNumberForRegion(parsedNumber, countryCode)
	if !isValid {
		return false, errors.New(fmt.Sprintf(
			"uable to validate the mobile number for this number: %s, country code: %s, and calling code: %s",
			mobileNumber, countryCode, callingCode),
		)
	}

	return true, nil
}
