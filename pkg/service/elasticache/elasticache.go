package elasticache

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elasticache"
)

func init() {
	service.Register("elasticache", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeCacheClustersPagesWithContext(context.Context, *elasticache.DescribeCacheClustersInput, func(*elasticache.DescribeCacheClustersOutput, bool) bool, ...request.Option) error
	DescribeReplicationGroupsPagesWithContext(context.Context, *elasticache.DescribeReplicationGroupsInput, func(*elasticache.DescribeReplicationGroupsOutput, bool) bool, ...request.Option) error
	DescribeCacheSubnetGroupsPagesWithContext(context.Context, *elasticache.DescribeCacheSubnetGroupsInput, func(*elasticache.DescribeCacheSubnetGroupsOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	elasticacheapi := elasticache.New(p, opts...)
	return api.Endpoint{
		"DescribeCacheClusters":     &DescribeCacheClusters{elasticacheapi},
		"DescribeReplicationGroups": &DescribeReplicationGroups{elasticacheapi},
		"DescribeCacheSubnetGroups": &DescribeCacheSubnetGroups{elasticacheapi},
	}
}
