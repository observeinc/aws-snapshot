package route53

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversionListBuckets(t *testing.T) {
	var o ListHostedZonesOutput
	apitest.TestSource(t, &apitest.TestSourceConfig{
		Source:     &o,
		InputFile:  "testdata/listhostedzones-input.json",
		OutputFile: "testdata/listhostedzones-expect.json",
		Update:     *update,
	})
}
