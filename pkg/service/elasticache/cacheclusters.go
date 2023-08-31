package elasticache

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

type DescribeCacheClustersOutput struct {
	*elasticache.DescribeCacheClustersOutput
}

func (o *DescribeCacheClustersOutput) Records() (records []*api.Record) {
	for _, cc := range o.CacheClusters {
		records = append(records, &api.Record{
			ID:   cc.ARN,
			Data: cc,
		})
	}
	return
}

type DescribeCacheClusters struct {
	API
}

var _ api.RequestBuilder = &DescribeCacheClusters{}

// New implements api.RequestBuilder
func (fn *DescribeCacheClusters) New(name string, config interface{}) ([]api.Request, error) {
	var input elasticache.DescribeCacheClustersInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countCacheClusters int
		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeCacheClustersPagesWithContext(ctx, &input, func(output *elasticache.DescribeCacheClustersOutput, last bool) bool {
			if r.Stats {
				countCacheClusters += len(output.CacheClusters)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeCacheClustersOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countCacheClusters})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
