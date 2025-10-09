package capture

import (
	"encoding/json"
	"fmt"
	"os"
	"probe/internal/capture/layer"
	"probe/internal/models"
	"probe/pkg/storage"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Capturer struct {
	interfaceName string
	handle        *pcap.Handle
	storage       storage.Storage
	running       bool
	mu            sync.RWMutex
	domainMap     map[string]string // IP到域名的映射
	domainMu      sync.RWMutex      // 保护domainMap的互斥锁
}

func NewCapturer(interfaceName string) (*Capturer, error) {
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("打开网络接口失败: %v", err)
	}

	// 设置过滤器，只捕获HTTP/HTTPS流量
	err = handle.SetBPFFilter("tcp port 80 or tcp port 443 or tcp port 8080 or tcp port 3000 or udp port 53")
	if err != nil {
		return nil, fmt.Errorf("设置BPF过滤器失败: %v", err)
	}

	return &Capturer{
		interfaceName: interfaceName,
		handle:        handle,
		running:       false,
		domainMap:     make(map[string]string),
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
		// 分层输出数据包信息
		c.processPacket(packet)
	}

	return nil
}

// processPacket 分层输出数据包的详细信息
func (c *Capturer) processPacket(packet gopacket.Packet) {
	// 解析网络层
	networkLayer := packet.NetworkLayer()
	if networkLayer == nil {
		return
	}

	// 解析传输层
	transportLayer := packet.TransportLayer()
	if transportLayer == nil {
		// 添加调试日志，查看没有传输层的数据包
		// fmt.Printf("无传输层数据包: %v\n", packet.Dump())
		return
	}

	applicationLayer := packet.ApplicationLayer()
	if applicationLayer == nil {
		return
	}

	var packetInfo models.PacketInfo
	// 输出数据包元信息
	metaInfo := layer.ExtractPacketMetadataInfo(packet)
	packetInfo.Metadata = metaInfo

	// 处理链路层
	if linkLayer := packet.LinkLayer(); linkLayer != nil {
		packetInfo.LinkLayer = layer.ExtractLinkLayerInfo(linkLayer)
	}

	// 特别处理网络层，尤其是IPv6
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		packetInfo.NetworkLayer = layer.ExtractNetworkLayerInfo(networkLayer)
	}

	// 处理传输层
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		packetInfo.TransportLayer = layer.ExtractTransportLayerInfo(transportLayer)
	}

	// 处理应用层
	if appLayer := packet.ApplicationLayer(); appLayer != nil {
		packetInfo.ApplicationLayer = layer.ExtractApplicationLayerInfo(appLayer)
		c.collectDomainFromDNS(packet, &packetInfo)
	}

	// 处理错误层
	if errLayer := packet.ErrorLayer(); errLayer != nil {
		packetInfo.ErrorLayer = layer.ExtractErrorLayerInfo(errLayer)
	}

	if packetInfo.ApplicationLayer.Domain != "" && strings.Contains(packetInfo.ApplicationLayer.Domain, "code") {
		savePacketInfoToJSON(&packetInfo)
	}
}

// savePacketInfoToJSON 将数据包信息保存到JSON文件中
func savePacketInfoToJSON(packetInfo *models.PacketInfo) {
	// 创建文件名，使用时间戳确保唯一性
	filename := fmt.Sprintf("packet_%d.json", time.Now().UnixNano())

	// 将packetInfo序列化为JSON
	data, err := json.MarshalIndent(packetInfo, "", "  ")
	if err != nil {
		fmt.Printf("序列化数据包信息失败: %v\n", err)
		return
	}

	// 写入文件
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}

	fmt.Printf("数据包信息已保存到文件: %s\n", filename)
}

// collectDomainFromDNS 从DNS数据包中收集域名信息
func (c *Capturer) collectDomainFromDNS(packet gopacket.Packet, packetInfo *models.PacketInfo) {
	// 检查是否是DNS数据包
	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns, _ := dnsLayer.(*layers.DNS)

		// 添加调试日志
		fmt.Printf("DNS数据包捕获: OpCode=%v, QR=%v, 问题数=%d, 答案数=%d\n",
			dns.OpCode, dns.QR, len(dns.Questions), len(dns.Answers))

		if dns.OpCode == layers.DNSOpCodeQuery && len(dns.Questions) > 0 {
			// 从DNS查询中提取域名
			domain := string(dns.Questions[0].Name)
			packetInfo.ApplicationLayer.Domain = domain
			// DNS请求中不生成FullURL，因为我们不知道协议类型和完整路径
			// 只保留域名信息，前端可以根据需要自行构建URL
			fmt.Printf("DNS查询域名: %s\n", domain)

			// 如果是DNS响应且包含答案
			if dns.QR && len(dns.Answers) > 0 {
				fmt.Printf("DNS响应，答案数: %d\n", len(dns.Answers))
				for i, answer := range dns.Answers {
					// 记录域名和IP的映射关系，可用于后续数据包的域名解析
					if answer.IP != nil {
						fmt.Printf("DNS答案[%d]: 域名=%s, IP=%s\n", i, string(answer.Name), answer.IP.String())
						// 存储域名和IP的映射关系
						c.domainMu.Lock()
						c.domainMap[answer.IP.String()] = string(answer.Name)
						c.domainMu.Unlock()

						// 这里可以存储域名和IP的映射关系，供其他数据包使用
						// 简化起见，我们只在当前数据包中记录
						if answer.IP.String() == packetInfo.NetworkLayer.SrcIP || answer.IP.String() == packetInfo.NetworkLayer.DstIP {
							packetInfo.ApplicationLayer.Domain = string(answer.Name)
							fmt.Printf("匹配IP地址，设置域名: %s\n", packetInfo.ApplicationLayer.Domain)
						}
					}
				}
			}
		}
	} else {
		// 检查是否是UDP 53端口的数据包但没有DNS层
		if (packetInfo.TransportLayer.SrcPort == 53 || packetInfo.TransportLayer.DstPort == 53) &&
			(packetInfo.NetworkLayer.DstIP == "139.155.132.244" || packetInfo.NetworkLayer.SrcIP == "139.155.132.244" ||
				packetInfo.NetworkLayer.DstIP == "192.168.222.127" || packetInfo.NetworkLayer.SrcIP == "192.168.222.127") {
			fmt.Printf("UDP 53端口数据包但无DNS层: SrcIP=%s, DstIP=%s, SrcPort=%d, DstPort=%d\n",
				packetInfo.NetworkLayer.SrcIP, packetInfo.NetworkLayer.DstIP, packetInfo.TransportLayer.SrcPort, packetInfo.TransportLayer.DstPort)
		}
	}
}
