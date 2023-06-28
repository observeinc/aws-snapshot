package cloudformation

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type DescribeStacksOutput struct {
	*cloudformation.DescribeStacksOutput
}

func (o *DescribeStacksOutput) Records() (records []*api.Record) {
	for _, s := range o.Stacks {
		records = append(records, &api.Record{
			ID:   s.StackId,
			Data: s,
		})
	}
	return
}

type DescribeStacks struct {
	API
}

var _ api.RequestBuilder = &DescribeStacks{}

// New implements api.RequestBuilder
func (fn *DescribeStacks) New(name string, config interface{}) ([]api.Request, error) {
	var input cloudformation.DescribeStacksInput

	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		err := fn.DescribeStacksPagesWithContext(ctx, &input, func(output *cloudformation.DescribeStacksOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeStacksOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		if aerr, ok := err.(awserr.Error); !ok || aerr.Code() != "AccessDenied" {
			outerErr = err
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
