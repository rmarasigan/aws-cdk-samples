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

// handler function receives an event record data as input, validates the data, and sends
// the data as input via Step Function that will trigger the "sendEmail" Lambda Function.
func handler(ctx context.Context, data json.RawMessage) error {
	var (
		transaction     = new(schema.Transaction)
		stateMachineARN = os.Getenv("STATE_MACHINE_ARN")
	)

	// Check if the STATE_MACHINE_ARN is configured
	if stateMachineARN == "" {
		err := errors.New("step functions SATE_MACHINE_ARN environment variable is not set")
		utility.Error(err, "EnvError", "Step Functions STATE_MACHINE_ARN is not configured on the environment")

		return err
	}

	// Unmarshal the request
	err := json.Unmarshal([]byte(data), transaction)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal JSON-encoded data", utility.KVP{Key: "data", Value: data})
		return err
	}

	// Set the default values and validate
	transaction.SetDefaultValues()
	err = transaction.IsValid()
	if err != nil {
		utility.Error(err, "TransactionError", "transaction details is/are not valid", utility.KVP{Key: "data", Value: data})
		return err
	}

	input, err := json.Marshal(transaction)
	if err != nil {
		utility.Error(err, "JSONError", "failed to marshal the request", utility.KVP{Key: "data", Value: data})
		return err
	}

	// Start the Step Function State Macine and send the data as the input
	err = awswrapper.SFnStartExecution(ctx, stateMachineARN, string(input))
	if err != nil {
		utility.Error(err, "SFnError", "failed to start the step function",
			utility.KVP{Key: "state_machine", Value: stateMachineARN}, utility.KVP{Key: "data", Value: data})

		return err
	}

	return nil
}
