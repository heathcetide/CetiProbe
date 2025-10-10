package proxy

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"probe/internal/models"
	"probe/pkg/storage"
)

// DataExporter 数据导出器
type DataExporter struct {
	store storage.FlowStorage
}

// NewDataExporter 创建数据导出器
func NewDataExporter(store storage.FlowStorage) *DataExporter {
	return &DataExporter{store: store}
}

// ExportToJSON 导出为JSON格式
func (de *DataExporter) ExportToJSON(filename string) error {
	flows := de.store.GetAll()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(flows)
}

// ExportToCSV 导出为CSV格式
func (de *DataExporter) ExportToCSV(filename string) error {
	flows := de.store.GetAll()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入CSV头部
	header := []string{
		"ID", "Scheme", "RemoteAddr", "StartAt", "EndAt", "LatencyMs",
		"Method", "URL", "Host", "Path", "Query",
		"StatusCode", "Status", "ResponseLength",
		"DNSLookupTime", "TCPConnectTime", "TLSHandshakeTime", "TTFB", "ContentTransferTime", "TotalTime",
		"TLSVersion", "CipherSuite", "IsSecure",
		"ClientIP", "ServerIP", "Country", "Region", "City", "ISP",
		"MIMEType", "Encoding", "Compression", "IsCompressed", "IsText", "IsJSON", "IsImage",
		"ErrorType", "ErrorMessage", "IsTimeout", "RetryCount",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	// 写入数据行
	for _, flow := range flows {
		record := de.flowToCSVRecord(flow)
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// flowToCSVRecord 将Flow转换为CSV记录
func (de *DataExporter) flowToCSVRecord(flow *models.Flow) []string {
	record := make([]string, 0, 30)

	// 基本信息
	record = append(record, flow.ID, flow.Scheme, flow.RemoteAddr,
		flow.StartAt.Format(time.RFC3339), flow.EndAt.Format(time.RFC3339),
		fmt.Sprintf("%d", flow.LatencyMs))

	// 请求信息
	if flow.Request != nil {
		record = append(record, flow.Request.Method, flow.Request.URL,
			flow.Request.Host, flow.Request.Path, flow.Request.Query)
	} else {
		record = append(record, "", "", "", "", "")
	}

	// 响应信息
	if flow.Response != nil {
		record = append(record, fmt.Sprintf("%d", flow.Response.StatusCode),
			flow.Response.Status, fmt.Sprintf("%d", flow.Response.Length))
	} else {
		record = append(record, "", "", "")
	}

	// 性能指标
	if flow.Performance != nil {
		record = append(record,
			fmt.Sprintf("%d", flow.Performance.DNSLookupTime),
			fmt.Sprintf("%d", flow.Performance.TCPConnectTime),
			fmt.Sprintf("%d", flow.Performance.TLSHandshakeTime),
			fmt.Sprintf("%d", flow.Performance.TTFB),
			fmt.Sprintf("%d", flow.Performance.ContentTransferTime),
			fmt.Sprintf("%d", flow.Performance.TotalTime))
	} else {
		record = append(record, "", "", "", "", "", "")
	}

	// TLS信息
	if flow.TLS != nil {
		record = append(record, flow.TLS.Version, flow.TLS.CipherSuite,
			fmt.Sprintf("%t", flow.TLS.IsSecure))
	} else {
		record = append(record, "", "", "")
	}

	// 网络信息
	if flow.Network != nil {
		record = append(record, flow.Network.ClientIP, flow.Network.ServerIP,
			flow.Network.Country, flow.Network.Region, flow.Network.City, flow.Network.ISP)
	} else {
		record = append(record, "", "", "", "", "", "")
	}

	// 内容信息
	if flow.Content != nil {
		record = append(record, flow.Content.MIMEType, flow.Content.Encoding,
			flow.Content.Compression, fmt.Sprintf("%t", flow.Content.IsCompressed),
			fmt.Sprintf("%t", flow.Content.IsText), fmt.Sprintf("%t", flow.Content.IsJSON),
			fmt.Sprintf("%t", flow.Content.IsImage))
	} else {
		record = append(record, "", "", "", "", "", "", "")
	}

	// 错误信息
	if flow.Error != nil {
		record = append(record, flow.Error.Type, flow.Error.Message,
			fmt.Sprintf("%t", flow.Error.IsTimeout), fmt.Sprintf("%d", flow.Error.RetryCount))
	} else {
		record = append(record, "", "", "", "")
	}

	return record
}

// ExportPerformanceReport 导出性能报告
func (de *DataExporter) ExportPerformanceReport(filename string) error {
	flows := de.store.GetAll()

	// 计算性能统计
	stats := de.calculatePerformanceStats(flows)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(stats)
}

// calculatePerformanceStats 计算性能统计
func (de *DataExporter) calculatePerformanceStats(flows []*models.Flow) map[string]interface{} {
	stats := map[string]interface{}{
		"total_flows":         len(flows),
		"performance_metrics": map[string]interface{}{},
		"tls_stats":           map[string]interface{}{},
		"content_stats":       map[string]interface{}{},
		"error_stats":         map[string]interface{}{},
	}

	// 性能指标统计
	var totalLatency, dnsTime, tcpTime, tlsTime, ttfb, contentTime int64
	var latencyCount, dnsCount, tcpCount, tlsCount, ttfbCount, contentCount int

	for _, flow := range flows {
		if flow.LatencyMs > 0 {
			totalLatency += flow.LatencyMs
			latencyCount++
		}

		if flow.Performance != nil {
			if flow.Performance.DNSLookupTime > 0 {
				dnsTime += flow.Performance.DNSLookupTime
				dnsCount++
			}
			if flow.Performance.TCPConnectTime > 0 {
				tcpTime += flow.Performance.TCPConnectTime
				tcpCount++
			}
			if flow.Performance.TLSHandshakeTime > 0 {
				tlsTime += flow.Performance.TLSHandshakeTime
				tlsCount++
			}
			if flow.Performance.TTFB > 0 {
				ttfb += flow.Performance.TTFB
				ttfbCount++
			}
			if flow.Performance.ContentTransferTime > 0 {
				contentTime += flow.Performance.ContentTransferTime
				contentCount++
			}
		}
	}

	// 计算平均值
	if latencyCount > 0 {
		stats["performance_metrics"].(map[string]interface{})["average_latency"] = float64(totalLatency) / float64(latencyCount)
	}
	if dnsCount > 0 {
		stats["performance_metrics"].(map[string]interface{})["average_dns_time"] = float64(dnsTime) / float64(dnsCount)
	}
	if tcpCount > 0 {
		stats["performance_metrics"].(map[string]interface{})["average_tcp_time"] = float64(tcpTime) / float64(tcpCount)
	}
	if tlsCount > 0 {
		stats["performance_metrics"].(map[string]interface{})["average_tls_time"] = float64(tlsTime) / float64(tlsCount)
	}
	if ttfbCount > 0 {
		stats["performance_metrics"].(map[string]interface{})["average_ttfb"] = float64(ttfb) / float64(ttfbCount)
	}
	if contentCount > 0 {
		stats["performance_metrics"].(map[string]interface{})["average_content_time"] = float64(contentTime) / float64(contentCount)
	}

	return stats
}
