package elasticbeanstalk

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

type DescribeApplicationsOutput struct {
	*elasticbeanstalk.DescribeApplicationsOutput
}

func (o *DescribeApplicationsOutput) Records() (records []*api.Record) {
	for _, r := range o.DescribeApplicationsOutput.Applications {
		records = append(records, &api.Record{
			ID:   r.ApplicationName,
			Data: r,
		})
	}
	return
}

type DescribeApplications struct {
	API
}

var _ api.RequestBuilder = &DescribeApplications{}

func (fn *DescribeApplications) New(name string, config interface{}) ([]api.Request, error) {
	var input elasticbeanstalk.DescribeApplicationsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		// AWS has a quota of 75 Applications by default
		output, err := fn.DescribeApplicationsWithContext(ctx, &input)
		if err != nil {
			return err
		}

		return api.SendRecords(ctx, ch, name, &DescribeApplicationsOutput{output})
	}

	return []api.Request{call}, nil
}
