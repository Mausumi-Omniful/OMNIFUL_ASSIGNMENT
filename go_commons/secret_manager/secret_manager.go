package secret_manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/omniful/go_commons/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/omniful/go_commons/log"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	defaultTTL = 24 * time.Hour
	minTTL     = 24 * time.Hour
)

// Common errors
var (
	errEmptyRegion     = errors.New("region cannot be empty")
	errEmptySecretName = errors.New("secret name cannot be empty")
	errEmptyFieldName  = errors.New("field name cannot be empty")
	errFieldNotFound   = errors.New("field not found in secret")
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// Validate checks if the credentials are valid
func (c *Credentials) Validate() error {
	if len(c.AccessKeyID) == 0 {
		return errors.New("access key ID cannot be empty")
	}
	if len(c.SecretAccessKey) == 0 {
		return errors.New("secret access key cannot be empty")
	}
	return nil
}

// Config holds the configuration for SecretManager
type Config struct {
	region      string
	ttl         time.Duration
	credentials *Credentials
}

type Option func(*Config)

func WithRegion(region string) Option {
	return func(c *Config) {
		c.region = region
	}
}

// WithTTL sets the cache TTL
func WithTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.ttl = ttl
	}
}

func WithCredentials(creds *Credentials) Option {
	return func(c *Config) {
		c.credentials = creds
	}
}

// Validate checks if the Config is valid
func (c *Config) Validate() error {
	if len(c.region) == 0 {
		return errEmptyRegion
	}
	if c.credentials != nil {
		if err := c.credentials.Validate(); err != nil {
			return fmt.Errorf("invalid credentials: %w", err)
		}
	}
	return nil
}

// NewConfig creates a new Config with the provided options
// Example:
//
//	cfg := NewConfig(
//	    WithRegion("eu-central-1"),
//	    WithTTL(24 * time.Hour),
//	    WithCredentials(&Credentials{
//	        AccessKeyID:     "your-access-key",
//	        SecretAccessKey: "your-secret-key",
//	    }),
//	)
func NewConfig(ctx context.Context, opts ...Option) *Config {
	cfg := &Config{
		region: config.GetString(ctx, "cloudProvider.region"),
		ttl:    defaultTTL,
	}

	// Apply all options
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

type cachedSecret struct {
	data      map[string]interface{}
	expiresAt time.Time
}

// SecretFetcher defines the interface for secret management operations
type SecretFetcher interface {
	// GetSecret retrieves a secret from cache or Secrets Manager
	// Returns an error if the secret name is empty or if the secret cannot be retrieved
	GetSecret(ctx context.Context, secretName string) (map[string]interface{}, error)

	// GetSecretFieldValue retrieves a specific field from a secret
	// Returns an error if the secret name or field name is empty, or if the field is not found
	GetSecretFieldValue(ctx context.Context, secretName, field string) (interface{}, error)
}

// SecretManager implements the SecretFetcher interface
type SecretManager struct {
	*secretsmanager.SecretsManager
	cache sync.Map
	ttl   time.Duration
}

// Ensure SecretManager implements SecretFetcher interface
// this is a compile-time check
var _ SecretFetcher = (*SecretManager)(nil)

// NewSecretManager creates a new SecretManager instance with the provided options
// If no options are provided, default configuration will be used.
// The default configuration uses defaultTTL (24 hours) for cache TTL.
//
// Example:
//
//	sm, err := NewSecretManager(ctx,
//	    WithRegion("eu-central-1"),
//	    WithTTL(24 * time.Hour),
//	    WithCredentials(&Credentials{
//	        AccessKeyID:     "your-access-key",
//	        SecretAccessKey: "your-secret-key",
//	    }),
//	)
func NewSecretManager(ctx context.Context, opts ...Option) (SecretFetcher, error) {
	cfg := NewConfig(ctx, opts...)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	awsConfig := &aws.Config{
		Region: aws.String(cfg.region),
	}

	// Add credentials if provided
	if cfg.credentials != nil {
		awsConfig.Credentials = credentials.NewStaticCredentials(
			cfg.credentials.AccessKeyID,
			cfg.credentials.SecretAccessKey,
			cfg.credentials.SessionToken,
		)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		log.Errorf("failed to create session, %v", err)
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Validate TTL
	ttl := cfg.ttl
	if ttl < minTTL {
		log.Warnf("Configured TTL (%v) is less than minimum allowed (%v), using minimum TTL", ttl, minTTL)
		ttl = minTTL
	} else if ttl == 0 {
		log.Infof("No TTL configured, using default TTL of %v", defaultTTL)
		ttl = defaultTTL
	}

	return &SecretManager{
		SecretsManager: secretsmanager.New(sess),
		ttl:            ttl,
	}, nil
}

// GetSecret retrieves a secret from cache or Secrets Manager and parses it into a map.
func (sm *SecretManager) GetSecret(
	ctx context.Context,
	secretName string,
) (map[string]interface{}, error) {
	if len(secretName) == 0 {
		return nil, errEmptySecretName
	}

	if data, found := sm.getFromCache(secretName); found {
		return data, nil
	}

	secretString, err := sm.fetchSecretString(ctx, secretName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secret %s: %w", secretName, err)
	}

	var parsed map[string]interface{}
	if err = json.Unmarshal([]byte(secretString), &parsed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret %s into map: %w", secretName, err)
	}

	sm.cache.Store(secretName, cachedSecret{
		data:      parsed,
		expiresAt: time.Now().Add(sm.ttl),
	})

	return parsed, nil
}

// GetSecretFieldValue returns the raw value of a field in a secret.
func (sm *SecretManager) GetSecretFieldValue(
	ctx context.Context,
	secretName, field string,
) (interface{}, error) {
	if len(secretName) == 0 {
		return nil, errEmptySecretName
	}
	if len(field) == 0 {
		return nil, errEmptyFieldName
	}

	secret, err := sm.GetSecret(ctx, secretName)
	if err != nil {
		return nil, err
	}

	value, ok := secret[field]
	if !ok {
		return nil, fmt.Errorf("%w: %s in secret %s", errFieldNotFound, field, secretName)
	}

	return value, nil
}

// fetchSecretString retrieves the raw secret string from Secrets Manager.
func (sm *SecretManager) fetchSecretString(ctx context.Context, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := sm.SecretsManager.GetSecretValueWithContext(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret %s: %w", secretName, err)
	}

	if result.SecretString != nil {
		return *result.SecretString, nil
	}

	return "", fmt.Errorf("secret %s has no string value", secretName)
}

// getFromCache returns a cached secret if it exists and hasn't expired.
func (sm *SecretManager) getFromCache(secretName string) (map[string]interface{}, bool) {
	if cached, found := sm.cache.Load(secretName); found {
		if cs, ok := cached.(cachedSecret); ok {
			if time.Now().Before(cs.expiresAt) {
				return cs.data, true
			}

			// Expired, clean up
			sm.cache.Delete(secretName)
		}
	}
	return nil, false
}
