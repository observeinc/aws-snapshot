package apitest

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api"
)

type TestSourceConfig struct {
	api.Source
	InputFile  string
	OutputFile string
	Update     bool
}

// TestSource compares the result of converting a given struct to api.Records
func TestSource(t *testing.T, tt *TestSourceConfig) {
	t.Helper()
	var expected []*api.Record

	UnmarshalFile(t, tt.InputFile, &tt.Source)
	actual := tt.Source.Records()
	if tt.Update {
		MarshalFile(t, actual, tt.OutputFile)
	}
	UnmarshalFile(t, tt.OutputFile, &expected)

	if diff := Diff(actual, expected); diff != "" {
		t.Logf("input=%q, output=%q", tt.InputFile, tt.OutputFile)
		t.Fatalf(diff)
	}
}

func UnmarshalFile(t *testing.T, filename string, o interface{}) {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %s: %s", filename, err)
	}

	if err := json.Unmarshal(data, o); err != nil {
		t.Fatalf("failed to unmarshal %s: %s", filename, err)
	}
}

func MarshalFile(t *testing.T, o interface{}, filename string) {
	t.Helper()

	data, err := json.MarshalIndent(o, "  ", "  ")
	if err != nil {
		t.Fatalf("failed to marshal %s: %s", filename, err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		t.Fatalf("failed to write to %s: %s", filename, err)
	}
}
