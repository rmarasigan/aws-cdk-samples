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

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		coffee    = new(api.Coffee)
		body      = request.Body
		tablename = os.Getenv("TABLE_NAME")
	)

	if tablename == "" {
		err := errors.New("DynamoDB TABLE_NAME environment variable is not set")
		utility.Error(err, "EnvError", "DynamoDB TABLE_NAME is not configured on the environment")

		return api.BadRequest(err)
	}

	// Unmarshal the request body
	err := api.UnmarshalJSON([]byte(body), coffee)
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

	// Insert the record into the DynamoDB Table
	err = awswrapper.DynamoPutItem(ctx, tablename, *coffee)
	if err != nil {
		utility.Error(err, "DynamoDBError", "Failed to put item to the DynamoDB Table", utility.KVP{Key: "tablename", Value: tablename})
		return api.BadRequest(err)
	}

	return api.OKWithoutBody()
}
