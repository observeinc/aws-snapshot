package api

import (
	"testing"
)

func TestLimitCh(t *testing.T) {

	var cancelled bool
	ch := newLimitCh(100, 10, func() {
		cancelled = true
	})

	for i := 0; i <= 100; i++ {
		ch.In() <- &Record{}
	}
	ch.Close()

	var records []*Record
	for r := range ch.Out() {
		records = append(records, r)
	}

	if !cancelled {
		t.Error("cancel callback not invoked")
	}

	if len(records) != 10 {
		t.Error("received wrong number of records")
	}
}
