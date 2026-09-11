package broker

import (
	"sync"
	"testing"
	"time"
)

func TestEventBrokerPublishSubscribe(t *testing.T) {
	b := NewEventBroker(100)

	sub := b.Subscribe("sub-1", "telemetry:kinetic")
	defer b.Unsubscribe("sub-1", "telemetry:kinetic")

	b.Publish("telemetry:kinetic", map[string]interface{}{
		"joules": 4500.0,
	})

	select {
	case evt := <-sub.Channel:
		if evt.Topic != "telemetry:kinetic" {
			t.Fatalf("Expected topic telemetry:kinetic, got %s", evt.Topic)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("Timed out waiting for published event")
	}
}

func TestEventBrokerConcurrentHighThroughput(t *testing.T) {
	b := NewEventBroker(5000)

	sub := b.Subscribe("bench-sub", "city:stream")
	defer b.Unsubscribe("bench-sub", "city:stream")
	_ = sub // Subscription registered

	const totalMsgs = 10000
	const numGoroutines = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < totalMsgs/numGoroutines; j++ {
				b.Publish("city:stream", j)
			}
		}(i)
	}

	wg.Wait()

	total, _ := b.Metrics()
	if total < totalMsgs {
		t.Fatalf("Expected at least %d total events, got %d", totalMsgs, total)
	}
}
