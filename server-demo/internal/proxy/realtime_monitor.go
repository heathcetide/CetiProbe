package proxy

import (
	"sync"
	"time"
)

// RealtimeMonitor 实时监控器
type RealtimeMonitor struct {
	mu sync.RWMutex
	// 实时统计
	stats *RealtimeStats
	// 历史数据
	history    []*RealtimeStats
	maxHistory int
}

// RealtimeStats 实时统计
type RealtimeStats struct {
	Timestamp          time.Time `json:"timestamp"`
	TotalRequests      int64     `json:"total_requests"`
	ActiveRequests     int64     `json:"active_requests"`
	SuccessfulRequests int64     `json:"successful_requests"`
	FailedRequests     int64     `json:"failed_requests"`
	AverageLatency     float64   `json:"average_latency"`
	MaxLatency         int64     `json:"max_latency"`
	MinLatency         int64     `json:"min_latency"`
	TotalBytes         int64     `json:"total_bytes"`
	RequestsPerSecond  float64   `json:"requests_per_second"`
	BytesPerSecond     float64   `json:"bytes_per_second"`
}

// NewRealtimeMonitor 创建实时监控器
func NewRealtimeMonitor() *RealtimeMonitor {
	return &RealtimeMonitor{
		stats: &RealtimeStats{
			Timestamp: time.Now(),
		},
		history:    make([]*RealtimeStats, 0),
		maxHistory: 100, // 保留最近100个时间点的数据
	}
}

// UpdateStats 更新统计信息
func (rm *RealtimeMonitor) UpdateStats(totalRequests, activeRequests, successfulRequests, failedRequests int64,
	latency int64, bytes int64) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	now := time.Now()

	// 如果时间间隔超过1秒，保存历史数据
	if now.Sub(rm.stats.Timestamp) >= time.Second {
		// 保存当前统计到历史
		rm.history = append(rm.history, rm.stats)
		if len(rm.history) > rm.maxHistory {
			rm.history = rm.history[1:]
		}

		// 创建新的统计
		rm.stats = &RealtimeStats{
			Timestamp: now,
		}
	}

	// 更新统计信息
	rm.stats.TotalRequests = totalRequests
	rm.stats.ActiveRequests = activeRequests
	rm.stats.SuccessfulRequests = successfulRequests
	rm.stats.FailedRequests = failedRequests
	rm.stats.TotalBytes += bytes

	// 更新延迟统计
	if rm.stats.MinLatency == 0 || latency < rm.stats.MinLatency {
		rm.stats.MinLatency = latency
	}
	if latency > rm.stats.MaxLatency {
		rm.stats.MaxLatency = latency
	}

	// 计算平均延迟
	if rm.stats.TotalRequests > 0 {
		rm.stats.AverageLatency = float64(rm.stats.TotalRequests-1)*rm.stats.AverageLatency + float64(latency)
		rm.stats.AverageLatency /= float64(rm.stats.TotalRequests)
	}

	// 计算每秒请求数和字节数
	timeDiff := now.Sub(rm.stats.Timestamp).Seconds()
	if timeDiff > 0 {
		rm.stats.RequestsPerSecond = float64(rm.stats.TotalRequests) / timeDiff
		rm.stats.BytesPerSecond = float64(rm.stats.TotalBytes) / timeDiff
	}
}

// GetCurrentStats 获取当前统计信息
func (rm *RealtimeMonitor) GetCurrentStats() *RealtimeStats {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.stats
}

// GetHistory 获取历史数据
func (rm *RealtimeMonitor) GetHistory() []*RealtimeStats {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	// 返回历史数据的副本
	result := make([]*RealtimeStats, len(rm.history))
	copy(result, rm.history)
	return result
}

// GetStatsSummary 获取统计摘要
func (rm *RealtimeMonitor) GetStatsSummary() map[string]interface{} {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return map[string]interface{}{
		"current": rm.stats,
		"history": rm.history,
		"summary": map[string]interface{}{
			"total_requests": rm.stats.TotalRequests,
			"success_rate": func() float64 {
				if rm.stats.TotalRequests > 0 {
					return float64(rm.stats.SuccessfulRequests) / float64(rm.stats.TotalRequests) * 100
				}
				return 0
			}(),
			"average_latency":     rm.stats.AverageLatency,
			"requests_per_second": rm.stats.RequestsPerSecond,
			"bytes_per_second":    rm.stats.BytesPerSecond,
		},
	}
}
