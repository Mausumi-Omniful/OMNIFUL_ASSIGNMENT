package jwt

type UserType uint64

type TenantType uint64

type Environment string

const (
	Omniful  UserType = 1
	Tenant   UserType = 2
	Customer UserType = 3
)

const (
	SelfUsage       TenantType = 1
	SharedHub       TenantType = 2
	ShippingPartner TenantType = 3
	RabtFulfilment  TenantType = 4
)

const (
	Live    Environment = "live"
	Testing Environment = "testing"
)
