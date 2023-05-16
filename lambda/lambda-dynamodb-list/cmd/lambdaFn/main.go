package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/common"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function once triggered, it will process the received
// order information, and either insert/update the record into
// the DynamoDB Table that is configured.
func handler(ctx context.Context, data json.RawMessage) error {
	var (
		order     = new(schema.Order)
		tablename = os.Getenv("TABLE_NAME")
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		err := errors.New("dynamodb TABLE_NAME environment variable is not set")
		utility.Error(err, "EnvError", "DynamoDB TABLE_NAME is not configured on the environment")

		return err
	}

	// Unmarshal the received JSON-encoded data
	err := json.Unmarshal([]byte(data), order)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data", utility.KVP{Key: "data", Value: data})
		return err
	}

	for _, line := range order.OrderLine {
		referenceId := order.ReferenceId

		// Check if the order line already exist
		existing, index, err := common.ExistingOrderLine(ctx, tablename, referenceId, line.ItemID)
		if err != nil {
			utility.Error(err, "DynamoDBError", "failed to fetch the item to the DynamoDB", utility.KVP{Key: "order", Value: order})
			return err
		}

		if !existing {
			// Append the non-existing order line
			order.OrderLine = append([]schema.OrderLine{}, line)

			// Insert the object to the DynamoDB table
			err = awswrapper.DynamoUpdateItem(ctx, tablename, *order)
			if err != nil {
				utility.Error(err, "DynamoDBError", "failed to insert the item to the DynamoDB", utility.KVP{Key: "order", Value: order})
				return err
			}

		} else {
			// Update an existing item object from the DynamoDB table
			err = awswrapper.DynamoUpdateObject(ctx, tablename, referenceId, index, line)
			if err != nil {
				utility.Error(err, "DynamoDBError", "failed to update the object from the DynamoDB", utility.KVP{Key: "order", Value: order})
				return err
			}
		}
	}

	return nil
}
