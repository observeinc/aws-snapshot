package elasticbeanstalk

import (
	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
)

func init() {
	service.Register("elasticbeanstalk", api.ServiceFunc(New))
}

type API interface {
	DescribeApplicationsWithContext(ctx aws.Context, input *elasticbeanstalk.DescribeApplicationsInput, opts ...request.Option) (*elasticbeanstalk.DescribeApplicationsOutput, error)
	DescribeApplicationVersionsWithContext(ctx aws.Context, input *elasticbeanstalk.DescribeApplicationVersionsInput, opts ...request.Option) (*elasticbeanstalk.DescribeApplicationVersionsOutput, error)
	DescribeEnvironmentsWithContext(ctx aws.Context, input *elasticbeanstalk.DescribeEnvironmentsInput, opts ...request.Option) (*elasticbeanstalk.EnvironmentDescriptionsMessage, error)
	DescribeEnvironmentHealthWithContext(ctx aws.Context, input *elasticbeanstalk.DescribeEnvironmentHealthInput, opts ...request.Option) (*elasticbeanstalk.DescribeEnvironmentHealthOutput, error)
	DescribeInstancesHealthWithContext(ctx aws.Context, input *elasticbeanstalk.DescribeInstancesHealthInput, opts ...request.Option) (*elasticbeanstalk.DescribeInstancesHealthOutput, error)
}

func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	beanapi := elasticbeanstalk.New(p, opts...)
	return api.Endpoint{
		"DescribeApplications":        &DescribeApplications{beanapi},
		"DescribeApplicationVersions": &DescribeApplicationVersions{beanapi},
		"DescribeEnvironmentHealth":   &DescribeEnvironmentHealth{beanapi},
		"DescribeEnvironments":        &DescribeEnvironments{beanapi},
		"DescribeInstancesHealth":     &DescribeInstancesHealth{beanapi},
	}
}
