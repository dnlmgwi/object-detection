# Product Requirements Document (PRD)
## Real-Time Football Object Detection & Possession Tracking System

**Version:** 1.0
**Date:** 2025-11-18
**Status:** Draft

---

## 1. Executive Summary

### 1.1 Overview
A real-time football game analysis system that uses computer vision to detect players from two teams and track ball possession from live video feeds. The system provides configurable team identification and low-latency real-time tracking through a client-server architecture.

### 1.2 Goals
- Enable real-time object detection of football players and ball from live video feeds
- Track and calculate ball possession statistics for two teams
- Provide configurable team identification (colors, names)
- Achieve low-latency processing (<100ms per frame)
- Deliver insights through an intuitive web interface

### 1.3 Technology Stack
- **Backend:** Go (Golang) with GoCV for video processing
- **Object Detection:** YOLOv8/YOLO11 (ONNX format) via Ultralytics
- **Frontend:** React Router 7
- **Real-time Communication:** WebSocket
- **Development Approach:** Test-Driven Development (TDD)

---

## 2. Product Scope

### 2.1 In Scope
- Video feed ingestion (file upload, webcam, RTSP stream)
- Player detection and team classification
- Ball detection and tracking
- Possession calculation based on ball proximity to players
- Team configuration (colors, names, jersey detection)
- Real-time WebSocket streaming of detection results
- Web-based UI for configuration and visualization
- Performance metrics and latency monitoring

### 2.2 Out of Scope (v1.0)
- Multi-camera support
- Player identity tracking across frames
- Tactical analysis (formations, heat maps)
- Video recording/storage
- Mobile applications
- Authentication/multi-user support

---

## 3. User Stories

### 3.1 Core User Stories

**US-001: Configure Teams**
- As a user, I want to configure two teams with custom names and jersey colors, so that the system can differentiate between them during detection.

**US-002: Upload Video Feed**
- As a user, I want to upload a video file or connect to a live stream, so that the system can analyze the football game.

**US-003: View Real-Time Detection**
- As a user, I want to see detected players and ball overlaid on the video in real-time, so that I can verify detection accuracy.

**US-004: Track Possession**
- As a user, I want to see live possession statistics (percentage per team), so that I can analyze game dynamics.

**US-005: Monitor Performance**
- As a user, I want to see processing latency metrics, so that I can ensure the system meets real-time requirements.

---

## 4. Functional Requirements

### 4.1 Object Detection Engine

**FR-001: YOLO Model Integration**
- System SHALL use YOLOv8 or YOLO11 ONNX model for object detection
- System SHALL detect persons (football players) with confidence threshold ≥ 0.5
- System SHALL detect sports balls with confidence threshold ≥ 0.4
- System SHALL process frames at minimum 10 FPS for real-time analysis

**FR-002: Team Classification**
- System SHALL classify detected players into Team A or Team B based on jersey color
- System SHALL use HSV color space for robust color matching
- System SHALL handle color tolerance configuration (±20° hue range)
- System SHALL label unclassified players as "neutral" (referee, etc.)

**FR-003: Ball Detection**
- System SHALL track the ball position across frames
- System SHALL apply Kalman filtering for smooth ball trajectory
- System SHALL handle ball occlusion (predict position when not visible)

### 4.2 Possession Tracking

**FR-004: Possession Calculation**
- System SHALL determine possession based on proximity (player within 2m of ball)
- System SHALL update possession statistics every second
- System SHALL calculate cumulative possession percentage per team
- System SHALL detect possession changes and log timestamps

**FR-005: Possession Events**
- System SHALL emit events when possession changes teams
- System SHALL track possession duration per possession period
- System SHALL count total possession instances per team

### 4.3 Configuration Management

**FR-006: Team Configuration**
- System SHALL accept team names (max 50 characters)
- System SHALL accept primary jersey color (RGB/HEX format)
- System SHALL accept secondary jersey color (optional)
- System SHALL persist configuration across sessions

**FR-007: Detection Configuration**
- System SHALL allow confidence threshold adjustment (0.0-1.0)
- System SHALL allow possession proximity threshold adjustment (0.5m-5m)
- System SHALL allow frame processing rate configuration (1-30 FPS)

### 4.4 Video Input Management

**FR-008: Video Sources**
- System SHALL accept video file uploads (MP4, AVI, MOV)
- System SHALL accept webcam input (device ID)
- System SHALL accept RTSP stream URLs
- System SHALL validate video source before processing

**FR-009: Video Processing**
- System SHALL decode video frames using GoCV
- System SHALL resize frames to optimal detection size (640x640 or 1280x1280)
- System SHALL maintain aspect ratio during resize

### 4.5 Real-Time Communication

**FR-010: WebSocket Streaming**
- System SHALL establish WebSocket connection for real-time data
- System SHALL stream detection results at configured FPS
- System SHALL stream possession statistics every second
- System SHALL handle client reconnection gracefully

**FR-011: Data Format**
- Detection results SHALL include: frame_id, timestamp, players[], ball, possession
- Player objects SHALL include: bbox, team_id, confidence, color_match
- Ball object SHALL include: bbox, confidence, velocity (optional)

### 4.6 Frontend Interface

**FR-012: Video Display**
- UI SHALL render video frames with detection overlays
- UI SHALL draw bounding boxes with team-specific colors
- UI SHALL display player count per team on frame
- UI SHALL display ball position indicator

**FR-013: Statistics Dashboard**
- UI SHALL display real-time possession gauge (Team A vs Team B)
- UI SHALL show cumulative possession percentages
- UI SHALL show current possession holder
- UI SHALL display processing FPS and latency

**FR-014: Configuration Panel**
- UI SHALL provide form for team configuration
- UI SHALL provide video source selector
- UI SHALL provide advanced settings panel (thresholds, FPS)
- UI SHALL validate all inputs before submission

---

## 5. Non-Functional Requirements

### 5.1 Performance

**NFR-001: Latency**
- End-to-end latency (video input to UI display) SHALL be ≤ 200ms for 1080p video
- Frame processing time SHALL be ≤ 100ms per frame on modern CPU
- Frame processing time SHALL be ≤ 30ms per frame on GPU (CUDA enabled)

**NFR-002: Throughput**
- System SHALL process minimum 10 FPS continuously
- System SHALL handle video resolutions up to 1920x1080
- WebSocket SHALL support up to 10 concurrent clients

**NFR-003: Resource Usage**
- Backend SHALL use ≤ 4GB RAM during processing
- Backend SHALL use ≤ 2 CPU cores (without GPU)
- Frontend bundle size SHALL be ≤ 2MB (gzipped)

### 5.2 Scalability

**NFR-004: Horizontal Scaling**
- System architecture SHALL support multiple video streams (future)
- WebSocket server SHALL support connection pooling
- Detection service SHALL be stateless for horizontal scaling

### 5.3 Reliability

**NFR-005: Error Handling**
- System SHALL handle video stream interruptions gracefully
- System SHALL retry failed WebSocket connections (max 3 attempts)
- System SHALL log all errors with stack traces

**NFR-006: Availability**
- System SHALL recover from crashes automatically (process supervisor)
- System SHALL validate all inputs to prevent crashes

### 5.4 Usability

**NFR-007: User Experience**
- UI SHALL be responsive (desktop and tablet)
- Configuration changes SHALL apply within 2 seconds
- Error messages SHALL be user-friendly and actionable

### 5.5 Maintainability

**NFR-008: Code Quality**
- Backend code SHALL have ≥ 80% test coverage
- Frontend code SHALL have ≥ 70% test coverage
- Code SHALL follow language-specific style guides (gofmt, ESLint)

**NFR-009: Documentation**
- All public APIs SHALL have inline documentation
- README SHALL include setup and usage instructions
- Architecture diagrams SHALL be provided

---

## 6. System Architecture

### 6.1 High-Level Architecture

```
┌─────────────────┐          ┌─────────────────────────┐
│                 │          │                         │
│  React Frontend │◄────────►│   Go Backend Server     │
│  (React Router) │ WebSocket│                         │
│                 │          │  ┌──────────────────┐   │
└─────────────────┘          │  │ Detection Engine │   │
                             │  │   (GoCV + YOLO)  │   │
                             │  └──────────────────┘   │
                             │  ┌──────────────────┐   │
                             │  │ Possession       │   │
                             │  │ Tracker          │   │
                             │  └──────────────────┘   │
                             │  ┌──────────────────┐   │
                             │  │ Video Processor  │   │
                             │  └──────────────────┘   │
                             └─────────────────────────┘
                                        ▲
                                        │
                                  Video Source
                              (File/Webcam/RTSP)
```

### 6.2 Backend Components

**6.2.1 HTTP Server**
- Serves frontend static files
- Handles configuration API endpoints
- Handles video upload endpoint

**6.2.2 WebSocket Server**
- Manages client connections
- Broadcasts detection results
- Handles client subscriptions

**6.2.3 Video Processor**
- Captures frames from video source
- Manages frame queue
- Handles frame rate control

**6.2.4 Detection Engine**
- Loads YOLO model (ONNX)
- Runs inference on frames
- Filters detections by confidence
- Classifies players by team

**6.2.5 Possession Tracker**
- Calculates ball-player distances
- Determines possession holder
- Tracks possession statistics
- Emits possession events

**6.2.6 Configuration Manager**
- Stores team configurations
- Provides configuration to other components
- Persists configuration to file/database

### 6.3 Frontend Components

**6.3.1 Video Player Component**
- Renders video frames from WebSocket
- Overlays bounding boxes
- Displays team colors

**6.3.2 Statistics Dashboard**
- Shows possession gauge
- Displays real-time stats
- Shows performance metrics

**6.3.3 Configuration Panel**
- Team configuration form
- Video source selector
- Advanced settings

**6.3.4 WebSocket Client**
- Manages WebSocket connection
- Handles reconnection logic
- Dispatches events to components

### 6.4 Data Flow

1. User configures teams and video source via frontend
2. Frontend sends configuration to backend via HTTP POST
3. Backend starts video processing pipeline
4. Video Processor captures frames and queues them
5. Detection Engine processes frames and detects objects
6. Possession Tracker calculates possession from detections
7. Results are broadcast to all WebSocket clients
8. Frontend receives results and updates UI

---

## 7. API Specifications

### 7.1 REST API

**POST /api/config/teams**
```json
Request:
{
  "teamA": {
    "name": "Red Team",
    "primaryColor": "#FF0000",
    "secondaryColor": "#FFFFFF"
  },
  "teamB": {
    "name": "Blue Team",
    "primaryColor": "#0000FF",
    "secondaryColor": "#FFFFFF"
  }
}

Response:
{
  "success": true,
  "message": "Teams configured successfully"
}
```

**POST /api/config/video**
```json
Request:
{
  "source": "file",  // "file" | "webcam" | "rtsp"
  "path": "/path/to/video.mp4",  // or device ID or RTSP URL
  "fps": 15
}

Response:
{
  "success": true,
  "message": "Video source configured"
}
```

**POST /api/video/upload**
```
Content-Type: multipart/form-data
File: video.mp4

Response:
{
  "success": true,
  "filePath": "/uploads/video_12345.mp4"
}
```

**POST /api/control/start**
```json
Response:
{
  "success": true,
  "sessionId": "session_12345"
}
```

**POST /api/control/stop**
```json
Response:
{
  "success": true
}
```

**GET /api/status**
```json
Response:
{
  "running": true,
  "fps": 14.5,
  "latency": 87,
  "framesProcessed": 1234
}
```

### 7.2 WebSocket API

**Connection:** `ws://localhost:8080/ws`

**Message: Detection Result**
```json
{
  "type": "detection",
  "data": {
    "frameId": 1234,
    "timestamp": 1634567890123,
    "players": [
      {
        "id": "p1",
        "bbox": {"x": 100, "y": 150, "width": 50, "height": 120},
        "team": "teamA",
        "confidence": 0.92
      }
    ],
    "ball": {
      "bbox": {"x": 320, "y": 240, "width": 20, "height": 20},
      "confidence": 0.87
    },
    "possession": {
      "team": "teamA",
      "playerId": "p1",
      "duration": 5.2
    }
  }
}
```

**Message: Possession Update**
```json
{
  "type": "possession",
  "data": {
    "teamA": 58.5,
    "teamB": 41.5,
    "currentHolder": "teamA",
    "lastChange": 1634567885000
  }
}
```

**Message: Performance Metrics**
```json
{
  "type": "metrics",
  "data": {
    "fps": 14.8,
    "latency": 92,
    "cpuUsage": 45.2,
    "memoryUsage": 1024
  }
}
```

---

## 8. Data Models

### 8.1 Team Configuration
```go
type Team struct {
    ID            string `json:"id"`
    Name          string `json:"name"`
    PrimaryColor  Color  `json:"primaryColor"`
    SecondaryColor Color `json:"secondaryColor"`
}

type Color struct {
    R uint8 `json:"r"`
    G uint8 `json:"g"`
    B uint8 `json:"b"`
}
```

### 8.2 Detection Result
```go
type DetectionResult struct {
    FrameID    int64      `json:"frameId"`
    Timestamp  int64      `json:"timestamp"`
    Players    []Player   `json:"players"`
    Ball       *Ball      `json:"ball"`
    Possession *Possession `json:"possession"`
}

type Player struct {
    ID         string      `json:"id"`
    BBox       BoundingBox `json:"bbox"`
    Team       string      `json:"team"`
    Confidence float32     `json:"confidence"`
}

type Ball struct {
    BBox       BoundingBox `json:"bbox"`
    Confidence float32     `json:"confidence"`
    Velocity   *Vector2D   `json:"velocity,omitempty"`
}

type BoundingBox struct {
    X      int `json:"x"`
    Y      int `json:"y"`
    Width  int `json:"width"`
    Height int `json:"height"`
}
```

### 8.3 Possession Statistics
```go
type PossessionStats struct {
    TeamAPercentage float64 `json:"teamA"`
    TeamBPercentage float64 `json:"teamB"`
    CurrentHolder   string  `json:"currentHolder"`
    LastChange      int64   `json:"lastChange"`
    TotalDuration   float64 `json:"totalDuration"`
}
```

---

## 9. Testing Strategy

### 9.1 Unit Testing

**Backend (Go)**
- Test each package in isolation
- Mock external dependencies (GoCV, file I/O)
- Target: ≥80% coverage
- Tools: `testing`, `testify`, `mockery`

**Frontend (React)**
- Test components with React Testing Library
- Test hooks and utilities
- Target: ≥70% coverage
- Tools: Vitest, React Testing Library

### 9.2 Integration Testing

**Backend**
- Test HTTP endpoints with `httptest`
- Test WebSocket communication
- Test video processing pipeline with sample video

**Frontend**
- Test user flows (configure → start → view results)
- Test WebSocket integration with mock server

### 9.3 End-to-End Testing

- Test complete workflow with sample football video
- Verify possession tracking accuracy
- Measure latency and performance

### 9.4 Performance Testing

- Benchmark frame processing time
- Measure WebSocket throughput
- Profile memory usage under load

---

## 10. Development Phases

### Phase 1: Foundation (Week 1)
- [x] Set up Go project structure
- [ ] Set up React Router 7 project
- [ ] Integrate GoCV and YOLO model
- [ ] Basic video frame capture and processing
- [ ] Unit tests for core detection logic

### Phase 2: Detection Engine (Week 2)
- [ ] Implement YOLO inference
- [ ] Implement team color classification
- [ ] Implement ball detection and tracking
- [ ] Unit tests for detection engine
- [ ] Integration tests with sample video

### Phase 3: Possession Tracking (Week 3)
- [ ] Implement possession calculation algorithm
- [ ] Implement possession statistics tracking
- [ ] Implement possession events
- [ ] Unit tests for possession tracker

### Phase 4: Backend API (Week 4)
- [ ] Implement REST API endpoints
- [ ] Implement WebSocket server
- [ ] Implement configuration management
- [ ] API integration tests

### Phase 5: Frontend (Week 5)
- [ ] Build configuration panel
- [ ] Build video player with overlays
- [ ] Build statistics dashboard
- [ ] Implement WebSocket client
- [ ] Component tests

### Phase 6: Integration & Optimization (Week 6)
- [ ] End-to-end testing
- [ ] Performance optimization (latency reduction)
- [ ] UI/UX improvements
- [ ] Documentation

---

## 11. Success Metrics

### 11.1 Technical Metrics
- **Detection Accuracy:** ≥85% player detection rate
- **Possession Accuracy:** ≥90% possession determination accuracy
- **Latency:** ≤200ms end-to-end
- **Frame Rate:** ≥10 FPS sustained
- **Test Coverage:** ≥80% backend, ≥70% frontend

### 11.2 Quality Metrics
- Zero critical bugs in production
- All API endpoints documented
- Setup time ≤15 minutes for new developers

---

## 12. Dependencies

### 12.1 Backend Dependencies
```
- Go 1.21+
- GoCV 0.35+
- OpenCV 4.8+
- gorilla/websocket
- YOLOv8/v11 ONNX model (download from Ultralytics)
```

### 12.2 Frontend Dependencies
```
- Node.js 20+
- React 19
- React Router 7
- WebSocket client library
- Tailwind CSS (optional, for styling)
```

### 12.3 System Dependencies
```
- OpenCV libraries (libopencv-dev)
- Optional: CUDA for GPU acceleration
```

---

## 13. Risks & Mitigation

### 13.1 Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| YOLO model accuracy insufficient for football | High | Medium | Use pre-trained sports-specific model; fine-tune if needed |
| GoCV performance bottleneck | High | Low | Profile early; optimize hot paths; use GPU |
| WebSocket scalability issues | Medium | Low | Implement connection pooling; load test early |
| Color-based team detection fails | High | Medium | Add jersey number recognition fallback; allow manual correction |

### 13.2 Project Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Scope creep | Medium | High | Strict adherence to MVP scope; defer nice-to-haves |
| Dependency version conflicts | Low | Medium | Use dependency management tools (go.mod, package-lock.json) |
| Insufficient test coverage | Medium | Medium | Enforce coverage gates in CI; TDD approach |

---

## 14. Open Questions

1. **Q:** Should we support multiple simultaneous video streams?
   **A:** Not in v1.0, but architecture should allow for future expansion.

2. **Q:** How should we handle occlusion (players blocking each other)?
   **A:** YOLO handles partial occlusion; severe cases may result in missed detections (acceptable for v1.0).

3. **Q:** Should possession tracking account for ball in the air?
   **A:** Yes, track 3D position if possible; otherwise, use last known ground position.

4. **Q:** What video formats should be supported?
   **A:** MP4, AVI, MOV for uploads; RTSP for live streams.

---

## 15. Future Enhancements (Post-v1.0)

- Player identity tracking across frames (re-identification)
- Tactical analysis (formations, passing networks)
- Multi-camera support with camera calibration
- Video recording and playback with annotations
- Mobile app (iOS/Android)
- Advanced metrics (distance covered, sprint detection)
- Integration with existing sports analytics platforms
- AI-powered highlights generation

---

## 16. Glossary

- **Bounding Box:** Rectangular box around detected object (x, y, width, height)
- **Confidence Threshold:** Minimum probability for accepting a detection
- **FPS:** Frames Per Second
- **GoCV:** Go bindings for OpenCV
- **HSV:** Hue, Saturation, Value color space
- **Latency:** Time delay from input to output
- **ONNX:** Open Neural Network Exchange format
- **RTSP:** Real-Time Streaming Protocol
- **TDD:** Test-Driven Development
- **WebSocket:** Protocol for bidirectional real-time communication
- **YOLO:** You Only Look Once (object detection algorithm)

---

## 17. Approval & Sign-off

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Product Owner | - | - | - |
| Tech Lead | - | - | - |
| QA Lead | - | - | - |

---

**Document Version History:**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-18 | Claude | Initial draft |
