package rds

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/rds"
)

type DescribeDBClustersOutput struct {
	*rds.DescribeDBClustersOutput
}

func (o *DescribeDBClustersOutput) Records() (records []*api.Record) {
	for _, c := range o.DBClusters {
		records = append(records, &api.Record{
			ID:   c.DBClusterArn,
			Data: c,
		})
	}
	return
}

type DescribeDBClusters struct {
	API
}

var _ api.RequestBuilder = &DescribeDBClusters{}

// New implements api.RequestBuilder
func (fn *DescribeDBClusters) New(name string, config interface{}) ([]api.Request, error) {
	var input rds.DescribeDBClustersInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countDBClusters int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeDBClustersPagesWithContext(ctx, &input, func(output *rds.DescribeDBClustersOutput, last bool) bool {
			if r.Stats {
				countDBClusters += len(output.DBClusters)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeDBClustersOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})

		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countDBClusters})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
