package layer

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// LinkLayerInfo 存储链路层信息
type LinkLayerInfo struct {
	// 基本信息
	Timestamp time.Time `json:"timestamp"` // 数据包捕获时间

	// 链路层信息
	SrcMAC  string `json:"src_mac,omitempty"`  // 源MAC地址
	DstMAC  string `json:"dst_mac,omitempty"`  // 目标MAC地址
	EthType string `json:"eth_type,omitempty"` // 以太网类型
	Length  int    `json:"length,omitempty"`   // 帧长度
}

// ExtractLinkLayerInfo 提取链路层信息并填充到LinkLayerInfo结构体中
func ExtractLinkLayerInfo(linkLayer gopacket.LinkLayer) *LinkLayerInfo {
	info := &LinkLayerInfo{
		Timestamp: time.Now(),
	}

	if linkLayer == nil {
		return info
	}

	// 处理以太网层
	switch l := linkLayer.(type) {
	case *layers.Ethernet:
		info.SrcMAC = l.SrcMAC.String()
		info.DstMAC = l.DstMAC.String()
		info.EthType = l.EthernetType.String()
		info.Length = len(l.Contents)
	}

	return info
}
