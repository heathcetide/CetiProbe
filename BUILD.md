# Cetiprobe 跨平台构建指南

本文档详细说明如何为 macOS、Linux 和 Windows 平台构建 Cetiprobe 网络分析工具。

## 📋 目录

- [环境要求](#环境要求)
- [macOS 构建](#macos-构建)
- [Linux 构建](#linux-构建)
- [Windows 构建](#windows-构建)
- [构建产物](#构建产物)
- [分发说明](#分发说明)
- [故障排除](#故障排除)

## 🔧 环境要求

### 基础环境
- **Node.js**: >= 18.0.0
- **npm**: >= 8.0.0
- **Rust**: >= 1.70.0
- **Tauri CLI**: >= 2.0.0

### 平台特定要求

#### macOS
- Xcode Command Line Tools
- macOS >= 10.13

#### Linux
- `libwebkit2gtk-4.0-dev`
- `libssl-dev`
- `libgtk-3-dev`
- `libayatana-appindicator3-dev`
- `librsvg2-dev`

#### Windows
- Microsoft Visual Studio C++ Build Tools
- Windows SDK

## 🍎 macOS 构建

### 1. 安装依赖

```bash
# 安装 Xcode Command Line Tools
xcode-select --install

# 安装 Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# 安装 Tauri CLI
cargo install tauri-cli
```

### 2. 构建应用

```bash
# 克隆项目
git clone <repository-url>
cd cosmic-ui-main

# 安装前端依赖
npm install

# 构建前端
npm run build

# 构建 Tauri 应用
npm run tauri:build
```

### 3. 构建产物

构建完成后，您将获得：
- **应用程序包**: `src-tauri/target/release/bundle/macos/Cetiprobe.app`
- **DMG 安装包**: `src-tauri/target/release/bundle/dmg/Cetiprobe_0.0.0_aarch64.dmg`

## 🐧 Linux 构建

### 1. 安装系统依赖

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install libwebkit2gtk-4.0-dev \
    build-essential \
    curl \
    wget \
    libssl-dev \
    libgtk-3-dev \
    libayatana-appindicator3-dev \
    librsvg2-dev
```

#### Fedora/CentOS/RHEL
```bash
sudo dnf install webkit2gtk3-devel.x86_64 \
    openssl-devel \
    gtk3-devel \
    libappindicator-gtk3-devel \
    librsvg2-devel
```

#### Arch Linux
```bash
sudo pacman -S webkit2gtk \
    base-devel \
    curl \
    wget \
    openssl \
    appmenu-gtk-module \
    gtk3 \
    libappindicator-gtk3 \
    librsvg
```

### 2. 安装 Rust 和 Tauri

```bash
# 安装 Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# 安装 Tauri CLI
cargo install tauri-cli

# 添加 Linux 目标架构
rustup target add x86_64-unknown-linux-gnu
```

### 3. 构建应用

```bash
# 克隆项目
git clone <repository-url>
cd cosmic-ui-main

# 安装前端依赖
npm install

# 构建前端
npm run build

# 构建 Tauri 应用
npm run tauri:build
```

### 4. 构建产物

构建完成后，您将获得：
- **AppImage**: `src-tauri/target/release/bundle/appimage/Cetiprobe_0.0.0_x86_64.AppImage`
- **DEB 包**: `src-tauri/target/release/bundle/deb/cetiprobe_0.0.0_amd64.deb`
- **RPM 包**: `src-tauri/target/release/bundle/rpm/cetiprobe-0.0.0-1.x86_64.rpm`

## 🪟 Windows 构建

### 1. 安装系统依赖

#### 使用 Chocolatey (推荐)
```powershell
# 安装 Chocolatey
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 安装必要工具
choco install -y microsoft-visual-studio-build-tools
choco install -y nodejs
choco install -y rust
```

#### 手动安装
1. 下载并安装 [Microsoft Visual Studio Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/)
2. 下载并安装 [Node.js](https://nodejs.org/)
3. 下载并安装 [Rust](https://rustup.rs/)

### 2. 安装 Rust 和 Tauri

```powershell
# 安装 Rust (如果未安装)
# 从 https://rustup.rs/ 下载并运行 rustup-init.exe

# 安装 Tauri CLI
cargo install tauri-cli

# 添加 Windows 目标架构
rustup target add x86_64-pc-windows-gnu
```

### 3. 构建应用

```powershell
# 克隆项目
git clone <repository-url>
cd cosmic-ui-main

# 安装前端依赖
npm install

# 构建前端
npm run build

# 构建 Tauri 应用
npm run tauri:build
```

### 4. 构建产物

构建完成后，您将获得：
- **MSI 安装包**: `src-tauri/target/release/bundle/msi/Cetiprobe_0.0.0_x64_en-US.msi`
- **NSIS 安装包**: `src-tauri/target/release/bundle/nsis/Cetiprobe_0.0.0_x64-setup.exe`

## 🚀 一键构建脚本

### 创建构建脚本

#### `build-all.sh` (macOS/Linux)
```bash
#!/bin/bash

echo "🚀 开始构建 Cetiprobe..."

# 检查依赖
echo "📋 检查依赖..."
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装"
    exit 1
fi

if ! command -v cargo &> /dev/null; then
    echo "❌ Rust 未安装"
    exit 1
fi

# 安装前端依赖
echo "📦 安装前端依赖..."
npm install

# 构建前端
echo "🔨 构建前端..."
npm run build

# 构建 Tauri 应用
echo "🔨 构建 Tauri 应用..."
npm run tauri:build

echo "✅ 构建完成！"
echo "📁 构建产物位置："
echo "   - macOS: src-tauri/target/release/bundle/macos/"
echo "   - Linux: src-tauri/target/release/bundle/appimage/"
echo "   - Windows: src-tauri/target/release/bundle/msi/"
```

#### `build-all.ps1` (Windows)
```powershell
Write-Host "🚀 开始构建 Cetiprobe..." -ForegroundColor Green

# 检查依赖
Write-Host "📋 检查依赖..." -ForegroundColor Yellow
if (!(Get-Command node -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Node.js 未安装" -ForegroundColor Red
    exit 1
}

if (!(Get-Command cargo -ErrorAction SilentlyContinue)) {
    Write-Host "❌ Rust 未安装" -ForegroundColor Red
    exit 1
}

# 安装前端依赖
Write-Host "📦 安装前端依赖..." -ForegroundColor Yellow
npm install

# 构建前端
Write-Host "🔨 构建前端..." -ForegroundColor Yellow
npm run build

# 构建 Tauri 应用
Write-Host "🔨 构建 Tauri 应用..." -ForegroundColor Yellow
npm run tauri:build

Write-Host "✅ 构建完成！" -ForegroundColor Green
Write-Host "📁 构建产物位置：" -ForegroundColor Cyan
Write-Host "   - Windows: src-tauri/target/release/bundle/msi/" -ForegroundColor Cyan
```

## 📦 构建产物说明

### macOS
- **Cetiprobe.app**: 应用程序包，可直接拖拽到 Applications 文件夹
- **Cetiprobe_0.0.0_aarch64.dmg**: DMG 安装包，包含安装向导

### Linux
- **Cetiprobe_0.0.0_x86_64.AppImage**: 便携式应用程序，无需安装
- **cetiprobe_0.0.0_amd64.deb**: Debian/Ubuntu 安装包
- **cetiprobe-0.0.0-1.x86_64.rpm**: Red Hat/Fedora 安装包

### Windows
- **Cetiprobe_0.0.0_x64_en-US.msi**: Windows 安装包
- **Cetiprobe_0.0.0_x64-setup.exe**: NSIS 安装程序

## 🚀 分发说明

### 文件大小参考
- **macOS**: ~11MB (应用) / ~4MB (DMG)
- **Linux**: ~15MB (AppImage) / ~12MB (DEB/RPM)
- **Windows**: ~8MB (MSI) / ~6MB (NSIS)

### 分发建议
1. **GitHub Releases**: 上传所有平台的构建产物
2. **应用商店**: 考虑发布到 Mac App Store、Microsoft Store
3. **包管理器**: 发布到 Homebrew、Chocolatey、Snap

## 🔧 故障排除

### 常见问题

#### 1. 构建失败：缺少系统依赖
```bash
# Linux: 安装缺失的依赖
sudo apt install libwebkit2gtk-4.0-dev libssl-dev libgtk-3-dev

# macOS: 安装 Xcode Command Line Tools
xcode-select --install
```

#### 2. 图标文件错误
```bash
# 确保图标文件存在且格式正确
ls -la src-tauri/icons/
# 应该包含: 32x32.png, 128x128.png, 128x128@2x.png
```

#### 3. 前端构建失败
```bash
# 清理并重新安装依赖
rm -rf node_modules package-lock.json
npm install
npm run build
```

#### 4. Rust 编译错误
```bash
# 更新 Rust 工具链
rustup update
cargo clean
npm run tauri:build
```

### 调试模式构建

```bash
# 开发模式构建（更快，包含调试信息）
npm run tauri:dev

# 生产模式构建（优化，体积小）
npm run tauri:build
```

## 📚 相关资源

- [Tauri 官方文档](https://tauri.app/)
- [Rust 安装指南](https://rustup.rs/)
- [Node.js 下载](https://nodejs.org/)
- [跨平台构建最佳实践](https://tauri.app/v1/guides/building/)

---

## 🎯 快速开始

如果您只想快速构建当前平台的应用：

```bash
# 1. 安装依赖
npm install

# 2. 构建
npm run tauri:build

# 3. 查看结果
ls -la src-tauri/target/release/bundle/
```

构建完成后，您就可以将应用分发给用户了！🎉
