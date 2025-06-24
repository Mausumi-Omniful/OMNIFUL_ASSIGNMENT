package private

import (
	"context"

	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/omniful/go_commons/constants"
	error2 "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/jwt"
	"github.com/omniful/go_commons/newrelic"
	"github.com/omniful/go_commons/permissions"
	"github.com/omniful/go_commons/rules"
)

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

type UserDetails struct {
	UserID          string                       `json:"user_id"`
	UserType        jwt.UserType                 `json:"user_type"`
	UserName        string                       `json:"user_name"`
	UserEmail       string                       `json:"user_email"`
	TenantID        string                       `json:"tenant_id"`
	TenantType      jwt.TenantType               `json:"tenant_type"`
	TenantName      string                       `json:"tenant_name"`
	UserTimeZone    string                       `json:"user_time_zone"`
	RuleGroup       *rules.RuleGroup             `json:"rule_group"`
	Permissions     []permissions.PermissionType `json:"permissions"`
	Roles           []string                     `json:"roles"`
	Environment     jwt.Environment              `json:"environment"`
	CustomerDetails *CustomerDetails             `json:"customer_details"`
}

type CustomerDetails struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SellerID string `json:"seller_id"`
}

func GetUserID(ctx context.Context) (string, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.UserNotFoundInCtx, "user details not found in ctx")
	}
	return userDetails.UserID, error2.CustomError{}
}

func GetUserType(ctx context.Context) (jwt.UserType, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return 0, error2.NewCustomError(error2.UserTypeNotFoundInCtx, "user details not found in ctx")
	}
	return userDetails.UserType, error2.CustomError{}
}

func GetTenantID(ctx context.Context) (string, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.TenantIDNotFoundInCtx, "user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return "", error2.NewCustomError(error2.UserIsOmnifulUserNotTenant, "user is omniful user")
	}
	return userDetails.TenantID, error2.CustomError{}
}

func GetTenantType(ctx context.Context) (jwt.TenantType, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return 0, error2.NewCustomError(error2.TenantTypeNotFoundInCtx, "user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return 0, error2.NewCustomError(error2.UserIsOmnifulUserNotTenant, "user is omniful user")
	}
	return userDetails.TenantType, error2.CustomError{}
}

func GetTenantName(ctx context.Context) (string, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.TenantNameNotFoundInCtx, "user details not found in ctx")
	}

	if userDetails.UserType == jwt.Omniful {
		return "", error2.NewCustomError(error2.UserIsOmnifulUserNotTenant, "user is omniful user")
	}
	return userDetails.TenantName, error2.CustomError{}
}

func GetUserTimeZone(ctx context.Context) (string, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.RulesNotFoundInCtx, "user timezone not found in ctx")
	}
	return userDetails.UserTimeZone, error2.CustomError{}
}

func GetUserName(ctx context.Context) (string, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.RulesNotFoundInCtx, "user name not found in ctx")
	}
	return userDetails.UserName, error2.CustomError{}
}

func GetUserEmail(ctx context.Context) (string, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.RulesNotFoundInCtx, "user email not found in ctx")
	}
	return userDetails.UserEmail, error2.CustomError{}
}

func GetRuleGroup(ctx context.Context) (*rules.RuleGroup, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return nil, error2.NewCustomError(error2.RulesNotFoundInCtx, "rules not found")
	}
	return userDetails.RuleGroup, error2.CustomError{}
}

func GetUserPermissions(ctx context.Context) ([]permissions.PermissionType, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return nil, error2.NewCustomError(error2.PermissionsNotFoundInCtx, "permissions not found")
	}
	return userDetails.Permissions, error2.CustomError{}
}

func GetUserRoles(ctx context.Context) ([]string, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return nil, error2.NewCustomError(error2.PermissionsNotFoundInCtx, "roles not found")
	}
	return userDetails.Roles, error2.CustomError{}
}

func GetUserPlan(ctx context.Context) (*UserPlan, error2.CustomError) {
	userPlan, ok := ctx.Value(constants.PrivateUserPlan).(*UserPlan)
	if !ok {
		return nil, error2.NewCustomError(error2.PlanNotFoundInCtx, "plan not found")
	}

	return userPlan, error2.CustomError{}
}

func GetEnvironment(ctx context.Context) (jwt.Environment, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return "", error2.NewCustomError(error2.EnvironmentNotFoundInCtx, "environment not found")
	}
	return userDetails.Environment, error2.CustomError{}
}

func GetUserDetails(ctx context.Context) (*UserDetails, error2.CustomError) {
	userDetails, ok := ctx.Value(constants.PrivateUserDetails).(*UserDetails)
	if !ok {
		return nil, error2.NewCustomError(error2.UserDetailsNotFoundInCtx, "user details not found in ctx")
	}
	return userDetails, error2.CustomError{}
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
