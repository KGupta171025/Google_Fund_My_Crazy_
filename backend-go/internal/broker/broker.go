package broker

import (
	"sync"
	"sync/atomic"
	"time"
)

// Event is the generic envelope for pub/sub messages
type Event struct {
	Topic     string      `json:"topic"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// Subscription represents an active subscriber's channel
type Subscription struct {
	ID      string
	Topic   string
	Channel chan Event
}

// EventBroker coordinates ultra-high throughput event publishing and subscriptions
type EventBroker struct {
	mu           sync.RWMutex
	subscribers  map[string]map[string]*Subscription // Topic -> SubscriberID -> Sub
	bufferSize   int
	totalEvents  uint64
	droppedEvents uint64
}

// NewEventBroker initializes a high-concurrency event broker
func NewEventBroker(bufferSize int) *EventBroker {
	if bufferSize <= 0 {
		bufferSize = 4096
	}
	return &EventBroker{
		subscribers: make(map[string]map[string]*Subscription),
		bufferSize:  bufferSize,
	}
}

// Subscribe registers a subscriber for a given topic
func (b *EventBroker) Subscribe(subID, topic string) *Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = make(map[string]*Subscription)
	}

	sub := &Subscription{
		ID:      subID,
		Topic:   topic,
		Channel: make(chan Event, b.bufferSize),
	}

	b.subscribers[topic][subID] = sub
	return sub
}

// Unsubscribe removes an active subscriber
func (b *EventBroker) Unsubscribe(subID, topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if subs, ok := b.subscribers[topic]; ok {
		if sub, exists := subs[subID]; exists {
			close(sub.Channel)
			delete(subs, subID)
		}
		if len(subs) == 0 {
			delete(b.subscribers, topic)
		}
	}
}

// Publish broadcasts an event to all subscribers of a topic (non-blocking)
func (b *EventBroker) Publish(topic string, payload interface{}) {
	event := Event{
		Topic:     topic,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	atomic.AddUint64(&b.totalEvents, 1)

	b.mu.RLock()
	defer b.mu.RUnlock()

	if subs, ok := b.subscribers[topic]; ok {
		for _, sub := range subs {
			select {
			case sub.Channel <- event:
				// Successfully delivered
			default:
				// Channel buffer full; drop event to avoid blocking publishers
				atomic.AddUint64(&b.droppedEvents, 1)
			}
		}
	}

	// Also broadcast to wildcard topic "*"
	if wildSubs, ok := b.subscribers["*"]; ok {
		for _, sub := range wildSubs {
			select {
			case sub.Channel <- event:
			default:
				atomic.AddUint64(&b.droppedEvents, 1)
			}
		}
	}
}

// Metrics returns total and dropped event counts
func (b *EventBroker) Metrics() (total uint64, dropped uint64) {
	return atomic.LoadUint64(&b.totalEvents), atomic.LoadUint64(&b.droppedEvents)
}
