// Package eventbus provides a thread-safe, type-safe event bus implementation
// for event-driven architectures in Go applications.
//
// The eventbus package enables loosely coupled communication between components
// through a publish-subscribe pattern. Components can publish events and subscribe
// to specific event types without direct dependencies on each other.
//
// Key features:
//   - Thread-safe operations with mutex-based synchronization
//   - Type-safe event handling with Go interfaces
//   - Support for multiple independent event buses
//   - Zero external dependencies
//   - Simple and intuitive API
//
// Basic usage:
//
//	// Define your event type
//	type UserLoggedIn struct {
//	    UserID string
//	    Timestamp time.Time
//	}
//
//	func (e UserLoggedIn) GetType() eventbus.EventType {
//	    return "user:logged_in"
//	}
//
//	// Create an event bus
//	bus := eventbus.New()
//
//	// Subscribe to events
//	bus.Subscribe("user:logged_in", func(event eventbus.Event) {
//	    e := event.(UserLoggedIn)
//	    fmt.Printf("User %s logged in at %v\n", e.UserID, e.Timestamp)
//	})
//
//	// Publish events
//	bus.Publish(UserLoggedIn{
//	    UserID: "user123",
//	    Timestamp: time.Now(),
//	})
package eventbus

import "sync"

// EventType represents the type identifier for an event.
// It's used to match events with their subscribers.
type EventType string

// Event is the interface that all events must implement.
// Events should be immutable value types for thread safety.
type Event interface {
	// GetType returns the event's type identifier.
	// This is used to route events to the appropriate subscribers.
	GetType() EventType
}

// EventListener is a function that handles an event.
// Listeners are called synchronously when an event is published.
// Listeners should not block for long periods as they will delay
// other listeners and the publisher.
type EventListener func(Event)

// EventBus provides thread-safe publish-subscribe functionality
// for event-driven communication between components.
type EventBus interface {
	// Subscribe registers a listener for a specific event type.
	// Multiple listeners can subscribe to the same event type.
	// Listeners are called in the order they were registered.
	//
	// Example:
	//   bus.Subscribe("user:login", func(event Event) {
	//       fmt.Println("User logged in:", event)
	//   })
	Subscribe(eventType EventType, listener EventListener)

	// Publish sends an event to all registered listeners for that event type.
	// Listeners are called synchronously in registration order.
	// If no listeners are registered for the event type, the event is silently dropped.
	//
	// Example:
	//   bus.Publish(UserLoginEvent{UserID: "123"})
	Publish(event Event)
}

// eventBusImpl is the internal implementation of EventBus.
// It uses a mutex to ensure thread-safe access to the listeners map.
type eventBusImpl struct {
	listeners map[EventType][]EventListener
	mutex     sync.Mutex
}

// New creates a new event bus instance.
// Each event bus is independent and maintains its own set of subscribers.
//
// Example:
//
//	bus := eventbus.New()
func New() EventBus {
	return &eventBusImpl{
		listeners: make(map[EventType][]EventListener),
	}
}

// Subscribe registers a listener for a specific event type.
func (bus *eventBusImpl) Subscribe(eventType EventType, listener EventListener) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	bus.listeners[eventType] = append(bus.listeners[eventType], listener)
}

// Publish sends an event to all registered listeners for that event type.
func (bus *eventBusImpl) Publish(event Event) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	if listeners, ok := bus.listeners[event.GetType()]; ok {
		for _, listener := range listeners {
			listener(event)
		}
	}
}
