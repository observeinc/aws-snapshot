package cloudwatchlogs

import (
	"flag"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

var update = flag.Bool("update", false, "update result files")

func TestConversionListAccountAliases(t *testing.T) {
	var o ListAccountAliasesOutput
	apitest.TestSource(t, &apitest.TestSourceConfig{
		Source:     &o,
		InputFile:  "testdata/listaccountaliases-input.json",
		OutputFile: "testdata/listaccountaliases-expect.json",
		Update:     *update,
	})
}

func TestConversionGetAccountAuthorizationDetails(t *testing.T) {
	var o GetAccountAuthorizationDetailsOutput
	apitest.TestSource(t, &apitest.TestSourceConfig{
		Source:     &o,
		InputFile:  "testdata/getaccountauthorizationdetails-input.json",
		OutputFile: "testdata/getaccountauthorizationdetails-expect.json",
		Update:     *update,
	})
}
