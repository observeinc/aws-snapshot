package autoscaling

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type DescribeLaunchConfigurationsOutput struct {
	*autoscaling.DescribeLaunchConfigurationsOutput
}

func (o *DescribeLaunchConfigurationsOutput) Records() (records []*api.Record) {
	for _, o := range o.LaunchConfigurations {
		records = append(records, &api.Record{
			ID:   o.LaunchConfigurationARN,
			Data: o,
		})
	}
	return
}

type DescribeLaunchConfigurations struct {
	API
}

var _ api.RequestBuilder = &DescribeLaunchConfigurations{}

// New implements api.RequestBuilder
func (fn *DescribeLaunchConfigurations) New(name string, config interface{}) ([]api.Request, error) {
	var input autoscaling.DescribeLaunchConfigurationsInput

	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeLaunchConfigurationsPagesWithContext(ctx, &input, func(output *autoscaling.DescribeLaunchConfigurationsOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeLaunchConfigurationsOutput{output}); err != nil {
				innerErr = err

				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
