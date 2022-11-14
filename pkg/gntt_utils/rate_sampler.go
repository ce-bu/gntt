package gntt_utils

import (
	"sync/atomic"
	"time"
)

type RateSampler struct {
	count    uint64
	stop     chan bool
	interval time.Duration
	callback func(float64)
}

func NewSampler(interval time.Duration, callback func(float64)) *RateSampler {
	return &RateSampler{
		stop:     make(chan bool),
		callback: callback,
		interval: interval,
	}
}

func (rs *RateSampler) AddSample(n uint64) {
	atomic.AddUint64(&rs.count, n)
}

func (rs *RateSampler) Start() {
	go func() {
		pt := time.Now()
		pb := atomic.LoadUint64(&rs.count)
		for {
			select {
			case <-rs.stop:
				goto end
			case <-time.After(rs.interval):
				ct := time.Now()
				cb := atomic.LoadUint64(&rs.count)
				rate := float64(cb-pb) / float64(ct.Sub(pt).Seconds())
				pt = ct
				pb = cb
				rs.callback(rate)
			}
		}
	end:
		rs.stop <- true
	}()

}

func (rs *RateSampler) Stop() {
	rs.stop <- true
	<-rs.stop
}
