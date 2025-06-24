package email

// SesClientOption is a function that configures the SES client
type SesClientOption func(*sesClientConfig)

// sesClientConfig holds all configuration options for the SES client
type sesClientConfig struct {
	region       string
	accessKey    string
	accessSecret string
	emailSetting EmailSetting
}

// WithRegion sets the AWS region for the SES client
func WithRegion(region string) SesClientOption {
	return func(cfg *sesClientConfig) {
		cfg.region = region
	}
}

// WithCredentials sets the AWS credentials for the SES client
func WithCredentials(accessKey, accessSecret string) SesClientOption {
	return func(cfg *sesClientConfig) {
		cfg.accessKey = accessKey
		cfg.accessSecret = accessSecret
	}
}

// WithEmailSetting sets the email settings for the SES client
func WithEmailSetting(setting EmailSetting) SesClientOption {
	return func(cfg *sesClientConfig) {
		cfg.emailSetting = setting
	}
}
