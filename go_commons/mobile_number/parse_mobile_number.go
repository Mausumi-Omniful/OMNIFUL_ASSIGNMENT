package mobile_number

import (
	"errors"
	"fmt"
	"github.com/nyaruka/phonenumbers"
)

var ErrNumberNotValidForRegion = errors.New("mobile number not valid for provided region")

type MobileNumber struct {
	number             string // Number doesn't contain any formatting and prefixes. (Example: "9876543211","561234567")
	countryCode        string // CountryCode to which number belongs. (Example: "IN","SA","US")
	countryCallingCode int32  // CountryCallingCode attached before any number. (Example: "+91","+966")
}

func (pn *MobileNumber) DisplayNumber() string {
	if pn == nil {
		return ""
	}

	if pn.countryCallingCode > 0 {
		return fmt.Sprintf("+%d-%s", pn.countryCallingCode, pn.number)
	}

	return pn.number
}

func (pn *MobileNumber) Number() string {
	if pn == nil {
		return ""
	}

	return pn.number
}

func (pn *MobileNumber) CountryCode() string {
	if pn == nil {
		return ""
	}

	return pn.countryCode
}

func (pn *MobileNumber) CountryCallingCode() int32 {
	if pn == nil {
		return 0
	}

	return pn.countryCallingCode
}

// ParseNumber Only parses the number, doesn't validate whether number is correct or not
func ParseNumber(number, countryCode string) (*MobileNumber, error) {
	phoneNum, err := phonenumbers.Parse(number, countryCode)
	if err != nil {
		return nil, err
	}

	return getMobileNumber(phoneNum), nil
}

// ParseNumberWithValidation Parses the number and validates number against country code as well.
// Country Code is mandatory.
// ErrNumberNotValidForRegion error will be raised if validation fails.
func ParseNumberWithValidation(number, countryCode string) (*MobileNumber, error) {
	phoneNum, err := phonenumbers.Parse(number, countryCode)
	if err != nil {
		return nil, err
	}

	if !phonenumbers.IsValidNumberForRegion(phoneNum, countryCode) {
		return getMobileNumber(phoneNum), ErrNumberNotValidForRegion
	}

	return getMobileNumber(phoneNum), nil
}

func getMobileNumber(number *phonenumbers.PhoneNumber) *MobileNumber {
	if number == nil {
		return nil
	}

	return &MobileNumber{
		number:             phonenumbers.GetNationalSignificantNumber(number),
		countryCode:        phonenumbers.GetRegionCodeForNumber(number),
		countryCallingCode: number.GetCountryCode(),
	}
}
