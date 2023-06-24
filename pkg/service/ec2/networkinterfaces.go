package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeNetworkInterfacesOutput struct {
	*ec2.DescribeNetworkInterfacesOutput
}

func (o *DescribeNetworkInterfacesOutput) Records() (records []*api.Record) {
	for _, s := range o.NetworkInterfaces {
		records = append(records, &api.Record{
			ID:   s.NetworkInterfaceId,
			Data: s,
		})
	}
	return
}

type DescribeNetworkInterfaces struct {
	API
}

var _ api.RequestBuilder = &DescribeNetworkInterfaces{}

// New implements api.RequestBuilder
func (fn *DescribeNetworkInterfaces) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeNetworkInterfacesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeNetworkInterfacesPagesWithContext(ctx, &input, func(output *ec2.DescribeNetworkInterfacesOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeNetworkInterfacesOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
