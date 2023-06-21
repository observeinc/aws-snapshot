package apigateway

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/apigateway"
)

type GetStagesOutput struct {
	*apigateway.GetStagesOutput

	restApiId   *string
	restApiName *string
}

type GetStagesRecord struct {
	*apigateway.Stage
	RestApiId   *string
	RestApiName *string
}

func (o *GetStagesOutput) Records() (records []*api.Record) {
	for _, r := range o.Item {
		rWithRestAPIInfo := GetStagesRecord{
			Stage:       r,
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

type GetStages struct {
	API
}

var _ api.RequestBuilder = &GetStages{}

// New implements api.RequestBuilder
func (fn *GetStages) New(name string, config interface{}) ([]api.Request, error) {
	var input apigateway.GetRestApisInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var innerErr, outerErr error

		outerErr = fn.GetRestApisPagesWithContext(ctx, &input, func(output *apigateway.GetRestApisOutput, last bool) bool {
			for _, restApi := range output.Items {
				stagesInput := &apigateway.GetStagesInput{
					RestApiId: restApi.Id,
				}

				stagesOutput, err := fn.GetStagesWithContext(ctx, stagesInput)
				if err != nil {
					innerErr = err
					return false
				}

				source := &GetStagesOutput{
					GetStagesOutput: stagesOutput,
					restApiId:       restApi.Id,
					restApiName:     restApi.Name,
				}

				if !api.SendRecords(ctx, ch, name, source) {
					return false
				}
			}
			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
