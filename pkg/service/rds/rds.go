package rds

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/rds"
)

func init() {
	service.Register("rds", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeDBClustersPagesWithContext(context.Context, *rds.DescribeDBClustersInput, func(*rds.DescribeDBClustersOutput, bool) bool, ...request.Option) error
	DescribeDBInstancesPagesWithContext(context.Context, *rds.DescribeDBInstancesInput, func(*rds.DescribeDBInstancesOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	rdsapi := rds.New(p, opts...)
	return api.Endpoint{
		"DescribeDBClusters":  &DescribeDBClusters{rdsapi},
		"DescribeDBInstances": &DescribeDBInstances{rdsapi},
	}
}
