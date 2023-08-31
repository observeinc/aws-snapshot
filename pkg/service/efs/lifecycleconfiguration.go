package efs

import (
	"context"
	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/efs"
)

type DescribeLifecycleConfigurationOutput struct {
	FilesystemID *string
	*efs.DescribeLifecycleConfigurationOutput
}

func (o *DescribeLifecycleConfigurationOutput) Records() (records []*api.Record) {
	for _, l := range o.LifecyclePolicies {
		records = append(records, &api.Record{
			ID:   o.FilesystemID,
			Data: l,
		})
	}
	return
}

type DescribeLifecycleConfiguration struct {
	API
}

var _ api.RequestBuilder = &DescribeLifecycleConfiguration{}

// New implements api.RequestBuilder
func (fn *DescribeLifecycleConfiguration) New(name string, config interface{}) ([]api.Request, error) {
	var fsInput efs.DescribeFileSystemsInput
	var mtInput efs.DescribeLifecycleConfigurationInput
	if err := api.DecodeConfig(config, &mtInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countLifeCyclePolices int
		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeFileSystemsPagesWithContext(ctx, &fsInput, func(output *efs.DescribeFileSystemsOutput, last bool) bool {
			if r.Stats {
				// assumption we have 1 lifecycle per fs?
				countLifeCyclePolices += len(output.FileSystems)
			} else {
				for _, fs := range output.FileSystems {
					mtInput.FileSystemId = fs.FileSystemId
					output, err := fn.DescribeLifecycleConfigurationWithContext(ctx, &mtInput)
					if err != nil {
						innerErr = err
						return false
					}

					if err := api.SendRecords(ctx, ch, name, &DescribeLifecycleConfigurationOutput{FilesystemID: fs.FileSystemId, DescribeLifecycleConfigurationOutput: output}); err != nil {
						innerErr = err
						return false
					}
				}
			}
			return true
		})
		if r.Stats {
			innerErr := api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countLifeCyclePolices})
			if innerErr != nil {
				return innerErr
			}
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
