package organizations

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/organizations"
)

type DescribeOrganizationOutput struct {
	*organizations.DescribeOrganizationOutput
}

func (o *DescribeOrganizationOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.Organization.Arn,
		Data: o,
	})
	return
}

type DescribeOrganization struct {
	API
}

var _ api.RequestBuilder = &DescribeOrganization{}

// New implements api.RequestBuilder
func (fn *DescribeOrganization) New(name string, config interface{}) ([]api.Request, error) {
	var input organizations.DescribeOrganizationInput

	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		output, err := fn.DescribeOrganizationWithContext(ctx, &input)
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == organizations.ErrCodeAWSOrganizationsNotInUseException {
			// nothing to do here
			return nil
		} else if err != nil {
			return err
		}

		return api.SendRecords(ctx, ch, name, &DescribeOrganizationOutput{output})
	}

	return []api.Request{call}, nil
}
