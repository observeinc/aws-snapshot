package elasticbeanstalk

import (
	"context"
	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

type DescribeInstancesHealthOutput struct {
	*elasticbeanstalk.DescribeInstancesHealthOutput
	environmentId *string
}

type DescribeInstancesHealthRecord struct {
	*elasticbeanstalk.SingleInstanceHealth
	EnvironmentId *string
}

func (o *DescribeInstancesHealthOutput) Records() (records []*api.Record) {
	for _, r := range o.InstanceHealthList {
		records = append(records, &api.Record{
			ID: r.InstanceId,
			Data: DescribeInstancesHealthRecord{
				SingleInstanceHealth: r,
				EnvironmentId:        o.environmentId,
			},
		})
	}
	return
}

type DescribeInstancesHealth struct {
	API
}

var _ api.RequestBuilder = &DescribeInstancesHealth{}

func (fn *DescribeInstancesHealth) New(name string, config interface{}) ([]api.Request, error) {
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
			healthInput := elasticbeanstalk.DescribeInstancesHealthInput{
				EnvironmentId:  env.EnvironmentId,
				AttributeNames: []*string{aws.String("All")},
			}
			healthOutput, err := fn.DescribeInstancesHealthWithContext(ctx, &healthInput)
			if err != nil {
				panic(err)
			}

			source := &DescribeInstancesHealthOutput{
				DescribeInstancesHealthOutput: healthOutput,
				environmentId:                 env.EnvironmentId,
			}
			_ = api.SendRecords(ctx, ch, name, source)
		}

		return nil
	}

	return []api.Request{call}, nil
}
