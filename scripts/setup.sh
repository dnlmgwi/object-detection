#!/bin/bash

# Complete setup script for Football Object Detection system

set -e

echo "⚽ Football Object Detection - Setup Script"
echo "==========================================="
echo ""

# Check prerequisites
echo "🔍 Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ from https://golang.org/"
    exit 1
fi
echo "✅ Go $(go version | awk '{print $3}')"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 20+ from https://nodejs.org/"
    exit 1
fi
echo "✅ Node.js $(node --version)"

# Check npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed."
    exit 1
fi
echo "✅ npm $(npm --version)"

# Check OpenCV (optional for now)
if pkg-config --exists opencv4; then
    echo "✅ OpenCV $(pkg-config --modversion opencv4)"
else
    echo "⚠️  OpenCV not found. GoCV features may not work."
    echo "   Install with: sudo apt-get install libopencv-dev (Ubuntu/Debian)"
    echo "   Or: brew install opencv (macOS)"
fi

echo ""
echo "📦 Installing dependencies..."
echo ""

# Backend setup
echo "1️⃣  Setting up Go backend..."
cd server
go mod download
go mod tidy
cd ..
echo "✅ Backend dependencies installed"

# Frontend setup
echo ""
echo "2️⃣  Setting up React frontend..."
cd frontend
npm install
cd ..
echo "✅ Frontend dependencies installed"

# Create necessary directories
echo ""
echo "3️⃣  Creating directories..."
mkdir -p server/models
mkdir -p server/uploads
mkdir -p server/static
echo "✅ Directories created"

# Download model
echo ""
echo "4️⃣  Downloading YOLO model..."
if [ -f "scripts/download-model.sh" ]; then
    chmod +x scripts/download-model.sh
    ./scripts/download-model.sh n || echo "⚠️  Model download failed. Please download manually."
else
    echo "⚠️  Model download script not found."
    echo "   Please download YOLOv8n ONNX model manually and place in server/models/"
fi

# Run tests
echo ""
echo "5️⃣  Running tests..."
echo "Testing backend..."
cd server
go test ./... || echo "⚠️  Some backend tests failed"
cd ..

echo ""
echo "✅ Setup complete!"
echo ""
echo "🚀 Quick Start:"
echo "   1. make dev          - Start development servers"
echo "   2. make build        - Build for production"
echo "   3. make run          - Run production build"
echo ""
echo "📖 See README.md for detailed documentation"
echo ""
