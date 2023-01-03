package securityhub

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/securityhub"
)

type GetFindingsOutput struct {
	*securityhub.GetFindingsOutput
}

func (o *GetFindingsOutput) Records() (records []*api.Record) {
	for _, finding := range o.Findings {
		records = append(records, &api.Record{
			ID:   finding.Id,
			Data: finding,
		})
	}
	return records
}

type GetFindings struct {
	API
}

var _ api.RequestBuilder = &GetFindings{}

// New implements api.RequestBuilder
func (fn *GetFindings) New(name string, config interface{}) ([]api.Request, error) {
	var input securityhub.GetFindingsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	if input.Filters == nil {
		// skip action if no filters are provided
		return nil, nil
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.GetFindingsPagesWithContext(ctx, &input, func(output *securityhub.GetFindingsOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &GetFindingsOutput{output})
		})
	}

	return []api.Request{call}, nil
}
