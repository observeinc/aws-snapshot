package firehose

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/firehose"
)

func init() {
	service.Register("firehose", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListDeliveryStreamsWithContext(context.Context, *firehose.ListDeliveryStreamsInput, ...request.Option) (*firehose.ListDeliveryStreamsOutput, error)
	DescribeDeliveryStreamWithContext(context.Context, *firehose.DescribeDeliveryStreamInput, ...request.Option) (*firehose.DescribeDeliveryStreamOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	firehoseapi := firehose.New(p, opts...)
	return api.Endpoint{
		"DescribeDeliveryStream": &DescribeDeliveryStreams{firehoseapi},
	}
}
