package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-schedule-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) error {
	utility.Info("LambdaInvoke", "The lambda function is invoked by the scheduled EventBridge")
	return nil
}
