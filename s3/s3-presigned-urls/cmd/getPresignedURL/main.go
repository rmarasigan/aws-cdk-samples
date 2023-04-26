package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/s3/s3-presigned-urls/api"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/s3/s3-presigned-urls/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/s3/s3-presigned-urls/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// handler function receives the Amazon API Gateway event record data as input
// and returns a response body with 200 OK HTTP Status. The API response must
// include the "Access-Control-Allow-Origin", "Access-Control-Allow-Methods" and
// "Access-Control-Allow-Headers" headers.
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		bucket      = os.Getenv("BUCKET_NAME")
		response    = new(awswrapper.Response)
		key         = request.QueryStringParameters["key"]
		action      = request.QueryStringParameters["action"]
		contentType = request.QueryStringParameters["content_type"]
	)

	// Check if the bucket is set in the environment
	if bucket == "" {
		err := errors.New("'BUCKET_NAME' is not set on the environment")
		utility.Error(err, "ConfigError", "'BUCKET_NAME' is not configured on the environment")

		return api.InternalServerError()
	}

	// Check if the key parameter is set
	if key == "" {
		err := errors.New("'key' parameter is not set")
		utility.Error(err, "APIError", "'key' parameter is missing")

		return api.BadRequest(err)
	}

	switch action {
	case "upload":
		// Check if the content_type parameter is set
		if contentType == "" {
			err := errors.New("'content_type' parameter is not set")
			utility.Error(err, "APIError", "'content_type' parameter is missing", utility.KVP{Key: "key", Value: key})

			return api.BadRequest(err)
		}

		// Get presigned HTTP Request
		presign, err := awswrapper.S3PresignPutObject(ctx, bucket, key, contentType)
		if err != nil {
			utility.Error(err, "S3Error", "Failed to get the presigned URL for PUT", utility.KVP{Key: "key", Value: key})
			return api.InternalServerError()
		}

		response.URL = &presign.URL
		response.Method = &presign.Method

	default:
		// Get presigned HTTP Request
		presign, err := awswrapper.S3PresignGetObject(ctx, bucket, key)
		if err != nil {
			utility.Error(err, "S3Error", "Failed to get the presigned URL for GET", utility.KVP{Key: "key", Value: key})
			return api.InternalServerError()
		}

		response.URL = &presign.URL
		response.Method = &presign.Method
	}

	return api.OK(response)
}
