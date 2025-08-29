package libsocket

import (
	"sync"
)

type callback[T any] func(T)

// EventEmitterCallback 是一个简单的事件发射器（事件分发器）。
// 它将事件（类型为 K）映射到监听器（类型为 V 的只读通道）。
// ⚠️ 警告：强烈建议使用只读通道作为监听器，
// 以避免数据竞争或意外行为。
type EventEmitterCallback[K comparable, V any] struct {
	listeners map[K][]callback[V]
	lock      sync.RWMutex
}

// NewEventEmitter 创建一个新的 EventEmitterCallback 并返回它的指针。
func NewEventEmitter[K comparable, V any]() *EventEmitterCallback[K, V] {
	return &EventEmitterCallback[K, V]{
		listeners: make(map[K][]callback[V]),
	}
}

// On 为指定的事件注册一个新的监听器。
func (e *EventEmitterCallback[K, V]) On(event K, listener callback[V]) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.listeners[event] = append(e.listeners[event], listener)
}

// Emit 会同步触发所有注册在指定事件上的监听器，
// 并将事件数据发送到它们的通道中。
// 该方法会等待所有数据发送完成后才返回。
// 如果 EventEmitterCallback 已经被关闭，Emit 将不会向任何通道发送数据。
func (e *EventEmitterCallback[K, V]) Emit(event K, data V) {
	e.lock.RLock()
	defer e.lock.RUnlock()

	listeners, found := e.listeners[event]
	if !found {
		return
	}

	for _, listener := range listeners {
		listener(data)
	}
}

// Close 会关闭所有监听器的通道，并清空所有监听器， 以防止内存泄漏。
func (e *EventEmitterCallback[K, V]) Close() {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.listeners = make(map[K][]callback[V])
}
