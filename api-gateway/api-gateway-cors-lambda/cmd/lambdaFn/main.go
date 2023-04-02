package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-cors-lambda/api"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-cors-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body := request.Body
	var coffee api.Coffee

	err := api.UnmarshalJSON([]byte(body), &coffee)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data")
		return api.BadRequest(err)
	}

	err = coffee.ValidateRequest()
	if err != nil {
		utility.Error(err, "APIError", "Some field(s) is/are missing")
		return api.BadRequest(err)
	}

	return api.OKWithoutBody()
}
