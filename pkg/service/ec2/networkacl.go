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
	var networkACLCount int
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		r, _ := ctx.Value("runner_config").(api.Runner)

		outerErr = fn.DescribeNetworkAclsPagesWithContext(ctx, &input, func(output *ec2.DescribeNetworkAclsOutput, last bool) bool {
			if r.Stats {
				networkACLCount += len(output.NetworkAcls)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeNetworkAclsOutput{output}); innerErr != nil {
					return false
				}
			}
			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{Count: networkACLCount})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
