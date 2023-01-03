package organizations

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/organizations"
)

func init() {
	service.Register("organizations", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeOrganizationWithContext(context.Context, *organizations.DescribeOrganizationInput, ...request.Option) (*organizations.DescribeOrganizationOutput, error)
	ListAccountsPagesWithContext(context.Context, *organizations.ListAccountsInput, func(*organizations.ListAccountsOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	organizationsapi := organizations.New(p, opts...)
	return api.Endpoint{
		"DescribeOrganizations": &DescribeOrganization{organizationsapi},
		"ListAccounts":          &ListAccounts{organizationsapi},
	}
}
