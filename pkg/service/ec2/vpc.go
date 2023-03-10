package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeVpcsOutput struct {
	*ec2.DescribeVpcsOutput
}

func (o *DescribeVpcsOutput) Records() (records []*api.Record) {
	for _, v := range o.Vpcs {
		records = append(records, &api.Record{
			ID:   v.VpcId,
			Data: v,
		})
	}
	return
}

type DescribeVpcs struct {
	API
}

var _ api.RequestBuilder = &DescribeVpcs{}

// New implements api.RequestBuilder
func (fn *DescribeVpcs) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeVpcsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeVpcsPagesWithContext(ctx, &input, func(output *ec2.DescribeVpcsOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeVpcsOutput{output})
		})
	}

	return []api.Request{call}, nil
}
