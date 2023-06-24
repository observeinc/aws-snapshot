package secretsmanager

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type ListSecretsOutput struct {
	*secretsmanager.ListSecretsOutput
}

func (o *ListSecretsOutput) Records() (records []*api.Record) {
	for _, secret := range o.SecretList {
		records = append(records, &api.Record{
			ID:   secret.ARN,
			Data: secret,
		})
	}
	return
}

type ListSecrets struct {
	API
}

var _ api.RequestBuilder = &ListSecrets{}

// New implements api.RequestBuilder
func (fn *ListSecrets) New(name string, config interface{}) ([]api.Request, error) {
	var listSecretsInput secretsmanager.ListSecretsInput

	if err := api.DecodeConfig(config, &listSecretsInput); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.ListSecretsPagesWithContext(ctx, &listSecretsInput, func(output *secretsmanager.ListSecretsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &ListSecretsOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
