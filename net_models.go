// socket模型
package libsocket

import (
	"context"
)

type (
	Connection interface {
		Write(m Message) error
		Open(ctx context.Context) error
		Close()
		CloseErr() error
		CloseChan() CloseChan
	}

	ConnectionFactory func(ctx context.Context, recvChan chan<- Message) Connection
)
