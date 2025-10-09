package models

import (
	"probe/internal/capture/layer"
)

// PacketInfo 表示网络数据包的信息
type PacketInfo struct {
	ID               int64                       `json:"id"` // 数据包唯一标识
	Metadata         *layer.PacketMetadataInfo   `json:"metadata"`
	LinkLayer        *layer.LinkLayerInfo        `json:"linkLayer"`
	NetworkLayer     *layer.NetworkLayerInfo     `json:"networkLayer"`
	TransportLayer   *layer.TransportLayerInfo   `json:"transportLayer"`
	ApplicationLayer *layer.ApplicationLayerInfo `json:"applicationLayer"`
	ErrorLayer       *layer.ErrorLayerInfo       `json:"errorLayer"`
}

// ToString 返回数据包的字符串表示，调用各层的打印方法
func (p *PacketInfo) ToString() string {

	if p.Metadata != nil {
		layer.PrintPacketMetadataInfo(p.Metadata)
	}

	if p.LinkLayer != nil {
		layer.PrintLinkLayerInfo(p.LinkLayer)
	}

	if p.NetworkLayer != nil {
		layer.PrintNetworkLayerInfo(p.NetworkLayer)
	}

	if p.TransportLayer != nil {
		layer.PrintTransportLayerDetails(p.TransportLayer)
	}

	if p.ApplicationLayer != nil {
		layer.PrintApplicationLayerDetails(p.ApplicationLayer)
	}

	if p.ErrorLayer != nil {
		layer.PrintErrorLayerDetails(p.ErrorLayer)
	}

	return ""
}
