package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"

	"github.com/dnlmgwi/football-object-detection/internal/config"
	"github.com/dnlmgwi/football-object-detection/internal/detection"
	"github.com/dnlmgwi/football-object-detection/internal/models"
	"github.com/dnlmgwi/football-object-detection/internal/possession"
	ws "github.com/dnlmgwi/football-object-detection/internal/websocket"
)

// Server represents the API server
type Server struct {
	config   *config.Manager
	detector *detection.Engine
	tracker  *possession.Tracker
	hub      *ws.Hub

	mu            sync.RWMutex
	running       bool
	stopChan      chan struct{}
	videoCapture  *gocv.VideoCapture
	frameID       int64
	startTime     time.Time
	framesProcessed int64
}

// NewServer creates a new API server
func NewServer(cfg *config.Manager, detector *detection.Engine, tracker *possession.Tracker, hub *ws.Hub) *Server {
	return &Server{
		config:   cfg,
		detector: detector,
		tracker:  tracker,
		hub:      hub,
		running:  false,
		stopChan: make(chan struct{}),
	}
}

// RegisterRoutes registers all API routes
func (s *Server) RegisterRoutes(router *mux.Router) {
	api := router.PathPrefix("/api").Subrouter()

	// Configuration endpoints
	api.HandleFunc("/config/teams", s.handleUpdateTeams).Methods("POST")
	api.HandleFunc("/config/video", s.handleUpdateVideo).Methods("POST")
	api.HandleFunc("/config/detection", s.handleUpdateDetection).Methods("POST")
	api.HandleFunc("/config/possession", s.handleUpdatePossession).Methods("POST")
	api.HandleFunc("/config", s.handleGetConfig).Methods("GET")

	// Video upload endpoint
	api.HandleFunc("/video/upload", s.handleVideoUpload).Methods("POST")

	// Control endpoints
	api.HandleFunc("/control/start", s.handleStart).Methods("POST")
	api.HandleFunc("/control/stop", s.handleStop).Methods("POST")
	api.HandleFunc("/control/status", s.handleStatus).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", s.handleWebSocket)

	// Enable CORS
	router.Use(corsMiddleware)
}

// handleUpdateTeams handles team configuration updates
func (s *Server) handleUpdateTeams(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TeamA config.Team `json:"teamA"`
		TeamB config.Team `json:"teamB"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.config.UpdateTeams(req.TeamA, req.TeamB); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update teams")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Teams configured successfully",
	})
}

// handleUpdateVideo handles video configuration updates
func (s *Server) handleUpdateVideo(w http.ResponseWriter, r *http.Request) {
	var video config.VideoConfig

	if err := json.NewDecoder(r.Body).Decode(&video); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.config.UpdateVideo(video); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update video config")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Video source configured",
	})
}

// handleUpdateDetection handles detection configuration updates
func (s *Server) handleUpdateDetection(w http.ResponseWriter, r *http.Request) {
	var detection config.DetectionConfig

	if err := json.NewDecoder(r.Body).Decode(&detection); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.config.UpdateDetection(detection); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update detection config")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Detection config updated",
	})
}

// handleUpdatePossession handles possession configuration updates
func (s *Server) handleUpdatePossession(w http.ResponseWriter, r *http.Request) {
	var poss config.PossessionConfig

	if err := json.NewDecoder(r.Body).Decode(&poss); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.config.UpdatePossession(poss); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update possession config")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Possession config updated",
	})
}

// handleGetConfig returns current configuration
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := s.config.GetConfig()
	respondJSON(w, http.StatusOK, cfg)
}

// handleVideoUpload handles video file uploads
func (s *Server) handleVideoUpload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 500MB)
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	// Create uploads directory
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, 0755)

	// Generate unique filename
	filename := fmt.Sprintf("video_%d%s", time.Now().Unix(), filepath.Ext(header.Filename))
	filepath := filepath.Join(uploadDir, filename)

	// Create file
	dst, err := os.Create(filepath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"filePath": filepath,
		"filename": filename,
	})
}

// handleStart starts video processing
func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		respondError(w, http.StatusBadRequest, "Already running")
		return
	}

	cfg := s.config.GetConfig()
	videoPath := cfg.Video.Path
	if videoPath == "" {
		s.mu.Unlock()
		respondError(w, http.StatusBadRequest, "No video source configured")
		return
	}

	// Open video capture
	var capture *gocv.VideoCapture
	var err error

	if cfg.Video.Source == "webcam" {
		// Open webcam by device ID
		deviceID := 0
		fmt.Sscanf(videoPath, "%d", &deviceID)
		capture, err = gocv.OpenVideoCapture(deviceID)
	} else {
		// Open video file or RTSP stream
		capture, err = gocv.OpenVideoCapture(videoPath)
	}

	if err != nil {
		s.mu.Unlock()
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to open video: %v", err))
		return
	}

	s.videoCapture = capture
	s.running = true
	s.frameID = 0
	s.framesProcessed = 0
	s.startTime = time.Now()
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	// Reset tracker
	s.tracker.Reset()

	// Start processing in goroutine
	go s.processVideo()

	sessionID := fmt.Sprintf("session_%d", time.Now().Unix())
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"sessionId": sessionID,
		"message":   "Processing started",
	})
}

// handleStop stops video processing
func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		respondError(w, http.StatusBadRequest, "Not running")
		return
	}

	close(s.stopChan)
	s.running = false

	if s.videoCapture != nil {
		s.videoCapture.Close()
		s.videoCapture = nil
	}
	s.mu.Unlock()

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Processing stopped",
	})
}

// handleStatus returns current status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	running := s.running
	framesProcessed := s.framesProcessed
	startTime := s.startTime
	s.mu.RUnlock()

	fps := 0.0
	latency := int64(0)

	if running && framesProcessed > 0 {
		elapsed := time.Since(startTime).Seconds()
		if elapsed > 0 {
			fps = float64(framesProcessed) / elapsed
		}
		// Estimate latency (simplified)
		latency = int64(1000.0 / fps)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"running":         running,
		"fps":             fps,
		"latency":         latency,
		"framesProcessed": framesProcessed,
	})
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	ws.ServeWS(s.hub, conn)
}

// processVideo processes video frames
func (s *Server) processVideo() {
	cfg := s.config.GetConfig()
	targetFPS := cfg.Video.FPS
	frameDuration := time.Duration(1000/targetFPS) * time.Millisecond

	frame := gocv.NewMat()
	defer frame.Close()

	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	lastPossessionBroadcast := time.Now()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			// Read frame
			if ok := s.videoCapture.Read(&frame); !ok || frame.Empty() {
				log.Println("End of video or failed to read frame")
				s.handleStop(nil, nil)
				return
			}

			startTime := time.Now()

			// Detect objects
			players, ball, err := s.detector.Detect(frame)
			if err != nil {
				log.Printf("Detection error: %v", err)
				continue
			}

			// Update possession
			possession := s.tracker.Update(players, ball)

			// Increment frame counters
			s.mu.Lock()
			s.frameID++
			s.framesProcessed++
			frameID := s.frameID
			s.mu.Unlock()

			// Broadcast detection result
			result := models.DetectionResult{
				FrameID:    frameID,
				Timestamp:  time.Now().UnixMilli(),
				Players:    players,
				Ball:       ball,
				Possession: possession,
			}

			if err := s.hub.BroadcastDetection(result); err != nil {
				log.Printf("Broadcast error: %v", err)
			}

			// Broadcast possession stats every second
			if time.Since(lastPossessionBroadcast) >= time.Second {
				stats := s.tracker.GetStats()
				if err := s.hub.BroadcastPossession(stats); err != nil {
					log.Printf("Possession broadcast error: %v", err)
				}
				lastPossessionBroadcast = time.Now()

				// Also broadcast metrics
				latency := time.Since(startTime).Milliseconds()
				elapsed := time.Since(s.startTime).Seconds()
				fps := 0.0
				if elapsed > 0 {
					fps = float64(s.framesProcessed) / elapsed
				}

				metrics := models.PerformanceMetrics{
					FPS:     fps,
					Latency: latency,
				}
				s.hub.BroadcastMetrics(metrics)
			}
		}
	}
}

// Helper functions

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
