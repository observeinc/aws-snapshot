package organizations

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/organizations"
)

type ListAccountsOutput struct {
	*organizations.ListAccountsOutput
}

func (o *ListAccountsOutput) Records() (records []*api.Record) {
	for _, a := range o.Accounts {
		records = append(records, &api.Record{
			ID:   a.Arn,
			Data: a,
		})
	}
	return
}

type ListAccounts struct {
	API
}

var _ api.RequestBuilder = &ListAccounts{}

// New implements api.RequestBuilder
func (fn *ListAccounts) New(name string, config interface{}) ([]api.Request, error) {
	var input organizations.ListAccountsInput

	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.ListAccountsPagesWithContext(ctx, &input, func(output *organizations.ListAccountsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &ListAccountsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		if aerr, ok := outerErr.(awserr.Error); ok && aerr.Code() == organizations.ErrCodeAWSOrganizationsNotInUseException {
			// nothing to do here
			outerErr = nil
		}

		if aerr, ok := outerErr.(awserr.Error); ok && aerr.Code() == organizations.ErrCodeAccessDeniedException {
			// ask for forgiveness, we may have configured this outside of master account
			// TODO: log warning
			outerErr = nil
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
