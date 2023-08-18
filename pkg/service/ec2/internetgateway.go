package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeInternetGatewaysOutput struct {
	*ec2.DescribeInternetGatewaysOutput
}

type CountInternetGatewaysOutput struct {
	Count int 	`json:Count`
}

func (o *CountInternetGatewaysOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		Data: o,
	})
	return
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
		var outerErr, innerErr error
		r, _:= ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeInternetGatewaysPagesWithContext(ctx, &input, func(output *ec2.DescribeInternetGatewaysOutput, last bool) bool {
			if r.Stats {
				if err := api.SendRecords(ctx, ch, name, &CountInternetGatewaysOutput{len(output.InternetGateways)}); err != nil {
					innerErr = err
					return false
				}
			} else if err := api.SendRecords(ctx, ch, name, &DescribeInternetGatewaysOutput{output}); err != nil {
				innerErr = err
				return false
			}
			return true
		})
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
