package cloudwatchlogs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
)

func init() {
	service.Register("iam", api.ServiceFunc(New))
}

var ignoredPolicies = map[string]bool{
	"arn:aws:iam::aws:policy/aws-service-role/AWSSupportServiceRolePolicy": true,
}

type ListAccountAliasesOutput struct {
	*iam.ListAccountAliasesOutput
}

func (o *ListAccountAliasesOutput) Records() (records []*api.Record) {
	// ideally the ID would be the account number, but we can't get it here.
	// We'll need to extract it from the invokedFunctionArn
	records = append(records, &api.Record{Data: o})
	return
}

type GetAccountAuthorizationDetailsOutput struct {
	*iam.GetAccountAuthorizationDetailsOutput
}

func (o *GetAccountAuthorizationDetailsOutput) Records() (records []*api.Record) {
	for _, groupDetail := range o.GroupDetailList {
		records = append(records, &api.Record{
			ID:   groupDetail.Arn,
			Data: groupDetail,
		})
	}
	for _, roleDetail := range o.RoleDetailList {
		records = append(records, &api.Record{
			ID:   roleDetail.Arn,
			Data: roleDetail,
		})
	}
	for _, userDetail := range o.UserDetailList {
		records = append(records, &api.Record{
			ID:   userDetail.Arn,
			Data: userDetail,
		})
	}
	for _, policy := range o.Policies {
		if policy.Arn != nil && ignoredPolicies[*policy.Arn] {
			continue
		}
		records = append(records, &api.Record{
			ID:   policy.Arn,
			Data: policy,
		})
	}
	return
}

// API documents the subset of AWS API we actually call
type API interface {
	ListAccountAliasesPagesWithContext(context.Context, *iam.ListAccountAliasesInput, func(*iam.ListAccountAliasesOutput, bool) bool, ...request.Option) error
	GetAccountAuthorizationDetailsPagesWithContext(context.Context, *iam.GetAccountAuthorizationDetailsInput, func(*iam.GetAccountAuthorizationDetailsOutput, bool) bool, ...request.Option) error
}

type ListAccountAliases struct {
	API
}

var _ api.RequestBuilder = &ListAccountAliases{}

// New implements api.RequestBuilder
func (fn *ListAccountAliases) New(name string, config interface{}) ([]api.Request, error) {
	var input iam.ListAccountAliasesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.ListAccountAliasesPagesWithContext(ctx, &input, func(output *iam.ListAccountAliasesOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &ListAccountAliasesOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}

type GetAccountAuthorizationDetails struct {
	API
}

var _ api.RequestBuilder = &GetAccountAuthorizationDetails{}

// New implements api.RequestBuilder
func (fn *GetAccountAuthorizationDetails) New(name string, config interface{}) ([]api.Request, error) {
	var input iam.GetAccountAuthorizationDetailsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.GetAccountAuthorizationDetailsPagesWithContext(ctx, &input, func(output *iam.GetAccountAuthorizationDetailsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &GetAccountAuthorizationDetailsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	iamapi := iam.New(p, opts...)
	return api.Endpoint{
		"ListAccountAliases":             &ListAccountAliases{iamapi},
		"GetAccountAuthorizationDetails": &GetAccountAuthorizationDetails{iamapi},
	}
}
