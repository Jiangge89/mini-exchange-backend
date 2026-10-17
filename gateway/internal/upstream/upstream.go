package upstream

import (
	"sync"
	"time"
)

type Upstream struct {
	Name string
	URL  string

	mu          sync.RWMutex
	blockHeight uint64
	errCount    int64
	totalCount  int64
	latencySum  time.Duration
	latencyN    int64
}

func (u *Upstream) RecordBlockHeight(height uint64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if height > u.blockHeight {
		u.blockHeight = height
	}
}

func (u *Upstream) RecordError() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.errCount++
	u.totalCount++
}

func (u *Upstream) RecordSuccess(latency time.Duration) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.totalCount++
	u.latencySum += latency
	u.latencyN++
}

func (u *Upstream) Stats() (blockHeight uint64, errorRate float64, avgLatency time.Duration) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	blockHeight = u.blockHeight
	if u.totalCount > 0 {
		errorRate = float64(u.errCount) / float64(u.totalCount)
	}
	if u.latencyN > 0 {
		avgLatency = u.latencySum / time.Duration(u.latencyN)
	}
	return
}

type Entry struct {
	Name string
	URL  string
}

type Pool struct {
	upstreams []*Upstream
}

func NewPool(entries []Entry) *Pool {
	pool := &Pool{}
	for _, e := range entries {
		pool.upstreams = append(pool.upstreams, &Upstream{Name: e.Name, URL: e.URL})
	}
	return pool
}

func (p *Pool) All() []*Upstream {
	return p.upstreams
}
