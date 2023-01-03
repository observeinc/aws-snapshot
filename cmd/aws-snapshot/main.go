package main

/*

aws-snapshot is a utility to quickly exercise the polling code that is part of our
lambda.

*/

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/observeinc/aws-snapshot/pkg/api"
	"github.com/observeinc/aws-snapshot/pkg/service"
	_ "github.com/observeinc/aws-snapshot/pkg/service/all"

	"github.com/aws/aws-sdk-go/aws/session"
)

func loadManifest(filename string, dst *api.Manifest) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func realMain() error {
	var (
		maxConcurrentRequests = flag.Int("max-concurrent-requests", 10, "Maximum concurrent requests")
		bufferSize            = flag.Int("buffer-size", 100, "Length of buffer for records")
		manifestFile          = flag.String("manifest-file", "", "Manifest filename")
	)

	flag.Parse()

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable, // To enable loading AWS region from .aws/config
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	s := service.New(sess)

	var manifest api.Manifest
	if *manifestFile != "" {
		if err = loadManifest(*manifestFile, &manifest); err != nil {
			return fmt.Errorf("failed to load manifest from file: %w", err)
		}
	}

	reqs, err := s.Resolve(&manifest)
	if err != nil {
		return fmt.Errorf("failed to resolve manifest: %w", err)
	}

	runner := api.Runner{
		Requests:              reqs,
		MaxConcurrentRequests: *maxConcurrentRequests,
		BufferSize:            *bufferSize,
		Recorder: api.RecorderFunc(func(ctx context.Context, ch <-chan *api.Record) error {
			for {
				r, ok := <-ch
				if !ok {
					return nil
				}
				data, _ := json.Marshal(r)
				fmt.Println(string(data))
			}
		}),
	}

	return runner.Run(context.Background())
}

func main() {
	if err := realMain(); err != nil {
		panic(err)
	}
}
