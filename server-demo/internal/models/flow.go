package models

import (
	"net/http"
	"time"
)

// HTTPMessage 表示一次HTTP消息的通用部分
type HTTPMessage struct {
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
	Proto   string            `json:"proto"`
	Length  int               `json:"length"`
}

// HTTPRequest 表示HTTP请求
type HTTPRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Path    string            `json:"path"`
	Query   string            `json:"query"`
	Host    string            `json:"host"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
	Proto   string            `json:"proto"`
	Length  int               `json:"length"`
}

// HTTPResponse 表示HTTP响应
type HTTPResponse struct {
	Status     string            `json:"status"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
	Proto      string            `json:"proto"`
	Length     int               `json:"length"`
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	DNSLookupTime       int64 `json:"dns_lookup_time"`       // DNS解析时间(ms)
	TCPConnectTime      int64 `json:"tcp_connect_time"`      // TCP连接时间(ms)
	TLSHandshakeTime    int64 `json:"tls_handshake_time"`    // TLS握手时间(ms)
	TTFB                int64 `json:"ttfb"`                  // 首字节时间(ms)
	ContentTransferTime int64 `json:"content_transfer_time"` // 内容传输时间(ms)
	TotalTime           int64 `json:"total_time"`            // 总时间(ms)
}

// TLSInfo TLS信息
type TLSInfo struct {
	Version     string           `json:"version"`      // TLS版本
	CipherSuite string           `json:"cipher_suite"` // 加密套件
	Certificate *CertificateInfo `json:"certificate"`  // 证书信息
	IsSecure    bool             `json:"is_secure"`    // 是否安全连接
	Protocol    string           `json:"protocol"`     // 协议类型
}

// CertificateInfo 证书信息
type CertificateInfo struct {
	Subject         string    `json:"subject"`          // 证书主题
	Issuer          string    `json:"issuer"`           // 证书颁发者
	NotBefore       time.Time `json:"not_before"`       // 有效期开始
	NotAfter        time.Time `json:"not_after"`        // 有效期结束
	SerialNumber    string    `json:"serial_number"`    // 序列号
	Fingerprint     string    `json:"fingerprint"`      // 指纹
	IsSelfSigned    bool      `json:"is_self_signed"`   // 是否自签名
	IsValid         bool      `json:"is_valid"`         // 是否有效
	ValidationError string    `json:"validation_error"` // 验证错误
}

// ErrorInfo 错误信息
type ErrorInfo struct {
	Type       string `json:"type"`        // 错误类型
	Message    string `json:"message"`     // 错误消息
	Code       int    `json:"code"`        // 错误代码
	IsTimeout  bool   `json:"is_timeout"`  // 是否超时
	RetryCount int    `json:"retry_count"` // 重试次数
	IsNetwork  bool   `json:"is_network"`  // 是否网络错误
	IsDNS      bool   `json:"is_dns"`      // 是否DNS错误
	IsTLS      bool   `json:"is_tls"`      // 是否TLS错误
}

// ContentInfo 内容信息
type ContentInfo struct {
	MIMEType         string  `json:"mime_type"`         // MIME类型
	Encoding         string  `json:"encoding"`          // 字符编码
	Compression      string  `json:"compression"`       // 压缩算法
	IsCompressed     bool    `json:"is_compressed"`     // 是否压缩
	OriginalSize     int64   `json:"original_size"`     // 原始大小
	CompressedSize   int64   `json:"compressed_size"`   // 压缩后大小
	CompressionRatio float64 `json:"compression_ratio"` // 压缩比
	IsText           bool    `json:"is_text"`           // 是否文本
	IsJSON           bool    `json:"is_json"`           // 是否JSON
	IsXML            bool    `json:"is_xml"`            // 是否XML
	IsImage          bool    `json:"is_image"`          // 是否图片
	IsVideo          bool    `json:"is_video"`          // 是否视频
	IsAudio          bool    `json:"is_audio"`          // 是否音频
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	ClientIP    string `json:"client_ip"`    // 客户端IP
	ServerIP    string `json:"server_ip"`    // 服务器IP
	ClientPort  int    `json:"client_port"`  // 客户端端口
	ServerPort  int    `json:"server_port"`  // 服务器端口
	IsIPv6      bool   `json:"is_ipv6"`      // 是否IPv6
	IsLocalhost bool   `json:"is_localhost"` // 是否本地
	IsPrivate   bool   `json:"is_private"`   // 是否私有IP
	Country     string `json:"country"`      // 国家
	Region      string `json:"region"`       // 地区
	City        string `json:"city"`         // 城市
	ISP         string `json:"isp"`          // ISP
	ASN         string `json:"asn"`          // ASN
}

// Flow 表示一次完整的请求-响应流
type Flow struct {
	ID         string        `json:"id"`
	Scheme     string        `json:"scheme"`
	RemoteAddr string        `json:"remote_addr"`
	StartAt    time.Time     `json:"start_at"`
	EndAt      time.Time     `json:"end_at"`
	LatencyMs  int64         `json:"latency_ms"`
	Request    *HTTPRequest  `json:"request"`
	Response   *HTTPResponse `json:"response"`

	// 新增字段
	Performance *PerformanceMetrics `json:"performance,omitempty"`
	TLS         *TLSInfo            `json:"tls,omitempty"`
	Error       *ErrorInfo          `json:"error,omitempty"`
	Content     *ContentInfo        `json:"content,omitempty"`
	Network     *NetworkInfo        `json:"network,omitempty"`
}

// CopyHeaders 将http.Header转换为map[string]string（首值）
func CopyHeaders(h http.Header) map[string]string {
	m := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}
