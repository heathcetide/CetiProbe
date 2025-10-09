#!/bin/bash

echo "🚀 Setting up Cosmic UI Desktop with Tauri..."

# 检查是否安装了 Rust
if ! command -v cargo &> /dev/null; then
    echo "❌ Rust is not installed. Please install Rust first:"
    echo "   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
    exit 1
fi

# 检查是否安装了 Tauri CLI
if ! command -v tauri &> /dev/null; then
    echo "📦 Installing Tauri CLI..."
    cargo install tauri-cli
fi

# 安装前端依赖
echo "📦 Installing frontend dependencies..."
npm install

# 安装 Tauri 依赖
echo "📦 Installing Tauri dependencies..."
npm install @tauri-apps/cli @tauri-apps/api

echo "✅ Setup complete!"
echo ""
echo "🎯 Available commands:"
echo "   npm run tauri:dev    - Start development server"
echo "   npm run tauri:build  - Build desktop application"
echo ""
echo "🚀 Run 'npm run tauri:dev' to start the desktop app!"
