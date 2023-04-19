package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-cloudwatch/internal/utility"
)

const (
	AWS_REGION = "us-east-1"
)

var (
	sfnClient *sfn.Client
)

// initStepFunctionsClient intializes the Step Function
// client from the provided configuration.
func initStepFunctionsClient(ctx context.Context) {
	if sfnClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "SFnClientError", "failed to load the default config")
		return
	}

	sfnClient = sfn.NewFromConfig(cfg)
}

// SFnStartExecution initializes the Step Function client and starts the state
// machine execution.
func SFnStartExecution(ctx context.Context, sfnARN, data string) error {
	// Initialize the Step Function Client
	initStepFunctionsClient(ctx)

	var input = &sfn.StartExecutionInput{
		StateMachineArn: aws.String(sfnARN),
		Input:           aws.String(data),
	}

	// Start the state machine execution
	_, err := sfnClient.StartExecution(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
