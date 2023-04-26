package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-lambda/api"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-lambda/api/models"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function will be triggered by the AWS Step Function State
// Machine and will write the event into the CloudWatch Logs.
func handler(ctx context.Context, event json.RawMessage) error {
	var (
		transaction = new(models.Transaction)
	)

	err := api.UnmarshalJSON(event, transaction)
	if err != nil {
		api.ErrorLog(err, "Failed to unmarshal the event message")
		return err
	}

	// Log the received input from the Step Function
	utility.Info("TRANSACTION_PROCESS", "Received event from the Step Function", utility.KVP{Key: "transaction", Value: transaction})

	return nil
}
