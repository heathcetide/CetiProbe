# CetiProbe Server 使用文档（Windows / macOS / Linux）

> 版本：Alpha | 端口：Web `:8080`，代理默认 `:8899`

## 目录
- [1. 概览](#1-概览)
- [2. 快速开始](#2-快速开始)
- [3. 三种系统的环境准备](#3-三种系统的环境准备)
- [4. 启动与访问](#4-启动与访问)
- [5. Web UI 使用指南](#5-web-ui-使用指南)
- [6. PCAP 抓包](#6-pcap-抓包)
- [7. MITM 代理（HTTP/HTTPS）](#7-mitm-代理httphttps)
- [8. 证书导入（分系统）](#8-证书导入分系统)
- [9. 常用 API 参考](#9-常用-api-参考)
- [10. 故障排查](#10-故障排查)
- [11. 安全与合规](#11-安全与合规)
- [12. 架构与数据模型](#12-架构与数据模型)
- [13. 性能与资源建议](#13-性能与资源建议)
- [14. 路线图与可扩展方向](#14-路线图与可扩展方向)
- [15. FAQ](#15-faq)

---

## 1. 概览
- 采集方式：
  - PCAP 抓包：基于 libpcap/Npcap 对网卡流量进行分层解析（链路/网络/传输/应用）。
  - MITM 代理：作为 HTTP/HTTPS 代理，中间人解密并记录完整请求/响应（URL、Query、Headers、Body 等）。
- Web UI：
  - 状态与统计、抓包控制、代理控制、下载根证书、查看最近数据包（PCAP）与 HTTP Flows（代理采集）。
- 目标：从“能监听”到“能看懂”，支持路径/参数/请求体/响应体的可读化展示。

---

## 2. 快速开始
```bash
# 进入服务目录并启动
cd server/cmd
go run main.go

# 打开浏览器
http://localhost:8080/
```
- UI 中点击“代理(MITM)”→ 勾选 HTTPS → 启动代理
- 点击“下载根证书”，按第 8 节导入证书
- 将系统/浏览器 HTTP 与 HTTPS 代理设为 `127.0.0.1:8899`
- 访问任意站点，UI 的 “HTTP Flows（代理采集）” 将出现记录

如需底层抓包（PCAP），按第 6 节获取网卡名后在 UI 输入并开始。

---

## 3. 三种系统的环境准备

### 3.1 Windows
- 安装 Go 1.23+
- PCAP：安装 Npcap（安装选项勾选 “WinPcap API-compatible Mode”）
- 建议用管理员运行（抓包/端口权限）
- 防火墙放行 8080 与 8899

### 3.2 macOS
- `brew install go`
- PCAP 一般无需额外驱动（需 `sudo` 权限运行）
- 代理在“系统设置 → 网络 → 代理”处配置 HTTP/HTTPS 代理

### 3.3 Linux（以 Ubuntu/Debian 为例）
- `sudo apt install golang` 或安装官方 tar 包
- 安装 `libpcap`；抓包建议 root/sudo 或给二进制赋能
- 桌面环境网络代理或浏览器代理设置

---

## 4. 启动与访问
- 启动：`cd server/cmd && go run main.go`
- 访问 UI：`http://localhost:8080/`（跳转到 `/ui/`）

---

## 5. Web UI 使用指南
- 顶部卡片
  - 运行状态（PCAP）、统计信息
- 抓包控制（PCAP）
  - 输入网卡名 → 开始/停止；支持清空已抓数据
  - Windows 网卡名需用 `\\Device\\NPF_{GUID}`（见 `/api/interfaces`）
- 代理（MITM）控制
  - 设置监听地址（默认 `:8899`）
  - 勾选 HTTPS(MITM)，点击“启动代理”
  - “下载根证书”导入为受信任根
  - 设置系统/浏览器 HTTP 与 HTTPS 代理为 `127.0.0.1:8899`
- 数据表格
  - 最近数据包（PCAP）
  - HTTP Flows（代理采集）

---

## 6. PCAP 抓包
- 获取网卡名：`GET /api/interfaces`
  - Windows：复制 `Name`（形如 `\\Device\\NPF_{GUID}`），不要使用“以太网/WLAN”的友好名
  - macOS/Linux：可显示 `en0/eth0` 等
- 启动：`POST /api/start?iface=网卡名`
- 停止：`POST /api/stop`
- 查看：`GET /api/packets?limit=200`、`GET /api/stats`
- 清空：`DELETE /api/packets`
- 注意：Windows 需 Npcap + 管理员；macOS/Linux 需 root/sudo

---

## 7. MITM 代理（HTTP/HTTPS）
- 启动：`POST /api/proxy/start?addr=:8899&https=1`（或在 UI 勾选 HTTPS）
- 停止：`POST /api/proxy/stop`
- 状态：`GET /api/proxy/status`
- 证书：`GET /api/proxy/ca` 下载，`POST /api/proxy/ca/generate` 重新生成
- flows：
  - 列表：`GET /api/flows?limit=200`
  - 详情：`GET /api/flows/:id`
  - 解码详情：`GET /api/flows/:id?decoded=1`（返回 `response.body_text`，自动解压 gzip/deflate，文本类型可读）
- 建议验证路径：
  - HTTP：访问 `http://neverssl.com`
  - HTTPS（忽略证书用于通路验证）：`curl -x http://127.0.0.1:8899 -k https://example.com`
- Chrome/Edge 如仍不通，禁用 QUIC：`chrome://flags/#enable-quic` → Disabled

---

## 8. 证书导入（分系统）
> 必须导入根证书为“受信任的根证书颁发机构”，HTTPS 才能被解密与展示完整内容。

### 8.1 Windows
- 系统级（管理员 PowerShell）：
```powershell
certutil -addstore -f "Root" "C:\路径\proxy_root_ca.pem"
```
- 当前用户（非管理员）：
```powershell
certutil -user -addstore -f "Root" "C:\路径\proxy_root_ca.pem"
```
- 图形界面：`mmc` → 文件→添加/删除管理单元→证书→计算机账户→受信任的根证书颁发机构→证书→导入

### 8.2 macOS
- 钥匙串访问（Keychain）→ 系统或登录 → 导入 `proxy_root_ca.pem` → 设置为“始终信任”

### 8.3 Linux
- 系统（curl 等）：复制到 `/usr/local/share/ca-certificates/xxx.crt` → `sudo update-ca-certificates`
- Firefox：设置 → 隐私与安全 → 证书 → 导入 → 勾选信任

导入后重启浏览器。若曾导入旧 CA，先删除旧证书再导入最新证书。

---

## 9. 常用 API 参考
- PCAP：
  - `GET /api/interfaces` 网卡列表
  - `POST /api/start?iface=...` 开始；`POST /api/stop` 停止
  - `GET /api/status` 状态；`GET /api/packets?limit=200` 列表；`GET /api/stats` 统计
  - `DELETE /api/packets` 清空
- 代理/证书：
  - `GET /api/proxy/status`；`POST /api/proxy/start?addr=:8899&https=1`；`POST /api/proxy/stop`
  - `GET /api/proxy/ca` 下载；`POST /api/proxy/ca/generate` 重生成
- flows：
  - `GET /api/flows?limit=200` 列表
  - `GET /api/flows/:id` 详情
  - `GET /api/flows/:id?decoded=1` 解码详情（含 `response.body_text`）
  - `GET /api/flows/stats` 统计；`DELETE /api/flows` 清空

---

## 10. 故障排查
- PCAP 报错 `couldn't load wpcap.dll`（Windows）
  - 未安装 Npcap 或缺少权限；安装 Npcap（WinPcap 兼容），管理员运行，检查 `C:\Windows\System32\wpcap.dll`
- 代理无数据或握手失败（unknown certificate/EOF）
  - 同时设置了 HTTP 与 HTTPS 代理
  - 已导入并信任根证书；重启浏览器
  - 禁用 QUIC，避免走 H3 绕过代理
  - 检查防火墙放行 8899
- 证书错误 `ERR_CERT_COMMON_NAME_INVALID`
  - 使用动态签发证书（项目已实现），若仍出现，删除旧 CA 并导入最新 `proxy_root_ca.pem`
- 验证代理链路：
```bash
curl -x http://127.0.0.1:8899 -k https://example.com
```

---

## 11. 安全与合规
- MITM 仅用于本地调试，不得用于未经授权的网络流量。
- 根证书只导入受控设备，不需要时请删除。
- 流量数据默认存内存，注意敏感信息处理与访问控制。

---

## 12. 架构与数据模型
### 12.1 目录结构（关键）
- `cmd/main.go`：入口，路由与服务装配
- `internal/capture/`：PCAP 抓包与分层解析
- `internal/proxy/`：MITM 代理、CA 生成与动态签发
- `pkg/storage/`：内存存储（packet/flow）
- `pkg/utils/`：工具（env、http 解码等）
- `server/web/`：前端静态页面
- `server/docs/`：文档

### 12.2 核心组件
- Capturer（PCAP）：打开网卡→设置 BPF→解析各层→写入 `Storage`
- ProxyServer（MITM）：HTTP/HTTPS 代理→请求/响应钩子→Flow 存储
- CA 管理：根证书生成/加载、为 host 动态签发叶子证书（包含 SAN）
- Storage：
  - Packets：最近 N 条包、统计、清空
  - Flows：最近 N 条 HTTP 流、ID 查询、统计、清空

### 12.3 数据模型（简要）
- Packet（分层信息：Link/Network/Transport/Application/Error/Metadata）
- Flow：
  - Request：method、url、path、query、host、headers、body、proto
  - Response：status、status_code、headers、body、proto
  - 附带 start/end/latency_ms、remote_addr、scheme

### 12.4 请求路径（MITM）
1) 客户端 → 代理（CONNECT）
2) 代理使用 CA 动态为目标签发证书（含 SAN），与客户端握手
3) 发起到上游的 TLS/HTTP 请求
4) OnRequest/OnResponse 捕获并存储 Flow

---

## 13. 性能与资源建议
- 默认内存限长：Packets（1w）、Flows（2w）。可按需调大或改为持久化。
- 代理与抓包不建议在超低性能设备长时间运行。
- Windows 上 PCAP 抓包建议仅在需要时启用，避免系统开销。

---

## 14. 路线图与可扩展方向
- 解码增强：
  - Brotli（br）、GBK/GB2312 → UTF-8 自动转码
  - JSON/HTML/文本 Pretty-Print 与高亮
- Flows 检索：
  - 分页、域名/方法/路径/状态码/时间范围/关键字筛选
- 持久化：
  - SQLite/MySQL 存储历史、索引优化、归档策略
- 调试能力：
  - 重写（URL/Headers/Body）、断点挂起、回放、镜像/反向代理
  - 脚本（Python）管线与安全沙箱
- 观测与运维：
  - Prometheus 指标、优雅退出、结构化日志
- 架构优化：
  - 抽象 app/server/router 层，集中依赖注入与生命周期管理

---

## 15. FAQ
- Q：为何 UI 有“代理运行中”但没有数据？
  - A：多数站点是 HTTPS，需导入根证书并设置 HTTPS 代理；禁用 QUIC；确认未被其他代理/杀软拦截。
- Q：PCAP 网卡名怎么填？
  - A：Windows 使用 `GET /api/interfaces` 返回的 `\\Device\\NPF_{GUID}`；macOS/Linux 使用 `en0/eth0` 等。
- Q：`COMMON_NAME_INVALID` 怎么办？
  - A：删除旧 CA，下载并导入最新 CA；项目已实现 SAN 动态签发，正常不会再出现。
