package redshift

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/redshift"
)

type DescribeClustersOutput struct {
	*redshift.DescribeClustersOutput
}

func (o *DescribeClustersOutput) Records() (records []*api.Record) {
	for _, c := range o.Clusters {
		records = append(records, &api.Record{
			// XXX: api endpoint does not return an ARN
			ID:   c.ClusterIdentifier,
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
	var input redshift.DescribeClustersInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeClustersPagesWithContext(ctx, &input, func(output *redshift.DescribeClustersOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeClustersOutput{output})
		})
	}

	return []api.Request{call}, nil
}
