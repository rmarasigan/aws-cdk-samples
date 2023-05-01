package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sns/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sns/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sns/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, data json.RawMessage) error {
	var (
		event = new(schema.Event)
		topic = os.Getenv("TOPIC_ARN")
	)

	// Check if the SNS Topic is configured
	if topic == "" {
		err := errors.New("sns TOPIC_ARN environment is not set")
		utility.Error(err, "EnvError", "SNS TOPIC_ARN is not configured on the environment")

		return err
	}

	// Unmarshal the JSON-encoded data
	err := json.Unmarshal([]byte(data), event)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data", utility.KVP{Key: "data", Value: data})
		return err
	}

	switch event.Action {
	case "subscribe":
		// Subscribe the endpoint to the configured SNS Topic
		err = awswrapper.SNSSubscribe(ctx, topic, event.Content.User.Email)
		if err != nil {
			utility.Error(err, "SNSError", "Failed to subscribe an endpoint to the SNS Topic", utility.KVP{Key: "topic", Value: topic},
				utility.KVP{Key: "event", Value: event})
			return err
		}

	case "alert":
		// Publish the message to the configured SNS Topic
		err = awswrapper.SNSPublish(ctx, topic, event.Message())
		if err != nil {
			utility.Error(err, "SNSError", "Failed to publish message to the topic", utility.KVP{Key: "topic", Value: topic},
				utility.KVP{Key: "event", Value: event})
			return err
		}
	}

	return nil
}
