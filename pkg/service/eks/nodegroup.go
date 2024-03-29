package eks

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/eks"
)

type DescribeNodegroupOutput struct {
	*eks.DescribeNodegroupOutput
}

func (o *DescribeNodegroupOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		ID:   o.Nodegroup.NodegroupArn,
		Data: o,
	})
	return records
}

type DescribeNodegroup struct {
	API
}

var _ api.RequestBuilder = &DescribeNodegroup{}

// New implements api.RequestBuilder
func (fn *DescribeNodegroup) New(name string, config interface{}) ([]api.Request, error) {
	var input eks.ListClustersInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.ListClustersPagesWithContext(ctx, &input, func(output *eks.ListClustersOutput, last bool) bool {
			for _, clusterName := range output.Clusters {
				listNodegroupsInput := &eks.ListNodegroupsInput{
					ClusterName: clusterName,
				}

				err := fn.ListNodegroupsPagesWithContext(ctx, listNodegroupsInput, func(ngoutput *eks.ListNodegroupsOutput, last bool) bool {
					for _, nodegroupName := range ngoutput.Nodegroups {
						describeNodegroupInput := &eks.DescribeNodegroupInput{
							ClusterName:   clusterName,
							NodegroupName: nodegroupName,
						}
						describeNodegroupOutput, err := fn.DescribeNodegroupWithContext(ctx, describeNodegroupInput)
						if err != nil {
							innerErr = err
							return false
						}

						if err := api.SendRecords(ctx, ch, name, &DescribeNodegroupOutput{describeNodegroupOutput}); err != nil {
							innerErr = err
							return false
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

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
