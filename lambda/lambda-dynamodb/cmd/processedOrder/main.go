package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the DynamoDB event data as input, reads
// the 'referenceId' attribute that is coming from an INSERT Event,
// and update the record to the DynamoDB Table.
func handler(ctx context.Context, event events.DynamoDBEvent) error {
	var (
		records   = event.Records
		order     = schema.Order{}
		tablename = os.Getenv("TABLE_NAME")
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		err := errors.New("DynamoDB TABLE_NAME environment variable is not set")
		utility.Error(err, "EnvError", "DynamoDB TABLE_NAME is not configured on the environment")

		return err
	}

	if len(records) == 0 {
		utility.Info("DynamoDBEvent", "There are no new events")
		return nil
	}

	for _, record := range records {
		if record.EventName == "INSERT" {
			// Read the value of attribute 'referenceId'
			referenceId := record.Change.NewImage["referenceId"].String()

			err := awswrapper.DynamoUpdateItem(ctx, tablename, referenceId, string(order.Status.Processed()))
			if err != nil {
				utility.Error(err, "DynamoDBError", "Failed to update item to the DynamoDB Table",
					utility.KVP{Key: "tablename", Value: tablename}, utility.KVP{Key: "referenceId", Value: referenceId})

				return err
			}
		}
	}

	return nil
}
