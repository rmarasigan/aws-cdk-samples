package schema

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Transaction struct {
	ID          string   `json:"id"`
	OrderLineID string   `json:"order_line_id"`
	Customer    Customer `json:"customer"`
	Total       float64  `json:"total"`
	Change      float64  `json:"change"`
}
