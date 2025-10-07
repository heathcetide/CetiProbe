package storage

import (
	"sync"
	"time"
)

type Storage interface {
	StorePacket(packet *PacketInfo)
	GetPackets(limit int) []*PacketInfo
	GetPacketsByFilter(filter Filter) []*PacketInfo
	Clear()
	GetStats() Stats
}

type PacketInfo struct {
	Timestamp   time.Time `json:"timestamp"`
	SrcIP       string    `json:"src_ip"`
	DstIP       string    `json:"dst_ip"`
	SrcPort     uint16    `json:"src_port"`
	DstPort     uint16    `json:"dst_port"`
	Protocol    string    `json:"protocol"`
	Length      int       `json:"length"`
	Payload     []byte    `json:"payload,omitempty"`
	HTTPMethod  string    `json:"http_method,omitempty"`
	HTTPURL     string    `json:"http_url,omitempty"`
	HTTPStatus  string    `json:"http_status,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	ContentType string    `json:"content_type,omitempty"`
}

type Filter struct {
	Protocol   string    `json:"protocol"`
	SrcIP      string    `json:"src_ip"`
	DstIP      string    `json:"dst_ip"`
	Port       uint16    `json:"port"`
	HTTPMethod string    `json:"http_method"`
	SearchText string    `json:"search_text"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

type Stats struct {
	TotalPackets   int       `json:"total_packets"`
	HTTPPackets    int       `json:"http_packets"`
	HTTPSPackets   int       `json:"https_packets"`
	StartTime      time.Time `json:"start_time"`
	LastPacketTime time.Time `json:"last_packet_time"`
	UniqueIPs      int       `json:"unique_ips"`
	UniquePorts    int       `json:"unique_ports"`
}

type MemoryStorage struct {
	packets []*PacketInfo
	mu      sync.RWMutex
	stats   Stats
	ipSet   map[string]bool
	portSet map[uint16]bool
}

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
	m.stats.LastPacketTime = packet.Timestamp

	// 统计协议类型
	if packet.Protocol == "TCP" {
		if packet.DstPort == 80 || packet.SrcPort == 80 {
			m.stats.HTTPPackets++
		}
		if packet.DstPort == 443 || packet.SrcPort == 443 {
			m.stats.HTTPSPackets++
		}
	}

	// 统计唯一IP和端口
	m.ipSet[packet.SrcIP] = true
	m.ipSet[packet.DstIP] = true
	m.portSet[packet.SrcPort] = true
	m.portSet[packet.DstPort] = true

	m.stats.UniqueIPs = len(m.ipSet)
	m.stats.UniquePorts = len(m.portSet)
}

func (m *MemoryStorage) GetPackets(limit int) []*PacketInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || limit > len(m.packets) {
		limit = len(m.packets)
	}

	// 返回最新的数据包
	start := len(m.packets) - limit
	if start < 0 {
		start = 0
	}

	result := make([]*PacketInfo, limit)
	copy(result, m.packets[start:])
	return result
}

func (m *MemoryStorage) GetPacketsByFilter(filter Filter) []*PacketInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PacketInfo

	for _, packet := range m.packets {
		if m.matchesFilter(packet, filter) {
			result = append(result, packet)
		}
	}

	return result
}

func (m *MemoryStorage) matchesFilter(packet *PacketInfo, filter Filter) bool {
	// 协议过滤
	if filter.Protocol != "" && packet.Protocol != filter.Protocol {
		return false
	}

	// IP过滤
	if filter.SrcIP != "" && packet.SrcIP != filter.SrcIP {
		return false
	}
	if filter.DstIP != "" && packet.DstIP != filter.DstIP {
		return false
	}

	// 端口过滤
	if filter.Port != 0 && packet.SrcPort != filter.Port && packet.DstPort != filter.Port {
		return false
	}

	// HTTP方法过滤
	if filter.HTTPMethod != "" && packet.HTTPMethod != filter.HTTPMethod {
		return false
	}

	// 时间过滤
	if !filter.StartTime.IsZero() && packet.Timestamp.Before(filter.StartTime) {
		return false
	}
	if !filter.EndTime.IsZero() && packet.Timestamp.After(filter.EndTime) {
		return false
	}

	// 文本搜索
	if filter.SearchText != "" {
		searchText := filter.SearchText
		if !contains(packet.SrcIP, searchText) &&
			!contains(packet.DstIP, searchText) &&
			!contains(packet.HTTPMethod, searchText) &&
			!contains(packet.HTTPURL, searchText) &&
			!contains(packet.UserAgent, searchText) {
			return false
		}
	}

	return true
}

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

func (m *MemoryStorage) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.stats
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
