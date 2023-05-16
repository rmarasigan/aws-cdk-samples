package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-ses/api/schema"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-ses/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-ses/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function will be triggered by the AWS Step Function State
// Machine to send an email to the endpoint.
func handler(ctx context.Context, event json.RawMessage) error {
	var (
		transaction = new(schema.Transaction)
		identity    = os.Getenv("EMAIL_IDENTITY")
		email       = new(awswrapper.EmailConfiguration)
	)

	// Check if the email identity is configured
	if identity == "" {
		err := errors.New("ses EMAIL_IDENTITY environment variable is not set")
		utility.Error(err, "EnvError", "SES EMAIL_IDENTITY is not configured on the environment")

		return err
	}
	email.Sender = identity

	// Unmarshal the received event
	err := json.Unmarshal([]byte(event), transaction)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data", utility.KVP{Key: "event", Value: event})
		return err
	}

	subject, body := transaction.EmailContent()
	email.Body = body
	email.Subject = subject
	email.Recipients = append(email.Recipients, transaction.Customer.Email)

	// Check if the email configuration parameters are valid
	err = email.IsValid()
	if err != nil {
		utility.Error(err, "EmailCfgError", "not valid email configuration", utility.KVP{Key: "event", Value: event}, utility.KVP{Key: "email", Value: email})
		return err
	}

	// Send an email
	err = awswrapper.SESSendSimpleEmail(ctx, email)
	if err != nil {
		utility.Error(err, "SESError", "failed to send email", utility.KVP{Key: "event", Value: event}, utility.KVP{Key: "email", Value: email})
		return err
	}

	return nil
}
