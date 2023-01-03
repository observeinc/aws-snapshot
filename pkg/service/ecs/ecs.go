package ecs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func init() {
	service.Register("ecs", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListClustersPagesWithContext(context.Context, *ecs.ListClustersInput, func(*ecs.ListClustersOutput, bool) bool, ...request.Option) error
	ListContainerInstancesPagesWithContext(context.Context, *ecs.ListContainerInstancesInput, func(*ecs.ListContainerInstancesOutput, bool) bool, ...request.Option) error
	ListServicesPagesWithContext(context.Context, *ecs.ListServicesInput, func(*ecs.ListServicesOutput, bool) bool, ...request.Option) error
	ListTasksPagesWithContext(context.Context, *ecs.ListTasksInput, func(*ecs.ListTasksOutput, bool) bool, ...request.Option) error
	DescribeClustersWithContext(context.Context, *ecs.DescribeClustersInput, ...request.Option) (*ecs.DescribeClustersOutput, error)
	DescribeContainerInstancesWithContext(context.Context, *ecs.DescribeContainerInstancesInput, ...request.Option) (*ecs.DescribeContainerInstancesOutput, error)
	DescribeServicesWithContext(context.Context, *ecs.DescribeServicesInput, ...request.Option) (*ecs.DescribeServicesOutput, error)
	DescribeTasksWithContext(context.Context, *ecs.DescribeTasksInput, ...request.Option) (*ecs.DescribeTasksOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	ecsapi := ecs.New(p, opts...)
	return api.Endpoint{
		"DescribeClusters":           &DescribeClusters{ecsapi},
		"DescribeContainerInstances": &DescribeContainerInstances{ecsapi},
		"DescribeServices":           &DescribeServices{ecsapi},
		"DescribeTasks":              &DescribeTasks{ecsapi},
	}
}
