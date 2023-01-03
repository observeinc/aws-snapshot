package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeNatGatewaysOutput struct {
	*ec2.DescribeNatGatewaysOutput
}

func (o *DescribeNatGatewaysOutput) Records() (records []*api.Record) {
	for _, s := range o.NatGateways {
		records = append(records, &api.Record{
			ID:   s.NatGatewayId,
			Data: s,
		})
	}
	return
}

type DescribeNatGateways struct {
	API
}

var _ api.RequestBuilder = &DescribeNatGateways{}

// New implements api.RequestBuilder
func (fn *DescribeNatGateways) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeNatGatewaysInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeNatGatewaysPagesWithContext(ctx, &input, func(output *ec2.DescribeNatGatewaysOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeNatGatewaysOutput{output})
		})
	}

	return []api.Request{call}, nil
}
