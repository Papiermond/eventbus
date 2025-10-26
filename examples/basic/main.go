package main

import (
	"fmt"
	"time"

	"github.com/papiermond/eventbus"
)

// UserLoggedIn event
type UserLoggedIn struct {
	UserID    string
	Username  string
	Timestamp time.Time
}

func (e UserLoggedIn) GetType() eventbus.EventType {
	return "user:logged_in"
}

// UserLoggedOut event
type UserLoggedOut struct {
	UserID    string
	Timestamp time.Time
}

func (e UserLoggedOut) GetType() eventbus.EventType {
	return "user:logged_out"
}

// MessageSent event
type MessageSent struct {
	From    string
	To      string
	Content string
}

func (e MessageSent) GetType() eventbus.EventType {
	return "message:sent"
}

func main() {
	fmt.Println("=== EventBus Basic Example ===\n")

	// Create an event bus
	bus := eventbus.New()

	// Subscribe to user login events
	bus.Subscribe("user:logged_in", func(event eventbus.Event) {
		e := event.(UserLoggedIn)
		fmt.Printf("[Login Handler] User %s (%s) logged in at %s\n",
			e.Username, e.UserID, e.Timestamp.Format("15:04:05"))
	})

	// Subscribe to user logout events
	bus.Subscribe("user:logged_out", func(event eventbus.Event) {
		e := event.(UserLoggedOut)
		fmt.Printf("[Logout Handler] User %s logged out at %s\n",
			e.UserID, e.Timestamp.Format("15:04:05"))
	})

	// Multiple subscribers for the same event
	bus.Subscribe("user:logged_in", func(event eventbus.Event) {
		e := event.(UserLoggedIn)
		fmt.Printf("[Analytics] Recording login for user %s\n", e.UserID)
	})

	// Subscribe to message events
	bus.Subscribe("message:sent", func(event eventbus.Event) {
		e := event.(MessageSent)
		fmt.Printf("[Message] %s -> %s: %s\n", e.From, e.To, e.Content)
	})

	// Publish some events
	fmt.Println("Publishing events:\n")

	bus.Publish(UserLoggedIn{
		UserID:    "user-123",
		Username:  "alice",
		Timestamp: time.Now(),
	})

	time.Sleep(100 * time.Millisecond)

	bus.Publish(MessageSent{
		From:    "alice",
		To:      "bob",
		Content: "Hello, Bob!",
	})

	time.Sleep(100 * time.Millisecond)

	bus.Publish(UserLoggedOut{
		UserID:    "user-123",
		Timestamp: time.Now(),
	})

	fmt.Println("\n=== Example Complete ===")
}
