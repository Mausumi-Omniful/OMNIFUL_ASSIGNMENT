package email

import (
	"bytes"
	"errors"
	"html"
	"html/template"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/util"
	"gopkg.in/gomail.v2"
)

type SesClient struct {
	svc *ses.SES
}

func NewSesClient(region, accessKey, accessSecret string) (EmailClient, error) {
	if region == "" {
		return nil, errors.New("region is required")
	}

	config := &aws.Config{
		Region: aws.String(region),
	}

	// Only set credentials if they are provided
	if accessKey != "" && accessSecret != "" {
		config.Credentials = credentials.NewStaticCredentials(accessKey, accessSecret, "")
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	svc := ses.New(sess)
	return &SesClient{
		svc: svc,
	}, nil
}

type Recipient struct {
	ToEmails  []string
	CcEmails  []string
	BccEmails []string
}

type Attachment struct {
	Name string
	Data []byte
}

type Message struct {
	Subject      string
	Template     *template.Template
	TemplateData interface{}
}

// SendEmail sends email to specified email IDs
func (c *SesClient) SendEmail(fromEmail string, message Message, recipient Recipient) (err error) {
	// set to section
	var recipients []*string
	for _, r := range recipient.ToEmails {
		recipients = append(recipients, aws.String(r))
	}

	// set cc section
	var ccRecipients []*string
	if len(recipient.CcEmails) > 0 {
		for _, r := range recipient.CcEmails {
			ccRecipients = append(ccRecipients, aws.String(r))
		}
	}

	// set bcc section
	var bccRecipients []*string
	if len(recipient.BccEmails) > 0 {
		for _, r := range recipient.BccEmails {
			bccRecipients = append(bccRecipients, aws.String(r))
		}
	}

	buf := new(bytes.Buffer)
	if err = message.Template.Execute(buf, message.TemplateData); err != nil {
		return err
	}

	unescapedString := html.UnescapeString(buf.String())
	if util.ContainsScriptTag(unescapedString) {
		return errors.New("content contains script tag")
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses:  ccRecipients,
			ToAddresses:  recipients,
			BccAddresses: bccRecipients,
		},

		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(unescapedString),
				},
			},

			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(message.Subject),
			},
		},
		Source: aws.String(fromEmail),
	}

	_, err = c.svc.SendEmail(input)
	if err != nil {
		return
	}

	log.Infof("Email sent successfully to: ", recipient.ToEmails)
	return
}

// SendEmailWithAttachment sends email with attachment
func (c *SesClient) SendEmailWithAttachment(fromEmail string, message Message, attachments []Attachment, recipient Recipient) (err error) {
	msg := gomail.NewMessage()

	var recipients []*string
	for _, r := range recipient.ToEmails {
		recipients = append(recipients, aws.String(r))
	}

	var ccRecipients []*string
	if len(recipient.CcEmails) > 0 {
		for _, r := range recipient.CcEmails {
			ccRecipients = append(ccRecipients, aws.String(r))
		}
		msg.SetHeader("cc", recipient.CcEmails...)
	}

	var bccRecipients []*string
	if len(recipient.BccEmails) > 0 {
		for _, r := range recipient.BccEmails {
			bccRecipients = append(bccRecipients, aws.String(r))
		}
		msg.SetHeader("bcc", recipient.BccEmails...)
	}

	buf := new(bytes.Buffer)
	if err = message.Template.Execute(buf, message.TemplateData); err != nil {
		return err
	}

	unescapedString := html.UnescapeString(buf.String())
	if util.ContainsScriptTag(unescapedString) {
		return errors.New("content contains script tag")
	}

	msg.SetHeader("To", recipient.ToEmails...)
	msg.SetHeader("Subject", message.Subject)
	msg.SetBody("text/html", unescapedString)

	for _, v := range attachments {
		msg.Attach(v.Name, gomail.SetCopyFunc(func(w io.Writer) error {
			_, cusErr := w.Write(v.Data)
			return cusErr
		}))
	}

	// create a new buffer to add raw data
	var emailRaw bytes.Buffer
	_, err = msg.WriteTo(&emailRaw)
	if err != nil {
		return
	}

	m := ses.RawMessage{Data: emailRaw.Bytes()}
	input := &ses.SendRawEmailInput{Source: &fromEmail, Destinations: recipients, RawMessage: &m}

	_, err = c.svc.SendRawEmail(input)
	if err != nil {
		log.Errorf("Error sending mail - ", err)
		return
	}

	log.Infof("Email sent successfully to: ", recipient.ToEmails)
	return
}
