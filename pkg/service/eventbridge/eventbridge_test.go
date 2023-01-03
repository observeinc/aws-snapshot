package eventbridge

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {

	testcases := []*apitest.TestSourceConfig{
		{
			Source:     &ListEventBusesOutput{},
			InputFile:  "testdata/listeventbuses-input.json",
			OutputFile: "testdata/listeventbuses-expect.json",
			Update:     *update,
		},
		{
			Source:     &ListApiDestinationsOutput{},
			InputFile:  "testdata/listapidestinations-input.json",
			OutputFile: "testdata/listapidestinations-expect.json",
			Update:     *update,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.InputFile, func(t *testing.T) {
			apitest.TestSource(t, tt)
		})
	}

}
