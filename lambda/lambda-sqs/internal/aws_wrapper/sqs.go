package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sqs/internal/utility"
)

const (
	AWS_REGION       = "us-east-1"
	SQS_MSG_GROUP_ID = "process.order"
)

var (
	sqsClient *sqs.Client
)

// initSQSClient initializes the SQS Client from the
// provided configuration.
func initSQSClient(ctx context.Context) {
	if sqsClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "SQSClientError", "failed to load the default config")
		return
	}

	sqsClient = sqs.NewFromConfig(cfg)
}

// SQSSendMessage initializes the SQS client and delivers message to the specified queue.
func SQSSendMessage(ctx context.Context, queue, message string) error {
	// Initialize the SQSClient
	initSQSClient(ctx)

	var input = &sqs.SendMessageInput{
		QueueUrl:       aws.String(queue),
		MessageBody:    aws.String(message),
		MessageGroupId: aws.String(SQS_MSG_GROUP_ID),
	}

	_, err := sqsClient.SendMessage(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
