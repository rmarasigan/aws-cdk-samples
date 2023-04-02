package api

import (
	"fmt"
	"strings"
)

type Coffee struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

func (c Coffee) ValidateRequest() error {
	var err_msg []string

	if c.Name == "" {
		err_msg = append(err_msg, "name")
	}

	if c.Description == "" {
		err_msg = append(err_msg, "description")
	}

	if len(c.Ingredients) == 0 {
		err_msg = append(err_msg, "ingredients")
	}

	if len(err_msg) > 0 {
		return fmt.Errorf("required field(s) is/are empty: %s", strings.Join(err_msg, ", "))
	}

	return nil
}
