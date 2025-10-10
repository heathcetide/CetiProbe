package layer

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
)

// ApplicationLayerInfo 存储应用层信息
type ApplicationLayerInfo struct {
	// 基本信息
	Timestamp time.Time `json:"timestamp"` // 数据包捕获时间

	// 应用层数据
	Payload     string            `json:"payload,omitempty"`      // 原始载荷数据 (base64编码)
	HTTPMethod  string            `json:"http_method,omitempty"`  // HTTP方法 (GET, POST等)
	HTTPVersion string            `json:"http_version,omitempty"` // HTTP版本
	HTTPStatus  string            `json:"http_status,omitempty"`  // HTTP状态码文本
	StatusCode  int               `json:"status_code,omitempty"`  // HTTP状态码
	RequestURI  string            `json:"request_uri,omitempty"`  // 请求URI
	Headers     map[string]string `json:"headers,omitempty"`      // HTTP头部字段
	Body        []byte            `json:"body,omitempty"`         // HTTP消息体

	// 常用HTTP头部字段
	Host           string `json:"host,omitempty"`            // Host头部
	UserAgent      string `json:"user_agent,omitempty"`      // User-Agent头部
	ContentType    string `json:"content_type,omitempty"`    // Content-Type头部
	ContentLength  int    `json:"content_length,omitempty"`  // Content-Length头部
	Authorization  string `json:"authorization,omitempty"`   // Authorization头部
	Referer        string `json:"referer,omitempty"`         // Referer头部
	Server         string `json:"server,omitempty"`          // Server头部
	Cookie         string `json:"cookie,omitempty"`          // Cookie头部
	SetCookie      string `json:"set_cookie,omitempty"`      // Set-Cookie头部
	Accept         string `json:"accept,omitempty"`          // Accept头部
	AcceptLanguage string `json:"accept_language,omitempty"` // Accept-Language头部
	AcceptEncoding string `json:"accept_encoding,omitempty"` // Accept-Encoding头部
	Connection     string `json:"connection,omitempty"`      // Connection头部

	// URL解析字段
	Domain  string `json:"domain,omitempty"`   // 域名
	Path    string `json:"path,omitempty"`     // 路径
	Query   string `json:"query,omitempty"`    // 查询参数
	FullURL string `json:"full_url,omitempty"` // 完整URL
}

// ExtractApplicationLayerInfo 提取应用层信息并填充到ApplicationLayerInfo结构体中
func ExtractApplicationLayerInfo(appLayer gopacket.ApplicationLayer) *ApplicationLayerInfo {
	info := &ApplicationLayerInfo{
		Timestamp: time.Now(),
		Headers:   make(map[string]string),
	}

	if appLayer == nil {
		return info
	}

	// 获取原始载荷数据并转换为base64编码
	payload := appLayer.Payload()
	info.Payload = base64.StdEncoding.EncodeToString(payload)

	// 尝试解析HTTP数据
	if len(payload) > 0 {
		// 使用net/http库解析HTTP请求或响应
		reader := bytes.NewReader(payload)

		// 首先尝试解析为HTTP请求
		request, err := http.ReadRequest(bufio.NewReader(reader))
		if err == nil {
			// 成功解析为HTTP请求
			info.HTTPMethod = request.Method
			info.RequestURI = request.RequestURI
			info.HTTPVersion = request.Proto
			info.Host = request.Host
			info.Headers = make(map[string]string)

			// 提取所有头部字段
			for key, values := range request.Header {
				if len(values) > 0 {
					info.Headers[key] = values[0]
				}
			}

			// 提取常用头部字段
			info.UserAgent = request.Header.Get("User-Agent")
			info.ContentType = request.Header.Get("Content-Type")
			info.Authorization = request.Header.Get("Authorization")
			info.Referer = request.Header.Get("Referer")
			info.Cookie = request.Header.Get("Cookie")
			info.Accept = request.Header.Get("Accept")
			info.AcceptLanguage = request.Header.Get("Accept-Language")
			info.AcceptEncoding = request.Header.Get("Accept-Encoding")
			info.Connection = request.Header.Get("Connection")

			// 解析Content-Length
			if contentLength := request.Header.Get("Content-Length"); contentLength != "" {
				if length, err := strconv.Atoi(contentLength); err == nil {
					info.ContentLength = length
				}
			}

			// 读取请求体（如果有）
			if request.Body != nil {
				body, _ := io.ReadAll(request.Body)
				info.Body = body
			}
			// 解析URL组件
			info.Domain = extractDomainFromHost(request.Host)
			info.Path = extractPathFromURL(request.RequestURI)
			info.Query = extractQueryFromURL(request.RequestURI)

			// 智能检测协议类型
			scheme := "http"
			// 检查TLS
			if request.TLS != nil {
				scheme = "https"
			}
			// 检查端口号
			if strings.Contains(request.Host, ":") {
				parts := strings.Split(request.Host, ":")
				if len(parts) > 1 {
					port := parts[1]
					if port == "443" {
						scheme = "https"
					}
				}
			}
			// 检查Referer头
			if referer := request.Header.Get("Referer"); referer != "" {
				if strings.HasPrefix(referer, "https://") {
					scheme = "https"
				}
			}

			// 只有当我们有足够信息时才生成完整URL
			if info.Domain != "" && (info.Path != "" || request.RequestURI != "") {
				info.FullURL = generateFullURL(scheme, info.Domain, info.Path, info.Query)
			}

			return info
		}

		// 如果不是HTTP请求，尝试解析为HTTP响应
		// 重置reader
		reader = bytes.NewReader(payload)
		response, err := http.ReadResponse(bufio.NewReader(reader), nil)
		if err == nil {
			// 成功解析为HTTP响应
			info.HTTPVersion = response.Proto
			info.HTTPStatus = response.Status
			info.StatusCode = response.StatusCode
			info.Headers = make(map[string]string)

			// 提取所有头部字段
			for key, values := range response.Header {
				if len(values) > 0 {
					info.Headers[key] = values[0]
				}
			}

			// 提取常用头部字段
			info.Server = response.Header.Get("Server")
			info.ContentType = response.Header.Get("Content-Type")
			info.SetCookie = response.Header.Get("Set-Cookie")
			info.Connection = response.Header.Get("Connection")

			// 解析Content-Length
			if contentLength := response.Header.Get("Content-Length"); contentLength != "" {
				if length, err := strconv.Atoi(contentLength); err == nil {
					info.ContentLength = length
				}
			}

			// 读取响应体（如果有）
			if response.Body != nil {
				body, _ := io.ReadAll(response.Body)
				info.Body = body
			}

			return info
		}
	}

	return info
}

// 从Host中提取域名
func extractDomainFromHost(host string) string {
	if host == "" {
		return ""
	}
	// 移除端口号
	if idx := strings.Index(host, ":"); idx != -1 {
		return host[:idx]
	}
	return host
}

// 从URL中提取路径
func extractPathFromURL(url string) string {
	if url == "" {
		return ""
	}
	// 移除查询参数
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	// 移除片段标识符
	if idx := strings.Index(url, "#"); idx != -1 {
		url = url[:idx]
	}
	return url
}

// 从URL中提取查询参数
func extractQueryFromURL(url string) string {
	if url == "" {
		return ""
	}
	if idx := strings.Index(url, "?"); idx != -1 {
		if hashIdx := strings.Index(url, "#"); hashIdx != -1 {
			return url[idx+1 : hashIdx]
		}
		return url[idx+1:]
	}
	return ""
}

// 生成完整URL
func generateFullURL(scheme, domain, path, query string) string {
	if domain == "" {
		return ""
	}

	// 移除域名末尾的点（DNS格式常见）
	domain = strings.TrimSuffix(domain, ".")

	url := scheme + "://" + domain

	// 确保路径存在且格式正确
	if path == "" {
		path = "/"
		url += path
	} else {
		// 确保路径以/开头
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		url += path
	}

	// 添加查询参数
	if query != "" {
		url += "?" + query
	}

	return url
}

// PrintApplicationLayerDetails 打印应用层详细信息
func PrintApplicationLayerDetails(appInfo *ApplicationLayerInfo) {
	fmt.Println("  Application Layer 详细信息:")

	if appInfo.FullURL != "" {
		fmt.Printf("    完整URL: %s\n", appInfo.FullURL)
	} else if appInfo.Domain != "" {
		fmt.Printf("    域名: %s\n", appInfo.Domain)
	}

	if appInfo.HTTPMethod != "" {
		fmt.Printf("    HTTP方法: %s\n", appInfo.HTTPMethod)
	}

	if appInfo.RequestURI != "" {
		fmt.Printf("    请求URI: %s\n", appInfo.RequestURI)
	}

	if appInfo.HTTPVersion != "" {
		fmt.Printf("    HTTP版本: %s\n", appInfo.HTTPVersion)
	}

	if appInfo.HTTPStatus != "" {
		fmt.Printf("    HTTP状态: %s\n", appInfo.HTTPStatus)
	}

	if appInfo.StatusCode != 0 {
		fmt.Printf("    状态码: %d\n", appInfo.StatusCode)
	}

	if appInfo.Host != "" {
		fmt.Printf("    Host: %s\n", appInfo.Host)
	}

	if appInfo.UserAgent != "" {
		fmt.Printf("    User-Agent: %s\n", appInfo.UserAgent)
	}

	if appInfo.ContentType != "" {
		fmt.Printf("    Content-Type: %s\n", appInfo.ContentType)
	}

	if appInfo.ContentLength > 0 {
		fmt.Printf("    Content-Length: %d\n", appInfo.ContentLength)
	}

	if len(appInfo.Headers) > 0 {
		fmt.Println("    HTTP头部:")
		for key, value := range appInfo.Headers {
			fmt.Printf("      %s: %s\n", key, value)
		}
	}

	if len(appInfo.Body) > 0 {
		fmt.Printf("    消息体长度: %d 字节\n", len(appInfo.Body))
	}

	// 显示base64编码的载荷数据
	if appInfo.Payload != "" {
		// 解码base64数据以获取原始长度
		if decoded, err := base64.StdEncoding.DecodeString(appInfo.Payload); err == nil {
			fmt.Printf("    原始载荷: %d 字节的base64编码数据\n", len(decoded))

			// 如果数据较短，显示解码后的内容
			if len(decoded) <= 256 && isLikelyText(decoded) {
				fmt.Printf("    解码内容: %s\n", string(decoded))
			} else if len(decoded) > 0 {
				// 显示前几个字节的十六进制表示
				hexData := fmt.Sprintf("%x", decoded[:min(16, len(decoded))])
				fmt.Printf("    前%d字节十六进制: %s\n", min(16, len(decoded)), hexData)
			}
		} else {
			fmt.Printf("    原始载荷: %s (base64编码)\n", appInfo.Payload)
		}
	}
}

// isLikelyText 检查数据是否可能是文本
func isLikelyText(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// 统计可打印字符的比例
	printableCount := 0
	for _, b := range data {
		// 检查是否为可打印ASCII字符或常见空白字符
		if (b >= 32 && b <= 126) || b == 9 || b == 10 || b == 13 {
			printableCount++
		}
	}

	// 如果可打印字符比例超过70%，则认为是文本
	return float64(printableCount)/float64(len(data)) > 0.7
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
