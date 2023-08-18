package cloudwatchlogs

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func init() {
	service.Register("logs", api.ServiceFunc(New))
}

type DescribeLogGroupsOutput struct {
	*cloudwatchlogs.DescribeLogGroupsOutput
}

func (o *DescribeLogGroupsOutput) Records() (records []*api.Record) {
	for _, logGroup := range o.LogGroups {
		records = append(records, &api.Record{
			ID:   logGroup.Arn,
			Data: logGroup,
		})
	}
	return
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeLogGroupsPagesWithContext(context.Context, *cloudwatchlogs.DescribeLogGroupsInput, func(*cloudwatchlogs.DescribeLogGroupsOutput, bool) bool, ...request.Option) error
}

type DescribeLogGroups struct {
	API
}

var _ api.RequestBuilder = &DescribeLogGroups{}

// New implements api.RequestBuilder
func (fn *DescribeLogGroups) New(name string, config interface{}) ([]api.Request, error) {
	var input cloudwatchlogs.DescribeLogGroupsInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error
		var describeLogGroupsCount int

		r, _ := ctx.Value("runner_config").(api.Runner)
		outerErr = fn.DescribeLogGroupsPagesWithContext(ctx, &input, func(output *cloudwatchlogs.DescribeLogGroupsOutput, last bool) bool {
			if r.Stats {
				describeLogGroupsCount += len(output.LogGroups)
			} else {
				if innerErr = api.SendRecords(ctx, ch, name, &DescribeLogGroupsOutput{output}); innerErr != nil {
					return false
				}

			}
			return true
		})

		if outerErr == nil && r.Stats {
			innerErr = api.SendRecords(ctx, ch, name, &api.CountRecords{describeLogGroupsCount})
		}
		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	cwlogs := cloudwatchlogs.New(p, opts...)
	return api.Endpoint{
		"DescribeLogGroups": &DescribeLogGroups{cwlogs},
	}
}
