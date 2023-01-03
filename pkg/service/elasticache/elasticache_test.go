package elasticache

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeCacheClustersOutput{},
			InputFile:  "testdata/describecacheclusters-input.json",
			OutputFile: "testdata/describecacheclusters-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeReplicationGroupsOutput{},
			InputFile:  "testdata/describereplicationgroups-input.json",
			OutputFile: "testdata/describereplicationgroups-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeCacheSubnetGroupsOutput{},
			InputFile:  "testdata/describesubnetgroups-input.json",
			OutputFile: "testdata/describesubnetgroups-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
