package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb/internal/utility"
)

const (
	AWS_REGION = "us-east-1"
)

var (
	dynamoClient *dynamodb.Client
)

// initDynamoClient initializes the DynamoDB Client from the
// provided configuration.
func initDynamoClient(ctx context.Context) {
	if dynamoClient != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "DynamoClientError", "failed to load the default config")
		return
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
}

// DynamoPutItem initializes the DynamoDB client and inserts the item into the DynamoDB Table.
func DynamoPutItem(ctx context.Context, tablename string, order schema.Order) error {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	// Convert the record to map[string]types.AttributeValue
	// that is to be used in the PutItemInput
	value, err := attributevalue.MarshalMap(order)
	if err != nil {
		return err
	}

	var input = &dynamodb.PutItemInput{
		Item:      value,
		TableName: aws.String(tablename),
	}

	// Insert the item to the DynamoDB Table
	_, err = dynamoClient.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// DynamoUpdateItem intializes the DynamoDB client and updates an existing item's
// attribute.
func DynamoUpdateItem(ctx context.Context, tablename, key, status string) error {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	// Create an update expression
	update := expression.Set(expression.Name("status"), expression.Value(status))

	// Using the update to create a DynamoDB Expression
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	// Use the built expression to populate the DynamoDB Update
	// Item API input parameters
	var params = &dynamodb.UpdateItemInput{
		TableName: aws.String(tablename),
		Key: map[string]types.AttributeValue{
			"referenceId": &types.AttributeValueMemberS{Value: key},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	_, err = dynamoClient.UpdateItem(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
