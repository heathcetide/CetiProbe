package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/google/uuid"
	"probe/internal/models"
	"probe/pkg/storage"
)

// EnhancedProxyServer 增强版代理服务器
type EnhancedProxyServer struct {
	addr  string
	https bool
	srv   *http.Server
	mu    sync.Mutex
	store storage.FlowStorage

	// 新增组件
	perfCollector  *PerformanceCollector
	networkMonitor *NetworkMonitor
	errorCollector *ErrorCollector
	geoService     *GeoLocationService
	dnsResolver    *DNSResolver
}

// NewEnhancedProxyServer 创建增强版代理服务器
func NewEnhancedProxyServer(addr string, https bool, store storage.FlowStorage) *EnhancedProxyServer {
	return &EnhancedProxyServer{
		addr:           addr,
		https:          https,
		store:          store,
		perfCollector:  NewPerformanceCollector(),
		networkMonitor: NewNetworkMonitor(),
		errorCollector: NewErrorCollector(),
		geoService:     NewGeoLocationService(),
		dnsResolver:    NewDNSResolver(),
	}
}

// Start 启动增强版代理服务器
func (p *EnhancedProxyServer) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.srv != nil {
		return nil
	}

	gp := goproxy.NewProxyHttpServer()
	gp.Verbose = false

	// 捕获请求
	gp.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		flowID := uuid.NewString()
		ctx.UserData = flowID

		start := time.Now()

		// 开始性能收集
		perfMetrics := p.perfCollector.StartCollecting(flowID)

		// 创建增强的Flow对象
		flow := &models.Flow{
			ID:          flowID,
			Scheme:      req.URL.Scheme,
			RemoteAddr:  req.RemoteAddr,
			StartAt:     start,
			Request:     p.buildHTTPRequest(req),
			Performance: perfMetrics,
			Network:     p.buildNetworkInfo(req),
		}

		// 读取请求体
		var reqBody []byte
		if req.Body != nil {
			reqBody, _ = io.ReadAll(req.Body)
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewReader(reqBody))
			flow.Request.Body = reqBody
		}

		// 分析请求内容
		flow.Content = p.analyzeContent(reqBody, req.Header)

		// 异步收集DNS解析时间
		go p.collectDNSMetrics(flowID, req.Host)

		// 记录网络统计
		p.networkMonitor.RecordRequest(req.Host, true, 0, int64(len(reqBody)))

		p.store.Add(flow)
		return req, nil
	})

	// 捕获响应
	gp.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		flowID, _ := ctx.UserData.(string)
		if flowID == "" {
			return resp
		}

		flow := p.store.GetByID(flowID)
		if flow == nil {
			return resp
		}

		end := time.Now()
		flow.EndAt = end
		flow.LatencyMs = end.Sub(flow.StartAt).Milliseconds()

		// 更新性能指标
		if flow.Performance != nil {
			flow.Performance.TotalTime = flow.LatencyMs
			// 从性能收集器获取最新的指标
			if latestPerf := p.perfCollector.GetMetrics(flowID); latestPerf != nil {
				flow.Performance = latestPerf
				flow.Performance.TotalTime = flow.LatencyMs
			}
		}

		// 读取响应体
		var body []byte
		if resp != nil && resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(body))
		}

		// 构建响应
		if resp != nil {
			flow.Response = &models.HTTPResponse{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				Headers:    models.CopyHeaders(resp.Header),
				Body:       body,
				Proto:      resp.Proto,
				Length:     int(resp.ContentLength),
			}

			// 分析响应内容
			flow.Content = p.analyzeContent(body, resp.Header)
		}

		// 收集TLS信息
		if flow.Scheme == "https" {
			flow.TLS = p.collectTLSInfo(resp)
		}

		// 更新网络统计
		success := resp != nil && resp.StatusCode < 400
		p.networkMonitor.RecordRequest(flow.Request.Host, success, flow.LatencyMs, int64(len(body)))

		// 清理性能数据
		p.CleanupFlow(flowID)

		return resp
	})

	// HTTPS MITM（自签证书）
	if p.https {
		certPEM, keyPEM, _, err := EnsureCAExists()
		if err != nil {
			return err
		}
		caCert, caKey, err := ParseCA(certPEM, keyPEM)
		if err != nil {
			return err
		}

		tlsConfigForHost := func(host string, ctx *goproxy.ProxyCtx) (*tls.Config, error) {
			h := host
			if i := strings.IndexByte(h, ':'); i >= 0 {
				h = h[:i]
			}
			leafCertPEM, leafKeyPEM, err := SignHostCert(caCert, caKey, h, 24*time.Hour)
			if err != nil {
				return nil, err
			}
			pair, err := tls.X509KeyPair(leafCertPEM, leafKeyPEM)
			if err != nil {
				return nil, err
			}
			return &tls.Config{
				Certificates: []tls.Certificate{pair},
				MinVersion:   tls.VersionTLS12,
				ServerName:   h,
			}, nil
		}

		gp.Tr = &http.Transport{
			TLSClientConfig:   &tls.Config{MinVersion: tls.VersionTLS12},
			Proxy:             http.ProxyFromEnvironment,
			DialContext:       (&net.Dialer{Timeout: 30 * time.Second}).DialContext,
			ForceAttemptHTTP2: true,
		}

		gp.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			return &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: tlsConfigForHost}, host
		}))
	}

	p.srv = &http.Server{Addr: p.addr, Handler: gp}
	go p.srv.ListenAndServe()
	return nil
}

// buildHTTPRequest 构建HTTP请求
func (p *EnhancedProxyServer) buildHTTPRequest(req *http.Request) *models.HTTPRequest {
	return &models.HTTPRequest{
		Method:  req.Method,
		URL:     req.URL.String(),
		Path:    req.URL.Path,
		Query:   req.URL.RawQuery,
		Host:    req.Host,
		Headers: models.CopyHeaders(req.Header),
		Proto:   req.Proto,
		Length:  int(req.ContentLength),
	}
}

// buildNetworkInfo 构建网络信息
func (p *EnhancedProxyServer) buildNetworkInfo(req *http.Request) *models.NetworkInfo {
	clientIP, clientPort, _ := net.SplitHostPort(req.RemoteAddr)
	serverIP, serverPort, _ := net.SplitHostPort(req.Host)

	clientPortInt := 0
	serverPortInt := 0
	if clientPort != "" {
		fmt.Sscanf(clientPort, "%d", &clientPortInt)
	}
	if serverPort != "" {
		fmt.Sscanf(serverPort, "%d", &serverPortInt)
	}

	// 获取地理位置信息
	country, region, city, isp, asn := p.geoService.GetLocationInfo(clientIP)

	return &models.NetworkInfo{
		ClientIP:    clientIP,
		ServerIP:    serverIP,
		ClientPort:  clientPortInt,
		ServerPort:  serverPortInt,
		IsIPv6:      strings.Contains(clientIP, ":"),
		IsLocalhost: clientIP == "127.0.0.1" || clientIP == "::1",
		IsPrivate:   p.isPrivateIP(clientIP),
		Country:     country,
		Region:      region,
		City:        city,
		ISP:         isp,
		ASN:         asn,
	}
}

// analyzeContent 分析内容
func (p *EnhancedProxyServer) analyzeContent(body []byte, headers http.Header) *models.ContentInfo {
	contentInfo := &models.ContentInfo{
		OriginalSize: int64(len(body)),
	}

	// 分析MIME类型
	if contentType := headers.Get("Content-Type"); contentType != "" {
		contentInfo.MIMEType = strings.Split(contentType, ";")[0]
	} else {
		contentInfo.MIMEType = p.detectMIMEType(body)
	}

	// 分析编码
	if encoding := headers.Get("Content-Encoding"); encoding != "" {
		contentInfo.Compression = encoding
		contentInfo.IsCompressed = true
	}

	// 分析字符编码
	if contentType := headers.Get("Content-Type"); contentType != "" {
		if strings.Contains(contentType, "charset=") {
			parts := strings.Split(contentType, "charset=")
			if len(parts) > 1 {
				contentInfo.Encoding = strings.TrimSpace(parts[1])
			}
		}
	}

	// 内容类型检测
	contentInfo.IsText = p.isTextContent(contentInfo.MIMEType)
	contentInfo.IsJSON = p.isJSONContent(body)
	contentInfo.IsXML = p.isXMLContent(body)
	contentInfo.IsImage = p.isImageContent(contentInfo.MIMEType)
	contentInfo.IsVideo = p.isVideoContent(contentInfo.MIMEType)
	contentInfo.IsAudio = p.isAudioContent(contentInfo.MIMEType)

	return contentInfo
}

// collectTLSInfo 收集TLS信息
func (p *EnhancedProxyServer) collectTLSInfo(resp *http.Response) *models.TLSInfo {
	if resp == nil || resp.TLS == nil {
		return nil
	}

	tlsInfo := &models.TLSInfo{
		Version:     p.getTLSVersion(resp.TLS.Version),
		CipherSuite: p.getCipherSuite(resp.TLS.CipherSuite),
		IsSecure:    true,
		Protocol:    "TLS",
	}

	// 收集证书信息
	if len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		tlsInfo.Certificate = &models.CertificateInfo{
			Subject:      cert.Subject.String(),
			Issuer:       cert.Issuer.String(),
			NotBefore:    cert.NotBefore,
			NotAfter:     cert.NotAfter,
			SerialNumber: cert.SerialNumber.String(),
			Fingerprint:  fmt.Sprintf("%x", cert.Signature),
			IsSelfSigned: cert.Issuer.String() == cert.Subject.String(),
			IsValid:      time.Now().After(cert.NotBefore) && time.Now().Before(cert.NotAfter),
		}
	}

	return tlsInfo
}

// 辅助方法
func (p *EnhancedProxyServer) isPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.IsPrivate()
}

func (p *EnhancedProxyServer) detectMIMEType(body []byte) string {
	if len(body) == 0 {
		return "application/octet-stream"
	}

	// 简单的MIME类型检测
	if len(body) > 4 && body[0] == 0x89 && body[1] == 0x50 && body[2] == 0x4E && body[3] == 0x47 {
		return "image/png"
	}
	if len(body) > 2 && body[0] == 0xFF && body[1] == 0xD8 {
		return "image/jpeg"
	}
	if len(body) > 4 && string(body[:4]) == "<!DOCTYPE" {
		return "text/html"
	}
	if len(body) > 1 && body[0] == '{' {
		return "application/json"
	}
	if len(body) > 1 && body[0] == '<' {
		return "text/xml"
	}

	return "application/octet-stream"
}

func (p *EnhancedProxyServer) isTextContent(mimeType string) bool {
	return strings.HasPrefix(mimeType, "text/") ||
		strings.Contains(mimeType, "json") ||
		strings.Contains(mimeType, "xml")
}

func (p *EnhancedProxyServer) isJSONContent(body []byte) bool {
	if len(body) == 0 {
		return false
	}
	var jsonData interface{}
	return json.Unmarshal(body, &jsonData) == nil
}

func (p *EnhancedProxyServer) isXMLContent(body []byte) bool {
	if len(body) == 0 {
		return false
	}
	bodyStr := strings.TrimSpace(string(body))
	return strings.HasPrefix(bodyStr, "<") && strings.HasSuffix(bodyStr, ">")
}

func (p *EnhancedProxyServer) isImageContent(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

func (p *EnhancedProxyServer) isVideoContent(mimeType string) bool {
	return strings.HasPrefix(mimeType, "video/")
}

func (p *EnhancedProxyServer) isAudioContent(mimeType string) bool {
	return strings.HasPrefix(mimeType, "audio/")
}

func (p *EnhancedProxyServer) getTLSVersion(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

func (p *EnhancedProxyServer) getCipherSuite(suite uint16) string {
	// 简化的TLS加密套件检测
	switch suite {
	case 0x0005:
		return "TLS_RSA_WITH_RC4_128_SHA"
	case 0x000A:
		return "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
	case 0x002F:
		return "TLS_RSA_WITH_AES_128_CBC_SHA"
	case 0x0035:
		return "TLS_RSA_WITH_AES_256_CBC_SHA"
	case 0x003C:
		return "TLS_RSA_WITH_AES_128_CBC_SHA256"
	case 0x003D:
		return "TLS_RSA_WITH_AES_256_CBC_SHA256"
	case 0xC011:
		return "TLS_ECDHE_RSA_WITH_RC4_128_SHA"
	case 0xC012:
		return "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA"
	case 0xC013:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	case 0xC014:
		return "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	case 0xC027:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256"
	case 0xC028:
		return "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384"
	case 0xCCA8:
		return "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"
	case 0xC007:
		return "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA"
	case 0xC008:
		return "TLS_ECDHE_ECDSA_WITH_3DES_EDE_CBC_SHA"
	case 0xC009:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
	case 0xC00A:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
	case 0xC023:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256"
	case 0xC024:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384"
	case 0xCCA9:
		return "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305"
	default:
		return fmt.Sprintf("Unknown (0x%04x)", suite)
	}
}

// Stop 停止代理服务器
func (p *EnhancedProxyServer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.srv == nil {
		return nil
	}
	err := p.srv.Close()
	p.srv = nil
	return err
}

// IsRunning 检查代理服务器是否运行
func (p *EnhancedProxyServer) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.srv != nil
}

// collectDNSMetrics 收集DNS指标
func (p *EnhancedProxyServer) collectDNSMetrics(flowID, host string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, duration, err := p.dnsResolver.ResolveWithTiming(ctx, host)
	if err != nil {
		p.errorCollector.RecordError(flowID, "dns", err.Error(), 0, false, 0)
	} else {
		p.perfCollector.RecordDNSLookup(flowID, duration)
	}
}

// GetNetworkStats 获取网络统计信息
func (p *EnhancedProxyServer) GetNetworkStats() map[string]*NetworkStats {
	return p.networkMonitor.GetAllStats()
}

// GetPerformanceStats 获取性能统计信息
func (p *EnhancedProxyServer) GetPerformanceStats() map[string]*models.PerformanceMetrics {
	return p.perfCollector.metrics
}

// CleanupFlow 清理Flow相关数据
func (p *EnhancedProxyServer) CleanupFlow(flowID string) {
	p.perfCollector.Cleanup(flowID)
	p.errorCollector.Cleanup(flowID)
}
