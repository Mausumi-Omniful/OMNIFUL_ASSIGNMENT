package mobile_number

import (
	"github.com/nyaruka/phonenumbers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseNumber(t *testing.T) {
	var tests = []struct {
		input       string
		region      string
		err         error
		expectedNum MobileNumber
	}{
		{
			input: "9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "919876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+91-9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+919876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+91 9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "91 9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "09876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "9876543211", region: "SA", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "SA", countryCallingCode: 966},
		},
		{
			input: "919876543211", region: "SA", err: nil,
			expectedNum: MobileNumber{number: "919876543211", countryCode: "SA", countryCallingCode: 966},
		},
		{
			input: "+91-9876543211", region: "SA", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+919876543211", region: "SA", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+91 9876543211", region: "SA", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "91 9876543211", region: "SA", err: nil,
			expectedNum: MobileNumber{number: "919876543211", countryCode: "SA", countryCallingCode: 966},
		},
		{
			input: "09876543211", region: "SA", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "SA", countryCallingCode: 966},
		},
		{
			input: "9876543211", region: "", err: phonenumbers.ErrInvalidCountryCode,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "919876543211", region: "", err: phonenumbers.ErrInvalidCountryCode,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+91-9876543211", region: "", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+919876543211", region: "", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+91 9876543211", region: "", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "91 9876543211", region: "", err: phonenumbers.ErrInvalidCountryCode,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "09876543211", region: "", err: phonenumbers.ErrInvalidCountryCode,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
	}

	for _, tc := range tests {
		num, err := ParseNumber(tc.input, tc.region)

		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error(), "error mismatch for input %s, input region: %s", tc.input, tc.region)
		} else {
			assert.NoError(t, err, "unexpected error %s for input %s", err, tc.input, tc.region)
			assert.Equal(t, tc.expectedNum.Number(), num.Number(), "number mismatch for input %s, input region: %s", tc.input, tc.region)
			assert.Equal(t, tc.expectedNum.CountryCallingCode(), num.CountryCallingCode(), "country calling code mismatch for input %s, input region: %s", tc.input, tc.region)
			assert.Equal(t, tc.expectedNum.CountryCode(), num.CountryCode(), "country code mismatch for input %s, input region: %s", tc.input, tc.region)
		}
	}
}

func TestParseNumberWithValidation(t *testing.T) {
	var tests = []struct {
		input       string
		region      string
		err         error
		expectedNum MobileNumber
	}{
		{
			input: "9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "919876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+91-9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+919876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "+91 9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "91 9876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "09876543211", region: "IN", err: nil,
			expectedNum: MobileNumber{number: "9876543211", countryCode: "IN", countryCallingCode: 91},
		},
		{
			input: "9876543211", region: "SA", err: ErrNumberNotValidForRegion,
			expectedNum: MobileNumber{},
		},
		{
			input: "919876543211", region: "SA", err: ErrNumberNotValidForRegion,
			expectedNum: MobileNumber{},
		},
		{
			input: "+91-9876543211", region: "SA", err: ErrNumberNotValidForRegion,
			expectedNum: MobileNumber{},
		},
		{
			input: "+919876543211", region: "SA", err: ErrNumberNotValidForRegion,
			expectedNum: MobileNumber{},
		},
		{
			input: "+91 9876543211", region: "SA", err: ErrNumberNotValidForRegion,
			expectedNum: MobileNumber{},
		},
		{
			input: "91 9876543211", region: "SA", err: ErrNumberNotValidForRegion,
			expectedNum: MobileNumber{},
		},
		{
			input: "09876543211", region: "SA", err: ErrNumberNotValidForRegion,
			expectedNum: MobileNumber{},
		},
	}

	for _, tc := range tests {
		num, err := ParseNumberWithValidation(tc.input, tc.region)

		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error(), "error mismatch for input %s, input region: %s", tc.input, tc.region)
		} else {
			assert.NoError(t, err, "unexpected error %s for input %s", err, tc.input, tc.region)
			assert.Equal(t, tc.expectedNum.Number(), num.Number(), "number mismatch for input %s, input region: %s", tc.input, tc.region)
			assert.Equal(t, tc.expectedNum.CountryCallingCode(), num.CountryCallingCode(), "country calling code mismatch for input %s, input region: %s", tc.input, tc.region)
			assert.Equal(t, tc.expectedNum.CountryCode(), num.CountryCode(), "country code mismatch for input %s, input region: %s", tc.input, tc.region)
		}
	}
}
