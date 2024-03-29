package efs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/efs"
)

type DescribeAccessPointsOutput struct {
	*efs.DescribeAccessPointsOutput
}

func (o *DescribeAccessPointsOutput) Records() (records []*api.Record) {
	for _, p := range o.AccessPoints {
		records = append(records, &api.Record{
			ID:   p.AccessPointId,
			Data: p,
		})
	}
	return
}

type DescribeAccessPoints struct {
	API
}

var _ api.RequestBuilder = &DescribeAccessPoints{}

// New implements api.RequestBuilder
func (fn *DescribeAccessPoints) New(name string, config interface{}) ([]api.Request, error) {
	var input efs.DescribeAccessPointsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeAccessPointsPagesWithContext(ctx, &input, func(output *efs.DescribeAccessPointsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeAccessPointsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
