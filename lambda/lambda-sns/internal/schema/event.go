package schema

import (
	"fmt"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Content struct {
	User        User   `json:"user"`
	Application string `json:"application"`
	Message     string `json:"message"`
}

type Event struct {
	Action  string  `json:"action"`
	Content Content `json:"content"`
}

// Message returns the message body for Alert SNS Topic.
func (event *Event) Message() string {
	message := fmt.Sprintf("Hello,\n\nWe regret to inform you that your application, %s, encountered an error. The error message is as follows:\n\n\t%s\n\nIf you have any questions or concerns, please contact our support team at xyz@email.com.", event.Content.Application, event.Content.Message)

	return message
}
