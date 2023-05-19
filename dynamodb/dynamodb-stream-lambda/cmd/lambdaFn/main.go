package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/dynamodb/dynamodb-stream-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives a DynamoDB event data as input and writes some
// of the record data to CloudWatch logs.
func handler(ctx context.Context, event events.DynamoDBEvent) error {
	var records = event.Records

	if len(records) == 0 {
		utility.Info("DynamoDBStream", "there are no records found")
		return nil
	}

	for _, record := range records {
		utility.Info("DynamoDBStreamRecord", "Events from the DynamoDB Table", utility.KVP{Key: "record", Value: record})
	}

	return nil
}
