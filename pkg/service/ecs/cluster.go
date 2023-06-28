package ecs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ecs"
)

type DescribeClustersOutput struct {
	*ecs.DescribeClustersOutput
}

func (o *DescribeClustersOutput) Records() (records []*api.Record) {
	for _, c := range o.Clusters {
		records = append(records, &api.Record{
			ID:   c.ClusterArn,
			Data: c,
		})
	}
	return
}

type DescribeClusters struct {
	API
}

var _ api.RequestBuilder = &DescribeClusters{}

// New implements api.RequestBuilder
func (fn *DescribeClusters) New(name string, config interface{}) ([]api.Request, error) {
	var input ecs.ListClustersInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var innerErr, outerErr error

		outerErr = fn.ListClustersPagesWithContext(ctx, &input, func(output *ecs.ListClustersOutput, last bool) bool {
			describeClustersInput := &ecs.DescribeClustersInput{
				Clusters: output.ClusterArns,
			}

			describeClustersOutput, err := fn.DescribeClustersWithContext(ctx, describeClustersInput)
			if err != nil {
				innerErr = err
				return false
			}

			if err := api.SendRecords(ctx, ch, name, &DescribeClustersOutput{describeClustersOutput}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
