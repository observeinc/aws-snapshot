package api_test

import (
	"context"
	"errors"
	"fmt"
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
	case <-time.After(10 * time.Second):
		return nil
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
			func(ctx context.Context, ch chan<- *api.Record) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(10 * time.Second):
					return nil
				}
			},
		},
		Recorder:       &recorder,
		RequestTimeout: &timeout,
	}

	err := r.Run(context.Background())

	if err == nil {
		t.Fatal("expected error!")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(fmt.Sprintf("wrong error type, expected DeadlineExceeded, got %s", err))
	}

	if len(recorder.Records) != 0 {
		t.Fatal("no records expected")
	}
}
