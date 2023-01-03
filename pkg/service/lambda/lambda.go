package lambda

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func init() {
	service.Register("lambda", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListFunctionsPagesWithContext(context.Context, *lambda.ListFunctionsInput, func(*lambda.ListFunctionsOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	lambdaapi := lambda.New(p, opts...)
	return api.Endpoint{
		"ListFunctions": &ListFunctions{lambdaapi},
	}
}
