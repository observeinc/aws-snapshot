package ecs

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeClustersOutput{},
			InputFile:  "testdata/describeclusters-input.json",
			OutputFile: "testdata/describeclusters-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeServicesOutput{},
			InputFile:  "testdata/describeservices-input.json",
			OutputFile: "testdata/describeservices-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeTasksOutput{},
			InputFile:  "testdata/describetasks-input.json",
			OutputFile: "testdata/describetasks-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeTasksOutput{},
			InputFile:  "testdata/describecontainerinstances-input.json",
			OutputFile: "testdata/describecontainerinstances-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
