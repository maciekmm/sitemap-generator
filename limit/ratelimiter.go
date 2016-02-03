package limit

import (
	"errors"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	ticker  *time.Ticker
	burst   chan time.Time
	started uint32
}

//NewRateLimitedWebClient constructs RateLimitedWebClient with given rate and burst limit
func NewRateLimiter(rps int64, burst int) *RateLimiter {
	return &RateLimiter{
		ticker:  time.NewTicker(time.Second / time.Duration(rps)),
		burst:   make(chan time.Time, burst),
		started: 0,
	}
}

//Stop stops a RateLimitedWebClient
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
	close(rl.burst)
}

//Start starts a ticker on RateLimitedWebClient
func (rl *RateLimiter) Start() error {
	if atomic.LoadUint32(&rl.started) > 0 {
		return errors.New("Client's ticker already started")
	}

	atomic.StoreUint32(&rl.started, 1)

	go func() {
		for t := range rl.ticker.C {
			select {
			case rl.burst <- t:
			default:
			}
		}
	}()
	return nil
}

func (rl *RateLimiter) Wait() {
	<-rl.burst
}
