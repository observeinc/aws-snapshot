package sqs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
		var outerErr, innerErr error
		var countQueues int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ListQueuesPagesWithContext(ctx, &listQueuesInput, func(output *sqs.ListQueuesOutput, last bool) bool {
			if r.Stats {
				countQueues += len(output.QueueUrls)
			} else {
				for _, queueURL := range output.QueueUrls {
					output, err := fn.GetQueueAttributesWithContext(ctx, &sqs.GetQueueAttributesInput{
						QueueUrl:       queueURL,
						AttributeNames: []*string{aws.String("All")},
					})
					if err != nil {
						if ae, ok := err.(awserr.Error); ok && ae.Code() == sqs.ErrCodeQueueDoesNotExist {
							continue
						}

						innerErr = err
						return false
					}
					if innerErr = api.SendRecords(ctx, ch, name, &GetQueueAttributesOutput{GetQueueAttributesOutput: output}); innerErr != nil {
						return false
					}
				}
			}
			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{Count: countQueues})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
