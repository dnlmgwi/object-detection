.PHONY: help setup build test run clean docker

help:
	@echo "Football Object Detection - Makefile Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make setup          - Install all dependencies and download model"
	@echo "  make setup-backend  - Setup Go backend only"
	@echo "  make setup-frontend - Setup React frontend only"
	@echo "  make download-model - Download YOLOv8n model"
	@echo ""
	@echo "Development:"
	@echo "  make dev            - Run both backend and frontend in dev mode"
	@echo "  make dev-backend    - Run backend in dev mode"
	@echo "  make dev-frontend   - Run frontend in dev mode"
	@echo ""
	@echo "Build:"
	@echo "  make build          - Build both backend and frontend"
	@echo "  make build-backend  - Build Go backend"
	@echo "  make build-frontend - Build React frontend"
	@echo ""
	@echo "Test:"
	@echo "  make test           - Run all tests"
	@echo "  make test-backend   - Run Go tests"
	@echo "  make test-frontend  - Run React tests"
	@echo "  make coverage       - Generate test coverage report"
	@echo ""
	@echo "Run:"
	@echo "  make run            - Run production build"
	@echo "  make run-backend    - Run backend server"
	@echo ""
	@echo "Clean:"
	@echo "  make clean          - Remove build artifacts"
	@echo ""

# Setup
setup: setup-backend setup-frontend download-model
	@echo "✅ Setup complete!"

setup-backend:
	@echo "📦 Setting up Go backend..."
	cd server && go mod download
	cd server && go mod tidy

setup-frontend:
	@echo "📦 Setting up React frontend..."
	cd frontend && npm install

download-model:
	@echo "⬇️  Downloading YOLOv8n model..."
	@mkdir -p server/models
	@if [ ! -f server/models/yolov8n.onnx ]; then \
		echo "Please download the YOLOv8n ONNX model from:"; \
		echo "https://github.com/ultralytics/assets/releases"; \
		echo "And place it in server/models/yolov8n.onnx"; \
	else \
		echo "✅ Model already exists"; \
	fi

# Build
build: build-backend build-frontend
	@echo "✅ Build complete!"

build-backend:
	@echo "🔨 Building Go backend..."
	cd server && mkdir -p bin
	cd server && go build -o bin/server main.go

build-frontend:
	@echo "🔨 Building React frontend..."
	cd frontend && npm run build
	@echo "📦 Copying frontend build to server/static..."
	@mkdir -p server/static
	@cp -r frontend/build/client/* server/static/

# Test
test: test-backend test-frontend
	@echo "✅ All tests passed!"

test-backend:
	@echo "🧪 Running Go tests..."
	cd server && go test ./... -v

test-frontend:
	@echo "🧪 Running React tests..."
	cd frontend && npm test

coverage:
	@echo "📊 Generating coverage report..."
	cd server && go test ./... -coverprofile=coverage.out
	cd server && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: server/coverage.html"

# Development
dev:
	@echo "🚀 Starting development servers..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:5173"
	@make -j 2 dev-backend dev-frontend

dev-backend:
	@echo "🚀 Starting backend in dev mode..."
	cd server && go run main.go --model=./models/yolov8n.onnx

dev-frontend:
	@echo "🚀 Starting frontend in dev mode..."
	cd frontend && npm run dev

# Run
run: build
	@echo "🚀 Starting production server..."
	cd server && ./bin/server --model=./models/yolov8n.onnx

run-backend:
	@echo "🚀 Starting backend..."
	cd server && ./bin/server --model=./models/yolov8n.onnx

# Clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf server/bin
	rm -rf server/static
	rm -rf server/uploads
	rm -rf server/coverage.out
	rm -rf server/coverage.html
	rm -rf frontend/build
	rm -rf frontend/dist
	rm -rf frontend/.react-router
	@echo "✅ Clean complete!"

# Docker
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t football-detection .

docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 -v $(PWD)/server/models:/app/models football-detection

# Utility
fmt:
	@echo "✨ Formatting code..."
	cd server && go fmt ./...
	cd frontend && npm run format

lint:
	@echo "🔍 Linting code..."
	cd server && go vet ./...
	cd frontend && npm run lint
