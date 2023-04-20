package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/cloudwatch/cloudwatch-subscription-lambda/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/cloudwatch/cloudwatch-subscription-lambda/internal/trail"
	"github.com/rmarasigan/aws-cdk-samples/cloudwatch/cloudwatch-subscription-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event awswrapper.CloudWatchEvent) error {
	// Decode and decompressed the received CloudWatch Event
	data, err := event.DecodeData()
	if err != nil {
		utility.Error(err, "CWError", "Failed to decode and decompressed the received event from CloudWatch", utility.KVP{Key: "event", Value: event})
		return err
	}

	// Log the CloudWatch Event
	trail.Info("CloudWatch Event: %s", utility.EncodeJSON(data))

	return nil
}
