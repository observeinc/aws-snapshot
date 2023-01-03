package securityhub

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/securityhub"
)

func init() {
	service.Register("securityhub", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	GetFindingsPagesWithContext(ctx context.Context, input *securityhub.GetFindingsInput, fn func(*securityhub.GetFindingsOutput, bool) bool, opts ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	securityhubapi := securityhub.New(p, opts...)
	return api.Endpoint{
		"GetFindings": &GetFindings{securityhubapi},
	}
}
