package firehose

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/firehose"
)

type DescribeDeliveryStreamOutput struct {
	*firehose.DescribeDeliveryStreamOutput
}

func (o *DescribeDeliveryStreamOutput) Records() (records []*api.Record) {
	if desc := o.DeliveryStreamDescription; desc != nil {
		records = append(records, &api.Record{
			ID:   o.DeliveryStreamDescription.DeliveryStreamARN,
			Data: o.DeliveryStreamDescription,
		})
	}
	return
}

type DescribeDeliveryStreams struct {
	API
}

var _ api.RequestBuilder = &DescribeDeliveryStreams{}

// New implements api.RequestBuilder
func (fn *DescribeDeliveryStreams) New(name string, config interface{}) ([]api.Request, error) {
	var input firehose.ListDeliveryStreamsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {

		var lastPage bool
		var deliveryStreamCounts int

		r, _ := ctx.Value("runner_config").(api.Runner)
		for !lastPage {
			listOutput, err := fn.ListDeliveryStreamsWithContext(ctx, &input)
			if err != nil {
				return fmt.Errorf("failed to list streams: %w", err)
			}
			if len(listOutput.DeliveryStreamNames) == 0 {
				break
			}

			if r.Stats {
				deliveryStreamCounts += len(listOutput.DeliveryStreamNames)

				} else {

				for _, deliveryStreamName := range listOutput.DeliveryStreamNames {
					describeDeliveryStreamOutput, err := fn.DescribeDeliveryStreamWithContext(ctx, &firehose.DescribeDeliveryStreamInput{
						DeliveryStreamName: deliveryStreamName,
					})
					if err != nil {
						return fmt.Errorf("failed to describe stream %q: %w", *deliveryStreamName, err)
					}

					if err := api.SendRecords(ctx, ch, name, &DescribeDeliveryStreamOutput{describeDeliveryStreamOutput}); err != nil {
						return err
					}

					input.SetExclusiveStartDeliveryStreamName(*deliveryStreamName)
				}
			}


			if listOutput.HasMoreDeliveryStreams != nil {
				lastPage = !(*listOutput.HasMoreDeliveryStreams)
			}
		}
		if r.Stats {
			innerErr := api.SendRecords(ctx, ch, name, &api.CountRecords{deliveryStreamCounts})
			if innerErr != nil {
				return innerErr
			}

		}
		return nil
	}

	return []api.Request{call}, nil
}
