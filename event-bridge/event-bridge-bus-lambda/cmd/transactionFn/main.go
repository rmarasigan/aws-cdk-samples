package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-bus-lambda/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-bus-lambda/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/event-bridge/event-bridge-bus-lambda/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, data json.RawMessage) error {
	var (
		transaction    = new(schema.Transaction)
		CANCEL_SOURCE  = "transaction:cancel"
		PAYMENT_SOURCE = "transaction:payment"
		EVENT_BUS      = os.Getenv("EVENT_BUS_NAME")
	)

	// Check if the EventBridge Event Bus is configured
	if EVENT_BUS == "" {
		err := errors.New("eventbridge EVENT_BUS_NAME is not set")
		utility.Error(err, "EnvError", "EVENT_BUS_NAME is not configured on the environment")

		return err
	}

	// Unmarshal the JSON-encoded data
	err := json.Unmarshal([]byte(data), transaction)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data", utility.KVP{Key: "data", Value: data})
		return err
	}

	detail, err := transaction.Marshal()
	if err != nil {
		utility.Error(err, "JSONError", "Failed to marshal the transaction data", utility.KVP{Key: "data", Value: data})
		return err
	}

	switch transaction.Type {
	case "payment":
		// Send the event to the configured Event Bus and specific source (transaction:payment)
		err := awswrapper.EventBridgePutEvents(ctx, detail, PAYMENT_SOURCE, EVENT_BUS)
		if err != nil {
			utility.Error(err, "EVBError", "Failed to send events to the Event Bus", utility.KVP{Key: "source", Value: PAYMENT_SOURCE}, utility.KVP{Key: "detail", Value: detail})
			return err
		}

	case "cancel":
		// Send the event to the configured Event Bus and specific source (transaction:cancel)
		err := awswrapper.EventBridgePutEvents(ctx, detail, CANCEL_SOURCE, EVENT_BUS)
		if err != nil {
			utility.Error(err, "EVBError", "Failed to send events to the Event Bus", utility.KVP{Key: "source", Value: CANCEL_SOURCE}, utility.KVP{Key: "detail", Value: detail})
			return err
		}

	default:
		// Return an error as the Transaction type is unknown
		err := errors.New("invalid transaction 'type'")
		utility.Error(err, "TransactionError", "Unknown transaction 'type'", utility.KVP{Key: "type", Value: transaction.Type}, utility.KVP{Key: "detail", Value: detail})
		return err
	}

	return nil
}
