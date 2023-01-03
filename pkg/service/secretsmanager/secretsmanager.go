package secretsmanager

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func init() {
	service.Register("secretsmanager", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	ListSecretsPagesWithContext(context.Context, *secretsmanager.ListSecretsInput, func(*secretsmanager.ListSecretsOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	secretsmanagerapi := secretsmanager.New(p, opts...)
	return api.Endpoint{
		"ListSecrets": &ListSecrets{secretsmanagerapi},
	}
}
