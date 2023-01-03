package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeInternetGatewaysOutput struct {
	*ec2.DescribeInternetGatewaysOutput
}

func (o *DescribeInternetGatewaysOutput) Records() (records []*api.Record) {
	for _, s := range o.InternetGateways {
		records = append(records, &api.Record{
			ID:   s.InternetGatewayId,
			Data: s,
		})
	}
	return
}

type DescribeInternetGateways struct {
	API
}

var _ api.RequestBuilder = &DescribeInternetGateways{}

// New implements api.RequestBuilder
func (fn *DescribeInternetGateways) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeInternetGatewaysInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeInternetGatewaysPagesWithContext(ctx, &input, func(output *ec2.DescribeInternetGatewaysOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeInternetGatewaysOutput{output})
		})
	}

	return []api.Request{call}, nil
}
