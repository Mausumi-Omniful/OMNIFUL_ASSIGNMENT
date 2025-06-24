package email

import (
	"context"
	interservice_client "github.com/omniful/go_commons/interservice-client"
)

type EmailClient interface {
	SendEmail(fromEmail string, message Message, recipient Recipient) (err error)
	SendEmailWithAttachment(fromEmail string, message Message, attachments []Attachment, recipient Recipient) (err error)
}

type EmailClientV2 interface {
	SendEmail(ctx context.Context, tenantID string, message Message, recipient Recipient) (err error)
	SendEmailWithAttachment(ctx context.Context, tenantID string, message Message, attachments []Attachment, recipient Recipient) (err error)
}

type EmailSetting interface {
	GetEmailConfig(ctx context.Context, tenantID string) (res EmailConfig, err *interservice_client.Error)
}
