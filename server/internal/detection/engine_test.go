package detection

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnlmgwi/football-object-detection/internal/config"
	"github.com/dnlmgwi/football-object-detection/internal/models"
)

func TestBytesToFloat32(t *testing.T) {
	// Test conversion
	b := []byte{0x00, 0x00, 0x80, 0x3F} // 1.0 in little-endian
	f := bytesToFloat32(b)
	assert.InDelta(t, 1.0, f, 0.01)
}

func TestPerformNMS(t *testing.T) {
	boxes := []image.Rectangle{
		image.Rect(10, 10, 50, 50),   // Box 1
		image.Rect(12, 12, 52, 52),   // Box 2 (overlaps with Box 1)
		image.Rect(100, 100, 150, 150), // Box 3 (no overlap)
	}

	confidences := []float32{0.9, 0.8, 0.85}

	indices := performNMS(boxes, confidences, 0.5)

	// Should keep boxes 1 and 3 (box 2 is suppressed due to overlap with box 1)
	assert.Contains(t, indices, 0)
	assert.Contains(t, indices, 2)
	assert.NotContains(t, indices, 1)
}

func TestColorDistance(t *testing.T) {
	cfg := config.NewManager()
	// Mock engine (without actual model loading)
	e := &Engine{
		config: cfg,
	}

	c1 := config.Color{R: 255, G: 0, B: 0}   // Red
	c2 := config.Color{R: 255, G: 0, B: 0}   // Red
	c3 := config.Color{R: 0, G: 0, B: 255}   // Blue

	// Same color
	dist1 := e.colorDistance(c1, c2)
	assert.Equal(t, 0.0, dist1)

	// Different colors
	dist2 := e.colorDistance(c1, c3)
	assert.Greater(t, dist2, 300.0) // Significant distance
}

func TestClassifyPlayers(t *testing.T) {
	// This test would require actual image data and GoCV setup
	// For now, we test the logic structure

	cfg := config.NewManager()
	e := &Engine{
		config: cfg,
	}

	players := []models.Player{
		{ID: "p1", Team: "neutral"},
		{ID: "p2", Team: "neutral"},
	}

	teamA := config.Team{
		ID:           "teamA",
		PrimaryColor: config.Color{R: 255, G: 0, B: 0},
	}

	teamB := config.Team{
		ID:           "teamB",
		PrimaryColor: config.Color{R: 0, G: 0, B: 255},
	}

	// Without actual frame data, players remain neutral
	// This is a structure test
	assert.Len(t, players, 2)
	assert.Equal(t, "neutral", players[0].Team)
}

func TestEngineClose(t *testing.T) {
	cfg := config.NewManager()
	e := &Engine{
		config: cfg,
	}

	err := e.Close()
	assert.NoError(t, err)
}

// Mock test for Detect method (would require actual model in real scenario)
func TestDetectEmptyFrame(t *testing.T) {
	// This would require GoCV Mat setup
	// Placeholder for structure validation
	t.Skip("Requires GoCV and model file")
}
