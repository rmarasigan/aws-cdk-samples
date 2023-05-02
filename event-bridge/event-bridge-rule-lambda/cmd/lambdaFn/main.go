package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-rule-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives a PutObject event from the S3 Bucket and writes S3
// events information to Lambda's CloudWatch Logs.
func handler(ctx context.Context, event events.CloudWatchEvent) error {
	utility.Info("S3Event", "S3 PutObject Event", utility.KVP{Key: "event", Value: event})
	return nil
}
