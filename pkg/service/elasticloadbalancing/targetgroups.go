package elasticloadbalancing

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type DescribeTargetGroupsOutput struct {
	*elbv2.DescribeTargetGroupsOutput
	*elbv2.DescribeTagsOutput
}

func (o *DescribeTargetGroupsOutput) Records() (records []*api.Record) {
	tags := make(map[string][]*elbv2.Tag)
	if o.DescribeTagsOutput != nil {
		for _, desc := range o.DescribeTagsOutput.TagDescriptions {
			tags[*desc.ResourceArn] = desc.Tags
		}
	}

	for _, t := range o.TargetGroups {
		records = append(records, &api.Record{
			ID: t.TargetGroupArn,
			Data: struct {
				*elbv2.TargetGroup
				Tags []*elbv2.Tag
			}{
				TargetGroup: t,
				Tags:        tags[*t.TargetGroupArn],
			},
		})
	}
	return
}

type DescribeTargetGroups struct {
	ELBv2
}

var _ api.RequestBuilder = &DescribeTargetGroups{}

// New implements api.RequestBuilder
func (fn *DescribeTargetGroups) New(name string, config interface{}) ([]api.Request, error) {
	input := elbv2.DescribeTargetGroupsInput{
		PageSize: aws.Int64(20),
	}
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		return fn.DescribeTargetGroupsPagesWithContext(ctx, &input, func(output *elbv2.DescribeTargetGroupsOutput, last bool) bool {
			var describeTagsInput elbv2.DescribeTagsInput
			for _, targetGroup := range output.TargetGroups {
				describeTagsInput.ResourceArns = append(describeTagsInput.ResourceArns, targetGroup.TargetGroupArn)
			}
			describeTagsOutput, _ := fn.ELBv2.DescribeTagsWithContext(ctx, &describeTagsInput)
			return api.SendRecords(ctx, ch, name, &DescribeTargetGroupsOutput{
				DescribeTargetGroupsOutput: output,
				DescribeTagsOutput:         describeTagsOutput,
			})
		})
	}

	return []api.Request{call}, nil
}
