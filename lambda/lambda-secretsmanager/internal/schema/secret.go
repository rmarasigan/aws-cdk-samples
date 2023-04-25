package schema

import "encoding/json"

type Event struct {
	Action string `json:"action"`
	Secret Secret `json:"secret"`
}

type Secret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (secret *Secret) Marshal() ([]byte, error) {
	data, err := json.Marshal(secret)
	if err != nil {
		return nil, err
	}

	return data, nil
}
