package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/papiermond/eventbus"
)

// Player events
type PlayerJumpedEvent struct {
	PlayerID string
	Height   float64
}

func (e PlayerJumpedEvent) GetType() eventbus.EventType {
	return "player:jumped"
}

type PlayerDiedEvent struct {
	PlayerID string
	Cause    string
}

func (e PlayerDiedEvent) GetType() eventbus.EventType {
	return "player:died"
}

type PlayerRespawnedEvent struct {
	PlayerID string
	X, Y     float64
}

func (e PlayerRespawnedEvent) GetType() eventbus.EventType {
	return "player:respawned"
}

// World events
type LevelLoadedEvent struct {
	LevelName string
}

func (e LevelLoadedEvent) GetType() eventbus.EventType {
	return "world:level_loaded"
}

// Game systems
type AudioSystem struct {
	bus eventbus.EventBus
}

func NewAudioSystem(bus eventbus.EventBus) *AudioSystem {
	s := &AudioSystem{bus: bus}

	// Subscribe to events that trigger sounds
	bus.Subscribe("player:jumped", func(event eventbus.Event) {
		fmt.Println("  [Audio] Playing jump sound")
	})

	bus.Subscribe("player:died", func(event eventbus.Event) {
		fmt.Println("  [Audio] Playing death sound")
	})

	return s
}

type PhysicsSystem struct {
	bus eventbus.EventBus
}

func NewPhysicsSystem(bus eventbus.EventBus) *PhysicsSystem {
	s := &PhysicsSystem{bus: bus}

	bus.Subscribe("player:jumped", func(event eventbus.Event) {
		e := event.(PlayerJumpedEvent)
		fmt.Printf("  [Physics] Applying jump force (height: %.1f)\n", e.Height)
	})

	return s
}

type RenderSystem struct {
	bus eventbus.EventBus
}

func NewRenderSystem(bus eventbus.EventBus) *RenderSystem {
	s := &RenderSystem{bus: bus}

	bus.Subscribe("player:respawned", func(event eventbus.Event) {
		e := event.(PlayerRespawnedEvent)
		fmt.Printf("  [Render] Moved camera to (%.0f, %.0f)\n", e.X, e.Y)
	})

	bus.Subscribe("world:level_loaded", func(event eventbus.Event) {
		e := event.(LevelLoadedEvent)
		fmt.Printf("  [Render] Loading textures for level '%s'\n", e.LevelName)
	})

	return s
}

type Player struct {
	ID       string
	bus      eventbus.EventBus
	isAlive  bool
	position struct{ x, y float64 }
}

func NewPlayer(id string, bus eventbus.EventBus) *Player {
	p := &Player{
		ID:      id,
		bus:     bus,
		isAlive: true,
	}
	p.position.x = 0
	p.position.y = 0
	return p
}

func (p *Player) Jump(height float64) {
	if p.isAlive {
		fmt.Printf("\n[Player] %s jumps!\n", p.ID)
		p.bus.Publish(PlayerJumpedEvent{
			PlayerID: p.ID,
			Height:   height,
		})
	}
}

func (p *Player) Die(cause string) {
	if p.isAlive {
		fmt.Printf("\n[Player] %s died from %s\n", p.ID, cause)
		p.isAlive = false
		p.bus.Publish(PlayerDiedEvent{
			PlayerID: p.ID,
			Cause:    cause,
		})

		// Respawn after delay
		go func() {
			time.Sleep(500 * time.Millisecond)
			p.Respawn()
		}()
	}
}

func (p *Player) Respawn() {
	fmt.Printf("\n[Player] %s respawning...\n", p.ID)
	p.isAlive = true
	p.position.x = 0
	p.position.y = 0
	p.bus.Publish(PlayerRespawnedEvent{
		PlayerID: p.ID,
		X:        p.position.x,
		Y:        p.position.y,
	})
}

type World struct {
	bus eventbus.EventBus
}

func NewWorld(bus eventbus.EventBus) *World {
	return &World{bus: bus}
}

func (w *World) LoadLevel(name string) {
	fmt.Printf("\n[World] Loading level '%s'\n", name)
	w.bus.Publish(LevelLoadedEvent{LevelName: name})
}

func main() {
	fmt.Println("=== Game Event Bus Example ===")

	// Create event buses for different game domains
	playerBus := eventbus.New()
	worldBus := eventbus.New()

	// Initialize game systems
	NewAudioSystem(playerBus)
	NewPhysicsSystem(playerBus)
	NewRenderSystem(playerBus)
	NewRenderSystem(worldBus)

	// Create game objects
	player := NewPlayer("player-1", playerBus)
	world := NewWorld(worldBus)

	// Simulate game events
	world.LoadLevel("level-1")

	time.Sleep(200 * time.Millisecond)

	player.Jump(10.0)

	time.Sleep(200 * time.Millisecond)

	player.Die("falling")

	// Wait for respawn
	var wg sync.WaitGroup
	wg.Add(1)
	playerBus.Subscribe("player:respawned", func(event eventbus.Event) {
		wg.Done()
	})

	wg.Wait()

	player.Jump(8.0)

	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== Example Complete ===")
}
