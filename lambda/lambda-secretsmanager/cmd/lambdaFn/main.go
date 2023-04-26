package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-secretsmanager/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-secretsmanager/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-secretsmanager/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function once triggered, will either update the
// Secret Value or retrieve it.
func handler(ctx context.Context, data json.RawMessage) error {
	var (
		event     = new(schema.Event)
		secretARN = os.Getenv("SECRET_ARN")
	)

	// Check if the Secrets Manager ARN is configured
	if secretARN == "" {
		err := errors.New("secrets manager SECRET_ARN environment is not set")
		utility.Error(err, "EnvError", "Secrets Manager SECRET_ARN is not configured on the environment")

		return err
	}

	// Unmarshal the received JSON event
	err := json.Unmarshal([]byte(data), event)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data")
		return err
	}

	switch event.Action {
	case "get":
		// Get the secret value and print the output
		output, err := awswrapper.SMGetSecretValue(ctx, secretARN)
		if err != nil {
			utility.Error(err, "SecretsManagerError", "Failed to retrieve the secret value", utility.KVP{Key: "Secret", Value: event.Secret})
			return err
		}

		utility.Info("SMGetSecretValue", "Secret Manager Information", utility.KVP{Key: "CreatedDate", Value: &output.CreatedDate},
			utility.KVP{Key: "Name", Value: &output.Name}, utility.KVP{Key: "SecretString", Value: &output.SecretString})
		return nil

	default:
		// Update the secret value using the information coming
		// from the JSON-encoded data
		if event.Action == "update" {
			secretInfo, err := event.Secret.Marshal()
			if err != nil {
				utility.Error(err, "JSONError", "Failed to marshal the secret", utility.KVP{Key: "Secret", Value: event.Secret})
				return err
			}

			output, err := awswrapper.SMPutSecretValue(ctx, secretARN, string(secretInfo))
			if err != nil {
				utility.Error(err, "SecretsManagerError", "Failed to put/update the secret value", utility.KVP{Key: "Secret", Value: event.Secret})
				return err
			}

			utility.Info("SMPutSecretValue", "Finished updating the secret", utility.KVP{Key: "Name", Value: &output.Name}, utility.KVP{Key: "VersionId", Value: &output.VersionId})
			return nil

		} else {
			err := errors.New("incorrect type of 'action'")
			utility.Error(err, "SecretsManagerError", "Wrong type of 'action'", utility.KVP{Key: "Action", Value: event.Action})

			return err
		}
	}
}
