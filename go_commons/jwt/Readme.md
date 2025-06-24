# JWT Package

## Overview
The JWT package provides functionality for creating and managing JSON Web Tokens (JWT) in Go applications. It supports both signed and unsigned token creation with flexible claim management, making it suitable for various authentication and authorization scenarios.

## Features
- Create signed JWT tokens using RSA private keys
- Create unsigned JWT tokens for development/testing purposes
- Support for custom claims
- Built-in user and tenant type management
- Separate implementations for public and internal service authentication
- Rich user context management

## Package Structure

The package is organized into three main components:

### 1. Core JWT (Root Directory)
Contains the base types and token creation functions:

```go
type UserType uint64

const (
    Omniful UserType = 1
    Tenant  UserType = 2
)

type TenantType uint64

const (
    SelfUsage       TenantType = 1
    SharedHub       TenantType = 2
    ShippingPartner TenantType = 3
)
```

### 2. Public JWT (`public/` Directory)
Handles external/client-facing authentication:
- Uses RSA public key verification for enhanced security
- Designed for end-user authentication flows
- Contains rich user context information to minimize inter-service calls
- Includes user plan details and tenant information

Key components:
```go
type UserDetails struct {
    UserID       string
    UserType     jwt.UserType
    UserName     string
    UserEmail    string
    TenantID     string
    TenantType   jwt.TenantType
    TenantName   string
    UserTimeZone string
    TenantDomain string
}

type UserPlan struct {
    IsOmsEnabled         bool
    IsFleetEnabled       bool
    IsWmsEnabled         bool
    IsExternalWmsEnabled bool
    IsPosEnabled        bool
    IsOmniShipEnabled   bool
    IsTmsEnabled        bool
}
```

Helper functions for context management:
- `GetUserID(ctx)` - Retrieve user ID from context
- `GetUserType(ctx)` - Get user type (Omniful/Tenant)
- `GetTenantID(ctx)` - Get tenant identifier
- `GetUserPlan(ctx)` - Retrieve user's plan details
- And more context accessors for user details

### 3. Private JWT (`private/` Directory)
Handles inter-service communication:
- Uses `jwt.UnsafeAllowNoneSignatureType` for internal service verification
- Optimized for service-to-service communication
- Maintains minimal but sufficient user context
- Includes middleware for automatic JWT parsing and context injection

## Functions

### CreateJWTTokenWithSignature
Creates a signed JWT token using RSA private key.

```go
func CreateJWTTokenWithSignature(method jwt.SigningMethod, claims jwt.Claims, privateKey *rsa.PrivateKey) (string, error)
```

Example:
```go
package main

import (
    "crypto/rsa"
    "github.com/dgrijalva/jwt-go"
    "github.com/omniful/go_commons/jwt"
)

func main() {
    // Create custom claims
    claims := jwt.MapClaims{
        "user_id": 123,
        "type": jwt.Omniful,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    
    // Load your RSA private key
    var privateKey *rsa.PrivateKey // Load your private key
    
    // Create signed token
    token, err := jwt.CreateJWTTokenWithSignature(jwt.SigningMethodRS256, claims, privateKey)
    if err != nil {
        panic(err)
    }
}
```

### CreateJWTTokenWithoutSignature
Creates an unsigned JWT token, useful for testing and internal service communication.

```go
func CreateJWTTokenWithoutSignature(claims jwt.Claims) (string, error)
```

Example:
```go
package main

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/omniful/go_commons/jwt"
)

func main() {
    claims := jwt.MapClaims{
        "user_id": 123,
        "type": jwt.Tenant,
        "tenant_type": jwt.SharedHub,
    }
    
    token, err := jwt.CreateJWTTokenWithoutSignature(claims)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Unsigned token: %s\n", token)
}
```

## Usage Guidelines

### Public JWT Usage
1. Use for all external API endpoints
2. Always verify tokens using RSA public key
3. Leverage the rich user context to avoid unnecessary service calls
4. Check user plan details before allowing access to features
5. Handle tenant-specific logic using the provided tenant information

### Private JWT Usage
1. Use for service-to-service communication only
2. Implement the middleware in your service
3. Access user context through the provided helper functions
4. Be aware that these tokens use `UnsafeAllowNoneSignatureType`
5. Keep the service communication within trusted network boundaries

## Security Considerations
1. Always use `CreateJWTTokenWithSignature` with RSA keys for public endpoints
2. Keep your RSA private keys secure and never expose them
3. Use appropriate token expiration times
4. Include necessary claims like "exp" (expiration time) and "iat" (issued at)
5. Only use `CreateJWTTokenWithoutSignature` for internal service communication
6. Ensure your internal network is properly secured when using private JWT

## Dependencies
- github.com/dgrijalva/jwt-go: For JWT operations

## Directory Structure
```
jwt/
├── create.go      # Token creation functions
├── jwt.go         # Type definitions
├── private/       # Internal service JWT implementation
│   ├── middleware.go    # Service-to-service JWT middleware
│   └── user_detail.go   # Internal user context management
└── public/        # External JWT implementation
    ├── key.go          # RSA public key management
    └── user_detail.go  # Rich user context management
```

## Best Practices
1. Always validate tokens on the server side
2. Use appropriate user and tenant types based on your application needs
3. Implement proper error handling for token creation and validation
4. Rotate keys periodically for enhanced security
5. Use environment variables or secure vaults for storing sensitive keys
6. Keep internal service communication within trusted networks
7. Regularly audit JWT usage patterns and security measures
