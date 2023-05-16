package awswrapper

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

const (
	UTF8       = "UTF-8"
	ISO_8859_1 = "ISO-8859-1"
)

// EmailConfiguration holds the SES Email
// Configuration for the SendEmailInput.
type EmailConfiguration struct {
	Subject    string `json:"subject"`
	Sender     string
	Body       string   `json:"body"`
	Recipients []string `json:"recipients"`
	BCCAddress []string `json:"bcc_recipients,omitempty"`
	CCAddress  []string `json:"cc_recipients,omitempty"`
}

// setContent represents the content of the email and if the
// charset option is set to empty, it will use the UTF8 as its
// default charset.
func setContent(data, charset string) *types.Content {
	content := &types.Content{}
	content.Data = aws.String(data)

	if charset != "" {
		content.Charset = aws.String(charset)
	} else {
		content.Charset = aws.String(UTF8)
	}

	return content
}

// setSimpleTextContent contains the body of the message and the subject.
func (cfg *EmailConfiguration) setSimpleTextContent() *types.EmailContent {
	return &types.EmailContent{
		Simple: &types.Message{
			Subject: setContent(cfg.Subject, UTF8),
			Body: &types.Body{
				Text: setContent(cfg.Body, UTF8),
			},
		},
	}
}

// setSimpleDestination contains the recipients of the email message.
func (cfg *EmailConfiguration) setSimpleDestination() *types.Destination {
	destination := &types.Destination{}

	if len(cfg.BCCAddress) > 0 {
		destination.BccAddresses = cfg.BCCAddress
	}

	if len(cfg.CCAddress) > 0 {
		destination.CcAddresses = cfg.CCAddress
	}

	destination.ToAddresses = cfg.Recipients

	return destination
}

// IsValid checks if all the required fields are present.
func (cfg *EmailConfiguration) IsValid() error {
	var params []string

	if cfg.Subject == "" {
		params = append(params, "subject")
	}

	if cfg.Body == "" {
		params = append(params, "body")
	}

	if len(cfg.Recipients) == 0 {
		params = append(params, "recipients")
	}

	if len(params) > 0 {
		return fmt.Errorf("required parameter(s) is/are missing: %s", strings.Join(params, ", "))
	}

	return nil
}
