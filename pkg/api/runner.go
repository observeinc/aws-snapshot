package api

import (
	"context"
	"time"

	"github.com/go-logr/logr"
)

const defaultBufferSize = 100

type Runner struct {
	Requests              []Request
	Recorder              Recorder
	BufferSize            int
	MaxConcurrentRequests int
	MaxRecords            int
	ConcurrentRecorders   int
	RequestTimeout        *time.Duration
	Logger                *logr.Logger
}

// Pool runs num copies of Recorder.
func pool(fn Recorder, num int) Recorder {
	if num <= 1 {
		return fn
	}

	return RecorderFunc(func(ctx context.Context, ch <-chan *Record) error {
		errCh := make(chan error, num)
		defer close(errCh)

		for i := 0; i < num; i++ {
			go func() {
				errCh <- fn.ReadFrom(ctx, ch)
			}()
		}

		var errs []error
		for i := 0; i < num; i++ {
			if e := <-errCh; e != nil {
				errs = append(errs, e)
			}
		}

		return Join(errs...)
	})
}

// withSemaphore runs up to maxConcurrency Requests at a time.
func withSemaphore(fns []Request, maxConcurrency int, requestTimeout *time.Duration) Request {
	if maxConcurrency == 0 {
		maxConcurrency = len(fns)
	}

	return func(ctx context.Context, ch chan<- *Record) error {
		var (
			errCh = make(chan error, len(fns))
			sem   = make(chan struct{}, maxConcurrency)
		)

		defer close(errCh)
		defer close(sem)

		for _, fn := range fns {
			sem <- struct{}{}
			go func(ctx context.Context, fn Request) {
				defer func() {
					<-sem
				}()

				var cancel context.CancelFunc

				if requestTimeout != nil {
					ctx, cancel = context.WithTimeout(ctx, *requestTimeout)
					defer cancel()
				}

				errCh <- fn(ctx, ch)
			}(ctx, fn)
		}

		var errs []error
		for i := 0; i < len(fns); i++ {
			if e := <-errCh; e != nil {
				errs = append(errs, e)
			}
		}

		return Join(errs...)

	}
}

func (r *Runner) Run(ctx context.Context) error {
	if r.Logger != nil {
		ctx = logr.NewContext(ctx, *r.Logger)
	}

	requestFunc := withSemaphore(r.Requests, r.MaxConcurrentRequests, r.RequestTimeout)
	recorder := pool(r.Recorder, r.ConcurrentRecorders)

	if r.BufferSize < 1 {
		r.BufferSize = defaultBufferSize
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	var (
		limitCh   = newLimitCh(r.BufferSize, r.MaxRecords, cancelFunc)
		readErrCh = make(chan error, 1)
	)

	go func() {
		defer close(readErrCh)
		defer limitCh.Close()
		readErrCh <- requestFunc(ctx, limitCh.In())
	}()

	err := recorder.ReadFrom(ctx, limitCh.Out())
	if err != nil {
		cancelFunc()
		<-readErrCh
		return err
	}

	return <-readErrCh
}
