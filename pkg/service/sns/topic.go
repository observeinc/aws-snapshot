package sns

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/sns"
)

type GetTopicAttributesOutput struct {
	topicArn *string
	*sns.GetTopicAttributesOutput
}

func (o *GetTopicAttributesOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.topicArn,
		Data: o,
	})
	return
}

type GetTopicAttributes struct {
	API
}

var _ api.RequestBuilder = &GetTopicAttributes{}

// New implements api.RequestBuilder
func (fn *GetTopicAttributes) New(name string, config interface{}) ([]api.Request, error) {
	var listTopicsInput sns.ListTopicsInput

	if err := api.DecodeConfig(config, &listTopicsInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var countTopics int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ListTopicsPagesWithContext(ctx, &listTopicsInput, func(output *sns.ListTopicsOutput, last bool) bool {
			if r.Stats {
				countTopics += len(output.Topics)
			} else { 
			for _, topic := range output.Topics {
				if topic.TopicArn == nil {
					continue
				}

				output, err := fn.GetTopicAttributesWithContext(ctx, &sns.GetTopicAttributesInput{TopicArn: topic.TopicArn})
				if err != nil {
					innerErr = fmt.Errorf("failed to get %q: %w", *topic.TopicArn, err)
					return false
				}
				if innerErr = api.SendRecords(ctx, ch, name, &GetTopicAttributesOutput{
					topicArn:                 topic.TopicArn,
					GetTopicAttributesOutput: output,
				}); innerErr != nil {
					return false
				}
			}
		}
			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countTopics})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
