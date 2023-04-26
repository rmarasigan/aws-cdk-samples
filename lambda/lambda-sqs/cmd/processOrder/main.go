package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sqs/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sqs/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon SQS event message data as input
// and will write some of the messages to CloudWatch Logs.
func handler(ctx context.Context, event events.SQSEvent) error {
	var (
		order   = new(schema.Order)
		records = event.Records
	)

	if len(records) == 0 {
		utility.Info("SQS", "No records found")
		return nil
	}

	for _, record := range records {
		// Unmarshal the event message
		err := json.Unmarshal([]byte(record.Body), order)
		if err != nil {
			utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data", utility.KVP{Key: "data", Value: record.Body})
			return err
		}

		utility.Info("SQSMessage", "Received event from Amazon SQS", utility.KVP{Key: "order", Value: order},
			utility.KVP{Key: "messageId", Value: record.MessageId}, utility.KVP{Key: "eventSource", Value: record.EventSource})
	}

	return nil
}
