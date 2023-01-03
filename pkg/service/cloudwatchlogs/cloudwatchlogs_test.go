package cloudwatchlogs

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversion(t *testing.T) {
	var o DescribeLogGroupsOutput
	apitest.TestSource(t, &apitest.TestSourceConfig{
		Source:     &o,
		InputFile:  "testdata/input.json",
		OutputFile: "testdata/expect.json",
		Update:     *update,
	})
}
