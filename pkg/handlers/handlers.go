package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"go-aws-serverless/pkg/user"
)

var ErrorMethodNotAllowed = "method not allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(r events.APIGatewayProxyRequest, tableName string, dynamoClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	email := r.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dynamoClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

	result, err := user.FetchUsers(tableName, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func NewUser(r events.APIGatewayProxyRequest, tableName string, dynamoClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(r, tableName, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)
}

func UpdateUser(r events.APIGatewayProxyRequest, tableName string, dynamoClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(r, tableName, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, result)
}

func DeleteUser(r events.APIGatewayProxyRequest, tableName string, dynamoClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(r, tableName, dynamoClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
