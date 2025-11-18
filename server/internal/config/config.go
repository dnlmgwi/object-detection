package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Color represents an RGB color
type Color struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

// Team represents a football team configuration
type Team struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PrimaryColor   Color  `json:"primaryColor"`
	SecondaryColor Color  `json:"secondaryColor"`
}

// VideoConfig represents video source configuration
type VideoConfig struct {
	Source string  `json:"source"` // "file", "webcam", "rtsp"
	Path   string  `json:"path"`   // file path, device ID, or RTSP URL
	FPS    float64 `json:"fps"`
}

// DetectionConfig represents detection engine configuration
type DetectionConfig struct {
	ConfidenceThreshold float32 `json:"confidenceThreshold"`
	NMSThreshold        float32 `json:"nmsThreshold"`
	InputSize           int     `json:"inputSize"`
}

// PossessionConfig represents possession tracking configuration
type PossessionConfig struct {
	ProximityThreshold float64 `json:"proximityThreshold"` // meters
	ColorTolerance     float64 `json:"colorTolerance"`     // HSV hue tolerance
}

// Config holds all application configuration
type Config struct {
	TeamA      Team             `json:"teamA"`
	TeamB      Team             `json:"teamB"`
	Video      VideoConfig      `json:"video"`
	Detection  DetectionConfig  `json:"detection"`
	Possession PossessionConfig `json:"possession"`
}

// Manager manages application configuration
type Manager struct {
	mu     sync.RWMutex
	config Config
	path   string
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		config: Config{
			TeamA: Team{
				ID:             "teamA",
				Name:           "Team A",
				PrimaryColor:   Color{R: 255, G: 0, B: 0},
				SecondaryColor: Color{R: 255, G: 255, B: 255},
			},
			TeamB: Team{
				ID:             "teamB",
				Name:           "Team B",
				PrimaryColor:   Color{R: 0, G: 0, B: 255},
				SecondaryColor: Color{R: 255, G: 255, B: 255},
			},
			Video: VideoConfig{
				Source: "file",
				Path:   "",
				FPS:    15.0,
			},
			Detection: DetectionConfig{
				ConfidenceThreshold: 0.5,
				NMSThreshold:        0.4,
				InputSize:           640,
			},
			Possession: PossessionConfig{
				ProximityThreshold: 2.0,
				ColorTolerance:     20.0,
			},
		},
		path: "config.json",
	}
}

// GetConfig returns a copy of the current configuration
func (m *Manager) GetConfig() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// UpdateTeams updates team configurations
func (m *Manager) UpdateTeams(teamA, teamB Team) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.TeamA = teamA
	m.config.TeamB = teamB

	return m.save()
}

// UpdateVideo updates video configuration
func (m *Manager) UpdateVideo(video VideoConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Video = video

	return m.save()
}

// UpdateDetection updates detection configuration
func (m *Manager) UpdateDetection(detection DetectionConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Detection = detection

	return m.save()
}

// UpdatePossession updates possession configuration
func (m *Manager) UpdatePossession(possession PossessionConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Possession = possession

	return m.save()
}

// Load loads configuration from file
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Use default config
		}
		return err
	}

	return json.Unmarshal(data, &m.config)
}

// save persists configuration to file (must be called with lock held)
func (m *Manager) save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0644)
}
