# EventBus

A lightweight, thread-safe event bus library for Go that enables loosely coupled, event-driven architectures.

## Features

- **Thread-Safe**: All operations are protected by mutex locks, safe for concurrent use
- **Type-Safe**: Uses Go interfaces for compile-time type safety
- **Zero Dependencies**: No external dependencies beyond the Go standard library
- **Simple API**: Clean, intuitive interface with just two main methods
- **High Performance**: Minimal overhead, suitable for high-frequency events
- **Independent Buses**: Create multiple isolated event buses for different domains

## Installation

```bash
go get github.com/papiermond/eventbus
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/papiermond/eventbus"
)

// Define your event type
type UserLoggedIn struct {
    UserID   string
    Username string
}

// Implement the Event interface
func (e UserLoggedIn) GetType() eventbus.EventType {
    return "user:logged_in"
}

func main() {
    // Create an event bus
    bus := eventbus.New()

    // Subscribe to events
    bus.Subscribe("user:logged_in", func(event eventbus.Event) {
        e := event.(UserLoggedIn)
        fmt.Printf("User %s logged in!\n", e.Username)
    })

    // Publish an event
    bus.Publish(UserLoggedIn{
        UserID:   "123",
        Username: "alice",
    })
}
```

## Core Concepts

### Event

An event is any type that implements the `Event` interface:

```go
type Event interface {
    GetType() EventType
}
```

Events should be immutable value types for thread safety. Include all relevant data in the event struct.

### EventType

An `EventType` is a string identifier used to route events to subscribers:

```go
type EventType string
```

Convention: Use namespaced identifiers like `"domain:action"` (e.g., `"user:login"`, `"payment:completed"`).

### EventListener

A listener is a function that handles events:

```go
type EventListener func(Event)
```

Listeners are called synchronously in the order they were registered.

### EventBus

The event bus coordinates publishing and subscribing:

```go
type EventBus interface {
    Subscribe(eventType EventType, listener EventListener)
    Publish(event Event)
}
```

## Usage Examples

### Basic Publish/Subscribe

```go
bus := eventbus.New()

// Subscribe to an event
bus.Subscribe("order:created", func(event eventbus.Event) {
    order := event.(OrderCreatedEvent)
    fmt.Printf("Order %s created for $%.2f\n", order.ID, order.Total)
})

// Publish an event
bus.Publish(OrderCreatedEvent{
    ID:    "order-123",
    Total: 99.99,
})
```

### Multiple Subscribers

Multiple listeners can subscribe to the same event type:

```go
bus := eventbus.New()

// Logger
bus.Subscribe("user:login", func(event eventbus.Event) {
    log.Println("Login event:", event)
})

// Analytics
bus.Subscribe("user:login", func(event eventbus.Event) {
    analytics.Track("login", event)
})

// Notification
bus.Subscribe("user:login", func(event eventbus.Event) {
    notifications.Send("Welcome back!")
})
```

### Multiple Event Buses

Create separate buses for different domains:

```go
applicationBus := eventbus.New()
physicsBus := eventbus.New()
audioBus := eventbus.New()

// Each bus is independent
applicationBus.Subscribe("app:quit", handleQuit)
physicsBus.Subscribe("collision:detected", handleCollision)
audioBus.Subscribe("sound:play", playSound)
```

### Type Assertions

Safely extract event data with type assertions:

```go
bus.Subscribe("payment:completed", func(event eventbus.Event) {
    // Type assert to access specific fields
    payment, ok := event.(PaymentCompletedEvent)
    if !ok {
        log.Error("Unexpected event type")
        return
    }

    processPayment(payment.Amount, payment.Currency)
})
```

### Domain-Specific Event Naming

Use a consistent naming convention for events:

```go
// Domain:Action pattern
const (
    UserLoggedIn     eventbus.EventType = "user:logged_in"
    UserLoggedOut    eventbus.EventType = "user:logged_out"
    OrderCreated     eventbus.EventType = "order:created"
    OrderShipped     eventbus.EventType = "order:shipped"
    PaymentCompleted eventbus.EventType = "payment:completed"
)
```

## Advanced Patterns

### Organizing Events by Domain

```go
package events

import "github.com/papiermond/eventbus"

// User events
type UserLoggedIn struct {
    UserID string
}

func (e UserLoggedIn) GetType() eventbus.EventType {
    return "user:logged_in"
}

// Order events
type OrderCreated struct {
    OrderID string
    UserID  string
    Total   float64
}

func (e OrderCreated) GetType() eventbus.EventType {
    return "order:created"
}
```

### Event Bus Container

Group related buses together:

```go
type EventBuses struct {
    Application eventbus.EventBus
    Physics     eventbus.EventBus
    Audio       eventbus.EventBus
    UI          eventbus.EventBus
}

func NewEventBuses() *EventBuses {
    return &EventBuses{
        Application: eventbus.New(),
        Physics:     eventbus.New(),
        Audio:       eventbus.New(),
        UI:          eventbus.New(),
    }
}
```

### Decoupling Components

Event buses enable loose coupling between components:

```go
// Producer doesn't know about consumers
type OrderService struct {
    bus eventbus.EventBus
}

func (s *OrderService) CreateOrder(order Order) {
    // Business logic...
    s.bus.Publish(OrderCreatedEvent{ID: order.ID})
}

// Consumer doesn't know about producer
type EmailService struct {
    bus eventbus.EventBus
}

func (s *EmailService) Start() {
    s.bus.Subscribe("order:created", func(event eventbus.Event) {
        e := event.(OrderCreatedEvent)
        s.sendOrderConfirmation(e.ID)
    })
}
```

## Thread Safety

All operations are thread-safe and can be called from multiple goroutines:

```go
bus := eventbus.New()

// Safe to subscribe from multiple goroutines
go bus.Subscribe("event:type", handler1)
go bus.Subscribe("event:type", handler2)

// Safe to publish from multiple goroutines
go bus.Publish(event1)
go bus.Publish(event2)
```

## Performance Considerations

- **Listeners are called synchronously**: Long-running listeners will block the publisher
- **Use goroutines for async work**: Spawn goroutines inside listeners for non-blocking operations
- **Consider listener count**: More listeners = more overhead per publish

```go
// Async listener for slow operations
bus.Subscribe("heavy:task", func(event eventbus.Event) {
    go func() {
        // Long-running work won't block publisher
        processHeavyTask(event)
    }()
})
```

## Testing

The library includes comprehensive tests:

```bash
# Run tests
go test

# Run tests with coverage
go test -cover

# Run tests with race detection
go test -race

# Run benchmarks
go test -bench=.
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

MIT License - see LICENSE file for details.

## Examples

See the [examples](examples/) directory for complete working examples:

- [Basic Usage](examples/basic/main.go) - Simple publish/subscribe
- [Multiple Buses](examples/multiple_buses/main.go) - Domain separation
- [Game Events](examples/game/main.go) - Real-world game engine integration

## API Reference

See [GoDoc](https://pkg.go.dev/github.com/papiermond/eventbus) for detailed API documentation.
