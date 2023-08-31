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
		var countEnvironments int
		r, _ := ctx.Value("runner_config").(api.Runner)
		for {
			output, err := fn.DescribeEnvironmentsWithContext(ctx, &input)
			if err != nil {
				return err
			}
			if r.Stats {
				countEnvironments += len(output.Environments)
			} else {
				if err := api.SendRecords(ctx, ch, name, &DescribeEnvironmentsOutput{output}); err != nil {
					return err
				}
			}

			if output.NextToken == nil {
				break
			}
			input.NextToken = output.NextToken
		}

		if r.Stats {
			innerErr := api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countEnvironments})
			if innerErr != nil {
				return innerErr
			}
		}

		return nil
	}

	return []api.Request{call}, nil
}
