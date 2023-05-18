package api

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-cognito-lambda/internal/utility"
)

const (
	ALLOW_ORIGIN  = "*"
	ALLOW_METHOD  = "*"
	ALLOW_HEADERS = "Content-Type"
	CONTENT_TYPE  = "application/json"
)

// OK returns an API Gateway Response that contains a body and an HTTP
// OK status.
func OK(body any) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                   CONTENT_TYPE,
			"Access-Control-Allow-Origin":    ALLOW_ORIGIN,
			"Access-Control-Request-Methods": ALLOW_METHOD,
			"Access-Control-Allow-Headers":   ALLOW_HEADERS,
		},
		StatusCode: http.StatusOK,
		Body:       utility.EncodeJSON(body),
	}, nil
}
