package efs

import (
	"context"
	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/efs"
)

type DescribeReplicationConfigurationsOutput struct {
	*efs.DescribeReplicationConfigurationsOutput
}

func (o *DescribeReplicationConfigurationsOutput) Records() (records []*api.Record) {
	for _, r := range o.Replications {
		records = append(records, &api.Record{
			ID:   r.SourceFileSystemId,
			Data: r,
		})
	}
	return
}

type DescribeReplicationConfigurations struct {
	API
}

var _ api.RequestBuilder = &DescribeReplicationConfigurations{}

// New implements api.RequestBuilder
func (fn *DescribeReplicationConfigurations) New(name string, config interface{}) ([]api.Request, error) {
	var fsInput efs.DescribeFileSystemsInput
	var rcInput efs.DescribeReplicationConfigurationsInput
	if err := api.DecodeConfig(config, &rcInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeFileSystemsPagesWithContext(ctx, &fsInput, func(output *efs.DescribeFileSystemsOutput, last bool) bool {
			for _, fs := range output.FileSystems {
				rcInput.FileSystemId = fs.FileSystemId
				output, err := fn.DescribeReplicationConfigurationsWithContext(ctx, &rcInput)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok && aerr.Code() == efs.ErrCodeReplicationNotFound {
						continue
					}
					panic(err)
				}
				api.SendRecords(ctx, ch, name, &DescribeReplicationConfigurationsOutput{output})
			}
			return true
		})
	}

	return []api.Request{call}, nil
}
