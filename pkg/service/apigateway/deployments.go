package apigateway

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/apigateway"
)

type GetDeploymentsOutput struct {
	*apigateway.GetDeploymentsOutput

	restApiId   *string
	restApiName *string
}

type GetDeploymentsRecord struct {
	*apigateway.Deployment
	RestApiId   *string
	RestApiName *string
}

func (o *GetDeploymentsOutput) Records() (records []*api.Record) {
	for _, r := range o.Items {
		rWithRestAPIInfo := GetDeploymentsRecord{
			Deployment:  r,
			RestApiId:   o.restApiId,
			RestApiName: o.restApiName,
		}
		records = append(records, &api.Record{
			ID:   o.restApiId,
			Data: rWithRestAPIInfo,
		})
	}
	return records
}

type GetDeployments struct {
	API
}

var _ api.RequestBuilder = &GetDeployments{}

// New implements api.RequestBuilder
func (fn *GetDeployments) New(name string, config interface{}) ([]api.Request, error) {
	var input apigateway.GetRestApisInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var innerErr, outerErr error
		var countDeployments int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.GetRestApisPagesWithContext(ctx, &input, func(output *apigateway.GetRestApisOutput, last bool) bool {
			if r.Stats {
				countDeployments += len(output.Items)
			} else {
				for _, restApi := range output.Items {
					stagesInput := &apigateway.GetDeploymentsInput{
						RestApiId: restApi.Id,
					}

					deploymentsOutput, err := fn.GetDeploymentsWithContext(ctx, stagesInput)
					if err != nil {
						innerErr = err
						return false
					}

					source := &GetDeploymentsOutput{
						GetDeploymentsOutput: deploymentsOutput,
						restApiId:            restApi.Id,
						restApiName:          restApi.Name,
					}

					if err := api.SendRecords(ctx, ch, name, source); err != nil {
						innerErr = err

						return false
					}
				}
			}
			return true
		})
		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{countDeployments})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
