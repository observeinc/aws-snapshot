package dynamodb

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// AttributeValue is used to marshal dynamodb.AttributeValue with omitempty
type AttributeValue struct {
	B    []byte                     `json:"B,omitempty"`
	BOOL *bool                      `json:"BOOL,omitempty"`
	BS   [][]byte                   `json:"BS,omitempty"`
	L    []*AttributeValue          `json:"L,omitempty"`
	M    map[string]*AttributeValue `json:"M,omitempty"`
	N    *string                    `json:"N,omitempty"`
	NS   []*string                  `json:"NS,omitempty"`
	NULL *bool                      `json:"NULL,omitempty"`
	S    *string                    `json:"S,omitempty"`
	SS   []*string                  `json:"SS,omitempty"`
}

func (a *AttributeValue) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	a.B = av.B
	a.BOOL = av.BOOL
	a.BS = av.BS
	a.N = av.N
	a.NS = av.NS
	a.NULL = av.NULL
	a.S = av.S
	a.SS = av.SS

	if av.L != nil {
		if err := dynamodbattribute.UnmarshalList(av.L, &a.L); err != nil {
			return fmt.Errorf("failed to unmarshal embedded list: %w", err)
		}
	}
	if av.M != nil {
		if err := dynamodbattribute.UnmarshalMap(av.M, &a.M); err != nil {
			return fmt.Errorf("failed to unmarshal embedded map: %w", err)
		}
	}
	return nil
}

type ScanOutput struct {
	*dynamodb.ScanOutput
	TableArn  *string
	TableKeys []*string
}

// Records constructs valid observations from which we can build a resource
// Scan gives us individual items, but we need to know the key in order to
// makeresource. This leads to some additional bloat, but still simpler than
// stitching things back together from the table resource in OPAL
func (o *ScanOutput) Records() (records []*api.Record) {

	var items []map[string]AttributeValue
	_ = dynamodbattribute.UnmarshalListOfMaps(o.ScanOutput.Items, &items)
	for _, item := range items {

		keys := make(map[string]interface{})
		for _, keyName := range o.TableKeys {
			keys[*keyName] = item[*keyName]
		}

		records = append(records, &api.Record{
			ID: o.TableArn,
			Data: map[string]interface{}{
				"Keys":     keys,
				"NewImage": item,
			},
		})
	}
	return
}

type Scan struct {
	API
}

var _ api.RequestBuilder = &Scan{}

// New implements api.RequestBuilder
func (fn *Scan) New(name string, config interface{}) ([]api.Request, error) {
	var input dynamodb.ScanInput
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}

	if input.TableName == nil {
		// skip action if no table name is provided
		return nil, nil
	}

	call := func(ctx context.Context, ch chan<- *api.Record) error {
		describeTableOutput, err := fn.DescribeTableWithContext(ctx, &dynamodb.DescribeTableInput{
			TableName: input.TableName,
		})
		if err != nil {
			return fmt.Errorf("failed to describe table: %w", err)
		}

		var tableKeys []*string
		for _, element := range describeTableOutput.Table.KeySchema {
			tableKeys = append(tableKeys, element.AttributeName)
		}

		return fn.ScanPagesWithContext(ctx, &input, func(output *dynamodb.ScanOutput, last bool) bool {
			return api.SendRecords(ctx, ch, name, &ScanOutput{
				ScanOutput: output,
				TableArn:   describeTableOutput.Table.TableArn,
				TableKeys:  tableKeys,
			})
		})
	}

	return []api.Request{call}, nil
}
