package capture

import (
	"fmt"
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

// PacketInfo 已移动到 storage 包中，这里不再重复定义

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

		// 检查是否是HTTP流量
		if packetInfo.DstPort == 80 || packetInfo.SrcPort == 80 {
			c.parseHTTP(packet, packetInfo)
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

	// 检查是否是HTTP请求
	if len(payloadStr) > 4 && (payloadStr[:4] == "GET " || payloadStr[:4] == "POST" ||
		payloadStr[:4] == "PUT " || payloadStr[:4] == "HEAD" || payloadStr[:4] == "DELE") {
		c.parseHTTPRequest(payloadStr, packetInfo)
	}

	// 检查是否是HTTP响应
	if len(payloadStr) > 4 && payloadStr[:4] == "HTTP" {
		c.parseHTTPResponse(payloadStr, packetInfo)
	}
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

			switch key {
			case "User-Agent":
				packetInfo.UserAgent = value
			case "Content-Type":
				packetInfo.ContentType = value
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

			switch key {
			case "Content-Type":
				packetInfo.ContentType = value
			}
		}
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
