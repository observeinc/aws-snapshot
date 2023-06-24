package efs

import (
	"context"
	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/efs"
)

type DescribeFileSystemPolicyOutput struct {
	*efs.DescribeFileSystemPolicyOutput
}

func (o *DescribeFileSystemPolicyOutput) Records() (records []*api.Record) {
	records = []*api.Record{
		{
			ID:   o.FileSystemId,
			Data: o,
		},
	}
	return
}

type DescribeFileSystemPolicy struct {
	API
}

var _ api.RequestBuilder = &DescribeFileSystemPolicy{}

// New implements api.RequestBuilder
func (fn *DescribeFileSystemPolicy) New(name string, config interface{}) ([]api.Request, error) {
	var fsInput efs.DescribeFileSystemsInput
	var fspInput efs.DescribeFileSystemPolicyInput
	if err := api.DecodeConfig(config, &fspInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeFileSystemsPagesWithContext(ctx, &fsInput, func(output *efs.DescribeFileSystemsOutput, last bool) bool {
			for _, fs := range output.FileSystems {
				fspInput.FileSystemId = fs.FileSystemId
				output, err := fn.DescribeFileSystemPolicyWithContext(ctx, &fspInput)
				if err != nil {
					if ae, ok := err.(awserr.Error); ok && ae.Code() == efs.ErrCodePolicyNotFound {
						continue
					}

					innerErr = err
					return false
				}

				if err := api.SendRecords(ctx, ch, name, &DescribeFileSystemPolicyOutput{output}); err != nil {
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
