package models

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoundingBoxCenter(t *testing.T) {
	bbox := BoundingBox{X: 100, Y: 100, Width: 50, Height: 80}
	x, y := bbox.Center()

	assert.Equal(t, 125, x)
	assert.Equal(t, 140, y)
}

func TestBoundingBoxArea(t *testing.T) {
	bbox := BoundingBox{X: 0, Y: 0, Width: 50, Height: 80}
	area := bbox.Area()

	assert.Equal(t, 4000, area)
}

func TestBoundingBoxToRect(t *testing.T) {
	bbox := BoundingBox{X: 10, Y: 20, Width: 30, Height: 40}
	rect := bbox.ToRect()

	expected := image.Rect(10, 20, 40, 60)
	assert.Equal(t, expected, rect)
}

func TestPlayerModel(t *testing.T) {
	player := Player{
		ID:         "p1",
		BBox:       BoundingBox{X: 100, Y: 100, Width: 50, Height: 100},
		Team:       "teamA",
		Confidence: 0.95,
		ColorMatch: 0.87,
	}

	assert.Equal(t, "p1", player.ID)
	assert.Equal(t, "teamA", player.Team)
	assert.Equal(t, float32(0.95), player.Confidence)
}

func TestBallModel(t *testing.T) {
	ball := Ball{
		BBox:       BoundingBox{X: 320, Y: 240, Width: 20, Height: 20},
		Confidence: 0.88,
		Velocity:   &Vector2D{X: 5.0, Y: -2.0},
	}

	assert.Equal(t, 20, ball.BBox.Width)
	assert.NotNil(t, ball.Velocity)
	assert.Equal(t, 5.0, ball.Velocity.X)
}

func TestPossessionModel(t *testing.T) {
	possession := Possession{
		Team:     "teamA",
		PlayerID: "p1",
		Duration: 5.2,
		Distance: 1.5,
	}

	assert.Equal(t, "teamA", possession.Team)
	assert.Equal(t, "p1", possession.PlayerID)
	assert.Equal(t, 5.2, possession.Duration)
}

func TestDetectionResult(t *testing.T) {
	result := DetectionResult{
		FrameID:   123,
		Timestamp: 1634567890123,
		Players: []Player{
			{ID: "p1", Team: "teamA", Confidence: 0.9},
			{ID: "p2", Team: "teamB", Confidence: 0.85},
		},
		Ball: &Ball{
			BBox:       BoundingBox{X: 320, Y: 240, Width: 20, Height: 20},
			Confidence: 0.75,
		},
		Possession: &Possession{
			Team:     "teamA",
			PlayerID: "p1",
			Duration: 3.5,
		},
	}

	assert.Equal(t, int64(123), result.FrameID)
	assert.Len(t, result.Players, 2)
	assert.NotNil(t, result.Ball)
	assert.NotNil(t, result.Possession)
	assert.Equal(t, "teamA", result.Possession.Team)
}
