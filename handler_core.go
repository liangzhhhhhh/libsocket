// Package libsocket
// @description
// @author      梁志豪
// @datetime    2025/8/28 14:07
package libsocket

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"sync"
	"time"
)

type connectionHandler struct {
	client         Client
	logger         Logger                        // 日志管理
	messageHandler MessageHandler                // 消息处理
	emitter        emitter[EventType, EventType] // 监听器
	connFactory    ConnectionFactory             // 连接工厂
	send           chan Message                  // 发送的消息 缓冲
	recv           chan Message                  // 接收的消息 缓冲
	conn           Connection                    // 连接
	closeC         CloseChan                     // 外部发送指令关闭连接
	connOnce       sync.Once                     // 连接只执行一次
	closeOnce      sync.Once                     // 单次执行一次的关闭连接操作
	closeReason    error                         // 关闭原因
	connMu         sync.RWMutex
	*ReConnParam
	*ReopenParam
	*HeartBeatParam
}

// Connect 建立连接并启动读循环
func (h *connectionHandler) Connect(ctx context.Context) error {
	var err error
	h.connOnce.Do(func() {
		err = h.newConn(ctx)
		if err != nil {
			log.Printf("connect failed: %v", err)
			return
		}
		go h.run(ctx)
	})

	return err
}

func (h *connectionHandler) newConn(ctx context.Context) error {
	h.connMu.Lock()
	defer h.connMu.Unlock()
	recvChan := make(chan Message, 16) // 定义recv的内存地址
	var (
		attempts     = 0
		isReConnOpen = h.ReConnParam.isOpen
	)
	if h.conn != nil {
		h.conn.Close()
	}
	for {
		h.conn = h.connFactory(ctx, recvChan)
		if err := h.conn.Open(ctx); err != nil {
			attempts += 1
			if attempts > h.ReConnParam.maxRetryCount {
				return err
			}
			if errors.Is(err, ErrCannotConnect) {
				h.logger.Infof("cannot connect, reconnecting asap [%d次] due to: %s", attempts, err)
				if isReConnOpen {
					time.Sleep(time.Second)
					continue
				}
			}
			if isReConnOpen {
				ttw := h.backoff(attempts)
				h.logger.Infof("cannot connect after %s, waiting %s", err, ttw)
				time.Sleep(ttw)
				continue
			}
			// 理论上不会到这里
			return err
		}
		return nil
	}
}

func (h *connectionHandler) run(ctx context.Context) {
	var (
		innerCloseChan = h.conn.CloseChan()
		isReConnOpen   = h.ReConnParam.isOpen
		reConnRounds   = 0
		reOpenRounds   = 0
		then           = time.Now().UTC()
	)

	defer h.conn.Close()

	var reopenIntervalTicker <-chan time.Time
	if h.reopenIntervalTicker != nil {
		reopenIntervalTicker = h.reopenIntervalTicker.C
		defer h.reopenIntervalTicker.Stop()
	}
	var activeHeartBeat <-chan time.Time
	if h.activeHeartBeat != nil && h.activeHeartBeat.interval != nil {
		activeHeartBeat = h.activeHeartBeat.interval.C
		defer h.activeHeartBeat.interval.Stop()
	}

	for {
		select {
		case <-ctx.Done():
		case <-h.closeC:
			return
		case msg := <-h.recv:
			h.ReadMessage(msg)
		case msg := <-h.send:
			h.WriteMessage(msg)
		case <-activeHeartBeat:
			h.WriteMessage(basePingMessage)
		case <-reopenIntervalTicker:
			reOpenRounds++
			// Time to spawn a new conn. When a new one is opened, close the previous one. Order matters
			// to prevent data loss (duplicated data is preferred above lack of it)
			h.logger.Infof("spawning and opening #%d conn due to reopen trigger", reOpenRounds)
			err := h.newConn(ctx)
			if err != nil {
				h.Close()
				return
			}
			nextCloseChan := h.conn.CloseChan()
			innerCloseChan = nextCloseChan
		case <-innerCloseChan:
			h.conn.Close()
			h.closeReason = h.conn.CloseErr()

			if h.closeReason != nil {
				if errors.Is(h.closeReason, ErrConnectionClosed) ||
					errors.Is(h.closeReason, ErrTerminated) {
					if isReConnOpen {
						delta := time.Since(then)
						if delta > h.connDurationThreshold {
							// We assume that the connection was healthy for `connDurationThreshold` and that it
							// was terminated due to natural reasons, so we should try to reconnect asap
							reConnRounds = 0
						} else {
							reConnRounds++
						}
					}
				}
			}
			var err error
			if isReConnOpen {
				h.logger.Infof("retrying to connect at once due to %s", h.closeReason)
				err = h.newConn(ctx)
			}
			if err != nil || !isReConnOpen {
				h.Close()
				return
			}
			innerCloseChan = h.conn.CloseChan()
			then = time.Now().UTC()
			go h.emitter.Emit(EventReconnect, EventReconnect)
		}
	}
}

// Send 消息写入底层连接
func (h *connectionHandler) Send(m Message) {
	// 增加吞吐
	h.send <- m
}

func (h *connectionHandler) Recv(m Message) {
	// 增加吞吐
	h.recv <- m
}

// Close 主动关闭
func (h *connectionHandler) Close() {
	h.closeOnce.Do(func() {
		close(h.closeC)
	})
}

// CloseChan 暴露底层关闭通道
func (h *connectionHandler) CloseChan() CloseChan {
	return h.closeC
}

// CloseErr 暴露底层错误
func (h *connectionHandler) CloseErr() error {
	return h.closeReason
}

func (h *connectionHandler) WriteMessage(m Message) {
	if h.conn == nil {
		log.Println("connection not ready")
		return
	}
	// TODO: queue to buffer messages to send while reconnecting. Procrastinated as of now since
	if err := h.conn.Write(m); err != nil {
		log.Printf("send error: %v\n", err)
	}
}

func (h *connectionHandler) ReadMessage(msg Message) {

	if h.conn == nil {
		log.Println("connection not ready")
		return
	}
	// TODO: queue to buffer messages to send while reconnecting. Procrastinated as of now since
	h.messageHandler(h.client, msg)
}

func NewHandlerCore(logger Logger, client Client, connFactory ConnectionFactory, handler MessageHandler, emitter emitter[EventType, EventType], paramFuncs ...WithOptionalHandler) *connectionHandler {
	c := &connectionHandler{
		client:         client,
		logger:         logger,
		messageHandler: handler,
		connFactory:    connFactory,
		emitter:        emitter,
		send:           make(chan Message, 32),
		recv:           make(chan Message, 32),
		closeC:         make(CloseChan),
		ReConnParam:    NewDefaultReConnParam(),
		ReopenParam:    NewDefaultReopenParam(),
		HeartBeatParam: NewDefaultHeartBeat(),
	}
	for _, fun := range paramFuncs {
		fun(c)
	}
	return c
}

func WithReopenParam(p *ReopenParam) WithOptionalHandler {
	return func(h *connectionHandler) { h.ReopenParam = p }
}

func WithReConnParam(p *ReConnParam) WithOptionalHandler {
	return func(h *connectionHandler) { h.ReConnParam = p }
}

func WithHeartBeatParam(p *HeartBeatParam) WithOptionalHandler {
	return func(h *connectionHandler) { h.HeartBeatParam = p }
}
