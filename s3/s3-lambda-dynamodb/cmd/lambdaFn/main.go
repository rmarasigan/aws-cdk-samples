package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/s3/s3-lambda-dynamodb/api"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/s3/s3-lambda-dynamodb/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/s3/s3-lambda-dynamodb/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, events events.S3Event) error {
	var (
		tablename = os.Getenv("TABLE_NAME")
	)

	if tablename == "" {
		err := errors.New("DynamoDB TABLE_NAME environment variable is not set")
		utility.Error(err, "EnvError", "DynamoDB TABLE_NAME is not configured on the environment")

		return err
	}

	for _, record := range events.Records {
		var (
			item   = new(api.Item)
			key    = record.S3.Object.Key
			bucket = record.S3.Bucket.Name
		)

		// Fetch the object data from the S3 bucket
		data, err := awswrapper.S3GetObject(ctx, bucket, key)
		if err != nil {
			utility.Error(err, "S3Error", "Failed to get object from S3 bucket",
				utility.KVP{Key: "bucket", Value: bucket}, utility.KVP{Key: "key", Value: key})

			return err
		}

		// Unmarshal the raw data
		err = api.UnmarshalJSON(data, item)
		if err != nil {
			utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data",
				utility.KVP{Key: "bucket", Value: bucket}, utility.KVP{Key: "key", Value: key})

			return err
		}

		// Validate the raw data
		err = item.ValidateRequest()
		if err != nil {
			utility.Error(err, "ItemError", "Some field(s) is/are missing",
				utility.KVP{Key: "bucket", Value: bucket}, utility.KVP{Key: "key", Value: key})

			return err
		}

		// Insert the record into the DynamoDB Table
		err = awswrapper.DynamoPutItem(ctx, tablename, *item)
		if err != nil {
			utility.Error(err, "DynamoDBError", "Failed to put item to the DynamoDB Table",
				utility.KVP{Key: "bucket", Value: bucket}, utility.KVP{Key: "key", Value: key},
				utility.KVP{Key: "tablename", Value: tablename})

			return err
		}

		// Delete object after being processed
		err = awswrapper.S3DeleteObject(ctx, bucket, key)
		if err != nil {
			utility.Error(err, "S3Error", "Failed to delete/remove object from S3 bucket",
				utility.KVP{Key: "bucket", Value: bucket}, utility.KVP{Key: "key", Value: key})

			return err
		}
	}

	return nil
}
