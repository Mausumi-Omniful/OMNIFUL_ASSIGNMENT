package public

import (
	"context"
	"errors"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/omniful/go_commons/constants"
)

type CustomerDetails struct {
	ID         string
	Name       string
	TenantCode string
	SellerCode string
}

type CustomerJWTClaim struct {
	jwt2.StandardClaims
	CustomerDetails CustomerDetails `json:"customer_details"`
}

func GetCustomerID(ctx context.Context) (string, error) {
	customerDetails, ok := ctx.Value(constants.PublicCustomerDetails).(*CustomerDetails)
	if !ok {
		return "", errors.New("customer details not found in ctx")
	}

	return customerDetails.ID, nil
}

func GetCustomerName(ctx context.Context) (string, error) {
	customerDetails, ok := ctx.Value(constants.PublicCustomerDetails).(*CustomerDetails)
	if !ok {
		return "", errors.New("customer details not found in ctx")
	}
	
	return customerDetails.Name, nil
}

func GetCustomerTenantCode(ctx context.Context) (string, error) {
	customerDetails, ok := ctx.Value(constants.PublicCustomerDetails).(*CustomerDetails)
	if !ok {
		return "", errors.New("customer details not found in ctx")
	}

	return customerDetails.TenantCode, nil
}

func GetCustomerSellerCode(ctx context.Context) (string, error) {
	customerDetails, ok := ctx.Value(constants.PublicCustomerDetails).(*CustomerDetails)
	if !ok {
		return "", errors.New("customer details not found in ctx")
	}

	return customerDetails.SellerCode, nil
}
