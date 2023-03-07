package elasticbeanstalk

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeApplicationsOutput{},
			InputFile:  "testdata/describeapplications-input.json",
			OutputFile: "testdata/describeapplications-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeApplicationVersionsOutput{},
			InputFile:  "testdata/describeapplicationversions-input.json",
			OutputFile: "testdata/describeapplicationversions-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeEnvironmentsOutput{},
			InputFile:  "testdata/describeenvironments-input.json",
			OutputFile: "testdata/describeenvironments-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeInstancesHealthOutput{},
			InputFile:  "testdata/describeinstanceshealth-input.json",
			OutputFile: "testdata/describeinstanceshealth-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}
}
