package api

import (
	"fmt"
	"strings"
	"time"
)

type Item struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    Category `json:"category"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (i *Item) ValidateRequest() error {
	var err_msg []string
	i.ID = time.Now().Format("20060102150405")

	if i.Name == "" {
		err_msg = append(err_msg, "item name")
	}

	if i.Description == "" {
		err_msg = append(err_msg, "item description")
	}

	if i.Category.ID == "" {
		err_msg = append(err_msg, "category id")
	}

	if i.Category.Name == "" {
		err_msg = append(err_msg, "category name")
	}

	if len(err_msg) > 0 {
		return fmt.Errorf("required field(s) is/are empty: %s", strings.Join(err_msg, ", "))
	}

	return nil
}
