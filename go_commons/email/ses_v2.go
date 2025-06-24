package email

import (
	"context"
	"errors"
)

type SesClientV2 struct {
	sesClient    EmailClient
	emailSetting EmailSetting
}

func NewSesClientV2(opts ...SesClientOption) (EmailClientV2, error) {
	// Default configuration
	cfg := &sesClientConfig{}

	// Apply all options
	for _, opt := range opts {
		opt(cfg)
	}

	// Create the base SES client with credentials
	c, err := NewSesClient(cfg.region, cfg.accessKey, cfg.accessSecret)
	if err != nil {
		return nil, err
	}

	if cfg.emailSetting == nil {
		return nil, errors.New("email setting implementation is required")
	}

	return &SesClientV2{
		sesClient:    c,
		emailSetting: cfg.emailSetting,
	}, nil
}

func (s *SesClientV2) SendEmail(ctx context.Context, tenantID string, message Message, recipient Recipient) (err error) {
	emailConfig, interserviceErr := s.emailSetting.GetEmailConfig(ctx, tenantID)
	if interserviceErr != nil {
		err = errors.New(interserviceErr.Message)
		return
	}

	err = s.sesClient.SendEmail(emailConfig.FromEmail, message, recipient)
	if err != nil {
		return
	}

	return
}

func (s *SesClientV2) SendEmailWithAttachment(ctx context.Context, tenantID string, message Message, attachments []Attachment, recipient Recipient) (err error) {
	emailConfig, interserviceErr := s.emailSetting.GetEmailConfig(ctx, tenantID)
	if interserviceErr != nil {
		err = errors.New(interserviceErr.Message)
		return
	}

	err = s.sesClient.SendEmailWithAttachment(emailConfig.FromEmail, message, attachments, recipient)
	if err != nil {
		return
	}

	return
}
