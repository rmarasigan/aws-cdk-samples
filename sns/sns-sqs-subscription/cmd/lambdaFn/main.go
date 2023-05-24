package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/sns/sns-sqs-subscription/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives an Amazon SQS event message data as input, writes
// the event message data to the CloudWatch logs.
func handler(ctx context.Context, event events.SQSEvent) error {
	var records = event.Records

	// Check if there are events
	if len(records) == 0 {
		utility.Info("SQSEvent", "there are no records found")
		return nil
	}

	for _, record := range records {
		utility.Info("SQSEvent", "Events from the SNS Topic -> SQS Queue", utility.KVP{Key: "record", Value: record})
	}

	return nil
}
