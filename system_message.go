// Package libsocket
// @description
// @author      梁志豪
// @datetime    2025/8/28 15:26
package libsocket

var (
	basePingMessage = NewPingMessage([]byte("ping"))
	basePongMessage = NewPongMessage([]byte("pong"))
)
