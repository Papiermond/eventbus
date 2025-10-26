package main

import (
	"fmt"

	"github.com/papiermond/eventbus"
)

// Application events
type AppQuitEvent struct{}

func (e AppQuitEvent) GetType() eventbus.EventType {
	return "app:quit"
}

// Physics events
type CollisionEvent struct {
	Object1 string
	Object2 string
	Force   float64
}

func (e CollisionEvent) GetType() eventbus.EventType {
	return "physics:collision"
}

// Audio events
type SoundPlayEvent struct {
	SoundID string
	Volume  float32
}

func (e SoundPlayEvent) GetType() eventbus.EventType {
	return "audio:play"
}

// EventBuses holds all domain-specific event buses
type EventBuses struct {
	Application eventbus.EventBus
	Physics     eventbus.EventBus
	Audio       eventbus.EventBus
}

func NewEventBuses() *EventBuses {
	return &EventBuses{
		Application: eventbus.New(),
		Physics:     eventbus.New(),
		Audio:       eventbus.New(),
	}
}

func main() {
	fmt.Println("=== Multiple Event Buses Example ===\n")

	// Create separate buses for different domains
	buses := NewEventBuses()

	// Subscribe to application events
	buses.Application.Subscribe("app:quit", func(event eventbus.Event) {
		fmt.Println("[Application] Shutting down...")
	})

	// Subscribe to physics events
	buses.Physics.Subscribe("physics:collision", func(event eventbus.Event) {
		e := event.(CollisionEvent)
		fmt.Printf("[Physics] Collision: %s <-> %s (force: %.2f)\n",
			e.Object1, e.Object2, e.Force)
	})

	// Subscribe to audio events
	buses.Audio.Subscribe("audio:play", func(event eventbus.Event) {
		e := event.(SoundPlayEvent)
		fmt.Printf("[Audio] Playing sound '%s' at volume %.1f%%\n",
			e.SoundID, e.Volume*100)
	})

	// Multiple handlers for physics events
	buses.Physics.Subscribe("physics:collision", func(event eventbus.Event) {
		e := event.(CollisionEvent)
		if e.Force > 50.0 {
			// Trigger sound effect on strong collision
			buses.Audio.Publish(SoundPlayEvent{
				SoundID: "impact_heavy",
				Volume:  0.8,
			})
		}
	})

	// Publish events to different buses
	fmt.Println("Publishing events:\n")

	buses.Physics.Publish(CollisionEvent{
		Object1: "Player",
		Object2: "Wall",
		Force:   25.0,
	})

	fmt.Println()

	buses.Physics.Publish(CollisionEvent{
		Object1: "Player",
		Object2: "Enemy",
		Force:   75.0, // This will trigger audio
	})

	fmt.Println()

	buses.Audio.Publish(SoundPlayEvent{
		SoundID: "background_music",
		Volume:  0.5,
	})

	fmt.Println()

	buses.Application.Publish(AppQuitEvent{})

	fmt.Println("\n=== Example Complete ===")
}
