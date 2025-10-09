package layer

import (
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// extractNetworkLayerInfo 提取网络层信息并填充到NetworkLayerInfo结构体中
func ExtractNetworkLayerInfo(networkLayer gopacket.NetworkLayer) *NetworkLayerInfo {
	info := &NetworkLayerInfo{
		Timestamp: time.Now(),
	}

	// 特别处理IPv6层
	switch l := networkLayer.(type) {
	case *layers.IPv6:
		info.IPVersion = 6
		info.SrcIP = l.SrcIP.String()
		info.DstIP = l.DstIP.String()
		info.Length = int(l.Length)
		info.TTL = int(l.HopLimit)
		info.TrafficClass = int(l.TrafficClass)
		info.FlowLabel = int(l.FlowLabel)
		info.NextHeader = l.NextHeader.String()
		info.HopLimit = int(l.HopLimit)

		// 验证IP地址是否有效
		info.IsSrcIPValid = net.ParseIP(l.SrcIP.String()) != nil
		info.IsDstIPValid = net.ParseIP(l.DstIP.String()) != nil

		// 检查是否为环回地址
		info.IsSrcLoopback = l.SrcIP.IsLoopback()
		info.IsDstLoopback = l.DstIP.IsLoopback()

		// 检查是否为链路本地地址
		info.IsSrcLinkLocal = l.SrcIP.IsLinkLocalUnicast() || l.SrcIP.IsLinkLocalMulticast()
		info.IsDstLinkLocal = l.DstIP.IsLinkLocalUnicast() || l.DstIP.IsLinkLocalMulticast()

	case *layers.IPv4:
		info.IPVersion = 4
		info.SrcIP = l.SrcIP.String()
		info.DstIP = l.DstIP.String()
		info.Length = int(l.Length)
		info.TTL = int(l.TTL)
		info.IHL = int(l.IHL)
		info.TOS = int(l.TOS)
		info.Identifier = int(l.Id)
		info.Flags = int(l.Flags)
		info.FragOffset = int(l.FragOffset)
		info.Checksum = int(l.Checksum)
		info.Protocol = l.Protocol.String()

		// 处理IPv4选项和填充
		if len(l.Options) > 0 {
			// 将选项转换为字节切片存储
			optionsData := make([]byte, 0)
			for _, opt := range l.Options {
				// 将每个选项的原始数据添加到字节切片中
				optionsData = append(optionsData, opt.OptionData...)
			}
			info.Options = optionsData
		}

		if len(l.Padding) > 0 {
			info.Padding = make([]byte, len(l.Padding))
			copy(info.Padding, l.Padding)
		}

		// 验证IP地址是否有效
		info.IsSrcIPValid = net.ParseIP(l.SrcIP.String()) != nil
		info.IsDstIPValid = net.ParseIP(l.DstIP.String()) != nil

		// 检查是否为环回地址
		info.IsSrcLoopback = l.SrcIP.IsLoopback()
		info.IsDstLoopback = l.DstIP.IsLoopback()

		// 检查是否为链路本地地址
		info.IsSrcLinkLocal = l.SrcIP.IsLinkLocalUnicast() || l.SrcIP.IsLinkLocalMulticast()
		info.IsDstLinkLocal = l.DstIP.IsLinkLocalUnicast() || l.DstIP.IsLinkLocalMulticast()
	}

	return info
}

// NetworkLayerInfo 存储网络层信息
type NetworkLayerInfo struct {
	Timestamp time.Time `json:"timestamp"` // 数据包捕获时间

	// 网络层通用信息
	IPVersion int    `json:"ip_version"`    // IP版本 (4 or 6)
	SrcIP     string `json:"src_ip"`        // 源IP地址
	DstIP     string `json:"dst_ip"`        // 目标IP地址
	Protocol  string `json:"protocol"`      // 协议类型
	Length    int    `json:"length"`        // 数据包长度
	TTL       int    `json:"ttl,omitempty"` // IPv4 TTL or IPv6 HopLimit

	// IPv4 特有字段
	IHL        int    `json:"ihl,omitempty"`         // 首部长度
	TOS        int    `json:"tos,omitempty"`         // 服务类型
	Identifier int    `json:"identifier,omitempty"`  // 标识符
	Flags      int    `json:"flags,omitempty"`       // 标志位
	FragOffset int    `json:"frag_offset,omitempty"` // 片偏移
	Checksum   int    `json:"checksum,omitempty"`    // 首部校验和
	Options    []byte `json:"options,omitempty"`     // IPv4选项
	Padding    []byte `json:"padding,omitempty"`     // 填充字段

	// IPv6 特有字段
	TrafficClass int    `json:"traffic_class,omitempty"` // 流量类别
	FlowLabel    int    `json:"flow_label,omitempty"`    // 流标签
	NextHeader   string `json:"next_header,omitempty"`   // 下一报头
	HopLimit     int    `json:"hop_limit,omitempty"`     // 跳数限制

	// 地址类型标记
	IsSrcLoopback  bool `json:"is_src_loopback"`   // 源IP是否为环回地址
	IsDstLoopback  bool `json:"is_dst_loopback"`   // 目标IP是否为环回地址
	IsSrcLinkLocal bool `json:"is_src_link_local"` // 源IP是否为链路本地地址
	IsDstLinkLocal bool `json:"is_dst_link_local"` // 目标IP是否为链路本地地址
	IsSrcIPValid   bool `json:"is_src_ip_valid"`   // 源IP地址是否有效
	IsDstIPValid   bool `json:"is_dst_ip_valid"`   // 目标IP地址是否有效
}
