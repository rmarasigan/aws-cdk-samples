package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-sqs/api"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-sqs/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon SQS event message data as input,
// and writes the message data to CloudWatch Logs.
func handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var item = new(api.Item)

		err := api.UnmarshalJSON([]byte(record.Body), item)
		if err != nil {
			utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data")
			return err
		}

		utility.Info("SQSEvent", "Message to process", utility.KVP{Key: "Item", Value: item})
	}

	return nil
}
