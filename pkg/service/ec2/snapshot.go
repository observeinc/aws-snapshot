package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeSnapshotsOutput struct {
	*ec2.DescribeSnapshotsOutput
}

func (o *DescribeSnapshotsOutput) Records() (records []*api.Record) {
	for _, v := range o.Snapshots {
		records = append(records, &api.Record{
			ID:   v.SnapshotId,
			Data: v,
		})
	}
	return
}

type DescribeSnapshots struct {
	API
}

var _ api.RequestBuilder = &DescribeSnapshots{}

// New implements api.RequestBuilder
func (fn *DescribeSnapshots) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeSnapshotsInput

	// set max, otherwise we're force fed all results in one call
	input.SetMaxResults(500)

	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	if len(input.Filters) == 0 && len(input.OwnerIds) == 0 {
		// snapshots tend to accrue over time, only collect this data if explicitly filtered
		return nil, nil
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeSnapshotsPagesWithContext(ctx, &input, func(output *ec2.DescribeSnapshotsOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeSnapshotsOutput{output})
		})
	}

	return []api.Request{call}, nil
}
