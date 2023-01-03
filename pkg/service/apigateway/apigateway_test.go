package apigateway

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &GetRestApisOutput{},
			InputFile:  "testdata/getrestapis-input.json",
			OutputFile: "testdata/getrestapis-expect.json",
			Update:     *update,
		},
		{
			Source:     &GetDeploymentsOutput{},
			InputFile:  "testdata/getdeployments-input.json",
			OutputFile: "testdata/getdeployments-expect.json",
			Update:     *update,
		},
		{
			Source:     &GetStagesOutput{},
			InputFile:  "testdata/getstages-input.json",
			OutputFile: "testdata/getstages-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}
}
