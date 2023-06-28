package elasticbeanstalk

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

type DescribeApplicationVersionsOutput struct {
	*elasticbeanstalk.DescribeApplicationVersionsOutput
}

func (o *DescribeApplicationVersionsOutput) Records() (records []*api.Record) {
	for _, r := range o.ApplicationVersions {
		records = append(records, &api.Record{
			ID:   r.ApplicationName,
			Data: r,
		})
	}
	return
}

type DescribeApplicationVersions struct {
	API
}

var _ api.RequestBuilder = &DescribeApplicationVersions{}

func (fn *DescribeApplicationVersions) New(name string, config interface{}) ([]api.Request, error) {
	var input elasticbeanstalk.DescribeApplicationVersionsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		// AWS has a quota of 1,000 Application Versions by default
		for {
			output, err := fn.DescribeApplicationVersionsWithContext(ctx, &input)
			if err != nil {
				return err
			}

			if err := api.SendRecords(ctx, ch, name, &DescribeApplicationVersionsOutput{output}); err != nil {
				return err
			}

			if output.NextToken == nil {
				break
			}
			input.NextToken = output.NextToken
		}

		return nil
	}

	return []api.Request{call}, nil
}
