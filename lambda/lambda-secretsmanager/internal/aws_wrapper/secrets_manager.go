package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-secretsmanager/internal/utility"
)

const (
	AWS_REGION = "us-east-1"
)

var (
	secretManagerClient *secretsmanager.Client
)

// initSecretManagerClient initializes the Secrets Manager
// Client from the provided configuration.
func initSecretManagerClient(ctx context.Context) {
	if secretManagerClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "SecretsManagerError", "failed to load the default config")
		return
	}

	secretManagerClient = secretsmanager.NewFromConfig(cfg)
}

// SMPutSecretValue initializes the Secrets Manager client
// and creates a new version of the secret value.
func SMPutSecretValue(ctx context.Context, secretId, secret string) (*secretsmanager.PutSecretValueOutput, error) {
	// Initialize the SecretManager Client
	initSecretManagerClient(ctx)

	var input = &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretId),
		SecretString: aws.String(secret),
	}

	output, err := secretManagerClient.PutSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// SMGetSecretValue initializes the Secrets Manager client and
// retrieves the contents of the encrypted fields.
func SMGetSecretValue(ctx context.Context, secretId string) (*secretsmanager.GetSecretValueOutput, error) {
	// Initialize the SecretManager Client
	initSecretManagerClient(ctx)

	var input = &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	output, err := secretManagerClient.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}

	return output, nil
}
