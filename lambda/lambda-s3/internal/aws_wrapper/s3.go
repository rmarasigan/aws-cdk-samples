package awswrapper

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-s3/internal/utility"
)

const (
	AWS_REGION = "us-east-1"
)

var (
	s3Client *s3.Client
)

// initS3Client initializes the S3 Client from the
// provided configuration.
func initS3Client(ctx context.Context) {
	if s3Client != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "S3ClientError", "failed to load the default config")
		return
	}

	s3Client = s3.NewFromConfig(cfg)
}

// S3PutObject initializes the S3 client and uploads an object to the bucket.
func S3PutObject(ctx context.Context, bucket, key string, data []byte) error {
	// Initialize the S3Client
	initS3Client(ctx)

	var input = &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	// Upload the object to the S3 Bucket
	_, err := s3Client.PutObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// S3GetObject initializes the S3 client and retrieves the object from the
// S3 Bucket.
func S3GetObject(ctx context.Context, bucket, key string) ([]byte, error) {
	// Initialize the S3Client
	initS3Client(ctx)

	var input = &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Get the object from the S3 bucket
	output, err := s3Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	var buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// S3DeleteObject initializes the S3 client and removes the object from the
// S3 Bucket.
func S3DeleteObject(ctx context.Context, bucket, key string) error {
	// Initialize the S3Client
	initS3Client(ctx)

	var input = &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// Delete/remove the object from the S3 bucket
	_, err := s3Client.DeleteObject(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
