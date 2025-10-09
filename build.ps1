# Cetiprobe 快速构建脚本 (Windows PowerShell)
# 支持 Windows 平台构建

Write-Host "🚀 Cetiprobe 跨平台构建脚本" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green

# 检查必要工具
Write-Host "📋 检查构建环境..." -ForegroundColor Yellow

function Test-Command {
    param($Command)
    if (!(Get-Command $Command -ErrorAction SilentlyContinue)) {
        Write-Host "❌ $Command 未安装，请先安装 $Command" -ForegroundColor Red
        exit 1
    } else {
        Write-Host "✅ $Command 已安装" -ForegroundColor Green
    }
}

Test-Command "node"
Test-Command "npm"
Test-Command "cargo"

# 检查 Tauri CLI
try {
    cargo tauri --version | Out-Null
    Write-Host "✅ Tauri CLI 已安装" -ForegroundColor Green
} catch {
    Write-Host "📦 安装 Tauri CLI..." -ForegroundColor Yellow
    cargo install tauri-cli
}

Write-Host "✅ 环境检查完成" -ForegroundColor Green

# 安装依赖
Write-Host "📦 安装项目依赖..." -ForegroundColor Yellow
npm install

# 构建前端
Write-Host "🔨 构建前端应用..." -ForegroundColor Yellow
npm run build

# 构建 Tauri 应用
Write-Host "🔨 构建 Tauri 桌面应用..." -ForegroundColor Yellow
npm run tauri:build

Write-Host "✅ 构建完成！" -ForegroundColor Green
Write-Host ""
Write-Host "📁 构建产物位置：" -ForegroundColor Cyan
Write-Host "   🪟 Windows:" -ForegroundColor Cyan
Write-Host "      - MSI 安装包: src-tauri/target/release/bundle/msi/" -ForegroundColor White
Write-Host "      - NSIS 安装包: src-tauri/target/release/bundle/nsis/" -ForegroundColor White

Write-Host ""
Write-Host "🎉 构建成功！您现在可以分发应用程序了。" -ForegroundColor Green
Write-Host "📖 详细说明请查看 BUILD.md 文档" -ForegroundColor Cyan
