package eventbridge

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

func init() {
	service.Register("events", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListApiDestinationsWithContext(context.Context, *eventbridge.ListApiDestinationsInput, ...request.Option) (*eventbridge.ListApiDestinationsOutput, error)
	ListEventBusesWithContext(context.Context, *eventbridge.ListEventBusesInput, ...request.Option) (*eventbridge.ListEventBusesOutput, error)
	ListRulesWithContext(context.Context, *eventbridge.ListRulesInput, ...request.Option) (*eventbridge.ListRulesOutput, error)
	ListTargetsByRuleWithContext(context.Context, *eventbridge.ListTargetsByRuleInput, ...request.Option) (*eventbridge.ListTargetsByRuleOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	eventbridgeapi := eventbridge.New(p, opts...)
	return api.Endpoint{
		"ListApiDestinations": &ListApiDestinations{eventbridgeapi},
		"ListEventBuses":      &ListEventBuses{eventbridgeapi},
		"ListRules":           &ListRules{eventbridgeapi},
	}
}
