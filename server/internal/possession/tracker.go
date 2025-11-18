package possession

import (
	"math"
	"sync"
	"time"

	"github.com/dnlmgwi/football-object-detection/internal/config"
	"github.com/dnlmgwi/football-object-detection/internal/models"
)

// Tracker tracks ball possession
type Tracker struct {
	config *config.Manager

	mu                sync.RWMutex
	currentHolder     string  // "teamA", "teamB", or "none"
	currentPlayerID   string
	currentDuration   float64
	lastChangeTime    time.Time
	teamATotal        float64 // seconds
	teamBTotal        float64 // seconds
	totalTrackedTime  float64
	lastBallPosition  *models.BoundingBox
	possessionHistory []PossessionEvent
}

// PossessionEvent represents a possession change event
type PossessionEvent struct {
	Timestamp time.Time
	Team      string
	PlayerID  string
	Duration  float64
}

// NewTracker creates a new possession tracker
func NewTracker(cfg *config.Manager) *Tracker {
	return &Tracker{
		config:            cfg,
		currentHolder:     "none",
		lastChangeTime:    time.Now(),
		possessionHistory: make([]PossessionEvent, 0),
	}
}

// Update updates possession based on current detections
func (t *Tracker) Update(players []models.Player, ball *models.Ball) *models.Possession {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(t.lastChangeTime).Seconds()

	if ball == nil {
		// No ball detected, use last known position if available
		if t.lastBallPosition == nil {
			return nil
		}
		ball = &models.Ball{BBox: *t.lastBallPosition}
	} else {
		t.lastBallPosition = &ball.BBox
	}

	// Find closest player to ball
	cfg := t.config.GetConfig()
	proximityThreshold := cfg.Possession.ProximityThreshold

	closestPlayer := t.findClosestPlayer(players, ball, proximityThreshold)

	// Determine new holder
	newHolder := "none"
	newPlayerID := ""
	distance := math.MaxFloat64

	if closestPlayer != nil {
		newHolder = closestPlayer.Team
		newPlayerID = closestPlayer.ID
		distance = t.calculateDistance(closestPlayer.BBox, ball.BBox)
	}

	// Check for possession change
	if newHolder != t.currentHolder {
		t.handlePossessionChange(t.currentHolder, t.currentDuration, now)
		t.currentHolder = newHolder
		t.currentPlayerID = newPlayerID
		t.currentDuration = 0
		t.lastChangeTime = now
	} else {
		// Update duration
		t.currentDuration += elapsed
		t.totalTrackedTime += elapsed

		// Update team totals
		if t.currentHolder == "teamA" {
			t.teamATotal += elapsed
		} else if t.currentHolder == "teamB" {
			t.teamBTotal += elapsed
		}

		t.lastChangeTime = now
	}

	return &models.Possession{
		Team:     t.currentHolder,
		PlayerID: t.currentPlayerID,
		Duration: t.currentDuration,
		Distance: distance,
	}
}

// GetStats returns current possession statistics
func (t *Tracker) GetStats() models.PossessionStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	totalPossessionTime := t.teamATotal + t.teamBTotal
	teamAPercentage := 0.0
	teamBPercentage := 0.0

	if totalPossessionTime > 0 {
		teamAPercentage = (t.teamATotal / totalPossessionTime) * 100
		teamBPercentage = (t.teamBTotal / totalPossessionTime) * 100
	}

	return models.PossessionStats{
		TeamAPercentage: teamAPercentage,
		TeamBPercentage: teamBPercentage,
		CurrentHolder:   t.currentHolder,
		LastChange:      t.lastChangeTime.UnixMilli(),
		TotalDuration:   totalPossessionTime,
	}
}

// Reset resets possession tracking
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.currentHolder = "none"
	t.currentPlayerID = ""
	t.currentDuration = 0
	t.lastChangeTime = time.Now()
	t.teamATotal = 0
	t.teamBTotal = 0
	t.totalTrackedTime = 0
	t.lastBallPosition = nil
	t.possessionHistory = make([]PossessionEvent, 0)
}

// GetHistory returns possession change history
func (t *Tracker) GetHistory() []PossessionEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()

	history := make([]PossessionEvent, len(t.possessionHistory))
	copy(history, t.possessionHistory)
	return history
}

// findClosestPlayer finds the player closest to the ball within threshold
func (t *Tracker) findClosestPlayer(players []models.Player, ball *models.Ball, threshold float64) *models.Player {
	var closest *models.Player
	minDistance := math.MaxFloat64

	for i := range players {
		// Skip neutral players (referees, etc.)
		if players[i].Team == "neutral" {
			continue
		}

		distance := t.calculateDistance(players[i].BBox, ball.BBox)
		if distance < minDistance && distance <= threshold {
			minDistance = distance
			closest = &players[i]
		}
	}

	return closest
}

// calculateDistance calculates distance between player and ball (simplified 2D Euclidean)
func (t *Tracker) calculateDistance(playerBox, ballBox models.BoundingBox) float64 {
	px, py := playerBox.Center()
	bx, by := ballBox.Center()

	dx := float64(px - bx)
	dy := float64(py - by)

	// Convert pixel distance to meters (rough estimate)
	// Assume average player height ~1.8m and bbox height in pixels
	// This is a simplified conversion and should be calibrated
	pixelDistance := math.Sqrt(dx*dx + dy*dy)

	// Rough calibration: 100 pixels ≈ 1 meter (this should be configurable)
	metersPerPixel := 0.01
	return pixelDistance * metersPerPixel
}

// handlePossessionChange records a possession change event
func (t *Tracker) handlePossessionChange(prevHolder string, duration float64, timestamp time.Time) {
	if prevHolder == "none" {
		return
	}

	event := PossessionEvent{
		Timestamp: timestamp,
		Team:      prevHolder,
		PlayerID:  t.currentPlayerID,
		Duration:  duration,
	}

	t.possessionHistory = append(t.possessionHistory, event)

	// Keep only last 100 events to prevent memory growth
	if len(t.possessionHistory) > 100 {
		t.possessionHistory = t.possessionHistory[1:]
	}
}
