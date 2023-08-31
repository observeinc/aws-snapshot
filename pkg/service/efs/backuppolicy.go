package efs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/efs"
)

type DescribeBackupPolicyOutput struct {
	FilesystemID *string
	*efs.DescribeBackupPolicyOutput
}

func (o *DescribeBackupPolicyOutput) Records() (records []*api.Record) {
	records = []*api.Record{
		{
			ID:   o.FilesystemID,
			Data: o,
		},
	}
	return
}

type DescribeBackupPolicy struct {
	API
}

var _ api.RequestBuilder = &DescribeBackupPolicy{}

// New implements api.RequestBuilder
func (fn *DescribeBackupPolicy) New(name string, config interface{}) ([]api.Request, error) {
	var fsInput efs.DescribeFileSystemsInput
	var bcInput efs.DescribeBackupPolicyInput
	if err := api.DecodeConfig(config, &bcInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countFilesystems int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeFileSystemsPagesWithContext(ctx, &fsInput, func(output *efs.DescribeFileSystemsOutput, last bool) bool {
			if r.Stats {
				countFilesystems += len(output.FileSystems)
			} else {
				for _, fs := range output.FileSystems {
					bcInput.FileSystemId = fs.FileSystemId
					output, err := fn.DescribeBackupPolicyWithContext(ctx, &bcInput)
					if err != nil {
						if aerr, ok := err.(awserr.Error); ok && aerr.Code() == efs.ErrCodePolicyNotFound {
							continue
						}

						innerErr = err
						return false
					}

					if err := api.SendRecords(ctx, ch, name, &DescribeBackupPolicyOutput{FilesystemID: fs.FileSystemId, DescribeBackupPolicyOutput: output}); err != nil {
						innerErr = err
						return false
					}
				}
			}
			return true
		})
		if r.Stats {
			innerErr := api.SendRecords(ctx, ch, name, &api.CountRecords{countFilesystems})
			if innerErr != nil {
				return innerErr
			}
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
