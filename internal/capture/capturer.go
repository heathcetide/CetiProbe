package capture

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"probe/internal/storage"
)

type Capturer struct {
	interfaceName string
	handle        *pcap.Handle
	storage       storage.Storage
	running       bool
	mu            sync.RWMutex
}

func NewCapturer(interfaceName string, storage storage.Storage) (*Capturer, error) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("打开网络接口失败: %v", err)
	}

	// 设置过滤器，只捕获HTTP/HTTPS流量
	err = handle.SetBPFFilter("tcp port 80 or tcp port 443 or tcp port 8080 or tcp port 3000")
	if err != nil {
		return nil, fmt.Errorf("设置BPF过滤器失败: %v", err)
	}

	return &Capturer{
		interfaceName: interfaceName,
		handle:        handle,
		storage:       storage,
		running:       false,
	}, nil
}

func (c *Capturer) Start() error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("抓包器已经在运行")
	}
	c.running = true
	c.mu.Unlock()

	fmt.Printf("开始抓包，网络接口: %s\n", c.interfaceName)

	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	for packet := range packetSource.Packets() {
		c.processPacket(packet)
	}

	return nil
}

func (c *Capturer) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		c.running = false
		c.handle.Close()
		fmt.Println("抓包已停止")
	}
}

func (c *Capturer) processPacket(packet gopacket.Packet) {
	// 解析网络层
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	// 解析传输层
	transportLayer := packet.TransportLayer()
	if transportLayer == nil {
		return
	}

	packetInfo := &storage.PacketInfo{
		Timestamp: time.Now(),
		Length:    len(packet.Data()),
	}

	// 解析IP层
	switch networkLayer.LayerType() {
	case layers.LayerTypeIPv4:
		ipv4 := networkLayer.(*layers.IPv4)
		packetInfo.SrcIP = ipv4.SrcIP.String()
		packetInfo.DstIP = ipv4.DstIP.String()
		packetInfo.Protocol = "IPv4"
	case layers.LayerTypeIPv6:
		ipv6 := networkLayer.(*layers.IPv6)
		packetInfo.SrcIP = ipv6.SrcIP.String()
		packetInfo.DstIP = ipv6.DstIP.String()
		packetInfo.Protocol = "IPv6"
	}

	// 解析传输层
	switch transportLayer.LayerType() {
	case layers.LayerTypeTCP:
		tcp := transportLayer.(*layers.TCP)
		packetInfo.SrcPort = uint16(tcp.SrcPort)
		packetInfo.DstPort = uint16(tcp.DstPort)
		packetInfo.Protocol = "TCP"

		// 检查是否是HTTP流量 - 扩展端口检测
		if packetInfo.DstPort == 80 || packetInfo.SrcPort == 80 ||
			packetInfo.DstPort == 8080 || packetInfo.SrcPort == 8080 ||
			packetInfo.DstPort == 3000 || packetInfo.SrcPort == 3000 {
			c.parseHTTP(packet, packetInfo)
		}

		// 调试信息：如果目标IP是139.155.132.244，打印调试信息
		if packetInfo.DstIP == "139.155.132.244" || packetInfo.SrcIP == "139.155.132.244" {
			fmt.Printf("调试: 捕获到目标IP %s 的数据包，域名: %s, Host: %s\n",
				packetInfo.DstIP, packetInfo.Domain, packetInfo.Host)
		}
	}

	// 存储数据包信息
	c.storage.StorePacket(packetInfo)
}

func (c *Capturer) parseHTTP(packet gopacket.Packet, packetInfo *storage.PacketInfo) {
	// 获取应用层数据
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer == nil {
		return
	}

	payload := applicationLayer.Payload()
	if len(payload) == 0 {
		return
	}

	// 简单的HTTP解析
	payloadStr := string(payload)

	// 调试信息：打印原始payload的前100个字符
	if packetInfo.DstIP == "139.155.132.244" || packetInfo.SrcIP == "139.155.132.244" {
		preview := payloadStr
		if len(preview) > 100 {
			preview = preview[:100]
		}
		fmt.Printf("调试: HTTP payload预览: %s\n", preview)
	}

	// 检查是否是HTTP请求 - 增强检测逻辑
	if c.isHTTPRequest(payloadStr) {
		c.parseHTTPRequest(payloadStr, packetInfo)
	}

	// 检查是否是HTTP响应
	if c.isHTTPResponse(payloadStr) {
		c.parseHTTPResponse(payloadStr, packetInfo)
	}
}

// isHTTPRequest 检查是否是HTTP请求
func (c *Capturer) isHTTPRequest(payload string) bool {
	if len(payload) < 4 {
		return false
	}

	// 检查常见的HTTP方法
	httpMethods := []string{"GET ", "POST", "PUT ", "HEAD", "DELE", "PATCH", "OPTIONS", "TRACE"}
	for _, method := range httpMethods {
		if len(payload) >= len(method) && payload[:len(method)] == method {
			return true
		}
	}

	// 检查是否包含HTTP头部特征
	lines := splitLines(payload)
	if len(lines) >= 2 {
		// 第一行是请求行，第二行应该包含Host头部
		for _, line := range lines[1:] {
			if strings.HasPrefix(strings.ToLower(line), "host:") {
				return true
			}
		}
	}

	return false
}

// isHTTPResponse 检查是否是HTTP响应
func (c *Capturer) isHTTPResponse(payload string) bool {
	if len(payload) < 4 {
		return false
	}

	// 检查HTTP版本
	if strings.HasPrefix(payload, "HTTP/") {
		return true
	}

	return false
}

func (c *Capturer) parseHTTPRequest(payload string, packetInfo *storage.PacketInfo) {
	lines := splitLines(payload)
	if len(lines) == 0 {
		return
	}

	// 解析请求行
	requestLine := lines[0]
	parts := splitWords(requestLine)
	if len(parts) >= 3 {
		packetInfo.HTTPMethod = parts[0]
		packetInfo.HTTPURL = parts[1]

		// 解析URL，提取域名和路径
		c.parseURL(parts[1], packetInfo)
	}

	// 解析头部 - 增强解析逻辑
	for _, line := range lines[1:] {
		if line == "" {
			break
		}

		// 查找冒号位置
		colonIndex := findColon(line)
		if colonIndex > 0 {
			key := strings.TrimSpace(line[:colonIndex])
			value := strings.TrimSpace(line[colonIndex+1:])

			// 特别处理Host头部
			if strings.ToLower(key) == "host" {
				packetInfo.Host = value
				// 提取域名（去掉端口）
				colonIndex := strings.Index(value, ":")
				if colonIndex > 0 {
					packetInfo.Domain = value[:colonIndex]
				} else {
					packetInfo.Domain = value
				}
			} else {
				c.parseHeader(key, value, packetInfo)
			}
		}
	}
}

func (c *Capturer) parseHTTPResponse(payload string, packetInfo *storage.PacketInfo) {
	lines := splitLines(payload)
	if len(lines) == 0 {
		return
	}

	// 解析状态行
	statusLine := lines[0]
	parts := splitWords(statusLine)
	if len(parts) >= 2 {
		packetInfo.HTTPStatus = parts[1]
	}

	// 解析头部
	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		if colonIndex := findColon(line); colonIndex > 0 {
			key := line[:colonIndex]
			value := line[colonIndex+1:]
			if len(value) > 0 && value[0] == ' ' {
				value = value[1:]
			}

			c.parseHeader(key, value, packetInfo)
		}
	}
}

// parseURL 解析URL，提取域名和路径
func (c *Capturer) parseURL(url string, packetInfo *storage.PacketInfo) {
	// 处理相对URL
	if strings.HasPrefix(url, "/") {
		packetInfo.Path = url
		return
	}

	// 处理完整URL
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// 提取协议后的部分
		urlPart := url
		if strings.HasPrefix(url, "http://") {
			urlPart = url[7:]
		} else if strings.HasPrefix(url, "https://") {
			urlPart = url[8:]
		}

		// 查找第一个斜杠
		slashIndex := strings.Index(urlPart, "/")
		if slashIndex > 0 {
			packetInfo.Host = urlPart[:slashIndex]
			packetInfo.Path = urlPart[slashIndex:]
		} else {
			packetInfo.Host = urlPart
			packetInfo.Path = "/"
		}

		// 提取域名（去掉端口）
		colonIndex := strings.Index(packetInfo.Host, ":")
		if colonIndex > 0 {
			packetInfo.Domain = packetInfo.Host[:colonIndex]
		} else {
			packetInfo.Domain = packetInfo.Host
		}
	}
}

// parseHeader 解析HTTP头部
func (c *Capturer) parseHeader(key, value string, packetInfo *storage.PacketInfo) {
	switch key {
	case "Host":
		packetInfo.Host = value
		// 提取域名（去掉端口）
		colonIndex := strings.Index(value, ":")
		if colonIndex > 0 {
			packetInfo.Domain = value[:colonIndex]
		} else {
			packetInfo.Domain = value
		}
	case "User-Agent":
		packetInfo.UserAgent = value
	case "Content-Type":
		packetInfo.ContentType = value
	case "Referer":
		packetInfo.Referer = value
	case "Server":
		packetInfo.Server = value
	case "Set-Cookie":
		packetInfo.SetCookie = value
	case "Cookie":
		packetInfo.Cookie = value
	case "Authorization":
		packetInfo.Authorization = value
	case "Accept":
		packetInfo.Accept = value
	case "Accept-Language":
		packetInfo.AcceptLanguage = value
	case "Accept-Encoding":
		packetInfo.AcceptEncoding = value
	case "Connection":
		packetInfo.Connection = value
	case "Cache-Control":
		packetInfo.CacheControl = value
	case "Pragma":
		packetInfo.Pragma = value
	case "If-Modified-Since":
		packetInfo.IfModifiedSince = value
	case "If-None-Match":
		packetInfo.IfNoneMatch = value
	case "Range":
		packetInfo.Range = value
	case "Content-Length":
		packetInfo.ContentLength = value
	case "Transfer-Encoding":
		packetInfo.TransferEncoding = value
	case "Location":
		packetInfo.Location = value
	case "Last-Modified":
		packetInfo.LastModified = value
	case "ETag":
		packetInfo.ETag = value
	case "Expires":
		packetInfo.Expires = value
	case "Date":
		packetInfo.Date = value
	case "Age":
		packetInfo.Age = value
	case "Via":
		packetInfo.Via = value
	case "X-Forwarded-For":
		packetInfo.XForwardedFor = value
	case "X-Real-IP":
		packetInfo.XRealIP = value
	case "X-Requested-With":
		packetInfo.XRequestedWith = value
	}
}

// 辅助函数
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func splitWords(s string) []string {
	var words []string
	start := 0
	inWord := false
	for i, c := range s {
		if c == ' ' || c == '\t' {
			if inWord {
				words = append(words, s[start:i])
				inWord = false
			}
		} else {
			if !inWord {
				start = i
				inWord = true
			}
		}
	}
	if inWord {
		words = append(words, s[start:])
	}
	return words
}

func findColon(s string) int {
	for i, c := range s {
		if c == ':' {
			return i
		}
	}
	return -1
}
