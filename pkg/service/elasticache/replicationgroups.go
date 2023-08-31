package elasticache

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

type DescribeReplicationGroupsOutput struct {
	*elasticache.DescribeReplicationGroupsOutput
}

func (o *DescribeReplicationGroupsOutput) Records() (records []*api.Record) {
	for _, rg := range o.ReplicationGroups {
		records = append(records, &api.Record{
			ID:   rg.ARN,
			Data: rg,
		})
	}
	return
}

type DescribeReplicationGroups struct {
	API
}

var _ api.RequestBuilder = &DescribeReplicationGroups{}

// New implements api.RequestBuilder
func (fn *DescribeReplicationGroups) New(name string, config interface{}) ([]api.Request, error) {
	var input elasticache.DescribeReplicationGroupsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countReplicationGroups int
		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeReplicationGroupsPagesWithContext(ctx, &input, func(output *elasticache.DescribeReplicationGroupsOutput, last bool) bool {
			if r.Stats {
				countReplicationGroups += len(output.ReplicationGroups)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeReplicationGroupsOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})
		if r.Stats {
			innerErr := api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countReplicationGroups})
			if innerErr != nil {
				return innerErr
			}
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
