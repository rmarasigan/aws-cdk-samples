package api

import "encoding/json"

// UnmarshalJSON unmarshal the JSON-encoded data and stores the result in the value pointed to v.
func UnmarshalJSON(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}
