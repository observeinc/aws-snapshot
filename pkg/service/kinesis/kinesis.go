package kinesis

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func init() {
	service.Register("kinesis", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListStreamsPagesWithContext(context.Context, *kinesis.ListStreamsInput, func(*kinesis.ListStreamsOutput, bool) bool, ...request.Option) error
	DescribeStreamPagesWithContext(context.Context, *kinesis.DescribeStreamInput, func(*kinesis.DescribeStreamOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	kinesisapi := kinesis.New(p, opts...)
	return api.Endpoint{
		"DescribeStream": &DescribeStreams{kinesisapi},
	}
}
