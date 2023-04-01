package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/s3/s3-eventbridge-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	utility.Info("S3Events", "S3 Events Information",
		utility.KVP{Key: "ID", Value: event.ID}, utility.KVP{Key: "Source", Value: event.Source},
		utility.KVP{Key: "Detail", Value: event.Detail}, utility.KVP{Key: "DetailType", Value: event.DetailType})

	return nil
}
