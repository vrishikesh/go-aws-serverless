package user

import (
	"encoding/json"
	"fmt"
	"go-aws-serverless/pkg/validators"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email string, tableName string, dynamoClient *dynamodb.DynamoDB) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("unable to get record: %w", err)
	}

	user := User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal record: %w", err)
	}

	return &user, nil
}

func FetchUsers(tableName string, dynamoClient *dynamodb.DynamoDB) ([]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoClient.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan records")
	}

	users := []User{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal map")
	}

	return users, nil
}

func CreateUser(r events.APIGatewayProxyRequest, tableName string, dynamoClient *dynamodb.DynamoDB) (*User, error) {
	var u User

	if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
		return nil, fmt.Errorf("unable to unmarshal body")
	}

	if !validators.IsEmailValid(u.Email) {
		return nil, fmt.Errorf("invalid email")
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, fmt.Errorf("could not marshal item")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return nil, fmt.Errorf("could not put item")
	}

	return &u, nil
}

func UpdateUser(r events.APIGatewayProxyRequest, tableName string, dynamoClient *dynamodb.DynamoDB) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
		return nil, fmt.Errorf("could not unmarshal body")
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynamoClient)
	if currentUser == nil || len(currentUser.Email) == 0 {
		return nil, fmt.Errorf("user does not exist")
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal item")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return nil, fmt.Errorf("could not update item")
	}

	return &u, nil
}

func DeleteUser(r events.APIGatewayProxyRequest, tableName string, dynamoClient *dynamodb.DynamoDB) error {
	email := r.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("could not delete item")
	}

	return nil
}
