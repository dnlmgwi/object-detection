#!/bin/bash

# Install OpenCV and dependencies for Ubuntu/Debian
# This script installs OpenCV 4.x which is required for GoCV

set -e

echo "📦 Installing OpenCV and dependencies..."
echo "This may take 5-10 minutes..."
echo ""

# Update package list
echo "Updating package list..."
apt-get update -qq

# Install OpenCV and build dependencies
echo "Installing OpenCV..."
apt-get install -y -qq \
    libopencv-dev \
    pkg-config \
    build-essential \
    cmake

# Verify installation
echo ""
echo "✅ Verifying installation..."
if pkg-config --modversion opencv4 > /dev/null 2>&1; then
    echo "✅ OpenCV $(pkg-config --modversion opencv4) installed successfully"
elif pkg-config --modversion opencv > /dev/null 2>&1; then
    echo "✅ OpenCV $(pkg-config --modversion opencv) installed successfully"
else
    echo "❌ OpenCV installation verification failed"
    exit 1
fi

echo ""
echo "🎉 OpenCV installation complete!"
echo ""
echo "You can now run:"
echo "  make dev       - Start development servers"
echo "  make build     - Build the application"
