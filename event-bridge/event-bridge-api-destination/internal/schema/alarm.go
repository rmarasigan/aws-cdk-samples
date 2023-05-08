package schema

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

type Alarm struct {
	RequestID string `json:"request_id"`
	Source    string `json:"source"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// CreateAlarm creates a dummy alarm message with the
// AWS Request and Lambda Function ARN.
func (a *Alarm) CreateAlarm(ctx context.Context) (string, error) {
	lctx, _ := lambdacontext.FromContext(ctx)
	a.RequestID = lctx.AwsRequestID
	a.Source = lctx.InvokedFunctionArn
	a.Message = "This is a sample alarm message"
	a.Timestamp = time.Now().Format("2006-01-02 15:04:05")

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
