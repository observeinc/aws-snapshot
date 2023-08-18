package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeSecurityGroupsOutput struct {
	*ec2.DescribeSecurityGroupsOutput
}

func (o *DescribeSecurityGroupsOutput) Records() (records []*api.Record) {
	for _, s := range o.SecurityGroups {
		records = append(records, &api.Record{
			ID:   s.GroupId,
			Data: s,
		})
	}
	return
}

type DescribeSecurityGroups struct {
	API
}

var _ api.RequestBuilder = &DescribeSecurityGroups{}

// New implements api.RequestBuilder
func (fn *DescribeSecurityGroups) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeSecurityGroupsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		groupCount:=0
		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeSecurityGroupsPagesWithContext(ctx, &input, func(output *ec2.DescribeSecurityGroupsOutput, last bool) bool {
			if r.Stats {
				groupCount += len(output.SecurityGroups)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeSecurityGroupsOutput{output}); innerErr != nil {
					return false
				}
			}
			return true
		})
		if outerErr != nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{groupCount})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
