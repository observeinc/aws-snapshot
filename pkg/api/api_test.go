package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/api/apitest"

	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	errSentinel = errors.New("this is a test")
)

type testRequestConfig struct {
	ErrBuild   error `json:"errBuild,omitempty"`
	ErrRequest error `json:"errRequest,omitempty"`
}

type TestRequestBuilder struct{}

func (t *TestRequestBuilder) New(name string, v interface{}) ([]api.Request, error) {
	var config testRequestConfig
	if err := api.DecodeConfig(v, &config); err != nil {
		return nil, err
	}

	if config.ErrBuild != nil {
		return nil, config.ErrBuild
	}

	var reqs []api.Request

	reqs = append(reqs, func(ctx context.Context, ch chan<- *api.Record) error {
		if config.ErrRequest != nil {
			return config.ErrRequest
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- &api.Record{Action: name}:
			return nil
		}
	})

	return reqs, nil
}

func TestEndpoint(t *testing.T) {
	builder := &TestRequestBuilder{}

	e := api.Endpoint{
		"DescribeA": builder,
		"DescribeB": builder,
		"DescribeC": builder,
		"ListD":     builder,
	}

	type TestResult struct {
		NumRequests int
		Records     []*api.Record
		ErrBuild    error
		ErrRun      error
	}

	testcases := map[string]struct {
		Manifest api.Manifest
		Expect   TestResult
	}{
		"Wildcard All": {
			// default manifest includes all
			Expect: TestResult{
				NumRequests: 4,
				Records: []*api.Record{
					{Action: "DescribeA"},
					{Action: "DescribeB"},
					{Action: "DescribeC"},
					{Action: "ListD"},
				},
			},
		},
		"Partial Wildcard": {
			Manifest: api.Manifest{
				Include: []string{
					"Describe*",
				},
			},
			Expect: TestResult{
				NumRequests: 3,
				Records: []*api.Record{
					{Action: "DescribeA"},
					{Action: "DescribeB"},
					{Action: "DescribeC"},
				},
			},
		},
		"Exclude": {
			Manifest: api.Manifest{
				Exclude: []string{
					"Describe*",
				},
			},
			Expect: TestResult{
				NumRequests: 1,
				Records: []*api.Record{
					{Action: "ListD"},
				},
			},
		},
		"Include With Exclude": {
			Manifest: api.Manifest{
				Include: []string{
					"Describe*",
				},
				Exclude: []string{
					"DescribeB",
				},
			},
			Expect: TestResult{
				NumRequests: 2,
				Records: []*api.Record{
					{Action: "DescribeA"},
					{Action: "DescribeC"},
				},
			},
		},
		"Build Error": {
			Manifest: api.Manifest{
				Include: []string{
					"DescribeA",
				},
				Overrides: []*api.Config{
					{
						Action: "DescribeA",
						Config: map[string]interface{}{
							"ErrBuild": errSentinel,
						},
					},
				},
			},
			Expect: TestResult{
				ErrBuild: errSentinel,
			},
		},
		"Run Error": {
			Manifest: api.Manifest{
				Include: []string{
					"Describe*",
				},
				Overrides: []*api.Config{
					{
						Action: "DescribeA",
						Config: map[string]interface{}{
							"ErrRequest": errSentinel,
						},
					},
				},
			},
			Expect: TestResult{
				NumRequests: 3,
				ErrRun:      errSentinel,
			},
		},
		"Config Error": {
			Manifest: api.Manifest{
				Include: []string{
					"Describe*",
				},
				Overrides: []*api.Config{
					{
						Action: "DescribeA",
						Config: map[string]interface{}{
							"NoSuchAttr": false,
						},
					},
				},
			},
			Expect: TestResult{
				// TODO: look specifically for config error
				ErrBuild: cmpopts.AnyError,
			},
		},
	}

	for name, tt := range testcases {
		t.Run(name, func(t *testing.T) {
			var actual TestResult

			manifest := tt.Manifest
			reqs, err := e.Resolve(&manifest)

			if err != nil {
				actual.ErrBuild = err
			} else {
				actual.NumRequests = len(reqs)

				var r apitest.Recorder
				runner := api.Runner{
					Requests:              reqs,
					MaxConcurrentRequests: 1,
					Recorder:              &r,
				}
				if err := runner.Run(context.Background()); err != nil {
					actual.ErrRun = err
				} else {
					actual.Records = r.Records
				}
			}

			if diff := apitest.Diff(tt.Expect, actual); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
