package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-s3/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-s3/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-s3/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, data json.RawMessage) error {
	var (
		event  = new(schema.Event)
		bucket = os.Getenv("BUCKET_NAME")
	)

	if len(bucket) == 0 {
		err := errors.New("s3 bucket BUCKET_NAME environment is not set")
		utility.Error(err, "EnvError", "S3 Bucket BUCKENT_NAME is not configured on the environment")

		return err
	}

	// Unmarshal the received JSON event
	err := json.Unmarshal([]byte(data), event)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data", utility.KVP{Key: "event", Value: event})
		return err
	}

	// Check if the key exist in the payload
	if len(event.Key) == 0 {
		err := errors.New("'key' field is not set")
		utility.Error(err, "JSONError", "The 'key' field is not set in the payload", utility.KVP{Key: "event", Value: event})

		return err
	}

	switch event.Action {
	case "get":
		// Retrieve the object data and print the object data
		object, err := awswrapper.S3GetObject(ctx, bucket, event.Key)
		if err != nil {
			utility.Error(err, "S3Error", "Failed to retrieve object", utility.KVP{Key: "key", Value: event.Key})
			return err
		}

		utility.Info("S3GetObject", "Retrieving the object data", utility.KVP{Key: "object", Value: string(object)})
		return nil

	case "delete":
		// Delete the object from the S3 bucket
		err := awswrapper.S3DeleteObject(ctx, bucket, event.Key)
		if err != nil {
			utility.Error(err, "S3Error", "Failed to remove the object", utility.KVP{Key: "key", Value: event.Key})
			return err
		}

		utility.Info("S3DeleteObject", "Removing the object from the bucket", utility.KVP{Key: "bucket", Value: bucket},
			utility.KVP{Key: "key", Value: event.Key})
		return nil

	default:
		// Create a new order object using the information
		// coming from the JSON-encoded data
		if event.Action == "put" {
			event.Order.Status = event.Order.Status.Received()
			event.Order.Timestamp = time.Now().Format("02 Jan 2006 15:04:05")

			order, err := event.Order.Marshal()
			if err != nil {
				utility.Error(err, "JSONError", "Failed to marshal the order", utility.KVP{Key: "key", Value: event.Key},
					utility.KVP{Key: "order", Value: event.Order})
				return err
			}

			err = awswrapper.S3PutObject(ctx, bucket, event.Key, order)
			if err != nil {
				utility.Error(err, "S3Error", "Failed to upload object to the bucket", utility.KVP{Key: "bucket", Value: bucket},
					utility.KVP{Key: "key", Value: event.Key}, utility.KVP{Key: "order", Value: event.Order})
				return err
			}

			utility.Info("S3PutObject", "Successfully uploaded the object to the bucket", utility.KVP{Key: "bucket", Value: bucket},
				utility.KVP{Key: "key", Value: event.Key}, utility.KVP{Key: "order", Value: event.Order})
			return nil

		} else {
			err := errors.New("incorrect type of 'action'")
			utility.Error(err, "ActionError", "Wrong type of 'action'", utility.KVP{Key: "action", Value: event.Action},
				utility.KVP{Key: "key", Value: event.Key}, utility.KVP{Key: "order", Value: event.Order})

			return err
		}
	}
}
