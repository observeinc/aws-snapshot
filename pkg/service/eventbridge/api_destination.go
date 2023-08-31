package eventbridge

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type ListApiDestinationsOutput struct {
	*eventbridge.ListApiDestinationsOutput
}

func (o *ListApiDestinationsOutput) Records() (records []*api.Record) {
	for _, o := range o.ApiDestinations {
		records = append(records, &api.Record{
			ID:   o.ApiDestinationArn,
			Data: o,
		})
	}
	return
}

type ListApiDestinations struct {
	API
}

var _ api.RequestBuilder = &ListApiDestinations{}

// New implements api.RequestBuilder
func (fn *ListApiDestinations) New(name string, config interface{}) ([]api.Request, error) {
	var input eventbridge.ListApiDestinationsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}
	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var countApiDestination int
		r, _ := ctx.Value("runner_config").(api.Runner)
		for {
			output, err := fn.ListApiDestinationsWithContext(ctx, &input)
			if err != nil {
				return fmt.Errorf("failed to list api destinations: %w", err)
			}
			if r.Stats {
				countApiDestination += len(output.ApiDestinations)
			} else {

			if err := api.SendRecords(ctx, ch, name, &ListApiDestinationsOutput{output}); err != nil {
				return err
			}
		}

			if output.NextToken == nil {
				break
			}
			input.NextToken = output.NextToken
		}
		if r.Stats {
			innerErr := api.SendRecords(ctx, ch, name, &api.CountRecords{countApiDestination})
			if innerErr != nil {
				return innerErr
			}
		}
		return nil
	}

	return []api.Request{call}, nil
}
