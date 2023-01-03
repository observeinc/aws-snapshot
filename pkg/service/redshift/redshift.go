package redshift

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/redshift"
)

func init() {
	service.Register("redshift", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeClustersPagesWithContext(context.Context, *redshift.DescribeClustersInput, func(*redshift.DescribeClustersOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	redshiftapi := redshift.New(p, opts...)
	return api.Endpoint{
		"DescribeClusters": &DescribeClusters{redshiftapi},
	}
}
