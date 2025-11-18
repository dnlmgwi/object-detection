package websocket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dnlmgwi/football-object-detection/internal/models"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
	assert.False(t, hub.running)
}

func TestHubRun(t *testing.T) {
	hub := NewHub()

	// Start hub in goroutine
	go hub.Run()

	// Wait for hub to start
	time.Sleep(10 * time.Millisecond)

	hub.mu.RLock()
	running := hub.running
	hub.mu.RUnlock()

	assert.True(t, running)

	// Shutdown
	hub.Shutdown()
}

func TestGetClientCount(t *testing.T) {
	hub := NewHub()

	count := hub.GetClientCount()
	assert.Equal(t, 0, count)
}

func TestBroadcastDetection(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Shutdown()

	result := models.DetectionResult{
		FrameID:   123,
		Timestamp: time.Now().UnixMilli(),
		Players: []models.Player{
			{ID: "p1", Team: "teamA", Confidence: 0.9},
		},
		Ball: &models.Ball{
			BBox:       models.BoundingBox{X: 100, Y: 100, Width: 20, Height: 20},
			Confidence: 0.85,
		},
	}

	err := hub.BroadcastDetection(result)
	assert.NoError(t, err)
}

func TestBroadcastPossession(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Shutdown()

	stats := models.PossessionStats{
		TeamAPercentage: 60.0,
		TeamBPercentage: 40.0,
		CurrentHolder:   "teamA",
		LastChange:      time.Now().UnixMilli(),
		TotalDuration:   100.0,
	}

	err := hub.BroadcastPossession(stats)
	assert.NoError(t, err)
}

func TestBroadcastMetrics(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Shutdown()

	metrics := models.PerformanceMetrics{
		FPS:         15.5,
		Latency:     85,
		CPUUsage:    45.2,
		MemoryUsage: 1024 * 1024 * 512,
	}

	err := hub.BroadcastMetrics(metrics)
	assert.NoError(t, err)
}

func TestBroadcastError(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Shutdown()

	err := hub.BroadcastError("Test error message")
	assert.NoError(t, err)
}

func TestShutdown(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(10 * time.Millisecond)

	hub.Shutdown()

	hub.mu.RLock()
	running := hub.running
	clientCount := len(hub.clients)
	hub.mu.RUnlock()

	assert.False(t, running)
	assert.Equal(t, 0, clientCount)
}

func TestClientRegistration(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Shutdown()

	// Create mock client
	client := &Client{
		hub:  hub,
		send: make(chan []byte, 256),
	}

	// Register client
	hub.register <- client

	// Wait for registration
	time.Sleep(10 * time.Millisecond)

	count := hub.GetClientCount()
	assert.Equal(t, 1, count)

	// Unregister client
	hub.unregister <- client

	// Wait for unregistration
	time.Sleep(10 * time.Millisecond)

	count = hub.GetClientCount()
	assert.Equal(t, 0, count)
}
