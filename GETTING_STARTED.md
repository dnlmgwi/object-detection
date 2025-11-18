# Getting Started - YOLO Model Required

The football object detection system requires a YOLOv8 ONNX model to run.

## Quick Setup

### Option 1: Download Pre-converted Model (Recommended for Testing)

Due to environment constraints, manually download a YOLOv8n ONNX model:

1. Visit: https://github.com/ultralytics/assets/releases
2. Download `yolov8n.onnx` (or any YOLOv8 variant)
3. Place it in: `server/models/yolov8n.onnx`

### Option 2: Using Python Script (If ultralytics is available)

```bash
# Install ultralytics
pip3 install ultralytics

# Run the download script
./scripts/download-model.sh n
```

### Option 3: Export from PyTorch Yourself

```python
from ultralytics import YOLO

# Load model
model = YOLO('yolov8n.pt')

# Export to ONNX
model.export(format='onnx', imgsz=640)

# Move to: server/models/yolov8n.onnx
```

## Model Sizes

- `yolov8n.onnx` - Nano (~6MB) - Fastest, good for CPU
- `yolov8s.onnx` - Small (~22MB) - Better accuracy
- `yolov8m.onnx` - Medium (~50MB) - Balanced
- `yolov8l.onnx` - Large (~88MB) - Higher accuracy
- `yolov8x.onnx` - Extra Large (~136MB) - Best accuracy

## After Getting the Model

Once you have the model in `server/models/yolov8n.onnx`:

```bash
# Test the backend
cd server
go run main.go --model=./models/yolov8n.onnx

# Or use make
cd ..
make dev
```

## Current Status

- ✅ OpenCV 4.6.0 installed
- ✅ Go backend code ready
- ✅ React frontend ready
- ⏳ YOLO model needed (download manually for now)

## For Testing Without a Model

The system will fail to start without a model, but you can:

1. Run frontend only: `cd frontend && npm run dev`
2. Test backend compilation: `cd server && go build main.go`
3. Run unit tests: `cd server && go test ./...` (tests don't require the model)
