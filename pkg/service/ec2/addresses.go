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

type CountInstancesAddressOutput struct {
	Count int 	`json:Count`
}

func (o *CountInstancesAddressOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		Data: o,
	})
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
	//fmt.Println("Doing EC2 Stats %d", r.Stats )
	call := func(ctx context.Context, ch chan<- *api.Record) error {
		r, ok := ctx.Value("runner_config").(api.Runner)

		if !ok {
			return nil
		} else if ok && r.Stats {
			output, err := fn.DescribeAddressesWithContext(ctx, &input)
			if err != nil {
				return err
			}
			return api.SendRecords(ctx, ch, name, &CountInstancesAddressOutput{len(output.Addresses)})
		} else {
			output, err := fn.DescribeAddressesWithContext(ctx, &input)
			if err != nil {
				return err
			}

			return api.SendRecords(ctx, ch, name, &DescribeAddressesOutput{output})
		}
	}

	return []api.Request{call}, nil
}
