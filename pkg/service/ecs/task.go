package ecs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type DescribeTasksOutput struct {
	*ecs.DescribeTasksOutput
}

func (o *DescribeTasksOutput) Records() (records []*api.Record) {
	for _, s := range o.Tasks {
		records = append(records, &api.Record{
			ID:   s.TaskArn,
			Data: s,
		})
	}
	return
}

type DescribeTasks struct {
	API
}

var _ api.RequestBuilder = &DescribeTasks{}

// New implements api.RequestBuilder
func (fn *DescribeTasks) New(name string, config interface{}) ([]api.Request, error) {

	var input ecs.ListClustersInput

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.ListClustersPagesWithContext(ctx, &input, func(output *ecs.ListClustersOutput, last bool) bool {
			for _, clusterArn := range output.ClusterArns {
				// we can now describe up to 10 tasks per nested page
				listTasksInput := &ecs.ListTasksInput{
					Cluster:    clusterArn,
					MaxResults: aws.Int64(10),
				}

				// run nested query
				err := fn.ListTasksPagesWithContext(ctx, listTasksInput, func(output *ecs.ListTasksOutput, last bool) bool {
					if len(output.TaskArns) == 0 {
						return true
					}

					describeTasksInput := &ecs.DescribeTasksInput{
						Cluster: clusterArn,
						Tasks:   output.TaskArns,
					}

					describeTasksOutput, err := fn.DescribeTasksWithContext(ctx, describeTasksInput)
					if err != nil {
						innerErr = err
						return false
					}

					if err := api.SendRecords(ctx, ch, name, &DescribeTasksOutput{describeTasksOutput}); err != nil {
						innerErr = err
						return false
					}

					return true
				})

				if innerErr = api.FirstError(err, innerErr); innerErr != nil {
					return false
				}
			}
			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
