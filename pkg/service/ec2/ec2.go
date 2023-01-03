package ec2

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func init() {
	service.Register("ec2", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeAddressesWithContext(context.Context, *ec2.DescribeAddressesInput, ...request.Option) (*ec2.DescribeAddressesOutput, error)
	DescribeFlowLogsPagesWithContext(context.Context, *ec2.DescribeFlowLogsInput, func(*ec2.DescribeFlowLogsOutput, bool) bool, ...request.Option) error
	DescribeInstancesPagesWithContext(context.Context, *ec2.DescribeInstancesInput, func(*ec2.DescribeInstancesOutput, bool) bool, ...request.Option) error
	DescribeInternetGatewaysPagesWithContext(context.Context, *ec2.DescribeInternetGatewaysInput, func(*ec2.DescribeInternetGatewaysOutput, bool) bool, ...request.Option) error
	DescribeNatGatewaysPagesWithContext(context.Context, *ec2.DescribeNatGatewaysInput, func(*ec2.DescribeNatGatewaysOutput, bool) bool, ...request.Option) error
	DescribeNetworkAclsPagesWithContext(context.Context, *ec2.DescribeNetworkAclsInput, func(*ec2.DescribeNetworkAclsOutput, bool) bool, ...request.Option) error
	DescribeNetworkInterfacesPagesWithContext(context.Context, *ec2.DescribeNetworkInterfacesInput, func(*ec2.DescribeNetworkInterfacesOutput, bool) bool, ...request.Option) error
	DescribeRouteTablesPagesWithContext(context.Context, *ec2.DescribeRouteTablesInput, func(*ec2.DescribeRouteTablesOutput, bool) bool, ...request.Option) error
	DescribeSecurityGroupsPagesWithContext(context.Context, *ec2.DescribeSecurityGroupsInput, func(*ec2.DescribeSecurityGroupsOutput, bool) bool, ...request.Option) error
	DescribeSnapshotsPagesWithContext(context.Context, *ec2.DescribeSnapshotsInput, func(*ec2.DescribeSnapshotsOutput, bool) bool, ...request.Option) error
	DescribeSubnetsPagesWithContext(context.Context, *ec2.DescribeSubnetsInput, func(*ec2.DescribeSubnetsOutput, bool) bool, ...request.Option) error
	DescribeVolumesPagesWithContext(context.Context, *ec2.DescribeVolumesInput, func(*ec2.DescribeVolumesOutput, bool) bool, ...request.Option) error
	DescribeVolumeStatusPagesWithContext(context.Context, *ec2.DescribeVolumeStatusInput, func(*ec2.DescribeVolumeStatusOutput, bool) bool, ...request.Option) error
	DescribeVpcsPagesWithContext(context.Context, *ec2.DescribeVpcsInput, func(*ec2.DescribeVpcsOutput, bool) bool, ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	ec2api := ec2.New(p, opts...)
	return api.Endpoint{
		"DescribeAddresses":         &DescribeAddresses{ec2api},
		"DescribeFlowLogs":          &DescribeFlowLogs{ec2api},
		"DescribeInstances":         &DescribeInstances{ec2api},
		"DescribeInternetGateways":  &DescribeInternetGateways{ec2api},
		"DescribeNatGateways":       &DescribeNatGateways{ec2api},
		"DescribeNetworkAcls":       &DescribeNetworkAcls{ec2api},
		"DescribeNetworkInterfaces": &DescribeNetworkInterfaces{ec2api},
		"DescribeRouteTables":       &DescribeRouteTables{ec2api},
		"DescribeSnapshots":         &DescribeSnapshots{ec2api},
		"DescribeSecurityGroups":    &DescribeSecurityGroups{ec2api},
		"DescribeSubnets":           &DescribeSubnets{ec2api},
		"DescribeVolumes":           &DescribeVolumes{ec2api},
		"DescribeVolumeStatus":      &DescribeVolumeStatus{ec2api},
		"DescribeVpcs":              &DescribeVpcs{ec2api},
	}
}
