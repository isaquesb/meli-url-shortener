package dynamoDb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var Client *dynamodb.Client

func NewClient(opts DispatcherOptions) (*dynamodb.Client, error) {
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
