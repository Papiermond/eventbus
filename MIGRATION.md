# EventBus Library Migration Guide

This document explains how the eventbus system was extracted from the goGameSDL project into a standalone library.

## What Was Done

### 1. Created Standalone Library

The eventbus implementation was extracted to a new repository at `../eventbus/`:

```
eventbus/
├── go.mod                  # Module definition
├── eventbus.go            # Core implementation with full documentation
├── eventbus_test.go       # Comprehensive tests (11 tests + benchmarks)
├── README.md              # Complete usage documentation
├── LICENSE                # MIT License
├── .gitignore            # Standard Go .gitignore
└── examples/             # Working examples
    ├── basic/            # Basic publish/subscribe
    ├── multiple_buses/   # Domain separation pattern
    └── game/             # Game engine integration example
```

### 2. Library Features

The standalone library includes:

- **Core functionality**: Thread-safe event bus with publish/subscribe
- **Zero dependencies**: Only uses Go standard library
- **Comprehensive tests**: 11 unit tests + 3 benchmarks, all passing
- **Full documentation**: Package docs, README, examples
- **Type safety**: Generic interfaces for compile-time safety
- **High performance**: Minimal overhead, mutex-based synchronization

### 3. Updated Game Project

The game project now uses the external library:

**go.mod changes:**
```go
require (
    github.com/papiermond/eventbus v0.0.0
)

replace github.com/papiermond/eventbus => ../eventbus
```

**internal/eventbus/eventbus.go:**
```go
import extbus "github.com/papiermond/eventbus"

// Re-export external library types
type Event = extbus.Event
type EventType = extbus.EventType
type EventListener = extbus.EventListener
type EventBus = extbus.EventBus

// Game-specific container
type EventBuses struct {
    Application EventBus
    Input       EventBus
    Physics     EventBus
    // ... etc
}
```

### 4. Backward Compatibility

The migration maintains 100% backward compatibility:

- All imports remain unchanged: `import "game/internal/eventbus"`
- All APIs remain identical
- All existing code works without modification
- All tests pass (51 tests total)

## Using the Library in Other Projects

### Installation

```bash
go get github.com/papiermond/eventbus
```

### Basic Usage

```go
package main

import "github.com/papiermond/eventbus"

type MyEvent struct {
    Data string
}

func (e MyEvent) GetType() eventbus.EventType {
    return "my:event"
}

func main() {
    bus := eventbus.New()

    bus.Subscribe("my:event", func(event eventbus.Event) {
        e := event.(MyEvent)
        println("Received:", e.Data)
    })

    bus.Publish(MyEvent{Data: "Hello"})
}
```

## Publishing the Library

To make this library publicly available:

1. **Create GitHub repository**:
   ```bash
   cd ../eventbus
   git init
   git add .
   git commit -m "Initial commit: EventBus library"
   git remote add origin https://github.com/papiermond/eventbus.git
   git push -u origin main
   ```

2. **Tag a release**:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

3. **Update game project** to use published version:
   ```bash
   cd ../goGameSDL
   # Remove replace directive from go.mod
   go get github.com/papiermond/eventbus@v0.1.0
   go mod tidy
   ```

## Benefits of Extraction

1. **Reusability**: Can be used in any Go project, not just games
2. **Focused development**: Library can evolve independently
3. **Better testing**: Library has its own comprehensive test suite
4. **Documentation**: Proper package documentation and examples
5. **Community**: Others can contribute improvements
6. **Versioning**: Semantic versioning for stable API

## Next Steps

- [ ] Create GitHub repository for eventbus library
- [ ] Publish to GitHub and tag v0.1.0
- [ ] Update game project to use published version
- [ ] Consider adding to pkg.go.dev
- [ ] Add CI/CD for automated testing
- [ ] Add badges (tests, coverage, go report)

## Notes

- The library uses Go 1.23 features
- All operations are thread-safe
- Listeners are called synchronously
- For async operations, spawn goroutines inside listeners
- The game project maintains its domain-specific `EventBuses` container
