package awswrapper

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/schema"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/utility"
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

// DynamoQueryObject initializes the DynamoDB client and retrieves the specific order.
func DynamoQueryObject(ctx context.Context, tablename, referenceId string) (*dynamodb.QueryOutput, error) {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	// Create a key expression
	key := expression.Key("referenceId").Equal(expression.Value(referenceId))

	// Build an expression to retrieve item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		return nil, err
	}

	// Build the query input parameter
	var input = &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	output, err := dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DynamoUpdateItem intializes the DynamoDB client and inserts a new item to the
// configured DynamoDB Table.
func DynamoUpdateItem(ctx context.Context, tablename string, order schema.Order) error {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	// Convert the record to []types.AttributeValue
	values, err := attributevalue.MarshalList(order.OrderLine)
	if err != nil {
		return err
	}

	// It will check if the "order_line" attribute already exist, otherwise
	// the value of ":empty_list" will be used and then the new ":lines" will
	// be appended to the list
	update := aws.String("SET order_line = list_append(if_not_exists(order_line, :empty_list), :lines)")

	// Use the built expression to populate the DynamoDB Update
	// Item API input parameters
	var params = &dynamodb.UpdateItemInput{
		TableName: aws.String(tablename),
		Key: map[string]types.AttributeValue{
			"referenceId": &types.AttributeValueMemberS{Value: order.ReferenceId},
		},
		UpdateExpression: update,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":lines":      &types.AttributeValueMemberL{Value: values},
			":empty_list": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
		},
	}

	_, err = dynamoClient.UpdateItem(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// DynamoUpdateObject initializes the DynamoDB client and update the object
// that is in the list.
func DynamoUpdateObject(ctx context.Context, tablename, referenceId string, index int, line schema.OrderLine) error {
	// Initialize the DynamoClient
	initDynamoClient(ctx)

	// It will update the specific object inside the list by providing the "index"
	update := aws.String(fmt.Sprintf("SET order_line[%v].price = :item_price, order_line[%v].quantity = :item_quantity", index, index))

	// Use the built expression to populate the DynamoDB Update
	// Item API input parameters
	var input = &dynamodb.UpdateItemInput{
		TableName: aws.String(tablename),
		Key: map[string]types.AttributeValue{
			"referenceId": &types.AttributeValueMemberS{Value: referenceId},
		},
		UpdateExpression: update,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":item_price":    &types.AttributeValueMemberN{Value: fmt.Sprint(line.Price)},
			":item_quantity": &types.AttributeValueMemberN{Value: fmt.Sprint(line.Quantity)},
		},
	}

	_, err := dynamoClient.UpdateItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
