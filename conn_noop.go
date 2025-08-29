package libsocket

import "context"

type (
	// emitter 是一个通用事件发射器接口。
	// 它接收一个键 (K) 和一个值 (V)，用于触发事件。
	emitter[K comparable, V any] interface {
		Emit(K, V)
	}

	// ConnectionHandler 定义了与连接交互的行为。
	ConnectionHandler interface {
		// Recv 在接收到来自服务器的消息时被调用。
		// 它负责处理来自服务器的入站数据流。
		Recv(m Message)

		// Send 在需要向服务器发送消息时被调用。
		// 它负责处理发往服务器的出站数据流。
		Send(m Message)

		// Connect 建立与服务器的连接。
		// 这是一个阻塞函数，只有当连接不再活跃时才会返回。
		Connect(ctx context.Context) error

		// CloseChan 返回一个 channel，当连接被关闭时会关闭该 channel。
		// 可以用来监听连接的关闭事件。
		CloseChan() CloseChan

		// CloseErr 返回导致连接关闭的错误原因。
		// 如果连接是正常关闭的，则应返回 nil。
		CloseErr() error

		// Close 关闭连接。
		// 它应确保清理与连接相关的所有资源。
		Close()
	}

	// ConnectionHandlerFactory 是一种函数类型，
	// 它接收一个 Client、一个 MessageHandler 和一个 EventEmitter，
	// 并返回一个 ConnectionHandler。
	ConnectionHandlerFactory func(Client, MessageHandler, emitter[EventType, EventType]) ConnectionHandler
)
