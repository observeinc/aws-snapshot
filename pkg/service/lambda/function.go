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

		outerErr = fn.ListFunctionsPagesWithContext(ctx, &input, func(output *lambda.ListFunctionsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &ListFunctionsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
