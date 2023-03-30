package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-s3/api"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-s3/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-s3/internal/utility"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	body := request.Body
	var coffee api.Coffee
	bucket := os.Getenv("BUCKET_NAME")

	err := api.UnmarshalJSON([]byte(body), &coffee)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to unmarshal JSON-encoded data")
		return api.BadRequest(err)
	}

	err = coffee.ValidateRequest()
	if err != nil {
		utility.Error(err, "APIError", "Some field(s) is/are missing")
		return api.BadRequest(err)
	}

	// Upload the request information to the S3 Bucket
	key := fmt.Sprintf("%s.json", time.Now().Format("2006-01-02-15-04-05"))
	err = awswrapper.S3PutObject(ctx, bucket, key, []byte(body))
	if err != nil {
		utility.Error(err, "S3Error", "failed to upload object to the s3 bucket",
			utility.KVP{Key: "bucket", Value: bucket}, utility.KVP{Key: "key", Value: key})

		return api.BadRequest(err)
	}

	return api.OKWithoutBody()
}
