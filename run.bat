@echo off
chcp 65001 >nul

echo 🔍 Probe - 网络抓包工具
echo ========================

REM 检查Go是否安装
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到Go，请先安装Go 1.21或更高版本
    pause
    exit /b 1
)

echo ✅ Go版本检查通过

REM 显示网络接口
echo.
echo 📡 可用的网络接口:
netsh interface show interface

echo.
set /p INTERFACE="请输入网络接口名称 (例如: Ethernet, Wi-Fi): "

if "%INTERFACE%"=="" (
    echo ❌ 错误: 必须指定网络接口名称
    pause
    exit /b 1
)

echo ✅ 网络接口: %INTERFACE%

REM 询问端口
set /p PORT="请输入Web服务器端口 (默认: 8080): "
if "%PORT%"=="" set PORT=8080

echo.
echo 🚀 启动Probe抓包工具...
echo 网络接口: %INTERFACE%
echo Web端口: %PORT%
echo Web界面: http://localhost:%PORT%
echo.
echo 按 Ctrl+C 停止程序
echo.

REM 运行程序
go run main.go -i "%INTERFACE%" -p "%PORT%" -v

pause
