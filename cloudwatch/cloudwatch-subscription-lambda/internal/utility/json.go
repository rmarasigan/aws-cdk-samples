package utility

import "encoding/json"

// EncodeJSON marshals the body to produce a string format JSON.
func EncodeJSON(body any) string {
	encodeJSON, err := json.Marshal(body)
	if err != nil {
		Error(err, "JSONError", "Failed to encode JSON")
	}

	return string(encodeJSON)
}
