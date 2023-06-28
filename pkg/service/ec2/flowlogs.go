package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeFlowLogsOutput struct {
	*ec2.DescribeFlowLogsOutput
}

func (o *DescribeFlowLogsOutput) Records() (records []*api.Record) {
	for _, s := range o.FlowLogs {
		records = append(records, &api.Record{
			ID:   s.FlowLogId,
			Data: s,
		})
	}
	return
}

type DescribeFlowLogs struct {
	API
}

var _ api.RequestBuilder = &DescribeFlowLogs{}

// New implements api.RequestBuilder
func (fn *DescribeFlowLogs) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeFlowLogsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeFlowLogsPagesWithContext(ctx, &input, func(output *ec2.DescribeFlowLogsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeFlowLogsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
