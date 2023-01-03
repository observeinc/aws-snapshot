package elasticbeanstalk

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

type DescribeEnvironmentsOutput struct {
	*elasticbeanstalk.EnvironmentDescriptionsMessage
}

func (o *DescribeEnvironmentsOutput) Records() (records []*api.Record) {
	for _, r := range o.EnvironmentDescriptionsMessage.Environments {
		records = append(records, &api.Record{
			ID:   r.EnvironmentId,
			Data: r,
		})
	}
	return
}

type DescribeEnvironments struct {
	API
}

var _ api.RequestBuilder = &DescribeEnvironments{}

func (fn *DescribeEnvironments) New(name string, config interface{}) ([]api.Request, error) {
	var input elasticbeanstalk.DescribeEnvironmentsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		// AWS has a quota of 200 Environments by default
		for {
			output, err := fn.DescribeEnvironmentsWithContext(ctx, &input)
			if err != nil {
				panic(err)
			}

			_ = api.SendRecords(ctx, ch, name, &DescribeEnvironmentsOutput{output})

			if output.NextToken == nil {
				break
			}
			input.NextToken = output.NextToken
		}

		return nil
	}

	return []api.Request{call}, nil
}
