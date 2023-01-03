package sqs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func init() {
	service.Register("sqs", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListQueuesPagesWithContext(context.Context, *sqs.ListQueuesInput, func(*sqs.ListQueuesOutput, bool) bool, ...request.Option) error
	GetQueueAttributesWithContext(context.Context, *sqs.GetQueueAttributesInput, ...request.Option) (*sqs.GetQueueAttributesOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	sqsapi := sqs.New(p, opts...)
	return api.Endpoint{
		"GetQueueAttributes": &GetQueueAttributes{sqsapi},
	}
}
