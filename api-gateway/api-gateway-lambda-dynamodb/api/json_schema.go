package api

import (
	"fmt"
	"strings"
	"time"
)

type Coffee struct {
	Key         string   `json:"key,omitempty" dynamodbav:"key"`
	Name        string   `json:"name" dynamodbav:"name"`
	Description string   `json:"description" dynamodbav:"description"`
	Ingredients []string `json:"ingredients" dynamodbav:"ingredients"`
}

func (c *Coffee) ValidateRequest() error {
	var err_msg []string
	c.Key = time.Now().Format("20060102150405")

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
