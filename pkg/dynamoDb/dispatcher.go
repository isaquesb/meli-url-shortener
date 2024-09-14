package dynamoDb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/isaquesb/meli-url-shortener/internal/events"
	"github.com/isaquesb/meli-url-shortener/internal/ports/output"
	"github.com/isaquesb/meli-url-shortener/internal/urls"
	"time"
)

type DispatcherOptions struct {
	Region string
	Host   string
}

type Dispatcher struct {
	writer *dynamodb.Client
}

func NewDispatcher(opts DispatcherOptions) (output.Dispatcher, error) {
	connection, err := CreateConnection(opts)
	if err != nil {
		return nil, err
	}
	return &Dispatcher{
		writer: connection,
	}, nil
}

func (d *Dispatcher) Dispatch(ctx context.Context, msg events.Event) error {
	tableName := "urls"
	evt := msg.(*urls.CreateEvent)
	item := map[string]types.AttributeValue{
		"short":      &types.AttributeValueMemberS{Value: string(evt.ShortCode)},
		"url":        &types.AttributeValueMemberS{Value: string(evt.Url)},
		"created_at": &types.AttributeValueMemberS{Value: evt.CreatedAt.Format(time.RFC3339)},
	}
	_, err := d.writer.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})

	return err
}

func (d *Dispatcher) Close() {
}

func CreateConnection(opts DispatcherOptions) (*dynamodb.Client, error) {
	isLocal := opts.Region == "local"
	region := opts.Region
	if isLocal {
		region = "us-west-2"
	}

	configOpts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if isLocal {
		configOpts = append(
			configOpts,
			config.WithEndpointResolver(aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					if service == dynamodb.ServiceID {
						return aws.Endpoint{
							URL: opts.Host,
						}, nil
					}
					return aws.Endpoint{}, fmt.Errorf("unknown aws service: %s", service)
				}),
			),
		)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), configOpts...)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}
