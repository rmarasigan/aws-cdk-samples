package schema

import (
	"encoding/json"
)

type Status string

func (Status) Received() Status {
	return "RECEIVED"
}

func (Status) Processed() Status {
	return "PROCESSED"
}

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	ReferenceID string `json:"referenceId"`
	Status      Status `json:"status,omitempty"`
	Item        []Item `json:"item"`
	Quantity    int64  `json:"quantity"`
}

func (order *Order) Marshal() ([]byte, error) {
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	return data, nil
}
