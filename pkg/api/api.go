/*
Package api provides a framework for scraping the AWS API.
*/
package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/go-logr/logr"
)

// Service (e.g. RDS, IAM, S3) provides the means of creating an Endpoint
// The client and config options provided end up determining region.
type Service interface {
	New(p client.ConfigProvider, opts ...*aws.Config) Endpoint
}

// ServiceFunc implements Service.
type ServiceFunc func(p client.ConfigProvider, opts ...*aws.Config) Endpoint

func (s ServiceFunc) New(p client.ConfigProvider, opts ...*aws.Config) Endpoint {
	return s(p, opts...)
}

// Record is submitted as an observation.
type Record struct {
	Timestamp *int64      `json:"timestamp,omitempty"` // Optional nanosecond timestamp.
	Action    string      `json:"action"`              // Action performed. Naming should match IAM policy.
	ID        *string     `json:"id,omitempty"`        // ID of resource contained in Data. Optional
	Data      interface{} `json:"data"`
}

// Request runs synchronously, publishing Records to a channel.
type Request func(context.Context, chan<- *Record) error

// Recorder reads from Record channel, does stuff.
type RecorderFunc func(context.Context, <-chan *Record) error

func (r RecorderFunc) ReadFrom(ctx context.Context, ch <-chan *Record) error {
	return r(ctx, ch)
}

type Recorder interface {
	ReadFrom(context.Context, <-chan *Record) error
}

// RequestBuilder translates a requested action and config into requests.
// A given action can translate into multiple requests in a variety of cases:
// - the action is wildcarded (e.g. Describe*)
// - the action spawns multiple sub-requests (e.g. DescribeLogStream needs to be run for every LogGroup).
type RequestBuilder interface {
	New(action string, config interface{}) ([]Request, error)
}

func prefixError(name string, reqs []Request) []Request {
	res := []Request{}

	for _, r := range reqs {
		rq := r

		res = append(res, func(ctx context.Context, ch chan<- *Record) (err error) {
			logger := logr.FromContextOrDiscard(ctx)
			logger.V(6).Info("request start", "action", name)
			defer func() {
				logger.V(6).Info("request complete", "action", name, "error", err)
			}()

			if err = rq(ctx, ch); err != nil {
				err = fmt.Errorf("failed to run %q: %w", name, err)
			}

			return
		})
	}

	return res
}

// Endpoint is a collection of RequestBuilders.
type Endpoint map[string]RequestBuilder

// New allows Endpoint to act as a RequestBuilder, with wildcard support on action.
func (e *Endpoint) New(action string, config interface{}) ([]Request, error) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to cast endpoint configuration: %#v", config)
	}

	if e == nil {
		return nil, nil
	}

	var reqs []Request
	for reqName, reqBuilder := range *e {
		if match(reqName, action) {
			// we pass on the reqName, rather than the pattern it matched
			req, err := reqBuilder.New(reqName, configMap[reqName])
			if err != nil {
				return nil, err
			}
			reqs = append(reqs, prefixError(reqName, req)...)
		}
	}
	return reqs, nil
}

// Filter returns an endpoint with subset of RequestBuilders.
func (e *Endpoint) Filter(pattern string) Endpoint {
	out := make(Endpoint)
	for k, v := range *e {
		if match(k, pattern) {
			out[k] = v
		}
	}
	return out
}

// match does the most rudimentary form of glob matching possible.
func match(name, pattern string) bool {
	if !strings.HasSuffix(pattern, "*") {
		return name == pattern
	}

	pattern = pattern[:len(pattern)-1]
	return strings.HasPrefix(name, pattern)
}

type Config struct {
	Action string                 `json:"action"`
	Config map[string]interface{} `json:"config"`
}

// Manifest describes desired actions and any necessary configuration overrides.
type Manifest struct {
	Include   []string  `json:"include"`   // Include is a list of actions
	Exclude   []string  `json:"exclude"`   // Exclude is a list of actions
	Overrides []*Config `json:"overrides"` // Overrides are key value pairs of configuration data.
}

// Resolve translates a configuration manifest into requests.
func (e *Endpoint) Resolve(m *Manifest) ([]Request, error) {
	if len(m.Include) == 0 {
		m.Include = append(m.Include, "*")
	}

	filtered := make(Endpoint)
	for _, pattern := range m.Include {
		for k, v := range e.Filter(pattern) {
			filtered[k] = v
		}
	}

	for _, pattern := range m.Exclude {
		for k := range filtered.Filter(pattern) {
			delete(filtered, k)
		}
	}

	overrides := make(map[string]interface{}, len(m.Overrides))
	for _, o := range m.Overrides {
		overrides[o.Action] = o.Config
	}

	return filtered.New("*", overrides)
}

type Source interface {
	Records() []*Record
}

func SendRecords(ctx context.Context, ch chan<- *Record, name string, s Source) error {
	ts := time.Now().UnixNano()

	for _, r := range s.Records() {
		r.Action = name
		r.Timestamp = &ts
		select {
		case ch <- r:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func FirstError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
