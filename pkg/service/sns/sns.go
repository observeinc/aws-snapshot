package sns

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sns"
)

func init() {
	service.Register("sns", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListTopicsPagesWithContext(context.Context, *sns.ListTopicsInput, func(*sns.ListTopicsOutput, bool) bool, ...request.Option) error
	GetTopicAttributesWithContext(context.Context, *sns.GetTopicAttributesInput, ...request.Option) (*sns.GetTopicAttributesOutput, error)
	ListSubscriptionsPagesWithContext(context.Context, *sns.ListSubscriptionsInput, func(*sns.ListSubscriptionsOutput, bool) bool, ...request.Option) error
	GetSubscriptionAttributesWithContext(context.Context, *sns.GetSubscriptionAttributesInput, ...request.Option) (*sns.GetSubscriptionAttributesOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	snsapi := sns.New(p, opts...)
	return api.Endpoint{
		"GetTopicAttributes":        &GetTopicAttributes{snsapi},
		"GetSubscriptionAttributes": &GetSubscriptionAttributes{snsapi},
	}
}
