package rds

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/rds"
)

type DescribeDBInstancesOutput struct {
	*rds.DescribeDBInstancesOutput
}

func (o *DescribeDBInstancesOutput) Records() (records []*api.Record) {
	for _, i := range o.DBInstances {
		records = append(records, &api.Record{
			ID:   i.DBInstanceArn,
			Data: i,
		})
	}
	return
}

type DescribeDBInstances struct {
	API
}

var _ api.RequestBuilder = &DescribeDBInstances{}

// New implements api.RequestBuilder
func (fn *DescribeDBInstances) New(name string, config interface{}) ([]api.Request, error) {
	var input rds.DescribeDBInstancesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeDBInstancesPagesWithContext(ctx, &input, func(output *rds.DescribeDBInstancesOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeDBInstancesOutput{output})
		})
	}

	return []api.Request{call}, nil
}
