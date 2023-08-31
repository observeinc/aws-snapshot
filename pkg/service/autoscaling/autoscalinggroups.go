package autoscaling

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type DescribeAutoScalingGroupsOutput struct {
	*autoscaling.DescribeAutoScalingGroupsOutput
}

func (o *DescribeAutoScalingGroupsOutput) Records() (records []*api.Record) {
	for _, o := range o.AutoScalingGroups {
		records = append(records, &api.Record{
			ID:   o.AutoScalingGroupARN,
			Data: o,
		})
	}
	return
}

type DescribeAutoScalingGroups struct {
	API
}

var _ api.RequestBuilder = &DescribeAutoScalingGroups{}

// New implements api.RequestBuilder
func (fn *DescribeAutoScalingGroups) New(name string, config interface{}) ([]api.Request, error) {
	var input autoscaling.DescribeAutoScalingGroupsInput

	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countAutoScaleGroups int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeAutoScalingGroupsPagesWithContext(ctx, &input, func(output *autoscaling.DescribeAutoScalingGroupsOutput, last bool) bool {

			if r.Stats {
				countAutoScaleGroups += len(output.AutoScalingGroups)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeAutoScalingGroupsOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countAutoScaleGroups})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
