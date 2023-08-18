package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeSubnetsOutput struct {
	*ec2.DescribeSubnetsOutput
}

func (o *DescribeSubnetsOutput) Records() (records []*api.Record) {
	for _, s := range o.Subnets {
		records = append(records, &api.Record{
			ID:   s.SubnetId,
			Data: s,
		})
	}
	return
}

type DescribeSubnets struct {
	API
}

var _ api.RequestBuilder = &DescribeSubnets{}

// New implements api.RequestBuilder
func (fn *DescribeSubnets) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeSubnetsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countSubnetsOutput int
		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeSubnetsPagesWithContext(ctx, &input, func(output *ec2.DescribeSubnetsOutput, last bool) bool {
			if r.Stats {
				countSubnetsOutput += len(output.Subnets)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeSubnetsOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})

		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countSubnetsOutput})
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
