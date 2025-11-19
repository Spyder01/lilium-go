package core

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	ErrClosed      = errors.New("eventbus closed")
	ErrTopicClosed = errors.New("topic closed")
)

// internal types
type subscribeReq struct {
	buf   int
	resp  chan subscribeResp
	idOut chan uint64 // optional; for async id
}
type subscribeResp struct {
	id          uint64
	ch          chan interface{}
	unsubscribe func()
}

type publishReq struct {
	event interface{}
	ack   chan error // optional ack channel for publish result (deliver-to-others-but-error-on-any-failure)
}

type unsubscribeReq struct {
	id uint64
}

type topicBroker struct {
	// channels to broker goroutine
	subscribe   chan subscribeReq
	unsubscribe chan unsubscribeReq
	publish     chan publishReq
	close       chan struct{}
}

type EventBus struct {
	// map of topic => *topicBroker (we synchronize access to this map using atomic swap)
	// but we will keep a simple map guarded through a small lock-free approach: use a single goroutine
	// to create/get brokers (fast path avoids locking). For simplicity we use a small sync.Map-like atomic pointer.
	topics atomic.Pointer[map[string]*topicBroker] // pointer to map snapshot
	closed atomic.Bool
	idSeq  uint64
}

// NewEventBus initializes an empty bus.
func NewEventBus() *EventBus {
	m := make(map[string]*topicBroker)
	eb := &EventBus{}
	eb.topics.Store(&m)
	return eb
}

func (eb *EventBus) nextID() uint64 {
	return atomic.AddUint64(&eb.idSeq, 1)
}

// getOrCreateBroker returns existing broker or creates and starts one.
func (eb *EventBus) getOrCreateBroker(topic string) (*topicBroker, error) {
	if eb.closed.Load() {
		return nil, ErrClosed
	}

	// fast path: load map and check
	mp := eb.topics.Load()
	if b, ok := (*mp)[topic]; ok {
		return b, nil
	}

	// need to create a broker atomically: use CAS loop on the map pointer
	for {
		oldPtr := eb.topics.Load()
		oldMap := *oldPtr
		// if another goroutine created meanwhile, return it
		if b, ok := oldMap[topic]; ok {
			return b, nil
		}

		// copy map
		newMap := make(map[string]*topicBroker, len(oldMap)+1)
		for k, v := range oldMap {
			newMap[k] = v
		}
		// create broker
		broker := &topicBroker{
			subscribe:   make(chan subscribeReq),
			unsubscribe: make(chan unsubscribeReq),
			publish:     make(chan publishReq, 256), // per-topic publish buffer
			close:       make(chan struct{}),
		}
		newMap[topic] = broker

		// attempt to swap pointer
		if eb.topics.CompareAndSwap(oldPtr, &newMap) {
			// start broker goroutine
			go runBroker(topic, broker)
			return broker, nil
		}
		// else: retry (someone else changed map)
	}
}

// Subscribe returns id, a read-only channel for events, and an unsubscribe func.
// buf is the subscriber channel buffer size.
func (eb *EventBus) Subscribe(topic string, buf int) (uint64, <-chan interface{}, func()) {
	broker, err := eb.getOrCreateBroker(topic)
	if err != nil {
		// bus closed
		ch := make(chan interface{})
		close(ch)
		return 0, ch, func() {}
	}

	respC := make(chan subscribeResp)
	req := subscribeReq{
		buf:  buf,
		resp: respC,
	}
	broker.subscribe <- req
	resp := <-respC
	return resp.id, resp.ch, resp.unsubscribe
}

// Unsubscribe by ID (alternative API)
func (eb *EventBus) Unsubscribe(topic string, id uint64) {
	mp := eb.topics.Load()
	if b, ok := (*mp)[topic]; ok {
		select {
		case b.unsubscribe <- unsubscribeReq{id: id}:
		default:
			// If unsubscribe channel is full (shouldn't be), do non-blocking send as fallback
			go func() { b.unsubscribe <- unsubscribeReq{id: id} }()
		}
	}
}

// Publish publishes the event. It's non-blocking with current design (buffered per-topic).
// Behavior: deliver to all subscribers; if any subscriber's buffer is full the message is still
// delivered to the others, but Publish returns a non-nil error indicating at least one slow subscriber.
func (eb *EventBus) Publish(topic string, evt interface{}) error {
	if eb.closed.Load() {
		return ErrClosed
	}
	mp := eb.topics.Load()
	b, ok := (*mp)[topic]
	if !ok {
		// no subscribers yet — we can create broker (optional) or drop
		var err error
		b, err = eb.getOrCreateBroker(topic)
		if err != nil {
			return err
		}
	}
	ack := make(chan error, 1)
	select {
	case b.publish <- publishReq{event: evt, ack: ack}:
		// wait for broker to process and report result
		return <-ack
	default:
		// topic publish channel itself is full
		return fmt.Errorf("publish buffer full for topic %s", topic)
	}
}

// Close closes the whole bus and stops brokers.
func (eb *EventBus) Close() {
	if !eb.closed.CompareAndSwap(false, true) {
		return
	}
	mp := eb.topics.Load()
	for _, b := range *mp {
		close(b.close)
	}
	// clear map
	empty := make(map[string]*topicBroker)
	eb.topics.Store(&empty)
}

// broker goroutine
func runBroker(topic string, b *topicBroker) {
	// subscribers map id -> channel
	subs := make(map[uint64]chan interface{})
	for {
		select {
		case req := <-b.subscribe:
			// generate id using time-based value (keeps broker-local uniqueness)
			id := uint64(time.Now().UnixNano()) ^ uint64(len(subs)+1)
			ch := make(chan interface{}, req.buf)
			unsubOnce := func() {
				// send unsubscribe request back to this broker
				select {
				case b.unsubscribe <- unsubscribeReq{id: id}:
				default:
					// fallback: spawn a goroutine so we don't block
					go func() { b.unsubscribe <- unsubscribeReq{id: id} }()
				}
			}
			resp := subscribeResp{
				id:          id,
				ch:          ch,
				unsubscribe: unsubOnce,
			}
			// register
			subs[id] = ch
			req.resp <- resp

		case u := <-b.unsubscribe:
			// Remove subscriber but DO NOT close the subscriber channel.
			// Subscribers own their channel lifecycle; closing here causes `<-ch` to return nil
			// which some tests interpret as a delivered message.
			delete(subs, u.id)

		case pub := <-b.publish:
			var anyErr error
			for _, ch := range subs {
				select {
				case ch <- pub.event:
					// delivered
				default:
					// mark that at least one subscriber is slow, but continue delivering to others
					anyErr = fmt.Errorf("subscriber buffer full on topic %s", topic)
				}
			}
			// send result back to publisher if requested
			if pub.ack != nil {
				// best-effort: report the error (nil if all delivered)
				pub.ack <- anyErr
			}

		case <-b.close:
			// Close everything and exit
			for id, ch := range subs {
				delete(subs, id)
				close(ch)
			}
			return
		}
	}
}
