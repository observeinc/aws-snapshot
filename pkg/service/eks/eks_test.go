package eks

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeClusterOutput{},
			InputFile:  "testdata/describecluster-input.json",
			OutputFile: "testdata/describecluster-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeNodegroupOutput{},
			InputFile:  "testdata/describenodegroup-input.json",
			OutputFile: "testdata/describenodegroup-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeFargateProfileOutput{},
			InputFile:  "testdata/describefargateprofile-input.json",
			OutputFile: "testdata/describefargateprofile-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}
}
