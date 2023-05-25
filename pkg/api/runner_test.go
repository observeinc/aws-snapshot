package api_test

import (
	"context"
	"testing"
	"time"

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

func testLongRequest(ctx context.Context, ch chan<- *api.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		time.Sleep(10 * time.Millisecond)
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

func TestRunnerTimeout(t *testing.T) {
	var recorder apitest.Recorder

	timeout := time.Second

	r := api.Runner{
		Requests: []api.Request{
			testLongRequest,
		},
		MaxConcurrentRequests: 1,
		ConcurrentRecorders:   1,
		Recorder:              &recorder,
		RequestTimeout:        &timeout,
	}

	if err := r.Run(context.Background()); err == nil {
		t.Fatal("expected timeout error!")
	}

	if len(recorder.Records) != 0 {
		t.Fatal("no records expected")
	}
}
