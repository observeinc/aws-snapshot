package dynamodb

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func init() {
	service.Register("dynamodb", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListTablesPagesWithContext(context.Context, *dynamodb.ListTablesInput, func(*dynamodb.ListTablesOutput, bool) bool, ...request.Option) error
	DescribeTableWithContext(context.Context, *dynamodb.DescribeTableInput, ...request.Option) (*dynamodb.DescribeTableOutput, error)
	ScanPagesWithContext(context.Context, *dynamodb.ScanInput, func(*dynamodb.ScanOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	dynamodbapi := dynamodb.New(p, opts...)
	return api.Endpoint{
		"DescribeTable": &DescribeTable{dynamodbapi},
		"Scan":          &Scan{dynamodbapi},
	}
}
