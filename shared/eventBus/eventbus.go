package eventBus

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"sync"
)

type Bus interface {
	Subscribe(topic string, handler interface{}) error
	Publish(ctx context.Context, topic string, args ...interface{})
}

// AsyncEventBus 异步事件总线
type AsyncEventBus struct {
	handlers map[string][]reflect.Value
	lock     sync.Mutex
}

// NewEventBus new
func NewEventBus() *AsyncEventBus {
	return &AsyncEventBus{
		handlers: map[string][]reflect.Value{},
		lock:     sync.Mutex{},
	}
}

// Subscribe 订阅
func (bus *AsyncEventBus) Subscribe(topic string, f interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	v := reflect.ValueOf(f)
	if v.Type().Kind() != reflect.Func {
		return fmt.Errorf("handler is not a function")
	}

	handler, ok := bus.handlers[topic]
	if !ok {
		handler = []reflect.Value{}
	}
	handler = append(handler, v)
	bus.handlers[topic] = handler

	return nil
}

// Publish 发布
// 这里异步执行，并且不会等待返回结果
func (bus *AsyncEventBus) Publish(ctx context.Context, topic string, args ...interface{}) {
	handlers, ok := bus.handlers[topic]
	if !ok {
		logx.WithContext(ctx).Debugf("Publish not found handlers in topic:", topic)
		return
	}

	params := make([]reflect.Value, len(args)+1)
	params[0] = reflect.ValueOf(ctx)
	for i, arg := range args {
		params[i+1] = reflect.ValueOf(arg)
	}

	for i := range handlers {
		//先不使用异步,异步ctx会超时,先不做这块
		//utils.Go(ctx, func() {
		handlers[i].Call(params)
		//})
	}
}
