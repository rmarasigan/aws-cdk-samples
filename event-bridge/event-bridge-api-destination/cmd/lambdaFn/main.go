package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-api-destination/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-api-destination/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-api-destination/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function once triggered, will create a new alarm information,
// and send the event to the configured EventBridge Bus.
func handler(ctx context.Context) error {
	var (
		alarm     = new(schema.Alarm)
		SOURCE    = "trigger:alarm"
		EVENT_BUS = os.Getenv("EVENT_BUS_NAME")
	)

	// Check if the EventBridge Event Bus is configured
	if EVENT_BUS == "" {
		err := errors.New("eventbridge EVENT_BUS_NAME is not set")
		utility.Error(err, "EnvError", "EVENT_BUS_NAME is not configured on the environment")

		return err
	}

	detail, err := alarm.CreateAlarm(ctx)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to marshal the alarm information")
		return err
	}

	// Send the event to the configured Event Bus and specific source
	err = awswrapper.EventBridgePutEvents(ctx, detail, SOURCE, EVENT_BUS)
	if err != nil {
		utility.Error(err, "EVBError", "Failed to send events to the Event Bus", utility.KVP{Key: "source", Value: SOURCE})
		return err
	}

	return nil
}
