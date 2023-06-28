package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeRouteTablesOutput struct {
	*ec2.DescribeRouteTablesOutput
}

func (o *DescribeRouteTablesOutput) Records() (records []*api.Record) {
	for _, s := range o.RouteTables {
		records = append(records, &api.Record{
			ID:   s.RouteTableId,
			Data: s,
		})
	}
	return
}

type DescribeRouteTables struct {
	API
}

var _ api.RequestBuilder = &DescribeRouteTables{}

// New implements api.RequestBuilder
func (fn *DescribeRouteTables) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeRouteTablesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeRouteTablesPagesWithContext(ctx, &input, func(output *ec2.DescribeRouteTablesOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeRouteTablesOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
