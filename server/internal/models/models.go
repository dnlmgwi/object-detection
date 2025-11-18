package models

import "image"

// BoundingBox represents a rectangular bounding box
type BoundingBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Center returns the center point of the bounding box
func (b BoundingBox) Center() (int, int) {
	return b.X + b.Width/2, b.Y + b.Height/2
}

// Area returns the area of the bounding box
func (b BoundingBox) Area() int {
	return b.Width * b.Height
}

// ToRect converts to image.Rectangle
func (b BoundingBox) ToRect() image.Rectangle {
	return image.Rect(b.X, b.Y, b.X+b.Width, b.Y+b.Height)
}

// Vector2D represents a 2D vector
type Vector2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Player represents a detected player
type Player struct {
	ID         string      `json:"id"`
	BBox       BoundingBox `json:"bbox"`
	Team       string      `json:"team"` // "teamA", "teamB", or "neutral"
	Confidence float32     `json:"confidence"`
	ColorMatch float64     `json:"colorMatch,omitempty"`
}

// Ball represents a detected ball
type Ball struct {
	BBox       BoundingBox `json:"bbox"`
	Confidence float32     `json:"confidence"`
	Velocity   *Vector2D   `json:"velocity,omitempty"`
}

// Possession represents current ball possession
type Possession struct {
	Team     string  `json:"team"`     // "teamA", "teamB", or "none"
	PlayerID string  `json:"playerId"` // ID of player with possession
	Duration float64 `json:"duration"` // duration in seconds
	Distance float64 `json:"distance"` // distance to ball in meters
}

// DetectionResult represents the result of object detection on a frame
type DetectionResult struct {
	FrameID    int64       `json:"frameId"`
	Timestamp  int64       `json:"timestamp"` // Unix timestamp in milliseconds
	Players    []Player    `json:"players"`
	Ball       *Ball       `json:"ball"`
	Possession *Possession `json:"possession"`
}

// PossessionStats represents possession statistics
type PossessionStats struct {
	TeamAPercentage float64 `json:"teamA"`
	TeamBPercentage float64 `json:"teamB"`
	CurrentHolder   string  `json:"currentHolder"` // "teamA", "teamB", or "none"
	LastChange      int64   `json:"lastChange"`    // Unix timestamp
	TotalDuration   float64 `json:"totalDuration"` // total time tracked in seconds
}

// PerformanceMetrics represents system performance metrics
type PerformanceMetrics struct {
	FPS         float64 `json:"fps"`
	Latency     int64   `json:"latency"`     // milliseconds
	CPUUsage    float64 `json:"cpuUsage"`    // percentage
	MemoryUsage int64   `json:"memoryUsage"` // bytes
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type string      `json:"type"` // "detection", "possession", "metrics", "error"
	Data interface{} `json:"data"`
}
