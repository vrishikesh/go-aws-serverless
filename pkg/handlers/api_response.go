package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(status int, body any) (*events.APIGatewayProxyResponse, error) {
	r := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	r.StatusCode = status
	stringBody, _ := json.Marshal(body)
	r.Body = string(stringBody)
	return &r, nil
}
