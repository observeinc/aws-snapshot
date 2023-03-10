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
		return fn.ListTablesPagesWithContext(ctx, &input, func(output *dynamodb.ListTablesOutput, last bool) bool {
			for _, tableName := range output.TableNames {
				describeTableInput := &dynamodb.DescribeTableInput{
					TableName: tableName,
				}

				describeTableOutput, err := fn.DescribeTableWithContext(ctx, describeTableInput)
				if err != nil {
					panic(err)
				}

				if !api.SendRecords(ctx, ch, name, &DescribeTableOutput{describeTableOutput}) {
					// failed to send records, stop handling tables
					return false
				}
			}
			return true
		})
	}

	return []api.Request{call}, nil
}
