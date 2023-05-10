package schema

type Order struct {
	ReferenceId string      `json:"referenceId" dynamodbav:"referenceId"`
	OrderLine   []OrderLine `json:"order_line" dynamodbav:"order_line"`
}

type OrderLine struct {
	ItemID   string  `json:"item_id" dynamodbav:"item_id"`
	Price    float64 `json:"price" dynamodbav:"price"`
	Quantity int64   `json:"quantity" dynamodbav:"quantity"`
}
