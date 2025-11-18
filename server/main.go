package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dnlmgwi/football-object-detection/internal/api"
	"github.com/dnlmgwi/football-object-detection/internal/config"
	"github.com/dnlmgwi/football-object-detection/internal/detection"
	"github.com/dnlmgwi/football-object-detection/internal/possession"
	"github.com/dnlmgwi/football-object-detection/internal/websocket"
	"github.com/gorilla/mux"
)

var (
	port      = flag.String("port", "8080", "Server port")
	modelPath = flag.String("model", "./models/yolov8n.onnx", "Path to YOLO model")
)

func main() {
	flag.Parse()

	// Initialize configuration
	cfg := config.NewManager()

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize detection engine
	detector, err := detection.NewEngine(*modelPath, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize detection engine: %v", err)
	}
	defer detector.Close()

	// Initialize possession tracker
	tracker := possession.NewTracker(cfg)

	// Initialize API server
	router := mux.NewRouter()
	apiServer := api.NewServer(cfg, detector, tracker, hub)
	apiServer.RegisterRoutes(router)

	// Serve static files for frontend
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	// Start HTTP server
	addr := ":" + *port
	log.Printf("Starting server on %s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	hub.Shutdown()
	detector.Close()
}
