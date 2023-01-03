package elasticbeanstalk

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

type DescribeEnvironmentHealthOutput struct {
	*elasticbeanstalk.DescribeEnvironmentHealthOutput
}

func (o *DescribeEnvironmentHealthOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.DescribeEnvironmentHealthOutput.EnvironmentName,
		Data: o.DescribeEnvironmentHealthOutput,
	})
	return
}

type DescribeEnvironmentHealth struct {
	API
}

var _ api.RequestBuilder = &DescribeEnvironmentHealth{}

func (fn *DescribeEnvironmentHealth) New(name string, config interface{}) ([]api.Request, error) {
	var envsInput elasticbeanstalk.DescribeEnvironmentsInput
	if err := api.DecodeConfig(config, &envsInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		// AWS has a quota of 200 Environments by default

		envsOutput, err := fn.DescribeEnvironmentsWithContext(ctx, &envsInput)
		if err != nil {
			panic(err)
		}
		for _, env := range envsOutput.Environments {
			healthInput := elasticbeanstalk.DescribeEnvironmentHealthInput{
				EnvironmentId:  env.EnvironmentId,
				AttributeNames: []*string{aws.String("All")},
			}
			healthOutput, err := fn.DescribeEnvironmentHealthWithContext(ctx, &healthInput)
			if err != nil {
				panic(err)
			}

			_ = api.SendRecords(ctx, ch, name, &DescribeEnvironmentHealthOutput{healthOutput})
		}

		return nil
	}

	return []api.Request{call}, nil
}
