package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/sns/sns-lambda-subscription/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon SNS event record data as input, and writes
// the record data to CloudWatch Logs.
func handler(ctx context.Context, event events.SNSEvent) error {
	var records = event.Records

	// Check if there are events
	if len(records) == 0 {
		utility.Info("SNSEvent", "there are no records found")
		return nil
	}

	for _, record := range records {
		utility.Info("SNSEvent", "Events from the SNS Topic", utility.KVP{Key: "record", Value: record})
	}

	return nil
}
