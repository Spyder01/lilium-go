package core

import (
	"sync"

	"github.com/spyder01/lilium-go/pkg/logger"
)

type Context struct {
	store     map[string]any
	mu        sync.RWMutex
	Bus       *EventBus
	isRunning bool
	Logger    *logger.Logger
	app       *Lilium
}

func (ctx *Context) Set(key string, val any) {
	ctx.mu.Lock()
	ctx.store[key] = val
	ctx.mu.Unlock()
}

func (ctx *Context) Exists(key string) bool {
	ctx.mu.RLock()
	_, ok := ctx.store[key]
	ctx.mu.RUnlock()
	return ok
}

func (ctx *Context) Get(key string) (any, bool) {
	ctx.mu.RLock()
	val, ok := ctx.store[key]
	ctx.mu.RUnlock()
	return val, ok
}

func (ctx *Context) Delete(key string) bool {
	ctx.mu.Lock()
	_, ok := ctx.store[key]
	if ok {
		delete(ctx.store, key)
	}
	ctx.mu.Unlock()
	return ok
}

func (ctx *Context) GetString(key string) (string, bool) {
	ctx.mu.RLock()
	v, ok := ctx.store[key]
	ctx.mu.RUnlock()
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func (ctx *Context) GetInt(key string) (int, bool) {
	ctx.mu.RLock()
	v, ok := ctx.store[key]
	ctx.mu.RUnlock()
	if !ok {
		return 0, false
	}
	i, ok := v.(int)
	return i, ok
}

func (ctx *Context) MustGet(key string) any {
	ctx.mu.RLock()
	val, ok := ctx.store[key]
	ctx.mu.RUnlock()
	if !ok {
		panic("context: missing key " + key)
	}
	return val
}

func (ctx *Context) GetOrDefault(key string, def any) any {
	ctx.mu.RLock()
	val, ok := ctx.store[key]
	ctx.mu.RUnlock()
	if !ok {
		return def
	}
	return val
}

func (ctx *Context) Update(key string, fn func(old any) any) {
	ctx.mu.Lock()
	old := ctx.store[key]
	ctx.store[key] = fn(old)
	ctx.mu.Unlock()
}

func (ctx *Context) Snapshot() map[string]any {
	ctx.mu.RLock()
	out := make(map[string]any, len(ctx.store))
	for k, v := range ctx.store {
		out[k] = v
	}
	ctx.mu.RUnlock()
	return out
}

func (ctx *Context) Clear() {
	ctx.mu.Lock()
	ctx.store = make(map[string]any)
	ctx.mu.Unlock()
}

func (ctx *Context) Start() {
	ctx.mu.Lock()
	ctx.isRunning = true
	ctx.app.isRunning = true
	ctx.mu.Unlock()
}

func (ctx *Context) Stop() {
	ctx.mu.Lock()
	ctx.isRunning = false
	ctx.app.isRunning = false
	ctx.mu.Unlock()
}

func (ctx *Context) Publish(topic string, data any) error {
	return ctx.Bus.Publish(topic, data)
}

func (ctx *Context) Subscribe(topic string, buf int) (<-chan any, func()) {
	_, ch, unsub := ctx.Bus.Subscribe(topic, buf)
	return ch, unsub
}

func (ctx *Context) SetLocal(key string, val any) {
	ctx.Set("local."+key, val)
}

func (ctx *Context) GetLocal(key string) (any, bool) {
	return ctx.Get("local." + key)
}

func (ctx *Context) Provide(key string, val any) {
	ctx.Set("di."+key, val)
}

func (ctx *Context) Resolve(key string) (any, bool) {
	return ctx.Get("di." + key)
}

func (ctx *Context) GetLogger() *logger.Logger {
	return ctx.app.Logger
}
