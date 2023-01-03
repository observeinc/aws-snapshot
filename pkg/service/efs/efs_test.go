package efs

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeAccessPointsOutput{},
			InputFile:  "testdata/describeaccesspoints-input.json",
			OutputFile: "testdata/describeaccesspoints-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeBackupPolicyOutput{},
			InputFile:  "testdata/describebackuppolicy-input.json",
			OutputFile: "testdata/describebackuppolicy-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeFileSystemPolicyOutput{},
			InputFile:  "testdata/describefilesystempolicy-input.json",
			OutputFile: "testdata/describefilesystempolicy-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeFileSystemsOutput{},
			InputFile:  "testdata/describefilesystems-input.json",
			OutputFile: "testdata/describefilesystems-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeLifecycleConfigurationOutput{},
			InputFile:  "testdata/describelifecycleconfiguration-input.json",
			OutputFile: "testdata/describelifecycleconfiguration-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeMountTargetsOutput{},
			InputFile:  "testdata/describemounttargets-input.json",
			OutputFile: "testdata/describemounttargets-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeReplicationConfigurationsOutput{},
			InputFile:  "testdata/describereplicationconfigurations-input.json",
			OutputFile: "testdata/describereplicationconfigurations-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
