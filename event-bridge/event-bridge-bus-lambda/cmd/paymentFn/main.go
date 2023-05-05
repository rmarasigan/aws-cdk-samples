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

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	utility.Info("PaymentEvent", "Transaction Payment Event", utility.KVP{Key: "event", Value: event})
	return nil
}
