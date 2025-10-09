package storage

import (
	"strings"
	"sync"
	"time"
)

// MemoryStorage 是基于内存的存储实现
type MemoryStorage struct {
	packets []*PacketInfo
	mu      sync.RWMutex
	stats   Stats
	ipSet   map[string]bool
	portSet map[uint16]bool
}

// NewMemoryStorage 创建一个新的内存存储实例
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		packets: make([]*PacketInfo, 0),
		stats: Stats{
			StartTime: time.Now(),
		},
		ipSet:   make(map[string]bool),
		portSet: make(map[uint16]bool),
	}
}

// StorePacket 存储一个数据包
func (m *MemoryStorage) StorePacket(packet *PacketInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 限制存储的数据包数量，避免内存溢出
	const maxPackets = 10000
	if len(m.packets) >= maxPackets {
		// 移除最旧的数据包
		m.packets = m.packets[1:]
	}

	m.packets = append(m.packets, packet)

	// 更新统计信息
	m.stats.TotalPackets++
	m.stats.LastPacketTime = packet.Metadata.CaptureTime

	// 统计协议类型
	if packet.TransportLayer.Protocol == "TCP" {
		if packet.TransportLayer.DstPort == 80 || packet.TransportLayer.SrcPort == 80 {
			m.stats.HTTPPackets++
		}
		if packet.TransportLayer.DstPort == 443 || packet.TransportLayer.SrcPort == 443 {
			m.stats.HTTPSPackets++
		}
	}

	// 统计唯一IP和端口
	m.ipSet[packet.NetworkLayer.SrcIP] = true
	m.ipSet[packet.NetworkLayer.DstIP] = true
	m.portSet[packet.TransportLayer.SrcPort] = true
	m.portSet[packet.TransportLayer.DstPort] = true

	m.stats.UniqueIPs = len(m.ipSet)
	m.stats.UniquePorts = len(m.portSet)
}

// GetPackets 获取指定数量的最新数据包
func (m *MemoryStorage) GetPackets(limit int) []*PacketInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 如果limit为0或负数，返回所有数据包
	if limit <= 0 || limit > len(m.packets) {
		limit = len(m.packets)
	}

	// 返回最新的limit个数据包
	start := len(m.packets) - limit
	return m.packets[start:]
}

// GetPacketsByFilter 根据过滤条件获取数据包
func (m *MemoryStorage) GetPacketsByFilter(filter Filter) []*PacketInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*PacketInfo, 0)

	// 遍历所有数据包
	for _, packet := range m.packets {
		if m.matchesFilter(packet, filter) {
			result = append(result, packet)
		}
	}

	return result
}

// matchesFilter 检查数据包是否匹配过滤条件
func (m *MemoryStorage) matchesFilter(packet *PacketInfo, filter Filter) bool {
	// 协议过滤
	if filter.Protocol != "" && packet.TransportLayer.Protocol != filter.Protocol {
		return false
	}

	// IP过滤
	if filter.SrcIP != "" && packet.NetworkLayer.SrcIP != filter.SrcIP {
		return false
	}
	if filter.DstIP != "" && packet.NetworkLayer.DstIP != filter.DstIP {
		return false
	}

	// 端口过滤
	if filter.Port != 0 && packet.TransportLayer.SrcPort != filter.Port && packet.TransportLayer.DstPort != filter.Port {
		return false
	}

	// HTTP方法过滤
	if filter.HTTPMethod != "" && packet.ApplicationLayer.HTTPMethod != filter.HTTPMethod {
		return false
	}

	// 时间过滤
	if !filter.StartTime.IsZero() && packet.Metadata.CaptureTime.Before(filter.StartTime) {
		return false
	}
	if !filter.EndTime.IsZero() && packet.Metadata.CaptureTime.After(filter.EndTime) {
		return false
	}

	// 主机名过滤
	if filter.Host != "" && packet.ApplicationLayer.Host != filter.Host {
		return false
	}

	// 域名过滤
	if filter.Domain != "" && packet.ApplicationLayer.Domain != filter.Domain {
		return false
	}

	// 路径过滤
	if filter.Path != "" && packet.ApplicationLayer.Path != filter.Path {
		return false
	}

	// User-Agent过滤
	if filter.UserAgent != "" && packet.ApplicationLayer.UserAgent != filter.UserAgent {
		return false
	}

	// Content-Type过滤
	if filter.ContentType != "" && packet.ApplicationLayer.ContentType != filter.ContentType {
		return false
	}

	// Referer过滤
	if filter.Referer != "" && packet.ApplicationLayer.Referer != filter.Referer {
		return false
	}

	// Server过滤
	if filter.Server != "" && packet.ApplicationLayer.Server != filter.Server {
		return false
	}

	// 文本搜索
	if filter.SearchText != "" {
		searchText := filter.SearchText
		// 在各种字段中搜索文本
		if !containsString(packet.NetworkLayer.SrcIP, searchText) &&
			!containsString(packet.NetworkLayer.DstIP, searchText) &&
			!containsString(packet.ApplicationLayer.Host, searchText) &&
			!containsString(packet.ApplicationLayer.Path, searchText) &&
			!containsString(packet.ApplicationLayer.UserAgent, searchText) &&
			!containsString(packet.ApplicationLayer.ContentType, searchText) &&
			!containsString(packet.ApplicationLayer.Referer, searchText) &&
			!containsString(packet.ApplicationLayer.Server, searchText) {
			return false
		}
	}

	return true
}

// containsString 检查字符串是否包含另一个字符串（不区分大小写）
func containsString(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Clear 清空所有数据包
func (m *MemoryStorage) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.packets = make([]*PacketInfo, 0)
	m.stats = Stats{
		StartTime: time.Now(),
	}
	m.ipSet = make(map[string]bool)
	m.portSet = make(map[uint16]bool)
}

// GetStats 获取统计信息
func (m *MemoryStorage) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.stats
}
