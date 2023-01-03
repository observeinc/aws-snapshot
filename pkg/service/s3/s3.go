package cloudwatchlogs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

func init() {
	service.Register("s3", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListBucketsWithContext(context.Context, *s3.ListBucketsInput, ...request.Option) (*s3.ListBucketsOutput, error)
	GetBucketTaggingWithContext(context.Context, *s3.GetBucketTaggingInput, ...request.Option) (*s3.GetBucketTaggingOutput, error)
	GetBucketLocationWithContext(context.Context, *s3.GetBucketLocationInput, ...request.Option) (*s3.GetBucketLocationOutput, error)
	GetBucketPolicyWithContext(context.Context, *s3.GetBucketPolicyInput, ...request.Option) (*s3.GetBucketPolicyOutput, error)
	GetBucketAclWithContext(context.Context, *s3.GetBucketAclInput, ...request.Option) (*s3.GetBucketAclOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	s3api := s3.New(p, opts...)
	return api.Endpoint{
		"ListBuckets": &ListBuckets{
			API:    s3api,
			Region: s3api.Client.Config.Region,
		},
	}
}
