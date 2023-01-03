package ecs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type DescribeServicesOutput struct {
	*ecs.DescribeServicesOutput
}

func (o *DescribeServicesOutput) Records() (records []*api.Record) {
	for _, s := range o.Services {
		records = append(records, &api.Record{
			ID:   s.ServiceArn,
			Data: s,
		})
	}
	return
}

type DescribeServices struct {
	API
}

var _ api.RequestBuilder = &DescribeServices{}

// New implements api.RequestBuilder
func (fn *DescribeServices) New(name string, config interface{}) ([]api.Request, error) {

	var input ecs.ListClustersInput

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.ListClustersPagesWithContext(ctx, &input, func(output *ecs.ListClustersOutput, last bool) bool {
			for _, clusterArn := range output.ClusterArns {

				// we can now describe up to 10 services per nested page
				listServicesInput := &ecs.ListServicesInput{
					Cluster:    clusterArn,
					MaxResults: aws.Int64(10),
				}

				// run nested query
				err := fn.ListServicesPagesWithContext(ctx, listServicesInput, func(output *ecs.ListServicesOutput, last bool) bool {
					if len(output.ServiceArns) == 0 {
						return true
					}

					describeServicesInput := &ecs.DescribeServicesInput{
						Cluster:  clusterArn,
						Services: output.ServiceArns,
					}

					describeServicesOutput, err := fn.DescribeServicesWithContext(ctx, describeServicesInput)
					if err != nil {
						panic(err)
					}
					return api.SendRecords(ctx, ch, name, &DescribeServicesOutput{describeServicesOutput})
				})
				if err != nil {
					panic(err)
				}
			}
			return true
		})
	}

	return []api.Request{call}, nil
}
