package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnlmgwi/football-object-detection/internal/config"
	"github.com/dnlmgwi/football-object-detection/internal/possession"
	ws "github.com/dnlmgwi/football-object-detection/internal/websocket"
)

func TestNewServer(t *testing.T) {
	cfg := config.NewManager()
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()

	server := NewServer(cfg, nil, tracker, hub)

	assert.NotNil(t, server)
	assert.False(t, server.running)
}

func TestHandleGetConfig(t *testing.T) {
	cfg := config.NewManager()
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()
	server := NewServer(cfg, nil, tracker, hub)

	router := mux.NewRouter()
	server.RegisterRoutes(router)

	req, err := http.NewRequest("GET", "/api/config", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response config.Config
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "teamA", response.TeamA.ID)
	assert.Equal(t, "teamB", response.TeamB.ID)
}

func TestHandleUpdateTeams(t *testing.T) {
	cfg := config.NewManager()
	cfg.path = "test_config.json"
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()
	server := NewServer(cfg, nil, tracker, hub)

	router := mux.NewRouter()
	server.RegisterRoutes(router)

	reqBody := map[string]interface{}{
		"teamA": map[string]interface{}{
			"id":   "teamA",
			"name": "Red Team",
			"primaryColor": map[string]uint8{
				"r": 255,
				"g": 0,
				"b": 0,
			},
		},
		"teamB": map[string]interface{}{
			"id":   "teamB",
			"name": "Blue Team",
			"primaryColor": map[string]uint8{
				"r": 0,
				"g": 0,
				"b": 255,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/config/teams", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
}

func TestHandleUpdateVideo(t *testing.T) {
	cfg := config.NewManager()
	cfg.path = "test_config.json"
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()
	server := NewServer(cfg, nil, tracker, hub)

	router := mux.NewRouter()
	server.RegisterRoutes(router)

	reqBody := map[string]interface{}{
		"source": "file",
		"path":   "/path/to/video.mp4",
		"fps":    30.0,
	}

	body, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/config/video", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	// Verify config was updated
	updatedCfg := cfg.GetConfig()
	assert.Equal(t, "file", updatedCfg.Video.Source)
	assert.Equal(t, 30.0, updatedCfg.Video.FPS)
}

func TestHandleStatus(t *testing.T) {
	cfg := config.NewManager()
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()
	server := NewServer(cfg, nil, tracker, hub)

	router := mux.NewRouter()
	server.RegisterRoutes(router)

	req, err := http.NewRequest("GET", "/api/control/status", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["running"].(bool))
	assert.Equal(t, float64(0), response["fps"].(float64))
}

func TestHandleStopWhenNotRunning(t *testing.T) {
	cfg := config.NewManager()
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()
	server := NewServer(cfg, nil, tracker, hub)

	router := mux.NewRouter()
	server.RegisterRoutes(router)

	req, err := http.NewRequest("POST", "/api/control/stop", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandleStartWithoutVideo(t *testing.T) {
	cfg := config.NewManager()
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()
	server := NewServer(cfg, nil, tracker, hub)

	router := mux.NewRouter()
	server.RegisterRoutes(router)

	req, err := http.NewRequest("POST", "/api/control/start", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCORSMiddleware(t *testing.T) {
	cfg := config.NewManager()
	tracker := possession.NewTracker(cfg)
	hub := ws.NewHub()
	server := NewServer(cfg, nil, tracker, hub)

	router := mux.NewRouter()
	server.RegisterRoutes(router)

	req, err := http.NewRequest("OPTIONS", "/api/config", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}
