package kms

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"
)

type DescribeKeyOutput struct {
	*kms.DescribeKeyOutput
}

func (o *DescribeKeyOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.KeyMetadata.Arn,
		Data: o,
	})
	return
}

type DescribeKey struct {
	API
}

var _ api.RequestBuilder = &DescribeKey{}

// New implements api.RequestBuilder
func (fn *DescribeKey) New(name string, config interface{}) ([]api.Request, error) {
	var listKeysInput kms.ListKeysInput

	if err := api.DecodeConfig(config, &listKeysInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var keyCount int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ListKeysPagesWithContext(ctx, &listKeysInput, func(output *kms.ListKeysOutput, last bool) bool {
			if r.Stats {
				keyCount = +len(output.Keys)
			} else {
				for _, entry := range output.Keys {
					output, err := fn.DescribeKeyWithContext(ctx, &kms.DescribeKeyInput{
						KeyId: entry.KeyId,
					})

					// According to https://github.com/awsdocs/aws-kms-developer-guide/blob/master/doc_source/key-policies.md#using-key-policies-in-aws-kms
					// "Unless the key policy explicitly allows it, you cannot use IAM policies to allow access to a KMS key."
					//
					// This means that AccessDeniedException errors like https://observe.atlassian.net/browse/OB-10958 should happen
					// for users that use KMS and set a more restrictive key policy.
					if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "AccessDeniedException" {
						continue
					}
					if err != nil {
						innerErr = err
						return false
					}

					if innerErr = api.SendRecords(ctx, ch, name, &DescribeKeyOutput{output}); innerErr != nil {
						return false
					}
				}
			}
			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{Count: keyCount})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
