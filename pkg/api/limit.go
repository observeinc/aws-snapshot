package api

import (
	"sync/atomic"
)

type limitCh struct {
	in    chan *Record
	out   chan *Record
	count atomic.Int64
	limit int64
}

func (l *limitCh) In() chan<- *Record {
	return l.in
}

func (l *limitCh) Out() <-chan *Record {
	return l.out
}

func (l *limitCh) Close() {
	close(l.in)
}

func (l *limitCh) run(cancel func()) {
	defer close(l.out)
	for {
		v, ok := <-l.in
		if !ok {
			return
		}
		l.out <- v
		if l.count.Add(1) == l.limit {
			cancel()
			return
		}
	}
}

func newLimitCh(bufferSize, limit int, cancel func()) *limitCh {
	if limit == 0 {
		ch := make(chan *Record, bufferSize)
		return &limitCh{
			in:  ch,
			out: ch,
		}
	}

	l := &limitCh{
		in:    make(chan *Record, bufferSize),
		out:   make(chan *Record, bufferSize),
		limit: int64(limit),
	}

	go l.run(cancel)
	return l
}
