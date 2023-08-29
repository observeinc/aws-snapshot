package rds

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/rds"
)

type DescribeDBInstancesOutput struct {
	*rds.DescribeDBInstancesOutput
}

func (o *DescribeDBInstancesOutput) Records() (records []*api.Record) {
	for _, i := range o.DBInstances {
		records = append(records, &api.Record{
			ID:   i.DBInstanceArn,
			Data: i,
		})
	}
	return
}

type DescribeDBInstances struct {
	API
}

var _ api.RequestBuilder = &DescribeDBInstances{}

// New implements api.RequestBuilder
func (fn *DescribeDBInstances) New(name string, config interface{}) ([]api.Request, error) {
	var input rds.DescribeDBInstancesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countDBInstances int
		r, _ := ctx.Value("runner_config").(api.Runner)

		outerErr = fn.DescribeDBInstancesPagesWithContext(ctx, &input, func(output *rds.DescribeDBInstancesOutput, last bool) bool {
			if r.Stats {
				countDBInstances += len(output.DBInstances)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeDBInstancesOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countDBInstances})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
