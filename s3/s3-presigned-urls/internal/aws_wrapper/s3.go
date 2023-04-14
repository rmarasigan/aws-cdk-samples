package awswrapper

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rmarasigan/aws-cdk-samples/s3/s3-presigned-urls/internal/utility"
)

const (
	AWS_REGION = "us-east-1"
)

var (
	presignClient *s3.PresignClient
)

type Response struct {
	Method *string `json:"method,omitempty"`
	URL    *string `json:"url"`
}

// initS3Client initializes the S3 Client and the
// presign client from the provided configuration.
func initS3PresignClient(ctx context.Context) {
	if presignClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "S3ClientError", "failed to load the default config")
		return
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient = s3.NewPresignClient(s3Client)
}

// S3PresignPutObject initializes the S3 presgin client and generate the
// a presigned HTTP Request.
func S3PresignPutObject(ctx context.Context, bucket, key, contentType string) (*v4.PresignedHTTPRequest, error) {
	// Initialize the S3PresignClient
	initS3PresignClient(ctx)

	var input = &s3.PutObjectInput{
		Key:         aws.String(key),
		Bucket:      aws.String(bucket),
		ContentType: aws.String(contentType),
	}

	var options = func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(300 * int64(time.Second))
	}

	// Generate a presigned HTTP request
	response, err := presignClient.PresignPutObject(ctx, input, options)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// S3PresignGetObject initializes the S3 presgin client and generate the
// a presigned HTTP Request.
func S3PresignGetObject(ctx context.Context, bucket, key string) (*v4.PresignedHTTPRequest, error) {
	// Initialize the S3PresignClient
	initS3PresignClient(ctx)

	var input = &s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("attachment"),
	}

	var options = func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(300 * int64(time.Second))
	}

	// Generate a presigned HTTP request
	response, err := presignClient.PresignGetObject(ctx, input, options)
	if err != nil {
		return nil, err
	}

	return response, nil
}
