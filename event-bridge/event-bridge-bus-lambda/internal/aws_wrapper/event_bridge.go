package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-bus-lambda/internal/utility"
)

const (
	AWS_REGION = "us-east-1"
)

var (
	evbClient *eventbridge.Client
)

// initEventBridgeClient intializes the EventBridge Service
// Client from the provided configuration.
func initEventBridgeClient(ctx context.Context) {
	if evbClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "EVBError", "failed to load the default config")
		return
	}

	evbClient = eventbridge.NewFromConfig(cfg)
}

// EventBridgePutEvents send custom events to the specified Amazon EventBridge Event
// Bus Name.
func EventBridgePutEvents(ctx context.Context, detail, source, eventBusName string) error {
	// Initialize the EventBridge Client
	initEventBridgeClient(ctx)

	var input = &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:       aws.String(detail),
				DetailType:   aws.String("transaction"),
				Source:       aws.String(source),
				EventBusName: aws.String(eventBusName),
			},
		},
	}

	_, err := evbClient.PutEvents(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
