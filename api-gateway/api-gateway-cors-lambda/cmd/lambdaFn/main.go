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

// handler function receives the Amazon API Gateway event record data as input,
// validates the request body, and responds with a 200 OK HTTP Status. The API response
// must include the "Access-Control-Allow-Origin" and "Access-Control-Allow-Methods" headers.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		coffee api.Coffee
		body   = request.Body
	)

	// Unmarshal the JSON-encoded request body
	err := api.UnmarshalJSON([]byte(body), &coffee)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data")
		return api.BadRequest(err)
	}

	// Validate the incoming request body
	err = coffee.ValidateRequest()
	if err != nil {
		utility.Error(err, "APIError", "Some field(s) is/are missing")
		return api.BadRequest(err)
	}

	return api.OKWithoutBody()
}
