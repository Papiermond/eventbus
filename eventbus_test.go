package eventbus

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Test event types
type testEvent struct {
	eventType EventType
	data      string
}

func (e testEvent) GetType() EventType {
	return e.eventType
}

type counterEvent struct {
	value int
}

func (e counterEvent) GetType() EventType {
	return "counter"
}

// TestNew verifies that New creates a valid event bus
func TestNew(t *testing.T) {
	bus := New()
	if bus == nil {
		t.Fatal("New() returned nil")
	}
}

// TestSubscribeAndPublish verifies basic subscribe and publish functionality
func TestSubscribeAndPublish(t *testing.T) {
	bus := New()
	received := false

	bus.Subscribe("test:event", func(event Event) {
		received = true
		e := event.(testEvent)
		if e.data != "test data" {
			t.Errorf("Expected data 'test data', got '%s'", e.data)
		}
	})

	bus.Publish(testEvent{
		eventType: "test:event",
		data:      "test data",
	})

	if !received {
		t.Error("Event was not received by subscriber")
	}
}

// TestMultipleSubscribers verifies that multiple subscribers receive the same event
func TestMultipleSubscribers(t *testing.T) {
	bus := New()
	var count atomic.Int32

	// Subscribe three listeners
	for i := 0; i < 3; i++ {
		bus.Subscribe("test:multi", func(event Event) {
			count.Add(1)
		})
	}

	bus.Publish(testEvent{eventType: "test:multi", data: "test"})

	if count.Load() != 3 {
		t.Errorf("Expected 3 listeners to be called, got %d", count.Load())
	}
}

// TestMultipleEventTypes verifies that different event types are routed correctly
func TestMultipleEventTypes(t *testing.T) {
	bus := New()
	event1Received := false
	event2Received := false

	bus.Subscribe("event:one", func(event Event) {
		event1Received = true
	})

	bus.Subscribe("event:two", func(event Event) {
		event2Received = true
	})

	bus.Publish(testEvent{eventType: "event:one", data: "test"})

	if !event1Received {
		t.Error("event:one was not received")
	}
	if event2Received {
		t.Error("event:two should not have been received")
	}
}

// TestPublishWithNoSubscribers verifies that publishing without subscribers doesn't panic
func TestPublishWithNoSubscribers(t *testing.T) {
	bus := New()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Publish panicked with no subscribers: %v", r)
		}
	}()

	bus.Publish(testEvent{eventType: "unsubscribed", data: "test"})
}

// TestSubscriberOrder verifies that subscribers are called in registration order
func TestSubscriberOrder(t *testing.T) {
	bus := New()
	var order []int
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		value := i
		bus.Subscribe("order:test", func(event Event) {
			mu.Lock()
			order = append(order, value)
			mu.Unlock()
		})
	}

	bus.Publish(testEvent{eventType: "order:test", data: "test"})

	if len(order) != 5 {
		t.Fatalf("Expected 5 calls, got %d", len(order))
	}

	for i := 0; i < 5; i++ {
		if order[i] != i {
			t.Errorf("Expected order[%d] = %d, got %d", i, i, order[i])
		}
	}
}

// TestConcurrentPublish verifies thread safety with concurrent publishing
func TestConcurrentPublish(t *testing.T) {
	bus := New()
	var count atomic.Int32

	bus.Subscribe("concurrent:test", func(event Event) {
		count.Add(1)
	})

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			bus.Publish(testEvent{eventType: "concurrent:test", data: "test"})
		}()
	}

	wg.Wait()

	if count.Load() != numGoroutines {
		t.Errorf("Expected %d events, got %d", numGoroutines, count.Load())
	}
}

// TestConcurrentSubscribe verifies thread safety with concurrent subscriptions
func TestConcurrentSubscribe(t *testing.T) {
	bus := New()
	var count atomic.Int32

	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			bus.Subscribe("concurrent:subscribe", func(event Event) {
				count.Add(1)
			})
		}()
	}

	wg.Wait()

	bus.Publish(testEvent{eventType: "concurrent:subscribe", data: "test"})

	if count.Load() != numGoroutines {
		t.Errorf("Expected %d subscribers to be called, got %d", numGoroutines, count.Load())
	}
}

// TestConcurrentSubscribeAndPublish verifies thread safety with concurrent operations
func TestConcurrentSubscribeAndPublish(t *testing.T) {
	bus := New()
	var count atomic.Int32

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Concurrent subscriptions
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			bus.Subscribe("mixed:test", func(event Event) {
				count.Add(1)
			})
		}()
	}

	// Small delay to allow some subscriptions
	time.Sleep(10 * time.Millisecond)

	// Concurrent publishes
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			bus.Publish(testEvent{eventType: "mixed:test", data: "test"})
		}()
	}

	wg.Wait()

	// We should have received some events (exact count depends on timing)
	if count.Load() == 0 {
		t.Error("No events were received")
	}
}

// TestMultipleBuses verifies that separate event buses are independent
func TestMultipleBuses(t *testing.T) {
	bus1 := New()
	bus2 := New()

	count1 := 0
	count2 := 0

	bus1.Subscribe("test:event", func(event Event) {
		count1++
	})

	bus2.Subscribe("test:event", func(event Event) {
		count2++
	})

	bus1.Publish(testEvent{eventType: "test:event", data: "test"})

	if count1 != 1 {
		t.Errorf("Bus1: expected 1 event, got %d", count1)
	}
	if count2 != 0 {
		t.Errorf("Bus2: expected 0 events, got %d", count2)
	}
}

// TestEventDataIntegrity verifies that event data is preserved
func TestEventDataIntegrity(t *testing.T) {
	bus := New()

	expected := "test data with special chars: !@#$%^&*()"
	var received string

	bus.Subscribe("data:test", func(event Event) {
		e := event.(testEvent)
		received = e.data
	})

	bus.Publish(testEvent{
		eventType: "data:test",
		data:      expected,
	})

	if received != expected {
		t.Errorf("Expected '%s', got '%s'", expected, received)
	}
}

// BenchmarkPublish benchmarks event publishing performance
func BenchmarkPublish(b *testing.B) {
	bus := New()
	bus.Subscribe("bench:test", func(event Event) {
		// Empty listener
	})

	event := testEvent{eventType: "bench:test", data: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(event)
	}
}

// BenchmarkSubscribe benchmarks subscription performance
func BenchmarkSubscribe(b *testing.B) {
	bus := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Subscribe("bench:test", func(event Event) {})
	}
}

// BenchmarkPublishMultipleListeners benchmarks publishing with many listeners
func BenchmarkPublishMultipleListeners(b *testing.B) {
	bus := New()

	// Add 100 listeners
	for i := 0; i < 100; i++ {
		bus.Subscribe("bench:multi", func(event Event) {})
	}

	event := testEvent{eventType: "bench:multi", data: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish(event)
	}
}
