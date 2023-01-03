package api_test

import (
	"context"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/api/apitest"
)

func testRequest(ctx context.Context, ch chan<- *api.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		ch <- &api.Record{}
	}
	return nil
}

func TestRunner(t *testing.T) {
	var recorder apitest.Recorder

	r := api.Runner{
		Requests: []api.Request{
			testRequest,
			testRequest,
			testRequest,
			testRequest,
			testRequest,
		},
		MaxConcurrentRequests: 2,
		ConcurrentRecorders:   2,
		Recorder:              &recorder,
	}

	if err := r.Run(context.Background()); err != nil {
		t.Fatal(err)
	}

	if len(recorder.Records) != 5 {
		t.Fatal("wrong number of records")
	}
}
