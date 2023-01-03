package cloudformation

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func init() {
	service.Register("cloudformation", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeStacksPagesWithContext(context.Context, *cloudformation.DescribeStacksInput, func(*cloudformation.DescribeStacksOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	cloudformationapi := cloudformation.New(p, opts...)
	return api.Endpoint{
		"DescribeStacks": &DescribeStacks{cloudformationapi},
	}
}
