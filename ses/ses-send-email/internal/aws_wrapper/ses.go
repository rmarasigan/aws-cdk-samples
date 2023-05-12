package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/rmarasigan/aws-cdk-samples/ses/ses-send-email/internal/utility"
)

var sesClient *sesv2.Client

const AWS_REGION = "us-east-1"

// initSESClient initializes the SES Client from the provided configuration.
func initSESClient(ctx context.Context) {
	if sesClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "SESClientError", "failed to load the default config")
		return
	}

	sesClient = sesv2.NewFromConfig(cfg)
}

// SESSendSimpleEmail sends a simple email where we specify the sender, recipient,
// and the message body.
func SESSendSimpleEmail(ctx context.Context, email *EmailConfiguration) error {
	// Initialize the SESClient
	initSESClient(ctx)

	var input = &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(email.Sender),
		Content:          email.setSimpleTextContent(),
		Destination:      email.setSimpleDestination(),
	}

	_, err := sesClient.SendEmail(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
