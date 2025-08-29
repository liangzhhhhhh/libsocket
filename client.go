package libsocket

import (
	"context"
)

type (
	// Client 接口定义了一个客户端的行为。这包括建立和关闭连接、发送消息以及管理事件通知。
	Client interface {
		// Open 与服务器建立连接
		Open(ctx context.Context) error
		// Send 向服务器发送一条消息
		Send(m Message)
		// Close 关闭与服务器的连接
		Close()
		// CloseChan 返回一个 channel，当连接关闭时会发出信号
		CloseChan() CloseChan
	}

	// WithOptional 定义了可选开启的功能
	WithOptionalHandler func(c *connectionHandler)

	// CloseChan 是一个 channel 类型，用于通知连接已关闭
	CloseChan chan struct{}

	// MessageHandler 定义了消息处理函数
	// 入参为 Client 和 Message
	MessageHandler func(Client, Message)

	// EventHandler 定义了事件处理函数
	// 入参为 Client 和 EventType
	EventHandler func(Client, EventType)

	// ClientFactory 是客户端工厂方法
	// 用于创建新的 Client 实例
	ClientFactory func(logger Logger, connFactory ConnectionFactory, params ...WithOptionalHandler) Client
)
