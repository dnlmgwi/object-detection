package possession

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dnlmgwi/football-object-detection/internal/config"
	"github.com/dnlmgwi/football-object-detection/internal/models"
)

func TestNewTracker(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	assert.NotNil(t, tracker)
	assert.Equal(t, "none", tracker.currentHolder)
}

func TestCalculateDistance(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	playerBox := models.BoundingBox{X: 100, Y: 100, Width: 50, Height: 100}
	ballBox := models.BoundingBox{X: 150, Y: 150, Width: 20, Height: 20}

	distance := tracker.calculateDistance(playerBox, ballBox)

	// Distance should be positive
	assert.Greater(t, distance, 0.0)
}

func TestFindClosestPlayer(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	players := []models.Player{
		{
			ID:   "p1",
			Team: "teamA",
			BBox: models.BoundingBox{X: 100, Y: 100, Width: 50, Height: 100},
		},
		{
			ID:   "p2",
			Team: "teamB",
			BBox: models.BoundingBox{X: 500, Y: 500, Width: 50, Height: 100},
		},
		{
			ID:   "p3",
			Team: "neutral",
			BBox: models.BoundingBox{X: 140, Y: 140, Width: 50, Height: 100},
		},
	}

	ball := &models.Ball{
		BBox: models.BoundingBox{X: 150, Y: 150, Width: 20, Height: 20},
	}

	closest := tracker.findClosestPlayer(players, ball, 5.0)

	// Should find p1 (teamA) as closest non-neutral player
	assert.NotNil(t, closest)
	assert.Equal(t, "p1", closest.ID)
	assert.Equal(t, "teamA", closest.Team)
}

func TestUpdatePossession(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	players := []models.Player{
		{
			ID:   "p1",
			Team: "teamA",
			BBox: models.BoundingBox{X: 100, Y: 100, Width: 50, Height: 100},
		},
	}

	ball := &models.Ball{
		BBox: models.BoundingBox{X: 120, Y: 120, Width: 20, Height: 20},
	}

	possession := tracker.Update(players, ball)

	assert.NotNil(t, possession)
	assert.Equal(t, "teamA", possession.Team)
	assert.Equal(t, "p1", possession.PlayerID)
	assert.Equal(t, "teamA", tracker.currentHolder)
}

func TestPossessionChange(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	playersA := []models.Player{
		{
			ID:   "p1",
			Team: "teamA",
			BBox: models.BoundingBox{X: 100, Y: 100, Width: 50, Height: 100},
		},
	}

	playersB := []models.Player{
		{
			ID:   "p2",
			Team: "teamB",
			BBox: models.BoundingBox{X: 100, Y: 100, Width: 50, Height: 100},
		},
	}

	ball := &models.Ball{
		BBox: models.BoundingBox{X: 120, Y: 120, Width: 20, Height: 20},
	}

	// First update - teamA gets possession
	tracker.Update(playersA, ball)
	assert.Equal(t, "teamA", tracker.currentHolder)

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Second update - teamB gets possession
	tracker.Update(playersB, ball)
	assert.Equal(t, "teamB", tracker.currentHolder)

	// Check history
	history := tracker.GetHistory()
	assert.Len(t, history, 1)
	assert.Equal(t, "teamA", history[0].Team)
}

func TestGetStats(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	// Simulate some possession
	tracker.teamATotal = 60.0 // 60 seconds
	tracker.teamBTotal = 40.0 // 40 seconds
	tracker.currentHolder = "teamA"

	stats := tracker.GetStats()

	assert.Equal(t, 60.0, stats.TeamAPercentage)
	assert.Equal(t, 40.0, stats.TeamBPercentage)
	assert.Equal(t, "teamA", stats.CurrentHolder)
}

func TestReset(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	// Set some state
	tracker.currentHolder = "teamA"
	tracker.teamATotal = 10.0
	tracker.teamBTotal = 5.0

	tracker.Reset()

	assert.Equal(t, "none", tracker.currentHolder)
	assert.Equal(t, 0.0, tracker.teamATotal)
	assert.Equal(t, 0.0, tracker.teamBTotal)
	assert.Len(t, tracker.possessionHistory, 0)
}

func TestUpdateWithNoBall(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	players := []models.Player{
		{
			ID:   "p1",
			Team: "teamA",
			BBox: models.BoundingBox{X: 100, Y: 100, Width: 50, Height: 100},
		},
	}

	// Update with no ball
	possession := tracker.Update(players, nil)

	// Should return nil when no ball and no last position
	assert.Nil(t, possession)
}

func TestPossessionHistoryLimit(t *testing.T) {
	cfg := config.NewManager()
	tracker := NewTracker(cfg)

	// Add more than 100 events
	for i := 0; i < 150; i++ {
		tracker.handlePossessionChange("teamA", 1.0, time.Now())
	}

	history := tracker.GetHistory()
	assert.LessOrEqual(t, len(history), 100)
}
