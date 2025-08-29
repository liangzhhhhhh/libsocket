// Package new
// @description
// @author      梁志豪
// @datetime    2025/8/28 13:38
package libsocket

import (
	"time"
)

type HeartBeatParam struct {
	isOpen           bool
	flag             TypeKeepAlive
	activeHeartBeat  *ActiveHeartBeat
	passiveHeartBeat *PassiveHeartBeat
}

type TypeKeepAlive int

// HeartStrategy
const (
	KEEP_ALIVE_ACTIVE  TypeKeepAlive = 1 << iota // 主动发ping
	KEEP_ALIVE_PASSIVE                           // 接收到ping发pong
)

func NewHeartBeat(isOpen bool) *HeartBeatParam {
	return &HeartBeatParam{
		isOpen: isOpen,
	}
}

func NewDefaultHeartBeat() *HeartBeatParam {
	return NewHeartBeat(false)
}

type ActiveHeartBeat struct {
	interval *time.Ticker
}

func (p *HeartBeatParam) ActiveActiveHeartBeat(beat *ActiveHeartBeat) {
	p.flag = p.flag | KEEP_ALIVE_ACTIVE
	p.activeHeartBeat = beat
}

func (p *HeartBeatParam) IsActiveActiveHeartBeat() bool {
	return p.flag&KEEP_ALIVE_ACTIVE != 0
}

func (beat *ActiveHeartBeat) GetInterval() *time.Ticker {
	return beat.interval
}

type PassiveHeartBeat struct{}

func (p *HeartBeatParam) ActivePassiveHeartBeat(beat *PassiveHeartBeat) {
	p.flag = p.flag | KEEP_ALIVE_PASSIVE
	p.passiveHeartBeat = beat
}

func (p *HeartBeatParam) IsActivePassiveHeartBeat() bool {
	return p.flag&KEEP_ALIVE_PASSIVE != 0
}

func (p *HeartBeatParam) IsOpen() bool {
	return p.isOpen
}
