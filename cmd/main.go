package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"go-aws-serverless/pkg/handlers"
)

const tableName = "users"
var dynamoClient *dynamodb.DynamoDB

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalln(err)
	}

	dynamoClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(r events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch r.HTTPMethod {
	case http.MethodGet:
		return handlers.GetUser(r, tableName, dynamoClient)
	case http.MethodPost:
		return handlers.NewUser(r, tableName, dynamoClient)
	case http.MethodPut:
		return handlers.UpdateUser(r, tableName, dynamoClient)
	case http.MethodDelete:
		return handlers.DeleteUser(r, tableName, dynamoClient)
	default:
		return handlers.UnhandledMethod()
	}
}
