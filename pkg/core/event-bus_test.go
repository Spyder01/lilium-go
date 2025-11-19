package core

import (
	"sync"
	"testing"
	"time"
)

// helper: wait for value or timeout
func recvOrFail(t *testing.T, ch <-chan interface{}, timeout time.Duration) interface{} {
	t.Helper()
	select {
	case v, ok := <-ch:
		if !ok {
			t.Fatalf("channel closed unexpectedly")
		}
		return v
	case <-time.After(timeout):
		t.Fatalf("timeout waiting for event")
		return nil
	}
}

func TestSubscribeAndPublish(t *testing.T) {
	bus := NewEventBus()
	defer bus.Close()

	_, ch, _ := bus.Subscribe("topic1", 8)

	err := bus.Publish("topic1", "hello")
	if err != nil {
		t.Fatalf("unexpected error on publish: %v", err)
	}

	got := recvOrFail(t, ch, 200*time.Millisecond)
	if got != "hello" {
		t.Fatalf("expected 'hello', got %v", got)
	}
}

func TestMultipleSubscribersReceiveEvents(t *testing.T) {
	bus := NewEventBus()
	defer bus.Close()

	_, ch1, _ := bus.Subscribe("topic", 8)
	_, ch2, _ := bus.Subscribe("topic", 8)

	bus.Publish("topic", 42)

	v1 := recvOrFail(t, ch1, 200*time.Millisecond)
	v2 := recvOrFail(t, ch2, 200*time.Millisecond)

	if v1 != 42 || v2 != 42 {
		t.Fatalf("expected both subscribers to receive 42, got %v and %v", v1, v2)
	}
}

func TestUnsubscribeStopsReceiving(t *testing.T) {
	bus := NewEventBus()
	defer bus.Close()

	_, ch, unsub := bus.Subscribe("topic", 8)

	// first publish works
	bus.Publish("topic", "A")
	if recvOrFail(t, ch, 200*time.Millisecond) != "A" {
		t.Fatalf("expected 'A'")
	}

	unsub() // remove subscriber

	// second publish should not be seen
	bus.Publish("topic", "B")

	select {
	case v := <-ch:
		t.Fatalf("did not expect a message after unsubscribe, got %v", v)
	case <-time.After(100 * time.Millisecond):
		// expected
	}

	// also test explicit Unsubscribe API
	id2, ch2, _ := bus.Subscribe("topic", 8)
	bus.Unsubscribe("topic", id2)
	bus.Publish("topic", 123)

	select {
	case v := <-ch2:
		t.Fatalf("unexpected msg after Unsubscribe(): %v", v)
	case <-time.After(100 * time.Millisecond):
		// ok
	}
}

func TestPublishBufferFull(t *testing.T) {
	bus := NewEventBus()
	defer bus.Close()

	// subscriber with small buffer
	_, ch, _ := bus.Subscribe("slow", 1)

	// first event fills buffer
	if err := bus.Publish("slow", "X"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// second event is dropped because subscriber buffer is full
	err := bus.Publish("slow", "Y")
	if err == nil {
		t.Fatalf("expected publish buffer full error")
	}

	got := recvOrFail(t, ch, 200*time.Millisecond)
	if got != "X" {
		t.Fatalf("expected X, got %v", got)
	}
}

func TestConcurrentPublishSubscribe(t *testing.T) {
	bus := NewEventBus()
	defer bus.Close()

	var wg sync.WaitGroup

	// many subscribers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, ch, _ := bus.Subscribe("concurrent", 16)

			// read some messages
			for j := 0; j < 20; j++ {
				_ = recvOrFail(t, ch, 500*time.Millisecond)
			}
		}()
	}

	time.Sleep(10 * time.Millisecond)

	// many publishers
	for p := 0; p < 5; p++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				if err := bus.Publish("concurrent", j); err != nil {
					// buffer full is okay during aggressive testing
				}
			}
		}()
	}

	wg.Wait()
}
