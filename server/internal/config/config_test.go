package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()

	assert.NotNil(t, manager)
	config := manager.GetConfig()

	// Check default values
	assert.Equal(t, "teamA", config.TeamA.ID)
	assert.Equal(t, "teamB", config.TeamB.ID)
	assert.Equal(t, float32(0.5), config.Detection.ConfidenceThreshold)
	assert.Equal(t, 640, config.Detection.InputSize)
	assert.Equal(t, 2.0, config.Possession.ProximityThreshold)
}

func TestUpdateTeams(t *testing.T) {
	manager := NewManager()
	manager.path = "test_config.json"
	defer os.Remove("test_config.json")

	teamA := Team{
		ID:           "teamA",
		Name:         "Red Team",
		PrimaryColor: Color{R: 200, G: 0, B: 0},
	}

	teamB := Team{
		ID:           "teamB",
		Name:         "Blue Team",
		PrimaryColor: Color{R: 0, G: 0, B: 200},
	}

	err := manager.UpdateTeams(teamA, teamB)
	require.NoError(t, err)

	config := manager.GetConfig()
	assert.Equal(t, "Red Team", config.TeamA.Name)
	assert.Equal(t, "Blue Team", config.TeamB.Name)
	assert.Equal(t, uint8(200), config.TeamA.PrimaryColor.R)
}

func TestUpdateVideo(t *testing.T) {
	manager := NewManager()
	manager.path = "test_config.json"
	defer os.Remove("test_config.json")

	video := VideoConfig{
		Source: "webcam",
		Path:   "0",
		FPS:    30.0,
	}

	err := manager.UpdateVideo(video)
	require.NoError(t, err)

	config := manager.GetConfig()
	assert.Equal(t, "webcam", config.Video.Source)
	assert.Equal(t, "0", config.Video.Path)
	assert.Equal(t, 30.0, config.Video.FPS)
}

func TestUpdateDetection(t *testing.T) {
	manager := NewManager()
	manager.path = "test_config.json"
	defer os.Remove("test_config.json")

	detection := DetectionConfig{
		ConfidenceThreshold: 0.7,
		NMSThreshold:        0.5,
		InputSize:           1280,
	}

	err := manager.UpdateDetection(detection)
	require.NoError(t, err)

	config := manager.GetConfig()
	assert.Equal(t, float32(0.7), config.Detection.ConfidenceThreshold)
	assert.Equal(t, 1280, config.Detection.InputSize)
}

func TestLoadSave(t *testing.T) {
	manager := NewManager()
	manager.path = "test_config.json"
	defer os.Remove("test_config.json")

	// Update config
	teamA := Team{
		ID:           "teamA",
		Name:         "Test Team A",
		PrimaryColor: Color{R: 100, G: 100, B: 100},
	}

	err := manager.UpdateTeams(teamA, manager.config.TeamB)
	require.NoError(t, err)

	// Create new manager and load
	manager2 := NewManager()
	manager2.path = "test_config.json"
	err = manager2.Load()
	require.NoError(t, err)

	config := manager2.GetConfig()
	assert.Equal(t, "Test Team A", config.TeamA.Name)
	assert.Equal(t, uint8(100), config.TeamA.PrimaryColor.R)
}

func TestLoadNonExistentFile(t *testing.T) {
	manager := NewManager()
	manager.path = "nonexistent.json"

	err := manager.Load()
	assert.NoError(t, err) // Should not error, use defaults
}

func TestConcurrentAccess(t *testing.T) {
	manager := NewManager()
	manager.path = "test_config.json"
	defer os.Remove("test_config.json")

	done := make(chan bool)

	// Concurrent readers
	for i := 0; i < 10; i++ {
		go func() {
			_ = manager.GetConfig()
			done <- true
		}()
	}

	// Concurrent writers
	for i := 0; i < 5; i++ {
		go func() {
			video := VideoConfig{
				Source: "file",
				Path:   "/test/path",
				FPS:    15.0,
			}
			_ = manager.UpdateVideo(video)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 15; i++ {
		<-done
	}
}
