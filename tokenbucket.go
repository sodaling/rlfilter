package rlfilter

import (
	"net/http"
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	rate, allowance, max, unit, lastCheck uint64
}

func (t *TokenBucket) Limit(handler WebHandler) WebHandler {
	return func(resp http.ResponseWriter, req *http.Request) {
		if !t.limit() {
			handler(resp, req)
		} else {
			http.Error(resp, "you have been ban", http.StatusForbidden)
		}
	}
}

func NewTokenBucket(rate int, per time.Duration) Limiter {
	// 单位时间发放令牌数
	if rate < 1 {
		rate = 1
	}
	// 单位时间
	nano := uint64(per)
	if nano < 1 {
		nano = uint64(time.Second)
	}

	return &TokenBucket{
		rate: uint64(rate),
		// 初始化时候装满桶
		allowance: uint64(rate) * nano,
		max:       uint64(rate) * nano,
		unit:      nano,
		lastCheck: unixNano(),
	}
}

func (t *TokenBucket) limit() bool {
	now := unixNano()
	passed := now - atomic.SwapUint64(&t.lastCheck, now)

	rate := atomic.LoadUint64(&t.rate)
	current := atomic.AddUint64(&t.allowance, rate*passed)
	if max := atomic.LoadUint64(&t.max); current > max {
		atomic.SwapUint64(&t.allowance, max)
		current = max
	}

	if current < t.unit {
		return true
	}
	atomic.AddUint64(&t.allowance, -t.unit)
	return false
}

func (t *TokenBucket) UpdateRate(rate int) {
	atomic.StoreUint64(&t.rate, uint64(rate))
	atomic.StoreUint64(&t.max, uint64(rate)*t.unit)
}

func unixNano() uint64 {
	return uint64(time.Now().UnixNano())
}

func (t *TokenBucket) Undo() {
	current := atomic.AddUint64(&t.allowance, t.unit)
	if max := atomic.LoadUint64(&t.max); current > max {
		atomic.StoreUint64(&t.allowance, max)
	}
}
