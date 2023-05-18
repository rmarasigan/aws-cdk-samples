package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-cognito-lambda/api"
)

type Message struct {
	Content any `json:"content"`
}

func main() {
	lambda.Start(handler)
}

// handler function once triggered by an Amazon API Gateway event, will return an HTTP OK status and
// a custom message body.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	message := new(Message)
	message.Content = "Hello from the Lambda function!"

	return api.OK(message)
}
