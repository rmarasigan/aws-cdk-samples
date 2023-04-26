package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-lambda/api"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-lambda/api/models"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-lambda/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon API Gateway event record data as input,
// validates the request body, and sends the data that will trigger the second
// Lambda function.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		body              = request.Body
		transaction       = new(models.Transaction)
		STATE_MACHINE_ARN = os.Getenv("STATE_MACHINE_ARN")
	)

	// Check if the STATE_MACHINE_ARN is configured
	if STATE_MACHINE_ARN == "" {
		err := errors.New("STATE_MACHINE_ARN is not set on the environment")
		utility.Error(err, "EnvError", "STATE_MACHINE_ARN is missing")

		return api.InternalServerError()
	}

	// Unmarshal the request
	err := api.UnmarshalJSON([]byte(body), transaction)
	if err != nil {
		api.ErrorLog(err, "Failed to unmarshal the JSON-encoded data", utility.KVP{Key: "body", Value: body})
		return api.InternalServerError()
	}

	// Validate the incoming request
	err = transaction.Validate()
	if err != nil {
		api.ErrorLog(err, "Some fields in the request body were missing/incorrect", utility.KVP{Key: "body", Value: body})
		return api.BadRequest(err)
	}
	transaction.SetDefaultTransactionValues()

	data, err := api.MarshalJSON(transaction)
	if err != nil {
		api.ErrorLog(err, "Failed to marshal the request", utility.KVP{Key: "transaction", Value: transaction})
		return api.InternalServerError()
	}

	// Start the Step Function State Machine and send the data as the input
	err = awswrapper.SFnStartExecution(ctx, STATE_MACHINE_ARN, string(data))
	if err != nil {
		utility.Error(err, "SFnError", "Failed to start the step function",
			utility.KVP{Key: "state_machine", Value: STATE_MACHINE_ARN}, utility.KVP{Key: "data", Value: string(data)})

		return api.InternalServerError()
	}

	return api.OKWithoutBody()
}
