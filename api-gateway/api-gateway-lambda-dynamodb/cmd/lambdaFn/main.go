package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-dynamodb/api"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-dynamodb/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-dynamodb/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon API Gateway event record data as input,
// validates the request body, saves the processed request to the DynamoDB, and
// responds with a 200 OK HTTP Status.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		body      = request.Body
		coffee    = new(api.Coffee)
		tablename = os.Getenv("TABLE_NAME")
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		err := errors.New("dynamodb TABLE_NAME environment variable is not set")
		utility.Error(err, "EnvError", "DynamoDB TABLE_NAME is not configured on the environment")

		return api.InternalServerError()
	}

	// Unmarshal the request body
	err := api.UnmarshalJSON([]byte(body), coffee)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal JSON-encoded data")
		return api.BadRequest(err)
	}

	// Validate the incoming request body
	err = coffee.ValidateRequest()
	if err != nil {
		utility.Error(err, "APIError", "some field(s) is/are missing")
		return api.BadRequest(err)
	}

	// Insert the record into the DynamoDB Table
	err = awswrapper.DynamoPutItem(ctx, tablename, *coffee)
	if err != nil {
		utility.Error(err, "DynamoDBError", "failed to put item to the DynamoDB Table", utility.KVP{Key: "tablename", Value: tablename})
		return api.BadRequest(err)
	}

	return api.OKWithoutBody()
}
