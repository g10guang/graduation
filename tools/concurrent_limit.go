package tools

import "sync/atomic"

// Control concurrency
type ConcurrentLimit struct {
	NoCopy
	lim int64  		// total threshold
	now int64		// current concurrency
	tmp int64		// use for take one test is success or not
}

func NewConcurrentLimit(limit int64) *ConcurrentLimit {
	return &ConcurrentLimit{
		lim: limit,
	}
}

// take a token from pool
// return success or not
func (c *ConcurrentLimit) TakeOne() bool {
	newVal := atomic.AddInt64(&c.tmp, 1)
	if newVal < atomic.LoadInt64(&c.lim) {
		return true
	}
	atomic.AddInt64(&c.tmp, -1)
	return false
}

// release a token into pool
func (c *ConcurrentLimit) ReleaseOne() {
	atomic.AddInt64(&c.now, -1)
	atomic.AddInt64(&c.tmp, -1)
}

func (c *ConcurrentLimit) UpdateLimit(limit int64) {
	atomic.StoreInt64(&c.lim, limit)
}

func (c *ConcurrentLimit) Now() int64 {
	return atomic.LoadInt64(&c.now)
}
