package sns

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
)

type GetSubscriptionAttributesOutput struct {
	subscriptionArn *string
	*sns.GetSubscriptionAttributesOutput
}

func (o *GetSubscriptionAttributesOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.subscriptionArn,
		Data: o,
	})
	return
}

type GetSubscriptionAttributes struct {
	API
}

var _ api.RequestBuilder = &GetSubscriptionAttributes{}

// New implements api.RequestBuilder
func (fn *GetSubscriptionAttributes) New(name string, config interface{}) ([]api.Request, error) {
	var listSubscriptionsInput sns.ListSubscriptionsInput

	if err := api.DecodeConfig(config, &listSubscriptionsInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.ListSubscriptionsPagesWithContext(ctx, &listSubscriptionsInput, func(output *sns.ListSubscriptionsOutput, last bool) bool {
			for _, subscription := range output.Subscriptions {
				if subscription.SubscriptionArn == nil {
					continue
				}

				// sometimes we get an invalid ARN due to pending confirmation
				if _, err := arn.Parse(*subscription.SubscriptionArn); err != nil {
					continue
				}

				output, err := fn.GetSubscriptionAttributesWithContext(ctx, &sns.GetSubscriptionAttributesInput{SubscriptionArn: subscription.SubscriptionArn})
				if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sns.ErrCodeNotFoundException {
					continue
				}

				if err != nil {
					innerErr = fmt.Errorf("failed to process %s: %w", *subscription.SubscriptionArn, err)
					return false
				}
				if err := api.SendRecords(ctx, ch, name, &GetSubscriptionAttributesOutput{
					subscriptionArn:                 subscription.SubscriptionArn,
					GetSubscriptionAttributesOutput: output,
				}); err != nil {
					innerErr = err
					return false
				}
			}
			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
