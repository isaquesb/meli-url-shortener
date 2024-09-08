package dynamoDb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"time"
)

var gSvc *dynamodb.Client

func Write(short, url string) {
	svc := GetConnection()

	tableName := "urls"

	item := map[string]types.AttributeValue{
		"short":      &types.AttributeValueMemberS{Value: short},
		"url":        &types.AttributeValueMemberS{Value: url},
		"created_at": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}
	_, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})

	if err != nil {
		log.Fatalf("Error on PutItem: %v", err)
	}
}

func GetConnection() *dynamodb.Client {
	if gSvc != nil {
		return gSvc
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"), // Região padrão (não faz muita diferença no DynamoDB Local)
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				if service == dynamodb.ServiceID {
					return aws.Endpoint{
						URL: "http://localhost:8000", // Endpoint do DynamoDB Local
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("serviço desconhecido: %s", service)
			}),
		),
	)

	if err != nil {
		log.Fatalf("Error on loading config: %v", err)
	}

	gSvc = dynamodb.NewFromConfig(cfg)

	log.Println("DynamoDB Connected!")

	return gSvc
}
