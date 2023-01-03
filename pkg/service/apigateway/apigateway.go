package apigateway

import (
	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/apigateway"
)

func init() {
	service.Register("apigateway", api.ServiceFunc(New))
}

type API interface {
	GetRestApisPagesWithContext(ctx aws.Context, input *apigateway.GetRestApisInput, fn func(*apigateway.GetRestApisOutput, bool) bool, opts ...request.Option) error
	GetStagesWithContext(ctx aws.Context, input *apigateway.GetStagesInput, opts ...request.Option) (*apigateway.GetStagesOutput, error)
	GetDeploymentsWithContext(ctx aws.Context, input *apigateway.GetDeploymentsInput, opts ...request.Option) (*apigateway.GetDeploymentsOutput, error)
}

func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	eksapi := apigateway.New(p, opts...)
	return api.Endpoint{
		"GetRestApis":    &GetRestApis{eksapi},
		"GetDeployments": &GetDeployments{eksapi},
		"GetStages":      &GetStages{eksapi},
	}
}
