package capture

import (
	"fmt"
	"probe/internal/capture/layer"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"golang.org/x/mod/sumdb/storage"
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
		printPacketDetails(packet)
	}

	return nil
}

// printPacketDetails 分层输出数据包的详细信息
func printPacketDetails(packet gopacket.Packet) {
	// 输出数据包元信息
	metaInfo := layer.ExtractPacketMetadataInfo(packet)
	layer.PrintPacketMetadataInfo(metaInfo)

	// 处理链路层
	if linkLayer := packet.LinkLayer(); linkLayer != nil {
		layer.ExtractLinkLayerInfo(linkLayer)
	}

	// 特别处理网络层，尤其是IPv6
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		layer.ExtractNetworkLayerInfo(networkLayer)
	}

	// 处理传输层
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		layer.ExtractTransportLayerInfo(transportLayer)
	}

	// 处理应用层
	if appLayer := packet.ApplicationLayer(); appLayer != nil {
		layer.ExtractApplicationLayerInfo(appLayer)
	}

	// 处理错误层
	if errLayer := packet.ErrorLayer(); errLayer != nil {
		errInfo := layer.ExtractErrorLayerInfo(errLayer)
		layer.PrintErrorLayerDetails(errInfo)
	}

	fmt.Println() // 添加空行分隔不同数据包
}
