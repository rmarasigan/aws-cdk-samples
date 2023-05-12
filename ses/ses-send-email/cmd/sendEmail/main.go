package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/ses/ses-send-email/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/ses/ses-send-email/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, data json.RawMessage) error {
	var (
		identity = os.Getenv("EMAIL_IDENTITY")
		email    = new(awswrapper.EmailConfiguration)
	)

	// Check if the email identity is configured
	if identity == "" {
		err := errors.New("ses EMAIL_IDENTITY environment variable is not set")
		utility.Error(err, "EnvError", "SES EMAIL_IDENTITY is not configured on the environment")

		return err
	}
	email.Sender = identity

	// Unmarshal the received data
	err := json.Unmarshal([]byte(data), email)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data", utility.KVP{Key: "data", Value: data})
		return err
	}

	// Check if the email configuration parameters are valid
	err = email.IsValid()
	if err != nil {
		utility.Error(err, "EmailCfgError", "not valid email configuration", utility.KVP{Key: "data", Value: data})
		return err
	}

	// Send an email
	err = awswrapper.SESSendSimpleEmail(ctx, email)
	if err != nil {
		utility.Error(err, "SESError", "failed to send email", utility.KVP{Key: "data", Value: data})
		return err
	}

	return nil
}
