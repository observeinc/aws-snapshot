package elasticache

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

type DescribeCacheSubnetGroupsOutput struct {
	*elasticache.DescribeCacheSubnetGroupsOutput
}

func (o *DescribeCacheSubnetGroupsOutput) Records() (records []*api.Record) {
	for _, sg := range o.CacheSubnetGroups {
		records = append(records, &api.Record{
			ID:   sg.ARN,
			Data: sg,
		})
	}
	return
}

type DescribeCacheSubnetGroups struct {
	API
}

var _ api.RequestBuilder = &DescribeCacheSubnetGroups{}

// New implements api.RequestBuilder
func (fn *DescribeCacheSubnetGroups) New(name string, config interface{}) ([]api.Request, error) {
	var input elasticache.DescribeCacheSubnetGroupsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeCacheSubnetGroupsPagesWithContext(ctx, &input, func(output *elasticache.DescribeCacheSubnetGroupsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeCacheSubnetGroupsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
