package rds

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeDBClustersOutput{},
			InputFile:  "testdata/describedbclusters-input.json",
			OutputFile: "testdata/describedbclusters-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeDBInstancesOutput{},
			InputFile:  "testdata/describedbinstances-input.json",
			OutputFile: "testdata/describedbinstances-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
