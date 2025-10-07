#!/bin/bash

# Probe 网络抓包工具启动脚本

echo "🔍 Probe - 网络抓包工具"
echo "========================"

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go，请先安装Go 1.21或更高版本"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ 错误: Go版本过低，需要1.21或更高版本，当前版本: $GO_VERSION"
    exit 1
fi

echo "✅ Go版本检查通过: $GO_VERSION"

# 检查网络接口
echo ""
echo "📡 可用的网络接口:"
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    ifconfig | grep -E "^[a-z]" | grep -v lo0 | awk '{print $1}' | sed 's/://'
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    ip link show | grep -E "^[0-9]+:" | awk '{print $2}' | sed 's/://'
else
    echo "❌ 不支持的操作系统: $OSTYPE"
    exit 1
fi

echo ""
read -p "请输入网络接口名称 (例如: en0, eth0): " INTERFACE

if [ -z "$INTERFACE" ]; then
    echo "❌ 错误: 必须指定网络接口名称"
    exit 1
fi

# 检查接口是否存在
if [[ "$OSTYPE" == "darwin"* ]]; then
    if ! ifconfig "$INTERFACE" &> /dev/null; then
        echo "❌ 错误: 网络接口 '$INTERFACE' 不存在"
        exit 1
    fi
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if ! ip link show "$INTERFACE" &> /dev/null; then
        echo "❌ 错误: 网络接口 '$INTERFACE' 不存在"
        exit 1
    fi
fi

echo "✅ 网络接口检查通过: $INTERFACE"

# 询问端口
read -p "请输入Web服务器端口 (默认: 8080): " PORT
if [ -z "$PORT" ]; then
    PORT="8080"
fi

# 检查端口是否被占用
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  警告: 端口 $PORT 已被占用"
    read -p "是否继续? (y/N): " CONTINUE
    if [[ ! "$CONTINUE" =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""
echo "🚀 启动Probe抓包工具..."
echo "网络接口: $INTERFACE"
echo "Web端口: $PORT"
echo "Web界面: http://localhost:$PORT"
echo ""
echo "按 Ctrl+C 停止程序"
echo ""

# 检查是否需要sudo权限
echo "🔐 检查权限..."
if ! go run main.go -i "$INTERFACE" -p "$PORT" -v 2>/dev/null; then
    echo "⚠️  需要管理员权限来访问网络接口"
    echo "🔑 请输入您的密码以继续..."
    sudo go run main.go -i "$INTERFACE" -p "$PORT" -v
else
    go run main.go -i "$INTERFACE" -p "$PORT" -v
fi
