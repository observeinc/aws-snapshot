package lambda

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/lambda"
)

type ListFunctionsOutput struct {
	*lambda.ListFunctionsOutput
}

func (o *ListFunctionsOutput) Records() (records []*api.Record) {
	for _, fn := range o.Functions {
		records = append(records, &api.Record{
			ID:   fn.FunctionArn,
			Data: fn,
		})
	}
	return
}

type ListFunctions struct {
	API
}

var _ api.RequestBuilder = &ListFunctions{}

// New implements api.RequestBuilder
func (fn *ListFunctions) New(name string, config interface{}) ([]api.Request, error) {
	var input lambda.ListFunctionsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countLambdaFunctions int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ListFunctionsPagesWithContext(ctx, &input, func(output *lambda.ListFunctionsOutput, last bool) bool {

			if r.Stats {
				countLambdaFunctions += len(output.Functions)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &ListFunctionsOutput{output}); innerErr != nil {
					return false
				}
			}

			return true
		})

		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countLambdaFunctions})
		}

		return api.FirstError(outerErr, innerErr)

	}

	return []api.Request{call}, nil
}
