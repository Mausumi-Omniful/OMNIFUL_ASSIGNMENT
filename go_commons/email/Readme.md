# Package email

## Overview
The email package provides a robust interface for sending emails using AWS Simple Email Service (SES). It supports both standard and multi-tenant email sending with features like attachments, HTML templates, and configurable recipients (To, CC, BCC).

## Key Components

### Interfaces
- `EmailClient`: Basic email sending interface
  - `SendEmail`: Send simple emails
  - `SendEmailWithAttachment`: Send emails with attachments

- `EmailClientV2`: Multi-tenant email sending interface with context support
  - `SendEmail`: Send emails with tenant context
  - `SendEmailWithAttachment`: Send emails with attachments and tenant context

- `EmailSetting`: Interface for retrieving email configurations
  - `GetEmailConfig`: Retrieve email configuration for a tenant

### Data Structures
- `Message`: Represents an email message
  ```go
  type Message struct {
      Subject      string
      Template     *template.Template
      TemplateData interface{}
  }
  ```

- `Recipient`: Defines email recipients
  ```go
  type Recipient struct {
      ToEmails  []string
      CcEmails  []string
      BccEmails []string
  }
  ```

- `Attachment`: Represents email attachments
  ```go
  type Attachment struct {
      Name string
      Data []byte
  }
  ```

- `EmailConfig`: Represents email configuration for a tenant
  ```go
  type EmailConfig struct {
      Workspace string
      FromEmail string
  }
  ```

## Usage Examples

### Basic Email Sending
```go
package main

import (
    "github.com/omniful/go_commons/email"
    "html/template"
)

func main() {
    // Initialize SES client with options
    sesClient, err := email.NewSesClientV2(
        email.WithRegion("us-west-2"),
        email.WithCredentials("your-access-key", "your-access-secret"),
    )
    if err != nil {
        panic(err)
    }

    // Create HTML template
    tmpl, err := template.New("email").Parse("<h1>Hello {{.Name}}!</h1>")
    if err != nil {
        panic(err)
    }

    // Create message
    message := email.Message{
        Subject:      "Welcome Email",
        Template:     tmpl,
        TemplateData: map[string]string{"Name": "John"},
    }

    // Define recipients
    recipient := email.Recipient{
        ToEmails:  []string{"user@example.com"},
        CcEmails:  []string{"cc@example.com"},
        BccEmails: []string{"bcc@example.com"},
    }

    // Send email
    err = sesClient.SendEmail(context.Background(), "tenant-123", message, recipient)
    if err != nil {
        panic(err)
    }
}
```

### Multi-tenant Email Sending with Attachments
```go
package main

import (
    "context"
    "github.com/omniful/go_commons/email"
    "html/template"
)

func main() {
    // Initialize SES client V2 with email settings
    sesClientV2, err := email.NewSesClientV2(
        email.WithRegion("us-west-2"),
        email.WithCredentials("your-access-key", "your-access-secret"),
        email.WithEmailSetting(emailSettingImpl),
    )
    if err != nil {
        panic(err)
    }

    // Create HTML template
    tmpl, err := template.New("email").Parse("<h1>Document Attached</h1>")
    if err != nil {
        panic(err)
    }

    // Create message
    message := email.Message{
        Subject:      "Document Delivery",
        Template:     tmpl,
        TemplateData: nil,
    }

    // Create attachment
    attachments := []email.Attachment{
        {
            Name: "document.pdf",
            Data: []byte("PDF content here"),
        },
    }

    // Define recipients
    recipient := email.Recipient{
        ToEmails: []string{"user@example.com"},
    }

    // Send email with attachment
    ctx := context.Background()
    err = sesClientV2.SendEmailWithAttachment(ctx, "tenant-123", message, attachments, recipient)
    if err != nil {
        panic(err)
    }
}
```

## Features
- AWS SES Integration
- HTML Template Support
- Multiple Recipients (To, CC, BCC)
- File Attachments
- Multi-tenant Support
- XSS Protection (blocks script tags in content)
- Configurable Email Settings per Tenant
- Flexible Client Configuration with Options Pattern

## Notes
- The package requires AWS credentials to be properly configured
- HTML templates are sanitized for security (script tags are blocked)
- Supports both single-tenant and multi-tenant architectures
- Uses `gomail.v2` for handling attachments
- Client initialization supports various configuration options:
  - `WithRegion`: Set AWS region
  - `WithCredentials`: Set AWS credentials
  - `WithEmailSetting`: Set email settings for multi-tenant support
