package layer

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TransportLayerInfo 存储传输层信息
type TransportLayerInfo struct {
	// 基本信息
	Timestamp time.Time `json:"timestamp"` // 数据包捕获时间

	// 传输层通用信息
	SrcPort    uint16 `json:"src_port"`              // 源端口号
	DstPort    uint16 `json:"dst_port"`              // 目标端口号
	Protocol   string `json:"protocol"`              // 协议类型 (TCP/UDP)
	Length     int    `json:"length,omitempty"`      // 数据长度
	SeqNumber  uint32 `json:"seq_number,omitempty"`  // 序列号 (TCP)
	AckNumber  uint32 `json:"ack_number,omitempty"`  // 确认号 (TCP)
	WindowSize uint16 `json:"window_size,omitempty"` // 窗口大小 (TCP)
	Checksum   uint16 `json:"checksum,omitempty"`    // 校验和
	UrgentPtr  uint16 `json:"urgent_ptr,omitempty"`  // 紧急指针 (TCP)

	// TCP标志位
	IsFIN bool `json:"is_fin,omitempty"` // 结束标志
	IsSYN bool `json:"is_syn,omitempty"` // 同步标志
	IsRST bool `json:"is_rst,omitempty"` // 重置标志
	IsPSH bool `json:"is_psh,omitempty"` // 推送标志
	IsACK bool `json:"is_ack,omitempty"` // 确认标志
	IsURG bool `json:"is_urg,omitempty"` // 紧急标志
	IsECE bool `json:"is_ece,omitempty"` // ECN-Echo标志
	IsCWR bool `json:"is_cwr,omitempty"` // 拥塞窗口减少标志

	// UDP特有字段
	UDPChecksum uint16 `json:"udp_checksum,omitempty"` // UDP校验和
	UDPLength   uint16 `json:"udp_length,omitempty"`   // UDP长度
}

// ExtractTransportLayerInfo 提取传输层信息并填充到TransportLayerInfo结构体中
func ExtractTransportLayerInfo(transportLayer gopacket.TransportLayer) *TransportLayerInfo {
	info := &TransportLayerInfo{
		Timestamp: time.Now(),
	}

	if transportLayer == nil {
		return info
	}

	// 处理TCP层
	switch l := transportLayer.(type) {
	case *layers.TCP:
		info.Protocol = "TCP"
		info.SrcPort = uint16(l.SrcPort)
		info.DstPort = uint16(l.DstPort)
		info.SeqNumber = uint32(l.Seq)
		info.AckNumber = uint32(l.Ack)
		info.WindowSize = uint16(l.Window)
		info.Checksum = uint16(l.Checksum)
		info.UrgentPtr = uint16(l.Urgent)

		// 设置TCP标志位
		info.IsFIN = l.FIN
		info.IsSYN = l.SYN
		info.IsRST = l.RST
		info.IsPSH = l.PSH
		info.IsACK = l.ACK
		info.IsURG = l.URG
		info.IsECE = l.ECE
		info.IsCWR = l.CWR

	case *layers.UDP:
		info.Protocol = "UDP"
		info.SrcPort = uint16(l.SrcPort)
		info.DstPort = uint16(l.DstPort)
		info.UDPLength = uint16(l.Length)
		info.UDPChecksum = uint16(l.Checksum)
	}

	return info
}

// PrintTransportLayerDetails 打印传输层详细信息
func PrintTransportLayerDetails(transInfo *TransportLayerInfo) {
	fmt.Println("  Transport Layer 详细信息:")
	fmt.Printf("    协议: %s\n", transInfo.Protocol)
	fmt.Printf("    源端口: %d\n", transInfo.SrcPort)
	fmt.Printf("    目标端口: %d\n", transInfo.DstPort)

	if transInfo.Protocol == "TCP" {
		fmt.Printf("    序列号: %d\n", transInfo.SeqNumber)
		fmt.Printf("    确认号: %d\n", transInfo.AckNumber)
		fmt.Printf("    窗口大小: %d\n", transInfo.WindowSize)
		fmt.Printf("    校验和: %d\n", transInfo.Checksum)
		fmt.Printf("    紧急指针: %d\n", transInfo.UrgentPtr)

		// 打印TCP标志位
		fmt.Print("    TCP标志位: ")
		flags := make([]string, 0)
		if transInfo.IsFIN {
			flags = append(flags, "FIN")
		}
		if transInfo.IsSYN {
			flags = append(flags, "SYN")
		}
		if transInfo.IsRST {
			flags = append(flags, "RST")
		}
		if transInfo.IsPSH {
			flags = append(flags, "PSH")
		}
		if transInfo.IsACK {
			flags = append(flags, "ACK")
		}
		if transInfo.IsURG {
			flags = append(flags, "URG")
		}
		if transInfo.IsECE {
			flags = append(flags, "ECE")
		}
		if transInfo.IsCWR {
			flags = append(flags, "CWR")
		}

		if len(flags) > 0 {
			fmt.Println(strings.Join(flags, ", "))
		} else {
			fmt.Println("无")
		}
	} else if transInfo.Protocol == "UDP" {
		fmt.Printf("    UDP长度: %d\n", transInfo.UDPLength)
		fmt.Printf("    UDP校验和: %d\n", transInfo.UDPChecksum)
	}
}
