package cloudfront

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func init() {
	service.Register("cloudfront", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListDistributionsPagesWithContext(context.Context, *cloudfront.ListDistributionsInput, func(*cloudfront.ListDistributionsOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	cloudfrontapi := cloudfront.New(p, opts...)
	return api.Endpoint{
		"ListDistributions": &ListDistributions{cloudfrontapi},
	}
}
