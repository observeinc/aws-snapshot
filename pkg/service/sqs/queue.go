package sqs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type GetQueueAttributesOutput struct {
	*sqs.GetQueueAttributesOutput
}

func (o *GetQueueAttributesOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.Attributes["QueueArn"],
		Data: o,
	})
	return
}

type GetQueueAttributes struct {
	API
}

var _ api.RequestBuilder = &GetQueueAttributes{}

// New implements api.RequestBuilder
func (fn *GetQueueAttributes) New(name string, config interface{}) ([]api.Request, error) {
	var listQueuesInput sqs.ListQueuesInput

	if err := api.DecodeConfig(config, &listQueuesInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.ListQueuesPagesWithContext(ctx, &listQueuesInput, func(output *sqs.ListQueuesOutput, last bool) bool {
			for _, queueURL := range output.QueueUrls {
				output, err := fn.GetQueueAttributesWithContext(ctx, &sqs.GetQueueAttributesInput{
					QueueUrl:       queueURL,
					AttributeNames: []*string{aws.String("All")},
				})
				if err != nil {
					panic(err)
				}
				if !api.SendRecords(ctx, ch, name, &GetQueueAttributesOutput{GetQueueAttributesOutput: output}) {
					return false
				}
			}
			return true
		})
	}

	return []api.Request{call}, nil
}
