package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeVolumesOutput struct {
	*ec2.DescribeVolumesOutput
}

func (o *DescribeVolumesOutput) Records() (records []*api.Record) {
	for _, v := range o.Volumes {
		records = append(records, &api.Record{
			ID:   v.VolumeId,
			Data: v,
		})
	}
	return
}

type DescribeVolumes struct {
	API
}

var _ api.RequestBuilder = &DescribeVolumes{}

// New implements api.RequestBuilder
func (fn *DescribeVolumes) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeVolumesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeVolumesPagesWithContext(ctx, &input, func(output *ec2.DescribeVolumesOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeVolumesOutput{output})
		})
	}

	return []api.Request{call}, nil
}

type DescribeVolumeStatusOutput struct {
	*ec2.DescribeVolumeStatusOutput
}

func (o *DescribeVolumeStatusOutput) Records() (records []*api.Record) {
	for _, v := range o.VolumeStatuses {
		records = append(records, &api.Record{
			ID:   v.VolumeId,
			Data: v,
		})
	}
	return
}

type DescribeVolumeStatus struct {
	API
}

var _ api.RequestBuilder = &DescribeVolumeStatus{}

// New implements api.RequestBuilder
func (fn *DescribeVolumeStatus) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeVolumeStatusInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeVolumeStatusPagesWithContext(ctx, &input, func(output *ec2.DescribeVolumeStatusOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeVolumeStatusOutput{output})
		})
	}

	return []api.Request{call}, nil
}
