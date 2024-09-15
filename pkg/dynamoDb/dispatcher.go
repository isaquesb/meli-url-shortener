package dynamoDb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/isaquesb/url-shortener/internal/events"
	"github.com/isaquesb/url-shortener/internal/ports/output"
	"github.com/isaquesb/url-shortener/internal/urls"
	"time"
)

type DispatcherOptions struct {
	Region string
	Host   string
}

type Dispatcher struct {
	client *dynamodb.Client
}

func NewDispatcher(client *dynamodb.Client) (output.Dispatcher, error) {
	return &Dispatcher{
		client: client,
	}, nil
}

func (d *Dispatcher) Dispatch(ctx context.Context, msg events.Event) error {
	switch evt := msg.(type) {
	case *urls.CreateEvent:
		return d.CreateUrl(ctx, evt)
	case *urls.VisitEvent:
		return d.IncreaseVisits(ctx, evt)
	case *urls.DeleteEvent:
		return d.DeleteUrl(ctx, evt)
	default:
		return fmt.Errorf("unhandled event type: %T", evt)
	}
}

func (d *Dispatcher) CreateUrl(ctx context.Context, evt *urls.CreateEvent) error {
	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("urls"),
		Item: map[string]types.AttributeValue{
			"short":      &types.AttributeValueMemberS{Value: string(evt.ShortCode)},
			"url":        &types.AttributeValueMemberS{Value: evt.Url},
			"visits":     &types.AttributeValueMemberN{Value: "0"},
			"created_at": &types.AttributeValueMemberS{Value: evt.CreatedAt.Format(time.RFC3339)},
		},
	})

	return err
}

func (d *Dispatcher) IncreaseVisits(ctx context.Context, evt *urls.VisitEvent) error {
	_, err := d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("urls"),
		Key: map[string]types.AttributeValue{
			"short": &types.AttributeValueMemberS{Value: string(evt.ShortCode)},
		},
		UpdateExpression: aws.String("set visits = visits + :inc"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	})

	return err
}

func (d *Dispatcher) DeleteUrl(ctx context.Context, evt *urls.DeleteEvent) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("urls"),
		Key: map[string]types.AttributeValue{
			"short": &types.AttributeValueMemberS{Value: string(evt.ShortCode)},
		},
	})
	return err
}

func (d *Dispatcher) Close() {}
