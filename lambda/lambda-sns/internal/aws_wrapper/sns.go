package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-sns/internal/utility"
)

const (
	AWS_REGION = "us-east-1"
)

var (
	snsClient *sns.Client
)

// initSNSClient initializes the Simple Notification
// Service Client from the provided configuration.
func initSNSClient(ctx context.Context) {
	if snsClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "SNSError", "failed to load the default config")
		return
	}

	snsClient = sns.NewFromConfig(cfg)
}

// SNSSubscribe subscribes an endpoint to the specified SNS Topic. The delivery
// of message is thru email address.
//
//  endpoint: an email address
func SNSSubscribe(ctx context.Context, topic, endpoint string) error {
	// Initialize the Simple Notification Service Client
	initSNSClient(ctx)

	var input = &sns.SubscribeInput{
		Protocol:              aws.String("email"),
		TopicArn:              aws.String(topic),
		Endpoint:              aws.String(endpoint),
		ReturnSubscriptionArn: true,
	}

	_, err := snsClient.Subscribe(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// SNSPublish publish a message to an SNS Topic. The delivery of message
// is thru email address.
func SNSPublish(ctx context.Context, topic, message string) error {
	// Initialize the Simple Notification Service Client
	initSNSClient(ctx)

	var input = &sns.PublishInput{
		TopicArn: aws.String(topic),
		Subject:  aws.String("ERROR: APPLICATION ALERT"),
		Message:  aws.String(message),
	}

	// Publish a message to the specified SNS Topic.
	_, err := snsClient.Publish(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
