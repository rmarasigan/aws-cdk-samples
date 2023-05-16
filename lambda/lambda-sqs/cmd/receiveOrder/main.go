package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sqs/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sqs/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sqs/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function once triggered, validates the JSON-encoded
// data and send the message to an SQS Queue.
func handler(ctx context.Context, data json.RawMessage) error {
	var (
		queue = os.Getenv("QUEUE_URL")
		order = new(schema.Order)
	)

	// Check if the QUEUE_URL is configured
	if len(queue) == 0 {
		err := errors.New("sqs QUEUE_URL environment variable is not set")
		utility.Error(err, "SQSError", "SQS QUEUE_URL is not configured on the environment")

		return err
	}

	// Unmarshal the received JSON event
	err := json.Unmarshal([]byte(data), order)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal JSON-encoded data", utility.KVP{Key: "data", Value: data})
		return err
	}

	order.Status = order.Status.Received()

	message, err := order.Marshal()
	if err != nil {
		utility.Error(err, "JSONError", "failed to marshal the order", utility.KVP{Key: "order", Value: order})
		return err
	}

	// Send the message to the queue
	err = awswrapper.SQSSendMessage(ctx, queue, string(message))
	if err != nil {
		utility.Error(err, "SQSError", "failed to send message", utility.KVP{Key: "queue", Value: queue}, utility.KVP{Key: "order", Value: order})
		return err
	}

	return nil
}
