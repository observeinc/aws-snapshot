package elasticloadbalancing

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeLoadBalancersOutput{},
			InputFile:  "testdata/describeloadbalancers-input.json",
			OutputFile: "testdata/describeloadbalancers-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeLoadBalancersOutputV2{},
			InputFile:  "testdata/describeloadbalancersv2-input.json",
			OutputFile: "testdata/describeloadbalancersv2-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeTargetGroupsOutput{},
			InputFile:  "testdata/describetargetgroups-input.json",
			OutputFile: "testdata/describetargetgroups-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
