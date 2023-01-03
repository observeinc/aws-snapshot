package ec2

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeFlowLogsOutput{},
			InputFile:  "testdata/describeflowlogs-input.json",
			OutputFile: "testdata/describeflowlogs-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeInstancesOutput{},
			InputFile:  "testdata/describeinstances-input.json",
			OutputFile: "testdata/describeinstances-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeInternetGatewaysOutput{},
			InputFile:  "testdata/describeinternetgateways-input.json",
			OutputFile: "testdata/describeinternetgateways-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeNetworkAclsOutput{},
			InputFile:  "testdata/describenetworkacls-input.json",
			OutputFile: "testdata/describenetworkacls-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeNetworkInterfacesOutput{},
			InputFile:  "testdata/describenetworkinterfaces-input.json",
			OutputFile: "testdata/describenetworkinterfaces-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeRouteTablesOutput{},
			InputFile:  "testdata/describeroutetables-input.json",
			OutputFile: "testdata/describeroutetables-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeSecurityGroupsOutput{},
			InputFile:  "testdata/describesecuritygroups-input.json",
			OutputFile: "testdata/describesecuritygroups-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeSnapshotsOutput{},
			InputFile:  "testdata/describesnapshots-input.json",
			OutputFile: "testdata/describesnapshots-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeSubnetsOutput{},
			InputFile:  "testdata/describesubnets-input.json",
			OutputFile: "testdata/describesubnets-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeVolumesOutput{},
			InputFile:  "testdata/describevolumes-input.json",
			OutputFile: "testdata/describevolumes-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeVolumeStatusOutput{},
			InputFile:  "testdata/describevolumestatus-input.json",
			OutputFile: "testdata/describevolumestatus-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeVpcsOutput{},
			InputFile:  "testdata/describevpcs-input.json",
			OutputFile: "testdata/describevpcs-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeNatGatewaysOutput{},
			InputFile:  "testdata/describenatgateways-input.json",
			OutputFile: "testdata/describenatgateways-expect.json",
			Update:     *update,
		},
		{
			Source:     &DescribeAddressesOutput{},
			InputFile:  "testdata/describeaddresses-input.json",
			OutputFile: "testdata/describeaddresses-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
