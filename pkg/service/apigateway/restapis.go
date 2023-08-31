package apigateway

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/apigateway"
)

type GetRestApisOutput struct {
	*apigateway.GetRestApisOutput
}

func (o *GetRestApisOutput) Records() (records []*api.Record) {
	for _, r := range o.Items {
		records = append(records, &api.Record{
			ID:   r.Id,
			Data: r,
		})
	}
	return records
}

type GetRestApis struct {
	API
}

var _ api.RequestBuilder = &GetRestApis{}

// New implements api.RequestBuilder
func (fn *GetRestApis) New(name string, config interface{}) ([]api.Request, error) {
	var input apigateway.GetRestApisInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countRestApis int
		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.GetRestApisPagesWithContext(ctx, &input, func(output *apigateway.GetRestApisOutput, last bool) bool {
			if r.Stats {
				countRestApis += len(output.Items)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &GetRestApisOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countRestApis})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
