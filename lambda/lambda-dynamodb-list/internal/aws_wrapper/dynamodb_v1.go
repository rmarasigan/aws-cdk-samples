package awswrapper

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/schema"
)

var (
	sess           *session.Session
	dynamoV1Client dynamodbiface.DynamoDBAPI
)

// getSession creates a new session or gets an existing session.
func getSession() {
	if sess != nil {
		return
	}

	sess = session.Must(session.NewSession(aws.NewConfig().WithRegion(AWS_REGION)))
}

// initDynamoV1Client initializes the DynamoDB Client from the
// provided configuration.
func initDynamoV1Client() {
	getSession()

	if dynamoV1Client != nil {
		return
	}

	dynamoV1Client = dynamodb.New(sess)
}

// DynamoV1QueryObject initializes the DynamoDB client and retrieves the specific order.
func DynamoV1QueryObject(ctx context.Context, tablename, referenceId string) (*dynamodb.QueryOutput, error) {
	// Initialize the DynamoClient
	initDynamoV1Client()

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

	output, err := dynamoV1Client.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// DynamoV1UpdateItem initializes the DynamoDB client and inserts a new item to the
// configured DynamoDB Table.
func DynamoV1UpdateItem(ctx context.Context, tablename string, order schema.Order) error {
	// Initialize the DynamoClient
	initDynamoV1Client()

	// It will check if the "order_line" attribute already exist, otherwise
	// the value of ":empty_list" will be used and then the new ":lines" will
	// be appended to the list
	update := aws.String("SET order_line = list_append(if_not_exists(order_line, :empty_list), :lines)")

	// Convert the record to []*dynamodb.AttributeValue
	values, err := dynamodbattribute.MarshalList(order.OrderLine)
	if err != nil {
		return err
	}

	// Create the Update Item parameters
	var input = &dynamodb.UpdateItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"referenceId": {S: aws.String(order.ReferenceId)},
		},
		UpdateExpression: update,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":lines":      {L: values},
			":empty_list": {L: []*dynamodb.AttributeValue{}},
		},
	}

	_, err = dynamoV1Client.UpdateItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

// DynamoV1UpdateObject initializes the DynamoDB client and update the
// object that is in the list.
func DynamoV1UpdateObject(ctx context.Context, tablename, referenceId string, index int, line schema.OrderLine) error {
	// Initialize the DynamoClient
	initDynamoV1Client()

	// It will update the specific object inside the list by providing the "index"
	update := aws.String(fmt.Sprintf("SET order_line[%v].price = :item_price, order_line[%v].quantity = :item_quantity", index, index))

	// Use the built expression to populate the DynamoDB Update
	// Item API input parameters
	var input = &dynamodb.UpdateItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"referenceId": {S: aws.String(referenceId)},
		},
		UpdateExpression: update,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":item_price":    {N: aws.String(fmt.Sprint(line.Price))},
			":item_quantity": {N: aws.String(fmt.Sprint(line.Quantity))},
		},
	}

	_, err := dynamoV1Client.UpdateItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
