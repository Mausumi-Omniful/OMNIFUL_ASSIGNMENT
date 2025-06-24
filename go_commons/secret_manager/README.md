# Secret Manager

A Go package for managing secrets with caching support. This package provides a simple interface to interact with AWS Secrets Manager while maintaining a local cache for improved performance.

## Features

- Caching with configurable TTL (minimum 24 hours)
- AWS Secrets Manager integration
- Configurable AWS credentials
- Thread-safe operations
- Functional options pattern for configuration
- Interface-based design for better testing and extensibility
- Automatic cache cleanup

## Installation

```bash
go get github.com/omniful/go_commons/secret_manager
```

## Usage

### Basic Usage

```go
import "github.com/omniful/go_commons/secret_manager"

// Create a new secret manager with default configuration
// Note: sm is of type SecretFetcher interface
var sm secret_manager.SecretFetcher
sm, err := secret_manager.NewSecretManager(ctx)
if err != nil {
    log.Fatal(err)
}

// Get a secret
secret, err := sm.GetSecret(ctx, "my-secret")
if err != nil {
    log.Fatal(err)
}

// Get a specific field from a secret
value, err := sm.GetSecretFieldValue(ctx, "my-secret", "api-key")
if err != nil {
    log.Fatal(err)
}
```

### Interface Design

The package uses an interface-based design for better testing and extensibility:

```go
// SecretFetcher defines the interface for secret management operations
type SecretFetcher interface {
    // GetSecret retrieves a secret from cache or Secrets Manager
    GetSecret(ctx context.Context, secretName string) (map[string]interface{}, error)

    // GetSecretFieldValue retrieves a specific field from a secret
    GetSecretFieldValue(ctx context.Context, secretName, field string) (interface{}, error)
}
```

This design allows you to:
- Mock the secret manager for testing
- Create alternative implementations if needed
- Use dependency injection in your applications

### Advanced Configuration

The secret manager can be configured using functional options:

```go
// Create a secret manager with custom configuration
var sm secret_manager.SecretFetcher
sm, err := secret_manager.NewSecretManager(ctx,
    secret_manager.WithRegion("eu-central-1"),
    secret_manager.WithTTL(48 * time.Hour),
    secret_manager.WithCredentials(&secret_manager.Credentials{
        AccessKeyID:     "your-access-key",
        SecretAccessKey: "your-secret-key",
        SessionToken:    "your-session-token", // Optional
    }),
)
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithRegion` | Sets the AWS region | From config |
| `WithTTL` | Sets the cache TTL | 24 hours |
| `WithCredentials` | Sets AWS credentials | Default AWS credential chain |

### Error Handling

The package defines several error types for common scenarios:

```go
// Check for specific errors
if errors.Is(err, secret_manager.ErrEmptySecretName) {
    // Handle empty secret name
}
if errors.Is(err, secret_manager.ErrFieldNotFound) {
    // Handle field not found
}
```

### Caching

The secret manager maintains a local cache of secrets with the following characteristics:

- Default TTL: 24 hours
- Minimum TTL: 24 hours
- Thread-safe operations using sync.Map
- Automatic cleanup of expired secrets

### Testing

The interface-based design makes it easy to mock the secret manager for testing:

```go
// Mock implementation for testing
type MockSecretFetcher struct {
    secrets map[string]map[string]interface{}
}

func (m *MockSecretFetcher) GetSecret(ctx context.Context, secretName string) (map[string]interface{}, error) {
    if secret, ok := m.secrets[secretName]; ok {
        return secret, nil
    }
    return nil, fmt.Errorf("secret not found: %s", secretName)
}

func (m *MockSecretFetcher) GetSecretFieldValue(ctx context.Context, secretName, field string) (interface{}, error) {
    secret, err := m.GetSecret(ctx, secretName)
    if err != nil {
        return nil, err
    }
    if value, ok := secret[field]; ok {
        return value, nil
    }
    return nil, fmt.Errorf("field not found: %s", field)
}

// Usage in tests
func TestMyFunction(t *testing.T) {
    mockSM := &MockSecretFetcher{
        secrets: map[string]map[string]interface{}{
            "test-secret": {
                "api-key": "test-key",
            },
        },
    }
    
    // Use mockSM in your tests
    result, err := myFunction(mockSM)
    // ... assertions
}
```

### Best Practices

1. **Error Handling**: Always check for errors when retrieving secrets
   ```go
   secret, err := sm.GetSecret(ctx, "my-secret")
   if err != nil {
       // Handle error appropriately
   }
   ```

2. **Context Usage**: Always pass a context with timeout
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   secret, err := sm.GetSecret(ctx, "my-secret")
   ```

3. **Credentials Management**: Use environment variables or AWS IAM roles when possible
   ```go
   // Use default credential chain
   sm, err := secret_manager.NewSecretManager(ctx,
       secret_manager.WithRegion("eu-central-1"),
   )
   ```

4. **Secret Naming**: Use consistent naming conventions for secrets
   ```go
   // Good
   secret, err := sm.GetSecret(ctx, "prod/api-keys/service-x")
   
   // Avoid
   secret, err := sm.GetSecret(ctx, "service-x-prod-key")
   ```

5. **Interface Usage**: Use the SecretFetcher interface type in your code
   ```go
   // Good
   func ProcessSecret(sm secret_manager.SecretFetcher) error {
       // ...
   }
   
   // Avoid
   func ProcessSecret(sm *secret_manager.SecretManager) error {
       // ...
   }
   ```

## Examples

### Basic Secret Retrieval

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/omniful/go_commons/secret_manager"
)

func main() {
    ctx := context.Background()
    
    // Create secret manager (using interface type)
    var sm secret_manager.SecretFetcher
    sm, err := secret_manager.NewSecretManager(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get a secret
    secret, err := sm.GetSecret(ctx, "my-secret")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Custom Configuration

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/omniful/go_commons/secret_manager"
)

func main() {
    ctx := context.Background()
    
    // Create secret manager with custom configuration
    var sm secret_manager.SecretFetcher
    sm, err := secret_manager.NewSecretManager(ctx,
        secret_manager.WithRegion("eu-central-1"),
        secret_manager.WithTTL(48 * time.Hour),
        secret_manager.WithCredentials(&secret_manager.Credentials{
            AccessKeyID:     "your-access-key",
            SecretAccessKey: "your-secret-key",
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Get a specific field
    value, err := sm.GetSecretFieldValue(ctx, "my-secret", "api-key")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Value: %v", value)
}
```
