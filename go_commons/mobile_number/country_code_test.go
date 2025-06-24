package mobile_number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCountryCallingCode(t *testing.T) {
	var tests = []struct {
		input        string
		err          error
		expectedCode string
	}{
		{
			input: "IN", err: nil,
			expectedCode: "+91",
		},
		{
			input: "SA", err: nil,
			expectedCode: "+966",
		},
		{
			input: "US", err: nil,
			expectedCode: "+1",
		},
		{
			input: "GB", err: nil,
			expectedCode: "+44",
		},
		{
			input: "AU", err: nil,
			expectedCode: "+61",
		},
		{
			input: "XX", err: ErrCountryCodeNotValid,
			expectedCode: "",
		},
		{
			input: "", err: ErrCountryCodeNotValid,
			expectedCode: "",
		},
	}

	for _, tc := range tests {
		got, err := GetCountryCallingCode(tc.input)

		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error(), "error mismatch for input %s", tc.input)
			assert.Equal(t, tc.expectedCode, got, "result mismatch for input %s", tc.input)
		} else {
			assert.NoError(t, err, "unexpected error %s for input %s", err, tc.input)
			assert.Equal(t, tc.expectedCode, got, "result mismatch for input %s", tc.input)
		}
	}
}
