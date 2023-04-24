package schema

import (
	"fmt"
	"time"
)

type Status string

func (Status) Received() Status {
	return "RECEIVED"
}

func (Status) Processed() Status {
	return "PROCESSED"
}

type Item struct {
	ID    string  `dynamodbav:"id"`
	Name  string  `dynamodbav:"name"`
	Price float64 `dynamodb:"price"`
}

type Order struct {
	ReferenceID string `dynamodbav:"referenceId"`
	Status      Status `dynamodbav:"status"`
	Item        []Item `dynamodbav:"item"`
	Quantity    int64  `dynamodbav:"quantity"`
	Timestamp   string `dynamodbav:"timestamp"`
}

// CreateOrder set the default values to the Order
// schema to create a new order.
func (order Order) CreateOrder() Order {
	order.ReferenceID = fmt.Sprintf("SMPLORDR-%s", time.Now().Format("02012006150405"))
	order.Status = order.Status.Received()

	order.Item = append(order.Item, Item{
		ID:    "SMPLITM-12345",
		Name:  "Sample Item",
		Price: 1234.50,
	})

	order.Quantity = 1
	order.Timestamp = time.Now().Format("02 Jan 2006 15:04:05")

	return order
}
