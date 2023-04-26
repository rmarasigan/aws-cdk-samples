package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function once triggered, will create a new order
// information, and  insert the record into the DynamoDB Table.
func handler(ctx context.Context) error {
	var (
		order     = new(schema.Order)
		tablename = os.Getenv("TABLE_NAME")
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		err := errors.New("DynamoDB TABLE_NAME environment variable is not set")
		utility.Error(err, "EnvError", "DynamoDB TABLE_NAME is not configured on the environment")

		return err
	}

	// Insert the record into the DynamoDB table
	err := awswrapper.DynamoPutItem(ctx, tablename, order.CreateOrder())
	if err != nil {
		utility.Error(err, "DynamoDBError", "Failed to put item to the DynamoDB Table", utility.KVP{Key: "tablename", Value: tablename})
		return err
	}

	return nil
}
