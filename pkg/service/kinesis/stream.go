package kinesis

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type DescribeStreamOutput struct {
	*kinesis.DescribeStreamOutput
}

func (o *DescribeStreamOutput) Records() (records []*api.Record) {
	if desc := o.StreamDescription; desc != nil {
		records = append(records, &api.Record{
			ID:   desc.StreamARN,
			Data: desc,
		})
	}
	return
}

type DescribeStreams struct {
	API
}

var _ api.RequestBuilder = &DescribeStreams{}

// New implements api.RequestBuilder
func (fn *DescribeStreams) New(name string, config interface{}) ([]api.Request, error) {
	var listStreamsInput kinesis.ListStreamsInput

	// Limit number of shards returned per call.
	describeStreamInput := kinesis.DescribeStreamInput{
		Limit: aws.Int64(50),
	}

	if err := api.DecodeConfig(config, &describeStreamInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countStreams int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ListStreamsPagesWithContext(ctx, &listStreamsInput, func(output *kinesis.ListStreamsOutput, last bool) bool {
			if r.Stats {
				countStreams += len(output.StreamNames)
			} else {
				for _, streamName := range output.StreamNames {
					describeStreamInput.StreamName = streamName
					err := fn.DescribeStreamPagesWithContext(ctx, &describeStreamInput, func(output *kinesis.DescribeStreamOutput, last bool) bool {
						if innerErr = api.SendRecords(ctx, ch, name, &DescribeStreamOutput{output}); innerErr != nil {
							return false
						}

						return true
					})

					if innerErr = api.FirstError(err, innerErr); innerErr != nil {
						return false
					}
				}
			}
			return true
		})

		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countStreams})
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
