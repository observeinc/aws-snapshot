package efs

import (
	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/efs"
)

func init() {
	// note: service name in policy does not match SDK code structure (elasticfilesystem vs efs)
	service.Register("elasticfilesystem", api.ServiceFunc(New))
}

type API interface {
	DescribeAccessPointsPagesWithContext(ctx aws.Context, input *efs.DescribeAccessPointsInput, fn func(*efs.DescribeAccessPointsOutput, bool) bool, opts ...request.Option) error
	DescribeBackupPolicyWithContext(ctx aws.Context, input *efs.DescribeBackupPolicyInput, opts ...request.Option) (*efs.DescribeBackupPolicyOutput, error)
	DescribeFileSystemsPagesWithContext(ctx aws.Context, input *efs.DescribeFileSystemsInput, fn func(*efs.DescribeFileSystemsOutput, bool) bool, opts ...request.Option) error
	DescribeFileSystemPolicyWithContext(ctx aws.Context, input *efs.DescribeFileSystemPolicyInput, opts ...request.Option) (*efs.DescribeFileSystemPolicyOutput, error)
	DescribeLifecycleConfigurationWithContext(ctx aws.Context, input *efs.DescribeLifecycleConfigurationInput, opts ...request.Option) (*efs.DescribeLifecycleConfigurationOutput, error)
	DescribeMountTargetsWithContext(ctx aws.Context, input *efs.DescribeMountTargetsInput, opts ...request.Option) (*efs.DescribeMountTargetsOutput, error)
	DescribeReplicationConfigurationsWithContext(ctx aws.Context, input *efs.DescribeReplicationConfigurationsInput, opts ...request.Option) (*efs.DescribeReplicationConfigurationsOutput, error)
}

func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	efsAPI := efs.New(p, opts...)
	return api.Endpoint{
		"DescribeAccessPoints":              &DescribeAccessPoints{efsAPI},
		"DescribeBackupPolicy":              &DescribeBackupPolicy{efsAPI},
		"DescribeFileSystems":               &DescribeFileSystems{efsAPI},
		"DescribeFileSystemPolicy":          &DescribeFileSystemPolicy{efsAPI},
		"DescribeLifecycleConfiguration":    &DescribeLifecycleConfiguration{efsAPI},
		"DescribeMountTargets":              &DescribeMountTargets{efsAPI},
		"DescribeReplicationConfigurations": &DescribeReplicationConfigurations{efsAPI},
	}
}
