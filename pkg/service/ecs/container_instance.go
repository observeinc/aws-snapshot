package ecs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type DescribeContainerInstancesOutput struct {
	*ecs.DescribeContainerInstancesOutput
}

func (o *DescribeContainerInstancesOutput) Records() (records []*api.Record) {
	for _, s := range o.ContainerInstances {
		records = append(records, &api.Record{
			ID:   s.ContainerInstanceArn,
			Data: s,
		})
	}
	return
}

type DescribeContainerInstances struct {
	API
}

var _ api.RequestBuilder = &DescribeContainerInstances{}

// New implements api.RequestBuilder
func (fn *DescribeContainerInstances) New(name string, config interface{}) ([]api.Request, error) {

	var input ecs.ListClustersInput

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var innerErr, outerErr error

		outerErr = fn.ListClustersPagesWithContext(ctx, &input, func(output *ecs.ListClustersOutput, last bool) bool {
			for _, clusterArn := range output.ClusterArns {
				// we can now describe up to 10 tasks per nested page
				listContainerInstancesInput := &ecs.ListContainerInstancesInput{
					Cluster:    clusterArn,
					MaxResults: aws.Int64(10),
				}

				// run nested query
				err := fn.ListContainerInstancesPagesWithContext(ctx, listContainerInstancesInput, func(output *ecs.ListContainerInstancesOutput, last bool) bool {
					if len(output.ContainerInstanceArns) == 0 {
						return true
					}

					describeContainerInstancesInput := &ecs.DescribeContainerInstancesInput{
						Cluster:            clusterArn,
						ContainerInstances: output.ContainerInstanceArns,
					}

					describeContainerInstancesOutput, err := fn.DescribeContainerInstancesWithContext(ctx, describeContainerInstancesInput)
					if err != nil {
						innerErr = err
						return false
					}

					if err := api.SendRecords(ctx, ch, name, &DescribeContainerInstancesOutput{describeContainerInstancesOutput}); err != nil {
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
