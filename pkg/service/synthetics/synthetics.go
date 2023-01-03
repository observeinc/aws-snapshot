package synthetics

import (
	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/synthetics"
)

func init() {
	service.Register("synthetics", api.ServiceFunc(New))
}

// API documents the subset of AWS API we actually call
type API interface {
	DescribeCanariesPagesWithContext(ctx aws.Context, input *synthetics.DescribeCanariesInput, fn func(*synthetics.DescribeCanariesOutput, bool) bool, opts ...request.Option) error
}

// New implements api.ServiceFunc
func New(p client.ConfigProvider, opts ...*aws.Config) api.Endpoint {
	syntheticsapi := synthetics.New(p, opts...)
	return api.Endpoint{
		"DescribeCanaries": &DescribeCanaries{syntheticsapi},
	}
}
