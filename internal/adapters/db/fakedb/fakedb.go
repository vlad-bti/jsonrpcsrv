package fakedb

import (
	"context"
	"sync"
)

type fakeDB struct {
	mu sync.Mutex
}

func NewFakeDB() *fakeDB {
	return &fakeDB{}
}

func (r *fakeDB) Begin(ctx context.Context) {
	r.mu.Lock()
}

func (r *fakeDB) Commit(ctx context.Context) {
	r.mu.Unlock()
}

func (r *fakeDB) Rollback(ctx context.Context) {
	r.mu.Unlock()
}
