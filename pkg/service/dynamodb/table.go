package dynamodb

import (
	"context"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DescribeTableOutput struct {
	*dynamodb.DescribeTableOutput
}

func (o *DescribeTableOutput) Records() (records []*api.Record) {
	if t := o.Table; t != nil {
		records = append(records, &api.Record{
			ID:   t.TableArn,
			Data: t,
		})
	}
	return
}

type DescribeTable struct {
	API
}

var _ api.RequestBuilder = &DescribeTable{}

// New implements api.RequestBuilder
func (fn *DescribeTable) New(name string, config interface{}) ([]api.Request, error) {
	var input dynamodb.ListTablesInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		var innerErr, outerErr error

		outerErr = fn.ListTablesPagesWithContext(ctx, &input, func(output *dynamodb.ListTablesOutput, last bool) bool {
			for _, tableName := range output.TableNames {
				describeTableInput := &dynamodb.DescribeTableInput{
					TableName: tableName,
				}

				describeTableOutput, err := fn.DescribeTableWithContext(ctx, describeTableInput)
				if err != nil {
					innerErr = err
					return false
				}

				if err := api.SendRecords(ctx, ch, name, &DescribeTableOutput{describeTableOutput}); err != nil {
					innerErr = err

					// failed to send records, stop handling tables
					return false
				}
			}
			return true
		})

		return api.FirstError(outerErr, innerErr)
	}

	return []api.Request{call}, nil
}
