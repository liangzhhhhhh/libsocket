// Package new
// @description
// @author      liangzh
// @datetime    2025/8/28 13:35
package libsocket

import (
	"time"
)

type ReopenParam struct {
	isOpen               bool
	reopenIntervalTicker *time.Ticker
}

func NewReOpenParam(isOpen bool, reopenIntervalTicker *time.Ticker) *ReopenParam {
	return &ReopenParam{
		isOpen:               isOpen,
		reopenIntervalTicker: reopenIntervalTicker,
	}
}

func NewDefaultReopenParam() *ReopenParam {
	return NewReOpenParam(false, nil)
}

func (p *ReopenParam) IsOpen() bool {
	return p.isOpen
}
