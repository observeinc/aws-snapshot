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

type CountNatGatewaysOutput struct {
	Count int `json:Count`
}

func (o *CountNatGatewaysOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		Data: o,
	})
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
		var outerErr, innerErr error
		var gatewayCount int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeNatGatewaysPagesWithContext(ctx, &input, func(output *ec2.DescribeNatGatewaysOutput, last bool) bool {
			if r.Stats {
				gatewayCount += len(output.NatGateways)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeNatGatewaysOutput{output}); innerErr != nil {
					return false
				}
			}
			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &CountNatGatewaysOutput{gatewayCount})
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
