// Package new
// @description
// @author      梁志豪
// @datetime    2025/8/28 13:33
package libsocket

import (
	"math"
	"time"
)

type ReConnParam struct {
	isOpen                bool // 是否开启重连
	maxRetryCount         int  // 最大重试次数
	connDurationThreshold time.Duration
	backoff               Backoff // 退避策略
}

func ExponentialBackoff(attempts int) float64 {
	return (math.Pow(2.0, float64(attempts)) - 1) / 2
}

func ExponentialBackoffSeconds(attempts int) time.Duration {
	return time.Duration(ExponentialBackoff(attempts)) * time.Second
}

type Backoff func(int) time.Duration

func NewReConnParam(isOpen bool, maxRetryCount int, connDurationThreshold time.Duration, backoff Backoff) *ReConnParam {
	if backoff == nil {
		backoff = ExponentialBackoffSeconds
	}
	return &ReConnParam{
		isOpen:                isOpen,
		maxRetryCount:         maxRetryCount,
		backoff:               backoff,
		connDurationThreshold: connDurationThreshold,
	}
}

func NewDefaultReConnParam() *ReConnParam {
	return NewReConnParam(false, 0, time.Duration(0), nil)
}

func (p *ReConnParam) IsOpen() bool {
	return p.isOpen
}
