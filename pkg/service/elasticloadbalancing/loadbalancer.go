package elasticloadbalancing

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type DescribeLoadBalancersOutput struct {
	*elb.DescribeLoadBalancersOutput
	*elb.DescribeTagsOutput
}

func (o *DescribeLoadBalancersOutput) Records() (records []*api.Record) {
	tags := make(map[string][]*elb.Tag)
	if o.DescribeTagsOutput != nil {
		for _, desc := range o.DescribeTagsOutput.TagDescriptions {
			tags[*desc.LoadBalancerName] = desc.Tags
		}
	}

	for _, lb := range o.LoadBalancerDescriptions {
		records = append(records, &api.Record{
			ID: lb.LoadBalancerName,
			Data: struct {
				*elb.LoadBalancerDescription
				Tags []*elb.Tag
			}{
				LoadBalancerDescription: lb,
				Tags:                    tags[*lb.LoadBalancerName],
			},
		})
	}
	return
}

type DescribeLoadBalancersOutputV2 struct {
	*elbv2.DescribeLoadBalancersOutput
	*elbv2.DescribeTagsOutput
}

func (o *DescribeLoadBalancersOutputV2) Records() (records []*api.Record) {
	tags := make(map[string][]*elbv2.Tag)
	if o.DescribeTagsOutput != nil {
		for _, desc := range o.DescribeTagsOutput.TagDescriptions {
			tags[*desc.ResourceArn] = desc.Tags
		}
	}

	for _, lb := range o.LoadBalancers {
		records = append(records, &api.Record{
			ID: lb.LoadBalancerArn,
			Data: struct {
				*elbv2.LoadBalancer
				Tags []*elbv2.Tag
			}{
				LoadBalancer: lb,
				Tags:         tags[*lb.LoadBalancerArn],
			},
		})
	}
	return
}

type DescribeLoadBalancers struct {
	ELBv2
	ELB
}

var _ api.RequestBuilder = &DescribeLoadBalancers{}

// New implements api.RequestBuilder
func (fn *DescribeLoadBalancers) New(name string, config interface{}) ([]api.Request, error) {
	if config != nil {
		return nil, fmt.Errorf("action %s is not configurable", name)
	}

	var (
		// DescribeTags takes at most 20 ARNs / names
		input   = elb.DescribeLoadBalancersInput{PageSize: aws.Int64(20)}
		inputv2 = elbv2.DescribeLoadBalancersInput{PageSize: aws.Int64(20)}
	)

	v1call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countLoadBalancers int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ELB.DescribeLoadBalancersPagesWithContext(ctx, &input, func(output *elb.DescribeLoadBalancersOutput, last bool) bool {
			if r.Stats {
				countLoadBalancers += len(output.LoadBalancerDescriptions)
			} else {
				var describeTagsInput elb.DescribeTagsInput
				for _, loadBalancerDescription := range output.LoadBalancerDescriptions {
					describeTagsInput.LoadBalancerNames = append(describeTagsInput.LoadBalancerNames, loadBalancerDescription.LoadBalancerName)
				}

				var describeTagsOutput *elb.DescribeTagsOutput
				var err error

				if len(describeTagsInput.LoadBalancerNames) > 0 {
					describeTagsOutput, err = fn.ELB.DescribeTagsWithContext(ctx, &describeTagsInput)

					if err != nil {
						innerErr = err
						return false
					}
				}

				if err = api.SendRecords(ctx, ch, name, &DescribeLoadBalancersOutput{
					DescribeLoadBalancersOutput: output,
					DescribeTagsOutput:          describeTagsOutput,
				}); err != nil {
					innerErr = err
					return false
				}
			}
			if outerErr == nil && r.Stats {
				innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countLoadBalancers})
			}
			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	v2call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countLoadBalancers int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ELBv2.DescribeLoadBalancersPagesWithContext(ctx, &inputv2, func(output *elbv2.DescribeLoadBalancersOutput, last bool) bool {
			var describeTagsInput elbv2.DescribeTagsInput
			if r.Stats {
				countLoadBalancers += len(output.LoadBalancers)
			} else {
				for _, loadBalancer := range output.LoadBalancers {
					describeTagsInput.ResourceArns = append(describeTagsInput.ResourceArns, loadBalancer.LoadBalancerArn)
				}

				var describeTagsOutput *elbv2.DescribeTagsOutput
				var err error

				if len(describeTagsInput.ResourceArns) > 0 {
					describeTagsOutput, err = fn.ELBv2.DescribeTagsWithContext(ctx, &describeTagsInput)

					if err != nil {
						innerErr = err
						return false
					}
				}

				if err = api.SendRecords(ctx, ch, name, &DescribeLoadBalancersOutputV2{
					DescribeLoadBalancersOutput: output,
					DescribeTagsOutput:          describeTagsOutput,
				}); err != nil {
					innerErr = err
					return false
				}
			}

			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countLoadBalancers})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{v1call, v2call}, nil
}
