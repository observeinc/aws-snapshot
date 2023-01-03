package eks

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/eks"
)

type DescribeClusterOutput struct {
	*eks.DescribeClusterOutput
}

func (o *DescribeClusterOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.Cluster.Arn,
		Data: o,
	})
	return records
}

type DescribeCluster struct {
	API
}

var _ api.RequestBuilder = &DescribeCluster{}

// New implements api.RequestBuilder
func (fn *DescribeCluster) New(name string, config interface{}) ([]api.Request, error) {
	var input eks.ListClustersInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.ListClustersPagesWithContext(ctx, &input, func(output *eks.ListClustersOutput, last bool) bool {
			for _, clusterName := range output.Clusters {
				describeClusterInput := &eks.DescribeClusterInput{
					Name: clusterName,
				}

				describeClusterOutput, err := fn.DescribeClusterWithContext(ctx, describeClusterInput)
				if err != nil {
					panic(err)
				}
				ok := api.SendRecords(ctx, ch, name, &DescribeClusterOutput{describeClusterOutput})
				if !ok {
					return false
				}
			}
			return true
		})
	}

	return []api.Request{call}, nil
}
