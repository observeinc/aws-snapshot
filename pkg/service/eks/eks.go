package eks

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
)

func init() {
	service.Register("eks", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListClustersPagesWithContext(context.Context, *eks.ListClustersInput, func(*eks.ListClustersOutput, bool) bool, ...request.Option) error
	ListNodegroupsPagesWithContext(context.Context, *eks.ListNodegroupsInput, func(*eks.ListNodegroupsOutput, bool) bool, ...request.Option) error
	ListFargateProfilesPagesWithContext(context.Context, *eks.ListFargateProfilesInput, func(*eks.ListFargateProfilesOutput, bool) bool, ...request.Option) error
	DescribeClusterWithContext(context.Context, *eks.DescribeClusterInput, ...request.Option) (*eks.DescribeClusterOutput, error)
	DescribeNodegroupWithContext(context.Context, *eks.DescribeNodegroupInput, ...request.Option) (*eks.DescribeNodegroupOutput, error)
	DescribeFargateProfileWithContext(context.Context, *eks.DescribeFargateProfileInput, ...request.Option) (*eks.DescribeFargateProfileOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	eksapi := eks.New(p, opts...)
	return api.Endpoint{
		"DescribeCluster":        &DescribeCluster{eksapi},
		"DescribeNodegroup":      &DescribeNodegroup{eksapi},
		"DescribeFargateProfile": &DescribeFargateProfile{eksapi},
	}
}
