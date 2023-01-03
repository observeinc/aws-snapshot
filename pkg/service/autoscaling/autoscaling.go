package autoscaling

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func init() {
	service.Register("autoscaling", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeAutoScalingGroupsPagesWithContext(context.Context, *autoscaling.DescribeAutoScalingGroupsInput, func(*autoscaling.DescribeAutoScalingGroupsOutput, bool) bool, ...request.Option) error
	DescribeLaunchConfigurationsPagesWithContext(context.Context, *autoscaling.DescribeLaunchConfigurationsInput, func(*autoscaling.DescribeLaunchConfigurationsOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	autoscalingapi := autoscaling.New(p, opts...)
	return api.Endpoint{
		"DescribeAutoScalingGroups":    &DescribeAutoScalingGroups{autoscalingapi},
		"DescribeLaunchConfigurations": &DescribeLaunchConfigurations{autoscalingapi},
	}
}
