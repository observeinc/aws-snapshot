package kms

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/kms"
)

func init() {
	service.Register("kms", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListKeysPagesWithContext(context.Context, *kms.ListKeysInput, func(*kms.ListKeysOutput, bool) bool, ...request.Option) error
	DescribeKeyWithContext(context.Context, *kms.DescribeKeyInput, ...request.Option) (*kms.DescribeKeyOutput, error)
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	kmsapi := kms.New(p, opts...)
	return api.Endpoint{
		"DescribeKey": &DescribeKey{kmsapi},
	}
}
