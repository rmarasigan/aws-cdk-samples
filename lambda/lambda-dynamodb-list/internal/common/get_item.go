package common

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	awswrapper "github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/aws_wrapper"
	"github.com/rmarasigan/aws-cdk-samples/lambda/lambda-dynamodb-list/internal/schema"
)

// ExistingOrderLine retrieves the specific order and loop through the order lines
// to check if the said order line exist or not.
func ExistingOrderLine(ctx context.Context, tablename, referenceId, item_id string) (bool, int, error) {
	var order schema.Order

	output, err := awswrapper.DynamoQueryObject(ctx, tablename, referenceId)
	if err != nil {
		return false, 0, err
	}

	if output.Count > 0 {
		// Unmarshal to the actual structure
		err = attributevalue.UnmarshalMap(output.Items[0], &order)
		if err != nil {
			return false, 0, err
		}

		for index, line := range order.OrderLine {
			if line.ItemID == item_id {
				return true, index, err
			}
		}
	}

	return false, 0, nil
}

func ExistingOrderLineV1(ctx context.Context, tablename, referenceId, item_id string) (bool, int, error) {
	var order schema.Order

	output, err := awswrapper.DynamoV1QueryObject(ctx, tablename, referenceId)
	if err != nil {
		return false, 0, err
	}

	if *output.Count > 0 {
		// Unmarshal to the actual structure
		err = dynamodbattribute.UnmarshalMap(output.Items[0], &order)
		if err != nil {
			return false, 0, err
		}

		for index, line := range order.OrderLine {
			if line.ItemID == item_id {
				return true, index, err
			}
		}
	}

	return false, 0, nil
}
