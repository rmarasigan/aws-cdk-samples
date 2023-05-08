package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-bus-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function will be triggered, if it received
// an event from an EventBus and writes events information
// to Lambda's CloudWatch Logs.
func handler(ctx context.Context, event events.CloudWatchEvent) error {
	utility.Info("CancelEvent", "Transaction Cancel Event", utility.KVP{Key: "event", Value: event})
	return nil
}
