package dynamodb

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"

	"github.com/aws/aws-sdk-go/aws"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &DescribeTableOutput{},
			InputFile:  "testdata/describetable-input.json",
			OutputFile: "testdata/describetable-expect.json",
			Update:     *update,
		},
		{
			Source: &ScanOutput{
				TableArn:  aws.String("arn:aws:dynamodb:us-west-1:123456789012:table/example"),
				TableKeys: []*string{aws.String("LockID")},
			},
			InputFile:  "testdata/scan-input.json",
			OutputFile: "testdata/scan-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
