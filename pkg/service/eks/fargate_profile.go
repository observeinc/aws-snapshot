package eks

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/eks"
)

type DescribeFargateProfileOutput struct {
	*eks.DescribeFargateProfileOutput
}

func (o *DescribeFargateProfileOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.FargateProfile.FargateProfileArn,
		Data: o,
	})
	return records
}

type DescribeFargateProfile struct {
	API
}

var _ api.RequestBuilder = &DescribeFargateProfile{}

// New implements api.RequestBuilder
func (fn *DescribeFargateProfile) New(name string, config interface{}) ([]api.Request, error) {
	var input eks.ListClustersInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var describeFargateCount int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.ListClustersPagesWithContext(ctx, &input, func(output *eks.ListClustersOutput, last bool) bool {

			for _, clusterName := range output.Clusters {
				listFargateProfilesInput := &eks.ListFargateProfilesInput{
					ClusterName: clusterName,
				}

				err := fn.ListFargateProfilesPagesWithContext(ctx, listFargateProfilesInput, func(fpoutput *eks.ListFargateProfilesOutput, last bool) bool {
					if r.Stats {
						describeFargateCount += len(fpoutput.FargateProfileNames)
					} else {
						for _, fargateProfileName := range fpoutput.FargateProfileNames {
							describeFargateProfileInput := &eks.DescribeFargateProfileInput{
								ClusterName:        clusterName,
								FargateProfileName: fargateProfileName,
							}
							describeFargateProfileOutput, err := fn.DescribeFargateProfileWithContext(ctx, describeFargateProfileInput)
							if err != nil {
								innerErr = err
								return false
							}

							if innerErr = api.SendRecords(ctx, ch, name, &DescribeFargateProfileOutput{describeFargateProfileOutput}); innerErr != nil {
								return false
							}
						}
					}
					return true
				})
				if innerErr = api.FirstError(err, innerErr); innerErr != nil {
					return false
				}
			}
			return true
		})

		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{Count: describeFargateCount})
		}

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
