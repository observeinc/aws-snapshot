package route53

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/route53"
)

func init() {
	service.Register("route53", api.ServiceFunc(New))
}

type ListHostedZonesOutput struct {
	*route53.ListHostedZonesOutput
}

func (o *ListHostedZonesOutput) Records() (records []*api.Record) {
	for _, z := range o.HostedZones {
		records = append(records, &api.Record{
			ID:   z.Id,
			Data: z,
		})
	}
	return
}

// API documents the subset of AWS API we actually call
type API interface {
	ListHostedZonesPagesWithContext(context.Context, *route53.ListHostedZonesInput, func(*route53.ListHostedZonesOutput, bool) bool, ...request.Option) error
}

type ListHostedZones struct {
	API
}

var _ api.RequestBuilder = &ListHostedZones{}

// New implements api.RequestBuilder
func (fn *ListHostedZones) New(name string, config interface{}) ([]api.Request, error) {
	var input route53.ListHostedZonesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.ListHostedZonesPagesWithContext(ctx, &input, func(output *route53.ListHostedZonesOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &ListHostedZonesOutput{output})
		})
	}

	return []api.Request{call}, nil
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	route53api := route53.New(p, opts...)
	return api.Endpoint{
		"ListHostedZones": &ListHostedZones{route53api},
	}
}
