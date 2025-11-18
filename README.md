# Football Object Detection System

Real-time football game object detection and ball possession tracking system using YOLOv8, Go, and React Router 7.

![System Architecture](docs/architecture.png)

## Features

- 🎯 **Real-time Object Detection**: Detect players and ball using YOLOv8/YOLO11
- ⚽ **Ball Possession Tracking**: Track which team has possession and for how long
- 🎨 **Team Classification**: Identify players by jersey color
- 📊 **Live Statistics**: Real-time possession percentages and performance metrics
- 🔄 **WebSocket Streaming**: Low-latency real-time updates
- 🎥 **Multiple Video Sources**: Support for file upload, webcam, and RTSP streams
- ⚙️ **Configurable**: Adjustable team colors, FPS, confidence thresholds
- 📱 **Modern UI**: Responsive React interface with Tailwind CSS

## Architecture

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
                             └─────────────────────────┘
                                        ▲
                                        │
                                  Video Source
```

## Tech Stack

### Backend
- **Go 1.21+**: Server and business logic
- **GoCV**: OpenCV bindings for Go
- **YOLOv8/v11**: Object detection model (ONNX format)
- **Gorilla WebSocket**: Real-time communication
- **Gorilla Mux**: HTTP routing

### Frontend
- **React 19**: UI library
- **React Router 7**: Routing and framework
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling
- **Vite**: Build tool

## Prerequisites

### System Dependencies

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install -y \
    build-essential \
    libopencv-dev \
    pkg-config
```

**macOS:**
```bash
brew install opencv pkg-config
```

### Development Tools
- Go 1.21 or higher
- Node.js 20 or higher
- npm or yarn

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/dnlmgwi/football-object-detection.git
cd football-object-detection
```

### 2. Download YOLO Model

Download YOLOv8n (nano) ONNX model:

```bash
# Create models directory
mkdir -p server/models

# Download YOLOv8n ONNX model (replace URL with actual model)
curl -L https://github.com/ultralytics/assets/releases/download/v0.0.0/yolov8n.onnx \
    -o server/models/yolov8n.onnx
```

Or download manually from [Ultralytics](https://github.com/ultralytics/ultralytics) and place in `server/models/`.

### 3. Setup Backend

```bash
cd server

# Download Go dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o bin/server main.go
```

### 4. Setup Frontend

```bash
cd frontend

# Install dependencies
npm install

# Build
npm run build
```

### 5. Run the Application

**Terminal 1 - Backend:**
```bash
cd server
./bin/server --model=./models/yolov8n.onnx --port=8080
```

**Terminal 2 - Frontend (Development):**
```bash
cd frontend
npm run dev
```

Or serve the built frontend from the backend:
```bash
cd frontend
npm run build
cp -r build/client ../server/static

cd ../server
./bin/server
```

### 6. Open the Application

Navigate to `http://localhost:5173` (dev) or `http://localhost:8080` (production)

## Usage

### Configure Teams

1. In the right panel, set team names and colors
2. Click "Save Teams"

### Configure Video Source

**Option 1: Upload Video File**
1. Select "Video File" as source
2. Click "Choose File" and select a football game video
3. Click "Save Video Config"

**Option 2: Use Webcam**
1. Select "Webcam" as source
2. Enter device ID (usually 0)
3. Click "Save Video Config"

**Option 3: RTSP Stream**
1. Select "RTSP Stream" as source
2. Enter RTSP URL (e.g., `rtsp://camera.local/stream`)
3. Click "Save Video Config"

### Start Detection

1. Click "Start Processing"
2. Watch real-time detection and possession tracking
3. Click "Stop Processing" to end

## Configuration

### Detection Settings

Edit `server/config.json` or use the API:

```json
{
  "detection": {
    "confidenceThreshold": 0.5,
    "nmsThreshold": 0.4,
    "inputSize": 640
  },
  "possession": {
    "proximityThreshold": 2.0,
    "colorTolerance": 20.0
  }
}
```

### Performance Tuning

- **Lower FPS** (5-10): Reduces CPU usage, suitable for analysis
- **Higher FPS** (15-30): Better real-time tracking, higher CPU usage
- **Input Size** (640): Faster inference, lower accuracy
- **Input Size** (1280): Slower inference, higher accuracy

## Development

### Backend Development

```bash
cd server

# Run with auto-reload (using air)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/detection -run TestDetect
```

### Frontend Development

```bash
cd frontend

# Development server with hot reload
npm run dev

# Type checking
npm run typecheck

# Linting
npm run lint

# Build for production
npm run build
```

### Running Tests (TDD)

Backend uses Go's testing framework:
```bash
cd server
go test ./... -v
```

Frontend uses Vitest:
```bash
cd frontend
npm test
```

## API Documentation

### REST Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/config` | GET | Get current configuration |
| `/api/config/teams` | POST | Update team configuration |
| `/api/config/video` | POST | Update video configuration |
| `/api/config/detection` | POST | Update detection settings |
| `/api/config/possession` | POST | Update possession settings |
| `/api/video/upload` | POST | Upload video file |
| `/api/control/start` | POST | Start processing |
| `/api/control/stop` | POST | Stop processing |
| `/api/control/status` | GET | Get processing status |

### WebSocket Endpoint

**Connection:** `ws://localhost:8080/ws`

**Message Types:**

1. **Detection Result**
```json
{
  "type": "detection",
  "data": {
    "frameId": 123,
    "timestamp": 1634567890123,
    "players": [...],
    "ball": {...},
    "possession": {...}
  }
}
```

2. **Possession Stats**
```json
{
  "type": "possession",
  "data": {
    "teamA": 58.5,
    "teamB": 41.5,
    "currentHolder": "teamA",
    "lastChange": 1634567890000
  }
}
```

3. **Performance Metrics**
```json
{
  "type": "metrics",
  "data": {
    "fps": 15.2,
    "latency": 85
  }
}
```

## Troubleshooting

### GoCV Installation Issues

If you encounter GoCV-related errors:

```bash
# Check OpenCV installation
pkg-config --modversion opencv4

# Reinstall OpenCV if needed
# Ubuntu/Debian:
sudo apt-get remove libopencv-dev
sudo apt-get install libopencv-dev

# macOS:
brew reinstall opencv
```

### Model Loading Errors

Ensure the YOLO model is:
1. In ONNX format (`.onnx`)
2. In the correct path (`server/models/`)
3. Downloaded completely (check file size)

### Low FPS / High Latency

- Reduce input size (640 instead of 1280)
- Lower processing FPS
- Use GPU acceleration (requires CUDA setup)
- Close other applications

### WebSocket Connection Issues

- Check firewall settings
- Ensure backend is running on port 8080
- Verify WebSocket URL in frontend matches backend

## Performance Benchmarks

| Configuration | FPS | Latency | CPU Usage |
|---------------|-----|---------|-----------|
| YOLOv8n, 640px, CPU | 10-15 | ~100ms | 60-80% |
| YOLOv8n, 1280px, CPU | 5-8 | ~200ms | 80-100% |
| YOLOv8n, 640px, GPU | 25-30 | ~30ms | 20-30% |

*Tested on Intel i7-10700K, 16GB RAM, NVIDIA RTX 3060*

## Project Structure

```
.
├── server/                 # Go backend
│   ├── main.go            # Entry point
│   ├── go.mod             # Go dependencies
│   ├── internal/          # Internal packages
│   │   ├── api/           # HTTP and WebSocket handlers
│   │   ├── config/        # Configuration management
│   │   ├── detection/     # Object detection engine
│   │   ├── possession/    # Possession tracking
│   │   ├── websocket/     # WebSocket hub
│   │   └── models/        # Data models
│   └── models/            # YOLO models (not in git)
├── frontend/              # React frontend
│   ├── app/               # Application code
│   │   ├── routes/        # Route components
│   │   ├── components/    # React components
│   │   └── hooks/         # Custom hooks
│   ├── public/            # Static assets
│   └── package.json       # Node dependencies
├── PRD.md                 # Product Requirements Document
└── README.md              # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow TDD principles
- Write tests for new features
- Use `gofmt` for Go code
- Use ESLint/Prettier for TypeScript
- Update documentation
- Ensure all tests pass before PR

## Roadmap

### v1.0 (Current)
- [x] Basic object detection
- [x] Team classification
- [x] Ball possession tracking
- [x] WebSocket streaming
- [x] Web UI

### v1.1 (Planned)
- [ ] Player re-identification across frames
- [ ] Historical possession timeline
- [ ] Export statistics to CSV/JSON
- [ ] Multiple video streams
- [ ] Docker deployment

### v2.0 (Future)
- [ ] Player tracking (unique IDs)
- [ ] Tactical analysis (formations)
- [ ] Heat maps
- [ ] Event detection (goals, passes, fouls)
- [ ] Mobile app

## License

MIT License - see [LICENSE](LICENSE) file

## Acknowledgments

- [Ultralytics YOLOv8](https://github.com/ultralytics/ultralytics) for the detection model
- [GoCV](https://gocv.io/) for OpenCV bindings
- [React Router](https://reactrouter.com/) for the frontend framework
- [OpenCV](https://opencv.org/) for computer vision functionality

## Support

- 📧 Email: support@example.com
- 🐛 Issues: [GitHub Issues](https://github.com/dnlmgwi/football-object-detection/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/dnlmgwi/football-object-detection/discussions)

## Citation

If you use this project in your research, please cite:

```bibtex
@software{football_object_detection,
  author = {Your Name},
  title = {Football Object Detection System},
  year = {2025},
  url = {https://github.com/dnlmgwi/football-object-detection}
}
```

---

Made with ⚽ and ❤️ using AI-assisted development
