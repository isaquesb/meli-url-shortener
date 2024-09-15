package config

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/isaquesb/meli-url-shortener/pkg/dynamoDb"
)

var dynamoClient *dynamodb.Client

func dynamoDbClient() (*dynamodb.Client, error) {
	var err error
	if dynamoClient == nil {
		dynamoClient, err = dynamoDb.NewClient(dynamoDb.DispatcherOptions{
			Region: GetEnv("DYNAMODB_REGION", "local"),
			Host:   GetEnv("DYNAMODB_HOST", "http://dynamodb-local:8000"),
		})
	}

	if err != nil {
		return nil, err
	}

	return dynamoClient, err
}
