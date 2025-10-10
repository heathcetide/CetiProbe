package layer

import (
	"fmt"
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

// PrintLinkLayerInfo 打印链路层信息
func PrintLinkLayerInfo(linkInfo *LinkLayerInfo) {
	fmt.Println("  Link Layer 详细信息:")
	fmt.Printf("    时间戳: %v\n", linkInfo.Timestamp)
	fmt.Printf("    源MAC: %s\n", linkInfo.SrcMAC)
	fmt.Printf("    目标MAC: %s\n", linkInfo.DstMAC)
	fmt.Printf("    以太网类型: %s\n", linkInfo.EthType)
	fmt.Printf("    长度: %d 字节\n", linkInfo.Length)
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
