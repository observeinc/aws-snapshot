package cloudformation

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

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
		return fn.DescribeStacksPagesWithContext(ctx, &input, func(output *cloudformation.DescribeStacksOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeStacksOutput{output})
		})
	}

	return []api.Request{call}, nil
}
