package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type DescribeInstancesOutput struct {
	*ec2.DescribeInstancesOutput
}

func (o *DescribeInstancesOutput) Records() (records []*api.Record) {
	type elem struct {
		Instance      interface{} `json:"instance,omitempty"`
		Group         interface{} `json:"group,omitempty"`
		OwnerId       *string     `json:"ownerId"`
		RequesterId   *string     `json:"requesterId"`
		ReservationId *string     `json:"reservationId"`
	}

	for _, r := range o.Reservations {
		for _, g := range r.Groups {
			records = append(records, &api.Record{
				ID: g.GroupId,
				Data: elem{
					Group:         g,
					OwnerId:       r.OwnerId,
					RequesterId:   r.RequesterId,
					ReservationId: r.ReservationId,
				},
			})
		}
		for _, i := range r.Instances {
			records = append(records, &api.Record{
				ID: i.InstanceId,
				Data: elem{
					Instance:      i,
					OwnerId:       r.OwnerId,
					RequesterId:   r.RequesterId,
					ReservationId: r.ReservationId,
				},
			})
		}
	}
	return
}

type DescribeInstances struct {
	API
}

var _ api.RequestBuilder = &DescribeInstances{}

// New implements api.RequestBuilder
func (fn *DescribeInstances) New(name string, config interface{}) ([]api.Request, error) {
	var input ec2.DescribeInstancesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeInstancesPagesWithContext(ctx, &input, func(output *ec2.DescribeInstancesOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &DescribeInstancesOutput{output})
		})
	}

	return []api.Request{call}, nil
}
