package storage

import (
	"probe/internal/models"
	"time"
)

// Storage 是存储接口，定义了所有存储实现必须支持的方法
type Storage interface {
	// StorePacket 存储一个数据包
	StorePacket(packet *models.PacketInfo)

	// GetPackets 获取指定数量的最新数据包
	GetPackets(limit int) []*models.PacketInfo

	// GetPacketsByFilter 根据过滤条件获取数据包
	GetPacketsByFilter(filter Filter) []*models.PacketInfo

	// Clear 清空所有数据包
	Clear()

	// GetStats 获取统计信息
	GetStats() Stats
}

// Filter 用于过滤数据包的条件
type Filter struct {
	Protocol    string    `json:"protocol"`
	SrcIP       string    `json:"src_ip"`
	DstIP       string    `json:"dst_ip"`
	Port        uint16    `json:"port"`
	HTTPMethod  string    `json:"http_method"`
	SearchText  string    `json:"search_text"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Host        string    `json:"host"`
	Domain      string    `json:"domain"`
	Path        string    `json:"path"`
	UserAgent   string    `json:"user_agent"`
	ContentType string    `json:"content_type"`
	Referer     string    `json:"referer"`
	Server      string    `json:"server"`
}

// Stats 存储统计信息
type Stats struct {
	TotalPackets   int       `json:"total_packets"`
	HTTPPackets    int       `json:"http_packets"`
	HTTPSPackets   int       `json:"https_packets"`
	StartTime      time.Time `json:"start_time"`
	LastPacketTime time.Time `json:"last_packet_time"`
	UniqueIPs      int       `json:"unique_ips"`
	UniquePorts    int       `json:"unique_ports"`
}
