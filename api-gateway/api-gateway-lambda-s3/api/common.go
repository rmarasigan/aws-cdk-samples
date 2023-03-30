package api

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rmarasigan/aws-cdk-samples/api-gateway/api-gateway-lambda-s3/internal/utility"
)

const (
	ALLOW_ORIGIN  = "*"
	ALLOW_METHOD  = "*"
	ALLOW_HEADERS = "Content-Type"
	CONTENT_TYPE  = "application/json"
)

type Body struct {
	ErrorMsg *string `json:"err_msg,omitempty"`
}

// OK returns an API Gateway Response that contains a body and an HTTP
// OK status.
func OK(body any) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                  CONTENT_TYPE,
			"Access-Control-Allow-Origin":   ALLOW_ORIGIN,
			"Access-Control-Request-Method": ALLOW_METHOD,
			"Access-Control-Allow-Headers":  ALLOW_HEADERS,
		},
		StatusCode: http.StatusOK,
		Body:       utility.EncodeJSON(body),
	}, nil
}

// OKWithoutBody returns an API Gateway Response with HTTP OK status without a body.
func OKWithoutBody() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                  CONTENT_TYPE,
			"Access-Control-Allow-Origin":   ALLOW_ORIGIN,
			"Access-Control-Request-Method": ALLOW_METHOD,
			"Access-Control-Allow-Headers":  ALLOW_HEADERS,
		},
		StatusCode: http.StatusOK,
	}, nil
}

// BadRequest returns an API Gateway Response with HTTP BadRequest status with
// an error message as the body.
func BadRequest(err error) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                  CONTENT_TYPE,
			"Access-Control-Allow-Origin":   ALLOW_ORIGIN,
			"Access-Control-Request-Method": ALLOW_METHOD,
			"Access-Control-Allow-Headers":  ALLOW_HEADERS,
		},
		StatusCode: http.StatusBadRequest,
		Body:       utility.EncodeJSON(Body{ErrorMsg: aws.String(err.Error())}),
	}, nil
}
