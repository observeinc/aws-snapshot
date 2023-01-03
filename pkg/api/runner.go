package api

import (
	"context"
)

const defaultBufferSize = 100

type Runner struct {
	Requests              []Request
	Recorder              Recorder
	BufferSize            int
	MaxConcurrentRequests int
	ConcurrentRecorders   int
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

		var err error
		for i := 0; i < num; i++ {
			if e := <-errCh; e != nil && err == nil {
				err = e
			}
		}

		return err
	})
}

// withSemaphore runs up to maxConcurrency Requests at a time.
func withSemaphore(fns []Request, maxConcurrency int) Request {
	if maxConcurrency == 0 {
		maxConcurrency = len(fns)
	}

	return func(ctx context.Context, ch chan<- *Record) error {
		var (
			errCh = make(chan error, len(fns))
			sem   = make(chan struct{}, maxConcurrency)
			err   error
		)

		defer close(errCh)
		defer close(sem)

		for _, fn := range fns {
			sem <- struct{}{}
			go func(fn Request) {
				defer func() {
					<-sem
				}()
				errCh <- fn(ctx, ch)
			}(fn)
		}

		for i := 0; i < len(fns); i++ {
			if e := <-errCh; e != nil && err == nil {
				err = e
			}
		}
		return err
	}
}

func (r *Runner) Run(ctx context.Context) error {
	requestFunc := withSemaphore(r.Requests, r.MaxConcurrentRequests)
	recorder := pool(r.Recorder, r.ConcurrentRecorders)

	if r.BufferSize < 1 {
		r.BufferSize = defaultBufferSize
	}

	var (
		requestCh = make(chan *Record, r.BufferSize)
		readErrCh = make(chan error, 1)
	)

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	go func() {
		defer close(readErrCh)
		defer close(requestCh)
		readErrCh <- requestFunc(ctx, requestCh)
	}()

	err := recorder.ReadFrom(ctx, requestCh)
	if err != nil {
		cancelFunc()
		<-readErrCh
		return err
	}

	return <-readErrCh
}
