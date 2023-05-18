package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-sqs/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/step-functions/step-functions-sqs/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon SQS event message data as input
// and will write the transaction message to CloudWatch Logs.
func handler(ctx context.Context, event events.SQSEvent) error {
	var (
		records     = event.Records
		transaction = new(schema.Transaction)
	)

	// Check if there are any event records
	if len(records) == 0 {
		utility.Info("SQS", "no record(s) found")
		return nil
	}

	for _, record := range records {
		// Unmarshal the event message
		err := json.Unmarshal([]byte(record.Body), transaction)
		if err != nil {
			utility.Error(err, "JSONError", "failed to unmarshal JSON-encoded data", utility.KVP{Key: "record", Value: record.Body})
			return err
		}

		utility.Info("SQSMessage", "recieved event from Amazon SQS", utility.KVP{Key: "transaction", Value: transaction})
	}

	return nil
}
