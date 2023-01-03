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
		return fn.ListClustersPagesWithContext(ctx, &input, func(output *eks.ListClustersOutput, last bool) bool {
			for _, clusterName := range output.Clusters {
				listFargateProfilesInput := &eks.ListFargateProfilesInput{
					ClusterName: clusterName,
				}

				err := fn.ListFargateProfilesPagesWithContext(ctx, listFargateProfilesInput, func(fpoutput *eks.ListFargateProfilesOutput, last bool) bool {
					for _, fargateProfileName := range fpoutput.FargateProfileNames {
						describeFargateProfileInput := &eks.DescribeFargateProfileInput{
							ClusterName:        clusterName,
							FargateProfileName: fargateProfileName,
						}
						describeFargateProfileOutput, err := fn.DescribeFargateProfileWithContext(ctx, describeFargateProfileInput)
						if err != nil {
							panic(err)
						}
						ok := api.SendRecords(ctx, ch, name, &DescribeFargateProfileOutput{describeFargateProfileOutput})
						if !ok {
							return false
						}
					}
					return true
				})
				if err != nil {
					panic(err)
				}
			}
			return true
		})
	}

	return []api.Request{call}, nil
}
