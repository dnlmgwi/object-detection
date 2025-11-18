#!/bin/bash

# Script to download YOLOv8 ONNX model
# Usage: ./scripts/download-model.sh [model-size]
# Model sizes: n (nano), s (small), m (medium), l (large), x (xlarge)

set -e

MODEL_SIZE="${1:-n}"
MODEL_DIR="server/models"
MODEL_NAME="yolov8${MODEL_SIZE}.onnx"
MODEL_PATH="${MODEL_DIR}/${MODEL_NAME}"

echo "🚀 Downloading YOLOv8${MODEL_SIZE} model..."

# Create models directory if it doesn't exist
mkdir -p "$MODEL_DIR"

# Check if model already exists
if [ -f "$MODEL_PATH" ]; then
    read -p "Model already exists. Download again? (y/N): " confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        echo "✅ Using existing model: $MODEL_PATH"
        exit 0
    fi
fi

# Download model using Python (ultralytics package)
echo "📦 Installing ultralytics (if not already installed)..."
pip3 install -q ultralytics

echo "⬇️  Downloading model..."
python3 << EOF
from ultralytics import YOLO

# Load model (will download if not exists)
model = YOLO('yolov8${MODEL_SIZE}.pt')

# Export to ONNX
model.export(format='onnx', imgsz=640)

print("✅ Model exported to ONNX format")
EOF

# Move model to correct location
if [ -f "yolov8${MODEL_SIZE}.onnx" ]; then
    mv "yolov8${MODEL_SIZE}.onnx" "$MODEL_PATH"
    echo "✅ Model saved to: $MODEL_PATH"
else
    echo "❌ Failed to download/export model"
    exit 1
fi

# Display model info
echo ""
echo "📊 Model Information:"
echo "   Name: YOLOv8${MODEL_SIZE}"
echo "   Path: $MODEL_PATH"
echo "   Size: $(du -h "$MODEL_PATH" | cut -f1)"
echo ""
echo "🎉 Download complete! You can now run the application."
