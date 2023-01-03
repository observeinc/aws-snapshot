package elasticloadbalancing

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func init() {
	service.Register("elasticloadbalancing", api.ServiceFunc(New))
}

type ELB interface {
	DescribeLoadBalancersPagesWithContext(context.Context, *elb.DescribeLoadBalancersInput, func(*elb.DescribeLoadBalancersOutput, bool) bool, ...request.Option) error
	DescribeTagsWithContext(context.Context, *elb.DescribeTagsInput, ...request.Option) (*elb.DescribeTagsOutput, error)
}

type ELBv2 interface {
	DescribeLoadBalancersPagesWithContext(context.Context, *elbv2.DescribeLoadBalancersInput, func(*elbv2.DescribeLoadBalancersOutput, bool) bool, ...request.Option) error
	DescribeTargetGroupsPagesWithContext(context.Context, *elbv2.DescribeTargetGroupsInput, func(*elbv2.DescribeTargetGroupsOutput, bool) bool, ...request.Option) error
	DescribeTagsWithContext(context.Context, *elbv2.DescribeTagsInput, ...request.Option) (*elbv2.DescribeTagsOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	elbv2api := elbv2.New(p, opts...)
	elbapi := elb.New(p, opts...)
	return api.Endpoint{
		"DescribeLoadBalancers": &DescribeLoadBalancers{
			ELBv2: elbv2api,
			ELB:   elbapi,
		},
		"DescribeTargetGroups": &DescribeTargetGroups{elbv2api},
	}
}
