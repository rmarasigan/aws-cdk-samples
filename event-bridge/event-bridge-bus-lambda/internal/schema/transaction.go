package schema

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ReferenceId string `json:"reference_id"`
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp,omitempty"`
}

func (t *Transaction) Marshal() (string, error) {
	t.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
