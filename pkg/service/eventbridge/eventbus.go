package eventbridge

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type ListEventBusesOutput struct {
	*eventbridge.ListEventBusesOutput
}

func (o *ListEventBusesOutput) Records() (records []*api.Record) {
	for _, o := range o.EventBuses {
		records = append(records, &api.Record{
			ID:   o.Arn,
			Data: o,
		})
	}
	return
}

type ListEventBuses struct {
	API
}

var _ api.RequestBuilder = &ListRules{}

// New implements api.RequestBuilder
func (fn *ListEventBuses) New(name string, config interface{}) ([]api.Request, error) {
	var input eventbridge.ListEventBusesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		for {
			output, err := fn.ListEventBusesWithContext(ctx, &input)
			if err != nil {
				return fmt.Errorf("failed to list event buses: %w", err)
			}

			if err := api.SendRecords(ctx, ch, name, &ListEventBusesOutput{output}); err != nil {
				return err
			}

			if output.NextToken == nil {
				break
			}
			input.NextToken = output.NextToken
		}
		return nil
	}

	return []api.Request{call}, nil
}
