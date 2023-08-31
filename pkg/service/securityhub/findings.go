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
		var outerErr, innerErr error
		var countFindings int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.GetFindingsPagesWithContext(ctx, &input, func(output *securityhub.GetFindingsOutput, last bool) bool {
			if r.Stats {
				countFindings += len(output.Findings)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &GetFindingsOutput{output}); innerErr != nil {
					return false
				}
			}
			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countFindings})
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
