package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeNetworkAclsOutput struct {
	*ec2.DescribeNetworkAclsOutput
}

func (o *DescribeNetworkAclsOutput) Records() (records []*api.Record) {
	for _, s := range o.NetworkAcls {
		records = append(records, &api.Record{
			ID:   s.NetworkAclId,
			Data: s,
		})
	}
	return
}

type DescribeNetworkAcls struct {
	API
}

var _ api.RequestBuilder = &DescribeNetworkAcls{}

// New implements api.RequestBuilder
func (fn *DescribeNetworkAcls) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeNetworkAclsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeNetworkAclsPagesWithContext(ctx, &input, func(output *ec2.DescribeNetworkAclsOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeNetworkAclsOutput{output})
		})
	}

	return []api.Request{call}, nil
}
