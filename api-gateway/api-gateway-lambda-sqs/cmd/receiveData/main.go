package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-sqs/api"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-sqs/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-sqs/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon API Gateway event record data as input,
// validates the request body, send the message to an SQS Queue, and responds
// with a 200 OK HTTP Status.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		body  = request.Body
		item  = new(api.Item)
		queue = os.Getenv("QUEUE_URL")
	)

	// Check if the SQS Queue is configured
	if queue == "" {
		err := errors.New("QUEUE_URL is not set on the environment")
		utility.Error(err, "EnvError", "QUEUE_URL is not configured on the environment")

		return api.InternalServerError()
	}

	// Unmarshal the request body
	err := api.UnmarshalJSON([]byte(body), item)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data")
		return api.InternalServerError()
	}

	// Validate the incoming request
	err = item.ValidateRequest()
	if err != nil {
		utility.Error(err, "APIError", "Some field(s) is/are missing")
		return api.BadRequest(err)
	}

	data, err := api.MarshalJSON(item)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to marshal JSON-encoded data")
		return api.InternalServerError()
	}

	// Send the received request to the configured SQS
	err = awswrapper.SQSSendMessage(ctx, queue, string(data))
	if err != nil {
		utility.Error(err, "SQSError", "Failed to send data to the configured SQS", utility.KVP{Key: "queue", Value: queue})
		return api.InternalServerError()
	}

	return api.OKWithoutBody()
}
