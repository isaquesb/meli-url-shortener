package dynamoDb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func NewRepository(client *dynamodb.Client) *Repository {
	return &Repository{
		client: client,
	}
}

type Repository struct {
	client *dynamodb.Client
}

func (r *Repository) getById(ctx context.Context, short string) (*dynamodb.GetItemOutput, error) {
	return r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("urls"),
		Key: map[string]types.AttributeValue{
			"short": &types.AttributeValueMemberS{Value: short},
		},
	})
}

func (r *Repository) UrlFromShort(ctx context.Context, short string) (string, error) {
	item, err := r.getById(ctx, short)

	if err != nil || item.Item == nil {
		return "", err
	}

	url := item.Item["url"]
	return url.(*types.AttributeValueMemberS).Value, nil
}

func (r *Repository) StatsFromShort(ctx context.Context, short string) (map[string]interface{}, error) {
	item, err := r.getById(ctx, short)

	if err != nil || item.Item == nil {
		return nil, err
	}

	visits := item.Item["visits"]

	return map[string]interface{}{
		"visits": visits.(*types.AttributeValueMemberN).Value,
	}, nil
}
