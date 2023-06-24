package efs

import (
	"context"
	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/efs"
)

type DescribeMountTargetsOutput struct {
	*efs.DescribeMountTargetsOutput
}

func (o *DescribeMountTargetsOutput) Records() (records []*api.Record) {
	for _, mt := range o.MountTargets {
		records = append(records, &api.Record{
			ID:   mt.MountTargetId,
			Data: mt,
		})
	}
	return
}

type DescribeMountTargets struct {
	API
}

var _ api.RequestBuilder = &DescribeMountTargets{}

// New implements api.RequestBuilder
func (fn *DescribeMountTargets) New(name string, config interface{}) ([]api.Request, error) {
	var fsInput efs.DescribeFileSystemsInput
	var mtInput efs.DescribeMountTargetsInput
	if err := api.DecodeConfig(config, &mtInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeFileSystemsPagesWithContext(ctx, &fsInput, func(output *efs.DescribeFileSystemsOutput, last bool) bool {
			for _, fs := range output.FileSystems {
				mtInput.FileSystemId = fs.FileSystemId
				output, err := fn.DescribeMountTargetsWithContext(ctx, &mtInput)
				if err != nil {
					innerErr = err
					return false
				}

				if err := api.SendRecords(ctx, ch, name, &DescribeMountTargetsOutput{output}); err != nil {
					innerErr = err
					return false
				}
			}
			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
