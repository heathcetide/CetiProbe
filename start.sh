#!/bin/bash

# Probe 智能启动脚本 - 自动处理权限问题

echo "🔍 Probe - 网络抓包工具"
echo "========================"

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go，请先安装Go 1.21或更高版本"
    exit 1
fi

# 检查网络接口
echo ""
echo "📡 可用的网络接口:"
ifconfig | grep -E "^[a-z]" | grep -v lo0 | awk '{print $1}' | sed 's/://'

echo ""
read -p "请输入网络接口名称 (例如: en0, en1): " INTERFACE

if [ -z "$INTERFACE" ]; then
    echo "❌ 错误: 必须指定网络接口名称"
    exit 1
fi

# 询问端口
read -p "请输入Web服务器端口 (默认: 8080): " PORT
if [ -z "$PORT" ]; then
    PORT="8080"
fi

echo ""
echo "🚀 启动Probe抓包工具..."
echo "网络接口: $INTERFACE"
echo "Web端口: $PORT"
echo "Web界面: http://localhost:$PORT"
echo ""

# 首先尝试普通权限运行
echo "🔐 尝试以普通权限运行..."
if timeout 3s go run main.go -i "$INTERFACE" -p "$PORT" -v 2>/dev/null; then
    echo "✅ 普通权限运行成功"
    go run main.go -i "$INTERFACE" -p "$PORT" -v
else
    echo "⚠️  普通权限失败，需要管理员权限"
    echo "🔑 请输入您的密码以继续..."
    sudo go run main.go -i "$INTERFACE" -p "$PORT" -v
fi
