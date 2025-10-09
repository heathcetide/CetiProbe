package layer

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
)

// PacketMetadataInfo 存储数据包元信息
type PacketMetadataInfo struct {
	// 数据包基本信息
	CaptureTime   time.Time `json:"capture_time"`   // 捕获时间
	DataSize      int       `json:"data_size"`      // 数据包数据大小(字节)
	WireLength    int       `json:"wire_length"`    // 线路上的原始长度
	CaptureLength int       `json:"capture_length"` // 实际捕获的长度
	Truncated     bool      `json:"truncated"`      // 是否被截断

	// 捕获相关信息
	InterfaceIndex  int  `json:"interface_index"`   // 接口索引
	CaptureLengthOk bool `json:"capture_length_ok"` // 捕获长度是否正常
}

// PrintPacketMetadataInfo 打印数据包元信息
func PrintPacketMetadataInfo(metaInfo *PacketMetadataInfo) {
	fmt.Println("Packet Metadata 详细信息:")
	fmt.Printf("  捕获时间: %v\n", metaInfo.CaptureTime)
	fmt.Printf("  数据大小: %d 字节\n", metaInfo.DataSize)
	fmt.Printf("  线路长度: %d 字节\n", metaInfo.WireLength)
	fmt.Printf("  捕获长度: %d 字节\n", metaInfo.CaptureLength)
	fmt.Printf("  是否被截断: %t\n", metaInfo.Truncated)
	fmt.Printf("  接口索引: %d\n", metaInfo.InterfaceIndex)
}

// ExtractPacketMetadataInfo 提取数据包元信息并填充到PacketMetadataInfo结构体中
func ExtractPacketMetadataInfo(packet gopacket.Packet) *PacketMetadataInfo {
	info := &PacketMetadataInfo{
		CaptureTime: time.Now(),
	}

	// 获取数据包元信息
	info.DataSize = len(packet.Data())

	// 获取捕获信息
	if metadata := packet.Metadata(); metadata != nil {
		info.CaptureTime = metadata.Timestamp
		info.WireLength = metadata.CaptureInfo.Length
		info.CaptureLength = metadata.CaptureInfo.CaptureLength
		info.InterfaceIndex = metadata.CaptureInfo.InterfaceIndex

		// 检查是否被截断
		info.Truncated = metadata.CaptureInfo.CaptureLength < metadata.CaptureInfo.Length
		info.CaptureLengthOk = metadata.CaptureInfo.CaptureLength >= 0
	}

	return info
}
