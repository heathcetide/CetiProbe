#!/bin/bash

# Cetiprobe 快速构建脚本
# 支持 macOS, Linux, Windows 跨平台构建

set -e

echo "🚀 Cetiprobe 跨平台构建脚本"
echo "================================"

# 检查操作系统
OS="unknown"
if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
    OS="windows"
fi

echo "🖥️  检测到操作系统: $OS"

# 检查必要工具
echo "📋 检查构建环境..."

check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 未安装，请先安装 $1"
        exit 1
    else
        echo "✅ $1 已安装"
    fi
}

check_command node
check_command npm
check_command cargo

# 检查 Tauri CLI
if ! cargo tauri --version &> /dev/null; then
    echo "📦 安装 Tauri CLI..."
    cargo install tauri-cli
fi

echo "✅ 环境检查完成"

# 安装依赖
echo "📦 安装项目依赖..."
npm install

# 构建前端
echo "🔨 构建前端应用..."
npm run build

# 构建 Tauri 应用
echo "🔨 构建 Tauri 桌面应用..."
npm run tauri:build

echo "✅ 构建完成！"
echo ""
echo "📁 构建产物位置："

case $OS in
    "macos")
        echo "   🍎 macOS:"
        echo "      - 应用程序: src-tauri/target/release/bundle/macos/Cetiprobe.app"
        echo "      - DMG 安装包: src-tauri/target/release/bundle/dmg/"
        ;;
    "linux")
        echo "   🐧 Linux:"
        echo "      - AppImage: src-tauri/target/release/bundle/appimage/"
        echo "      - DEB 包: src-tauri/target/release/bundle/deb/"
        echo "      - RPM 包: src-tauri/target/release/bundle/rpm/"
        ;;
    "windows")
        echo "   🪟 Windows:"
        echo "      - MSI 安装包: src-tauri/target/release/bundle/msi/"
        echo "      - NSIS 安装包: src-tauri/target/release/bundle/nsis/"
        ;;
esac

echo ""
echo "🎉 构建成功！您现在可以分发应用程序了。"
echo "📖 详细说明请查看 BUILD.md 文档"
