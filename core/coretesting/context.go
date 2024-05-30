package coretesting

import (
	"context"
	"time"

	"cosmossdk.io/core/store"
)

type dummyKey struct{}

func Context() context.Context {
	dummy := &dummyCtx{
		stores: map[string]store.KVStore{},
	}

	ctx := context.WithValue(context.Background(), dummyKey{}, dummy)
	return ctx
}

type dummyCtx struct {
	stores map[string]store.KVStore
}

func (d dummyCtx) Deadline() (deadline time.Time, ok bool) {
	panic("Deadline on dummy context")
}

func (d dummyCtx) Done() <-chan struct{} {
	panic("Done on dummy context")
}

func (d dummyCtx) Err() error {
	panic("Err on dummy context")
}

func (d dummyCtx) Value(key any) any {
	panic("Value on dummy context")
}

func unwrap(ctx context.Context) *dummyCtx {
	dummy := ctx.Value(dummyKey{})
	if dummy == nil {
		panic("invalid ctx without dummy")
	}

	return dummy.(*dummyCtx)
}
