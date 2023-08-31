package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeAddressesOutput struct {
	*ec2.DescribeAddressesOutput
}

func (o *DescribeAddressesOutput) Records() (records []*api.Record) {
	for _, a := range o.Addresses {
		records = append(records, &api.Record{
			ID:   a.AllocationId,
			Data: a,
		})
	}
	return
}

type DescribeAddresses struct {
	API
}

var _ api.RequestBuilder = &DescribeAddresses{}

// New implements api.RequestBuilder
func (fn *DescribeAddresses) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeAddressesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}
	call := func(ctx context.Context, ch chan<- *api.Record) error {
		r, _ := ctx.Value("runner_config").(api.Runner)

		output, err := fn.DescribeAddressesWithContext(ctx, &input)
		if err != nil {
			return err
		}
		if r.Stats {
			// Dont use intermediate variable because there is only one call
			return api.SendRecords(ctx, ch, name, &api.CountRecords{Count: len(output.Addresses)})
		} else {
			return api.SendRecords(ctx, ch, name, &DescribeAddressesOutput{output})
		}
	}

	return []api.Request{call}, nil
}
