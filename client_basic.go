package libsocket

import (
	"context"
)

// basicClient 是一种客户端实现，基于单一的连接套接字。
// 它只会将 websocket 的 "data" 消息转发给消息处理器（messageHandler），而 "ping"、"pong" 或 "close" 消息则会交由连接处理器（connection handlers）来处理。
// ⚠️ 重要提示：不要将其与 subscriber_static client 一起使用。该客户端旨在作为独立客户端使用，并计划在未来替代其他客户端。
type basicClient struct {
	// connectionHandlerFactory 是用于创建新连接处理器的工厂
	*connectionHandler

	// messageHandler 是用于处理传入消息的消息处理器
	messageHandler MessageHandler

	// eventHandler 是用于处理客户端事件的回调函数
	eventHandler func(Client, EventType)

	// eventEmitter 是事件分发器，用于管理事件监听与触发
	eventEmitter *EventEmitterCallback[EventType, EventType]
}

func (b *basicClient) createConnectionHandler(logger logger, connFactory ConnectionFactory, paramFuncs ...WithOptionalHandler) {
	handlerWrapper := func(cli Client, m Message) {
		if m.Type().IsData() {
			b.messageHandler(cli, m)
		} else {
			b.connectionHandler.Recv(m)
		}
	}

	b.connectionHandler = NewHandlerCore(logger, b, connFactory, handlerWrapper, b.eventEmitter, paramFuncs...)
}

func (b *basicClient) Open(ctx context.Context) error {

	b.eventEmitter.On(EventConnect, func(eventType EventType) {
		b.eventHandler(b, eventType)
	})

	b.eventEmitter.On(EventClose, func(eventType EventType) {
		b.eventHandler(b, eventType)
	})

	b.eventEmitter.On(EventReconnect, func(eventType EventType) {
		b.eventHandler(b, eventType)
	})

	if err := b.connectionHandler.Connect(ctx); err != nil {
		return err
	}

	return nil
}

func (b *basicClient) Send(m Message) {
	b.connectionHandler.Send(m)
}

func (b *basicClient) Close() {
	if b.eventEmitter != nil {
		b.eventEmitter.Close()
	}
	if b.connectionHandler != nil {
		b.connectionHandler.Close()
	}
}

func (b *basicClient) CloseChan() CloseChan {
	return b.connectionHandler.CloseChan()
}

func newBasicClient(
	logger logger,
	connFactory ConnectionFactory,
	messageHandler MessageHandler,
	eventHandler EventHandler,
	paramFuncs ...WithOptionalHandler,
) *basicClient {
	b := &basicClient{
		messageHandler: messageHandler,
		eventHandler:   eventHandler,
		eventEmitter:   NewEventEmitter[EventType, EventType](),
	}
	b.createConnectionHandler(logger, connFactory, paramFuncs...)
	return b
}

func NewBasicClientFactory(
	messageHandler MessageHandler,
	eventHandler EventHandler,
) ClientFactory {
	return func(logger logger, connFactory ConnectionFactory, paramFuncs ...WithOptionalHandler) Client {
		return newBasicClient(
			logger,
			connFactory,
			messageHandler,
			eventHandler,
			paramFuncs...,
		)
	}
}
