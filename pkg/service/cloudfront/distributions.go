package cloudfront

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/cloudfront"
)

type ListDistributionsOutput struct {
	*cloudfront.ListDistributionsOutput
}

func (o *ListDistributionsOutput) Records() (records []*api.Record) {
	for _, s := range o.DistributionList.Items {
		records = append(records, &api.Record{
			ID:   s.Id,
			Data: s,
		})
	}
	return
}

type ListDistributions struct {
	API
}

var _ api.RequestBuilder = &ListDistributions{}

// New implements api.RequestBuilder
func (fn *ListDistributions) New(name string, config interface{}) ([]api.Request, error) {
	var input cloudfront.ListDistributionsInput

	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.ListDistributionsPagesWithContext(ctx, &input, func(output *cloudfront.ListDistributionsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &ListDistributionsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
