package efs

import (
	"context"
	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/efs"
)

type DescribeFileSystemsOutput struct {
	*efs.DescribeFileSystemsOutput
}

func (o *DescribeFileSystemsOutput) Records() (records []*api.Record) {
	for _, f := range o.FileSystems {
		records = append(records, &api.Record{
			ID:   f.FileSystemId,
			Data: f,
		})
	}
	return
}

type DescribeFileSystems struct {
	API
}

var _ api.RequestBuilder = &DescribeFileSystems{}

// New implements api.RequestBuilder
func (fn *DescribeFileSystems) New(name string, config interface{}) ([]api.Request, error) {
	var input efs.DescribeFileSystemsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countFileSystems int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeFileSystemsPagesWithContext(ctx, &input, func(output *efs.DescribeFileSystemsOutput, last bool) bool {
			if r.Stats {
				countFileSystems += len(output.FileSystems)
			} else {
				if innerErr := api.SendRecords(ctx, ch, name, &DescribeFileSystemsOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})
		if r.Stats {
			innerErr := api.SendRecords(ctx, ch, name, &api.CountRecords{countFileSystems})
			if innerErr != nil {
				return innerErr
			}
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
