package api

import (
	"encoding/json"

	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-cors-lambda/internal/utility"
)

// MarshalJSON marshals the body.
func MarshalJSON(body any) []byte {
	encodeJSON, err := json.Marshal(body)
	if err != nil {
		utility.Error(err, "JSONError", "Failed to encode JSON")
	}

	return encodeJSON
}

// UnmarshalJSON unmarshal the JSON-encoded data and stores the result in the value pointed to v.
func UnmarshalJSON(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
