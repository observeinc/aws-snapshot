/*
Package apitest provides utilities to test api
*/
package apitest

import (
	"context"
	"encoding/json"
	"sort"
	"sync"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type sortedRecords []*api.Record

func (r sortedRecords) Len() int {
	return len(r)
}

func (r sortedRecords) Less(i, j int) bool {
	switch {
	case r[i] == nil:
		return true
	case r[j] == nil:
		return false
	case r[i].Action < r[j].Action:
		return true
	case r[i].Action > r[j].Action:
		return false
	}
	return true
}

func (r sortedRecords) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Records are equivalent if JSON representation is the same.
// This allows us to compare records with native structs against golden files
// which are unmarshaled as map[string]interface{}.
var cmpRecords = cmp.Transformer("Sort", func(in []*api.Record) []*api.Record {
	out := make(sortedRecords, len(in))
	for _, r := range in {
		var cp api.Record
		data, _ := json.Marshal(r)
		json.Unmarshal(data, &cp)
		out = append(out, &cp)
	}
	sort.Sort(out)
	return out
})

// Diff which handles Record comparison in sorted order
func Diff(a interface{}, b interface{}) string {
	return cmp.Diff(a, b, cmpRecords, cmpopts.EquateErrors())
}

type Recorder struct {
	Records []*api.Record
	sync.Mutex
}

func (r *Recorder) ReadFrom(ctx context.Context, ch <-chan *api.Record) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case record, ok := <-ch:
			if !ok {
				return nil
			}
			r.Lock()
			r.Records = append(r.Records, record)
			r.Unlock()
		}
	}
}
