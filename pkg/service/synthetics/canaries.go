package synthetics

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/synthetics"
)

type DescribeCanariesOutput struct {
	*synthetics.DescribeCanariesOutput
}

func (o *DescribeCanariesOutput) Records() (records []*api.Record) {
	for _, c := range o.Canaries {
		records = append(records, &api.Record{
			// XXX: api endpoint does not return an ARN
			ID:   c.Id,
			Data: c,
		})
	}
	return
}

type DescribeCanaries struct {
	API
}

var _ api.RequestBuilder = &DescribeCanaries{}

// New implements api.RequestBuilder
func (fn *DescribeCanaries) New(name string, config interface{}) ([]api.Request, error) {
	var input synthetics.DescribeCanariesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var outerErr, innerErr error

		outerErr = fn.DescribeCanariesPagesWithContext(ctx, &input, func(output *synthetics.DescribeCanariesOutput, last bool) bool {
			if err := api.SendRecords(ctx, ch, name, &DescribeCanariesOutput{output}); err != nil {
				innerErr = err
				return false
			}

			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
