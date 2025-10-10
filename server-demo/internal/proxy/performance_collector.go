package proxy

import (
	"context"
	"net"
	"sync"
	"time"

	"probe/internal/models"
)

// PerformanceCollector 性能指标收集器
type PerformanceCollector struct {
	mu sync.RWMutex
	// 存储每个请求的性能指标
	metrics map[string]*models.PerformanceMetrics
}

// NewPerformanceCollector 创建性能收集器
func NewPerformanceCollector() *PerformanceCollector {
	return &PerformanceCollector{
		metrics: make(map[string]*models.PerformanceMetrics),
	}
}

// StartCollecting 开始收集性能指标
func (pc *PerformanceCollector) StartCollecting(flowID string) *models.PerformanceMetrics {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	metrics := &models.PerformanceMetrics{
		DNSLookupTime:       -1, // -1表示未测量
		TCPConnectTime:      -1,
		TLSHandshakeTime:    -1,
		TTFB:                -1,
		ContentTransferTime: -1,
		TotalTime:           -1,
	}

	pc.metrics[flowID] = metrics
	return metrics
}

// RecordDNSLookup 记录DNS解析时间
func (pc *PerformanceCollector) RecordDNSLookup(flowID string, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if metrics, exists := pc.metrics[flowID]; exists {
		metrics.DNSLookupTime = duration.Milliseconds()
	}
}

// RecordTCPConnect 记录TCP连接时间
func (pc *PerformanceCollector) RecordTCPConnect(flowID string, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if metrics, exists := pc.metrics[flowID]; exists {
		metrics.TCPConnectTime = duration.Milliseconds()
	}
}

// RecordTLSHandshake 记录TLS握手时间
func (pc *PerformanceCollector) RecordTLSHandshake(flowID string, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if metrics, exists := pc.metrics[flowID]; exists {
		metrics.TLSHandshakeTime = duration.Milliseconds()
	}
}

// RecordTTFB 记录首字节时间
func (pc *PerformanceCollector) RecordTTFB(flowID string, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if metrics, exists := pc.metrics[flowID]; exists {
		metrics.TTFB = duration.Milliseconds()
	}
}

// RecordContentTransfer 记录内容传输时间
func (pc *PerformanceCollector) RecordContentTransfer(flowID string, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if metrics, exists := pc.metrics[flowID]; exists {
		metrics.ContentTransferTime = duration.Milliseconds()
	}
}

// RecordTotalTime 记录总时间
func (pc *PerformanceCollector) RecordTotalTime(flowID string, duration time.Duration) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if metrics, exists := pc.metrics[flowID]; exists {
		metrics.TotalTime = duration.Milliseconds()
	}
}

// GetMetrics 获取性能指标
func (pc *PerformanceCollector) GetMetrics(flowID string) *models.PerformanceMetrics {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	return pc.metrics[flowID]
}

// Cleanup 清理性能指标
func (pc *PerformanceCollector) Cleanup(flowID string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	delete(pc.metrics, flowID)
}

// DNSResolver DNS解析器
type DNSResolver struct {
	resolver *net.Resolver
}

// NewDNSResolver 创建DNS解析器
func NewDNSResolver() *DNSResolver {
	return &DNSResolver{
		resolver: &net.Resolver{
			PreferGo: true,
		},
	}
}

// ResolveWithTiming 带时间测量的DNS解析
func (dr *DNSResolver) ResolveWithTiming(ctx context.Context, host string) ([]string, time.Duration, error) {
	start := time.Now()
	ips, err := dr.resolver.LookupHost(ctx, host)
	duration := time.Since(start)
	return ips, duration, err
}

// NetworkMonitor 网络监控器
type NetworkMonitor struct {
	mu    sync.RWMutex
	stats map[string]*NetworkStats
}

// NetworkStats 网络统计
type NetworkStats struct {
	TotalRequests      int64 `json:"total_requests"`
	SuccessfulRequests int64 `json:"successful_requests"`
	FailedRequests     int64 `json:"failed_requests"`
	AverageLatency     int64 `json:"average_latency"`
	MaxLatency         int64 `json:"max_latency"`
	MinLatency         int64 `json:"min_latency"`
	TotalBytes         int64 `json:"total_bytes"`
	AverageBytes       int64 `json:"average_bytes"`
}

// NewNetworkMonitor 创建网络监控器
func NewNetworkMonitor() *NetworkMonitor {
	return &NetworkMonitor{
		stats: make(map[string]*NetworkStats),
	}
}

// RecordRequest 记录请求
func (nm *NetworkMonitor) RecordRequest(host string, success bool, latency int64, bytes int64) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	stats, exists := nm.stats[host]
	if !exists {
		stats = &NetworkStats{}
		nm.stats[host] = stats
	}

	stats.TotalRequests++
	if success {
		stats.SuccessfulRequests++
	} else {
		stats.FailedRequests++
	}

	// 更新延迟统计
	if stats.MinLatency == 0 || latency < stats.MinLatency {
		stats.MinLatency = latency
	}
	if latency > stats.MaxLatency {
		stats.MaxLatency = latency
	}

	// 计算平均延迟
	stats.AverageLatency = (stats.AverageLatency*(stats.TotalRequests-1) + latency) / stats.TotalRequests

	// 更新字节统计
	stats.TotalBytes += bytes
	stats.AverageBytes = stats.TotalBytes / stats.TotalRequests
}

// GetStats 获取统计信息
func (nm *NetworkMonitor) GetStats(host string) *NetworkStats {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	return nm.stats[host]
}

// GetAllStats 获取所有统计信息
func (nm *NetworkMonitor) GetAllStats() map[string]*NetworkStats {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	result := make(map[string]*NetworkStats)
	for host, stats := range nm.stats {
		result[host] = stats
	}
	return result
}

// ErrorCollector 错误收集器
type ErrorCollector struct {
	mu     sync.RWMutex
	errors map[string]*models.ErrorInfo
}

// NewErrorCollector 创建错误收集器
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make(map[string]*models.ErrorInfo),
	}
}

// RecordError 记录错误
func (ec *ErrorCollector) RecordError(flowID string, errType string, message string, code int, isTimeout bool, retryCount int) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.errors[flowID] = &models.ErrorInfo{
		Type:       errType,
		Message:    message,
		Code:       code,
		IsTimeout:  isTimeout,
		RetryCount: retryCount,
		IsNetwork:  errType == "network",
		IsDNS:      errType == "dns",
		IsTLS:      errType == "tls",
	}
}

// GetError 获取错误信息
func (ec *ErrorCollector) GetError(flowID string) *models.ErrorInfo {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	return ec.errors[flowID]
}

// Cleanup 清理错误信息
func (ec *ErrorCollector) Cleanup(flowID string) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	delete(ec.errors, flowID)
}

// GeoLocationService 地理位置服务
type GeoLocationService struct {
	// 这里可以集成真实的地理位置API
	// 目前使用模拟数据
}

// NewGeoLocationService 创建地理位置服务
func NewGeoLocationService() *GeoLocationService {
	return &GeoLocationService{}
}

// GetLocationInfo 获取地理位置信息
func (gls *GeoLocationService) GetLocationInfo(ip string) (country, region, city, isp, asn string) {
	// 模拟地理位置查询
	// 在实际应用中，这里应该调用真实的地理位置API
	if ip == "127.0.0.1" || ip == "::1" {
		return "本地", "本地", "本地", "本地网络", "本地ASN"
	}

	// 简单的IP分类
	if isPrivateIP(ip) {
		return "私有网络", "私有网络", "私有网络", "私有网络", "私有ASN"
	}

	// 模拟一些常见IP的地理位置
	switch {
	case ip == "8.8.8.8":
		return "美国", "加利福尼亚", "山景城", "Google", "AS15169"
	case ip == "1.1.1.1":
		return "美国", "加利福尼亚", "旧金山", "Cloudflare", "AS13335"
	default:
		return "未知", "未知", "未知", "未知", "未知"
	}
}

// isPrivateIP 检查是否为私有IP
func isPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.IsPrivate()
}
