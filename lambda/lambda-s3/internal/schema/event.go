package schema

type Event struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Order  Order  `json:"order,omitempty"`
}
