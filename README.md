# Probe - 网络抓包工具

一个用Go语言开发的网络抓包工具，类似于小黄鸟(HttpCanary)，用于捕获和分析网络流量。

## 功能特性

- 🔍 **实时网络流量捕获** - 捕获HTTP/HTTPS流量
- 📊 **Web界面展示** - 现代化的Web界面，实时显示抓包数据
- 🔎 **高级过滤功能** - 支持按IP、端口、协议、HTTP方法等过滤
- 📈 **统计信息** - 实时显示数据包统计信息
- 💾 **数据导出** - 支持JSON格式导出抓包数据
- 🎯 **HTTP协议解析** - 自动解析HTTP请求和响应
- 🔄 **实时更新** - 自动刷新数据，无需手动操作

## 系统要求

- Go 1.21 或更高版本
- macOS/Linux/Windows
- 管理员权限（用于网络接口访问）

## 安装和运行

### 1. 克隆项目
```bash
git clone <repository-url>
cd Probe
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 查看可用网络接口
```bash
# macOS
ifconfig

# Linux
ip addr show

# Windows
ipconfig
```

### 4. 运行程序
```bash
# 基本用法
go run main.go -i <网络接口名称>

# 示例
go run main.go -i en0

# 指定Web服务器端口
go run main.go -i en0 -p 8080

# 详细输出
go run main.go -i en0 -v
```

### 5. 访问Web界面
打开浏览器访问: http://localhost:8080

## 使用方法

### 命令行参数

- `-i`: 网络接口名称（必需）
- `-p`: Web服务器端口（默认: 8080）
- `-v`: 详细输出模式

### Web界面功能

1. **实时监控**: 页面会自动刷新显示最新的网络流量
2. **过滤功能**: 
   - 按源IP地址过滤
   - 按目标IP地址过滤
   - 按端口号过滤
   - 按协议类型过滤
   - 按HTTP方法过滤
   - 按文本内容搜索
3. **统计信息**: 显示总数据包数、HTTP/HTTPS包数、唯一IP数等
4. **数据管理**: 支持清空数据和导出数据

### 数据包信息

每个数据包包含以下信息：
- 时间戳
- 源IP地址和端口
- 目标IP地址和端口
- 协议类型
- 数据包长度
- HTTP方法（如果适用）
- HTTP URL（如果适用）
- HTTP状态码（如果适用）
- User-Agent（如果适用）
- Content-Type（如果适用）

## 技术架构

### 核心组件

1. **抓包器 (Capturer)**: 使用gopacket库捕获网络数据包
2. **存储模块 (Storage)**: 内存存储抓包数据
3. **Web服务器 (Server)**: 提供REST API和Web界面
4. **协议解析**: 解析HTTP/HTTPS协议

### 依赖库

- `github.com/google/gopacket`: 网络数据包捕获和解析
- `github.com/gorilla/mux`: HTTP路由
- `github.com/gorilla/websocket`: WebSocket支持（待实现）
- `github.com/sirupsen/logrus`: 日志记录

## 注意事项

1. **权限要求**: 程序需要管理员权限才能访问网络接口
2. **性能考虑**: 大量网络流量可能影响性能，建议在测试环境中使用
3. **数据安全**: 程序会捕获网络数据，请注意数据安全
4. **法律合规**: 请确保在合法范围内使用此工具

## 故障排除

### 常见问题

1. **权限不足**: 
   ```bash
   # 使用sudo运行
   sudo go run main.go -i en0
   
   # 或使用智能启动脚本
   ./start.sh
   ```

2. **网络接口不存在**: 检查接口名称是否正确
3. **端口被占用**: 更换Web服务器端口
4. **无法捕获数据**: 检查网络接口是否活跃

### 权限问题详解

在macOS上，网络抓包需要管理员权限，因为需要访问BPF (Berkeley Packet Filter) 设备。

**错误信息**: `Permission denied` 或 `cannot open BPF device`

**解决方案**:
1. 使用 `sudo` 运行程序
2. 使用提供的启动脚本 `./start.sh`
3. 检查网络接口是否活跃

### 调试模式

使用 `-v` 参数启用详细输出：
```bash
go run main.go -i en0 -v
```

## 开发计划

- [ ] WebSocket实时推送
- [ ] HTTPS流量解密
- [ ] 更多协议支持
- [ ] 数据包重放功能
- [ ] 更高级的过滤条件
- [ ] 数据包详情查看
- [ ] 性能优化

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 免责声明

本工具仅用于学习和研究目的。使用者需要遵守当地法律法规，不得用于非法用途。开发者不承担任何法律责任。
