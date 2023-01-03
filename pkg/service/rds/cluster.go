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
		return fn.DescribeDBClustersPagesWithContext(ctx, &input, func(output *rds.DescribeDBClustersOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeDBClustersOutput{output})
		})
	}

	return []api.Request{call}, nil
}
