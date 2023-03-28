package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-async-lambda/internal/trail"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-async-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) error {
	// Check if the request body is not empty
	if request.Body != "{}" {
		utility.OK("RestAPI", "Received API Gateway Request Body", utility.KVP{Key: "Body", Value: request.Body})
		return nil
	}

	trail.Info("Empty API Gateway Request Body")
	return nil
}
