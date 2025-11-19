package core

import (
	"sync"
	"testing"
)

func newTestContext() *Context {
	return &Context{
		store: make(map[string]any),
		Bus:   NewEventBus(),
	}
}

func TestSetAndGet(t *testing.T) {
	ctx := newTestContext()

	ctx.Set("name", "Lilium")

	val, ok := ctx.Get("name")
	if !ok {
		t.Fatalf("expected key to exist")
	}
	if val.(string) != "Lilium" {
		t.Fatalf("got %v, want %v", val, "Lilium")
	}
}

func TestExists(t *testing.T) {
	ctx := newTestContext()

	if ctx.Exists("missing") {
		t.Fatalf("Exists for missing key returned true")
	}

	ctx.Set("x", 1)

	if !ctx.Exists("x") {
		t.Fatalf("Exists for present key returned false")
	}
}

func TestDelete(t *testing.T) {
	ctx := newTestContext()
	ctx.Set("user", "alice")

	ok := ctx.Delete("user")
	if !ok {
		t.Fatalf("expected delete to return true")
	}

	if ctx.Exists("user") {
		t.Fatalf("expected key to be deleted")
	}

	if ctx.Delete("user") {
		t.Fatalf("delete on missing key should return false")
	}
}

func TestGetString(t *testing.T) {
	ctx := newTestContext()
	ctx.Set("name", "Alice")

	s, ok := ctx.GetString("name")
	if !ok || s != "Alice" {
		t.Fatalf("expected string Alice, got %v (ok=%v)", s, ok)
	}

	_, ok = ctx.GetString("missing")
	if ok {
		t.Fatalf("expected ok=false for missing key")
	}

	ctx.Set("num", 123)
	_, ok = ctx.GetString("num")
	if ok {
		t.Fatalf("expected type assertion to fail")
	}
}

func TestGetInt(t *testing.T) {
	ctx := newTestContext()
	ctx.Set("age", 30)

	i, ok := ctx.GetInt("age")
	if !ok || i != 30 {
		t.Fatalf("expected int 30, got %v (ok=%v)", i, ok)
	}

	_, ok = ctx.GetInt("missing")
	if ok {
		t.Fatalf("expected ok=false for missing key")
	}

	ctx.Set("str", "hello")
	_, ok = ctx.GetInt("str")
	if ok {
		t.Fatalf("expected type assertion to fail")
	}
}

func TestMustGet(t *testing.T) {
	ctx := newTestContext()
	ctx.Set("x", 7)

	if v := ctx.MustGet("x"); v.(int) != 7 {
		t.Fatalf("MustGet wrong value")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected MustGet to panic on missing key")
		}
	}()
	_ = ctx.MustGet("missing")
}

func TestGetOrDefault(t *testing.T) {
	ctx := newTestContext()

	if v := ctx.GetOrDefault("missing", 42); v.(int) != 42 {
		t.Fatalf("expected default value 42")
	}

	ctx.Set("x", "abc")
	if v := ctx.GetOrDefault("x", "zzz"); v.(string) != "abc" {
		t.Fatalf("expected stored value abc, got %v", v)
	}
}

func TestUpdate(t *testing.T) {
	ctx := newTestContext()
	ctx.Set("count", 1)

	ctx.Update("count", func(old any) any {
		return old.(int) + 10
	})

	val, _ := ctx.Get("count")
	if val.(int) != 11 {
		t.Fatalf("update failed, expected 11, got %v", val)
	}
}

func TestSnapshot(t *testing.T) {
	ctx := newTestContext()
	ctx.Set("a", 1)
	ctx.Set("b", 2)

	snap := ctx.Snapshot()

	if len(snap) != 2 {
		t.Fatalf("expected snapshot size 2")
	}

	snap["a"] = 100
	val, _ := ctx.Get("a")
	if val.(int) != 1 {
		t.Fatalf("snapshot should be a copy, original modified to %v", val)
	}
}

func TestClear(t *testing.T) {
	ctx := newTestContext()
	ctx.Set("a", 1)
	ctx.Set("b", 2)

	ctx.Clear()

	if ctx.Exists("a") || ctx.Exists("b") {
		t.Fatalf("expected all keys to be cleared")
	}
}

func TestSetLocalAndGetLocal(t *testing.T) {
	ctx := newTestContext()

	ctx.SetLocal("session", "xyz")

	v, ok := ctx.GetLocal("session")
	if !ok || v.(string) != "xyz" {
		t.Fatalf("local storage failed")
	}
}

func TestProvideAndResolve(t *testing.T) {
	ctx := newTestContext()

	ctx.Provide("db", "postgres")

	v, ok := ctx.Resolve("db")
	if !ok || v.(string) != "postgres" {
		t.Fatalf("DI storage failed")
	}
}

func TestPublishAndSubscribe(t *testing.T) {
	ctx := newTestContext()

	ch, unsub := ctx.Subscribe("test", 1)
	defer unsub()

	err := ctx.Publish("test", "hello")
	if err != nil {
		t.Fatalf("unexpected publish error: %v", err)
	}

	msg := <-ch
	if msg.(string) != "hello" {
		t.Fatalf("expected hello, got %v", msg)
	}
}

func TestConcurrencySafety(t *testing.T) {
	ctx := newTestContext()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			ctx.Set("k", i)
			_, _ = ctx.Get("k")
			ctx.Exists("k")
			ctx.Update("k", func(old any) any { return i })
		}(i)
	}

	wg.Wait()
}
