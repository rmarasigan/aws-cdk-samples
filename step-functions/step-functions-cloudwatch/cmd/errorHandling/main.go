package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-cloudwatch/internal/utility"

	awswrapper "github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-cloudwatch/internal/aws_wrapper"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the AWS CloudWatch Logs event record data as input,
// decode and decompressed the event, and writes the event information to the
// DynamoDB Table.
func handler(ctx context.Context, event awswrapper.CloudWatchEvent) error {
	var (
		STATE_MACHINE_ARN = os.Getenv("STATE_MACHINE_ARN")
	)

	// Check if the STATE_MACHINE_ARN is configured
	if STATE_MACHINE_ARN == "" {
		err := errors.New("step functions STATE_MACHINE_ARN environment is not set")
		utility.Error(err, "EnvError", "Step Functions STATE_MACHINE_ARN is missing")

		return err
	}

	// Decode and decompressed the received CloudWatch Event
	data, err := event.DecodeData()
	if err != nil {
		utility.Error(err, "CWError", "failed to decode and decompressed the received event from CloudWatch", utility.KVP{Key: "event", Value: event})
		return err
	}

	// Start the Step Function State Machine and send the input
	input := utility.EncodeJSON(data)
	err = awswrapper.SFnStartExecution(ctx, STATE_MACHINE_ARN, input)
	if err != nil {
		utility.Error(err, "SFnError", "failed to start the step function",
			utility.KVP{Key: "state_machine", Value: STATE_MACHINE_ARN}, utility.KVP{Key: "input", Value: input})

		return err
	}

	return nil
}
