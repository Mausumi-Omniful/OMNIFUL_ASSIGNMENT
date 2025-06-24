package public

import (
	"context"
	"errors"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/omniful/go_commons/constants"
	error2 "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/jwt"
	"github.com/omniful/go_commons/newrelic"
)

type UserDetails struct {
	UserID       string
	UserType     jwt.UserType
	UserName     string
	UserEmail    string
	TenantID     string
	TenantCode   string
	TenantType   jwt.TenantType
	TenantName   string
	UserTimeZone string
	TenantDomain string
	Environment  jwt.Environment
}

type JWTClaim struct {
	jwt2.StandardClaims
	UserDetails UserDetails `json:"user_details"`
	UserPlan    UserPlan    `json:"user_plan"`
}

type UserPlan struct {
	IsOmsEnabled         bool `json:"is_oms_enabled"`
	IsFleetEnabled       bool `json:"is_fleet_enabled"`
	IsWmsEnabled         bool `json:"is_wms_enabled"`
	IsExternalWmsEnabled bool `json:"is_external_wms_enabled"`
	IsPosEnabled         bool `json:"is_pos_enabled"`
	IsOmniShipEnabled    bool `json:"is_omni_ship_enabled"`
	IsTmsEnabled         bool `json:"is_tms_enabled"`
	IsTransitWmsEnabled  bool `json:"is_transit_wms_enabled"`
}

func GetUserID(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}
	return userDetails.UserID, nil
}

func GetUserType(ctx context.Context) (jwt.UserType, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return 0, errors.New("user details not found in ctx")
	}
	return userDetails.UserType, nil
}

func GetUserTimeZone(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}
	return userDetails.UserTimeZone, nil
}

func GetUserName(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}
	return userDetails.UserName, nil
}

func GetUserEmail(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}
	return userDetails.UserEmail, nil
}

func GetTenantID(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return "", errors.New("user is omniful user")
	}
	return userDetails.TenantID, nil
}

func GetTenantCode(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return "", errors.New("user is omniful user")
	}

	return userDetails.TenantCode, nil
}

func GetTenantType(ctx context.Context) (jwt.TenantType, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return 0, errors.New("user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return 0, errors.New("user is omniful user")
	}
	return userDetails.TenantType, nil
}

func GetTenantName(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return "", errors.New("user is omniful user")
	}
	return userDetails.TenantName, nil
}

func GetTenantDomain(ctx context.Context) (string, error) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", errors.New("user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return "", errors.New("user is omniful user")
	}

	return userDetails.TenantDomain, nil
}

func GetUserPlan(ctx context.Context) (*UserPlan, error2.CustomError) {
	userPlan, ok := ctx.Value(constants.PublicUserPlan).(*UserPlan)
	if !ok {
		return nil, error2.NewCustomError(error2.PlanNotFoundInCtx, "plan not found")
	}

	return userPlan, error2.CustomError{}
}

func GetEnvironment(ctx context.Context) (jwt.Environment, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PublicUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.EnvironmentNotFoundInCtx, "environment not found")
	}
	return userDetails.Environment, error2.CustomError{}
}

func AddCommonUserDetailAttributesInNewRelic(userDetails *UserDetails, ctx context.Context) {
	newrelic.AddAttributeWithContext(
		ctx,
		newrelic.Attribute{
			Name:  constants.AttributeTenantID,
			Value: userDetails.TenantID,
		},
		newrelic.Attribute{
			Name:  constants.AttributeTenantName,
			Value: userDetails.TenantName,
		},
		newrelic.Attribute{
			Name:  constants.AttributeTenantEnvironment,
			Value: userDetails.Environment,
		},
		newrelic.Attribute{
			Name:  constants.AttributeUserID,
			Value: userDetails.UserID,
		},
	)
}
