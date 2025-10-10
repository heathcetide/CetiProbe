package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"probe/internal/capture"
	pxy "probe/internal/proxy"
	"probe/pkg/storage"
	"probe/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/gopacket/pcap"
)

var (
	st        storage.Storage = storage.NewMemoryStorage()
	capMu     sync.Mutex
	capInst   *capture.Capturer
	currIF    string
	flowStore = storage.NewMemoryFlowStore()
	proxyMu   sync.Mutex
	proxyInst *pxy.EnhancedProxyServer
)

// generateInstallScript 生成安装脚本和指引
func generateInstallScript(osType, certPath string) (script string, instructions map[string]interface{}) {
	instructions = pxy.GetInstallInstructions(osType)

	switch osType {
	case "darwin":
		script = fmt.Sprintf("sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", certPath)
	case "windows":
		script = fmt.Sprintf("certutil -addstore -f ROOT %s", certPath)
	case "linux":
		script = fmt.Sprintf("sudo cp %s /usr/local/share/ca-certificates/probe-ca.crt && sudo update-ca-certificates", certPath)
	}

	return script, instructions
}

// getCertInstallInstructions 获取证书安装指引
func getCertInstallInstructions(osType string) map[string]interface{} {
	return pxy.GetInstallInstructions(osType)
}

func main() {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	// 静态前端页面
	r.Static("/ui", "../web")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/ui/")
	})

	api := r.Group("/api")
	{
		api.GET("/status", func(c *gin.Context) {
			capMu.Lock()
			running := false
			iface := currIF
			if capInst != nil {
				running = capInst.IsRunning()
			}
			capMu.Unlock()
			c.JSON(200, gin.H{"running": running, "iface": iface})
		})

		api.GET("/interfaces", func(c *gin.Context) {
			devs, err := pcap.FindAllDevs()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			type Iface struct {
				Name        string   `json:"name"`
				Description string   `json:"description"`
				Addresses   []string `json:"addresses"`
			}
			list := make([]Iface, 0, len(devs))
			for _, d := range devs {
				addrs := make([]string, 0, len(d.Addresses))
				for _, a := range d.Addresses {
					addrs = append(addrs, a.IP.String())
				}
				list = append(list, Iface{Name: d.Name, Description: d.Description, Addresses: addrs})
			}
			c.JSON(200, list)
		})

		api.POST("/start", func(c *gin.Context) {
			iface := c.Query("iface")
			if iface == "" {
				c.JSON(400, gin.H{"error": "缺少 iface 参数"})
				return
			}

			capMu.Lock()
			if capInst != nil && capInst.IsRunning() {
				capMu.Unlock()
				c.JSON(409, gin.H{"error": "已在运行"})
				return
			}
			// 创建新的抓包器
			cp, err := capture.NewCapturer(iface, st)
			if err != nil {
				capMu.Unlock()
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			capInst = cp
			currIF = iface
			capMu.Unlock()

			go func() {
				_ = cp.Start()
			}()

			c.JSON(200, gin.H{"ok": true})
		})

		api.POST("/stop", func(c *gin.Context) {
			capMu.Lock()
			if capInst == nil || !capInst.IsRunning() {
				capMu.Unlock()
				c.JSON(400, gin.H{"error": "未在运行"})
				return
			}
			err := capInst.Stop()
			capInst = nil
			capMu.Unlock()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"ok": true})
		})

		api.GET("/packets", func(c *gin.Context) {
			limitStr := c.Query("limit")
			limit := 100
			if limitStr != "" {
				if v, err := strconv.Atoi(limitStr); err == nil {
					limit = v
				}
			}
			packets := st.GetPackets(limit)
			c.JSON(200, packets)
		})

		api.GET("/stats", func(c *gin.Context) {
			c.JSON(200, st.GetStats())
		})

		api.DELETE("/packets", func(c *gin.Context) {
			st.Clear()
			c.JSON(200, gin.H{"ok": true})
		})

		// 代理控制 & flows
		api.GET("/proxy/status", func(c *gin.Context) {
			proxyMu.Lock()
			running := proxyInst != nil && proxyInst.IsRunning()
			proxyMu.Unlock()
			c.JSON(200, gin.H{"running": running})
		})

		api.POST("/proxy/start", func(c *gin.Context) {
			addr := c.Query("addr")
			if addr == "" {
				addr = ":8899"
			}
			https := c.Query("https") == "1"
			proxyMu.Lock()
			if proxyInst != nil && proxyInst.IsRunning() {
				proxyMu.Unlock()
				c.JSON(409, gin.H{"error": "代理已在运行"})
				return
			}
			ps := pxy.NewEnhancedProxyServer(addr, https, flowStore)
			proxyInst = ps
			proxyMu.Unlock()
			go ps.Start()
			c.JSON(200, gin.H{"ok": true, "addr": addr})
		})

		api.POST("/proxy/stop", func(c *gin.Context) {
			proxyMu.Lock()
			if proxyInst == nil || !proxyInst.IsRunning() {
				proxyMu.Unlock()
				c.JSON(400, gin.H{"error": "代理未在运行"})
				return
			}
			_ = proxyInst.Stop()
			proxyInst = nil
			proxyMu.Unlock()
			c.JSON(200, gin.H{"ok": true})
		})

		// CA 证书：下载与（重新）生成
		api.GET("/proxy/ca", func(c *gin.Context) {
			_, _, files, err := pxy.LoadCAFromDisk()
			if err != nil {
				// 若不存在则尝试生成
				_, _, files, err = pxy.GenerateCA()
				if err != nil {
					c.String(500, err.Error())
					return
				}
			}
			c.Header("Content-Disposition", "attachment; filename=proxy_root_ca.pem")
			c.Header("Content-Type", "application/x-pem-file")
			c.File(files.CertPath)
		})
		api.POST("/proxy/ca/generate", func(c *gin.Context) {
			_, _, _, err := pxy.GenerateCA()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"ok": true})
		})

		// 新增：自动安装证书端点
		api.POST("/proxy/ca/install", func(c *gin.Context) {
			os := c.Query("os") // darwin/windows/linux
			_, _, files, err := pxy.LoadCAFromDisk()
			if err != nil {
				_, _, files, err = pxy.GenerateCA()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
			}

			// 根据操作系统返回安装脚本或执行命令
			script, instructions := generateInstallScript(os, files.CertPath)
			c.JSON(200, gin.H{
				"script":       script,
				"instructions": instructions,
				"cert_path":    files.CertPath,
			})
		})

		// 新增：获取证书安装指引
		api.GET("/proxy/ca/instructions", func(c *gin.Context) {
			os := c.Query("os")
			instructions := getCertInstallInstructions(os)
			c.JSON(200, gin.H{"instructions": instructions})
		})

		api.GET("/flows", func(c *gin.Context) {
			limitStr := c.Query("limit")
			limit := 200
			if limitStr != "" {
				if v, err := strconv.Atoi(limitStr); err == nil {
					limit = v
				}
			}
			c.JSON(200, flowStore.GetAll(limit))
		})

		api.GET("/flows/stats", func(c *gin.Context) {
			c.JSON(200, flowStore.Stats())
		})

		api.GET("/flows/:id", func(c *gin.Context) {
			id := c.Param("id")
			f := flowStore.GetByID(id)
			if f == nil {
				c.JSON(404, gin.H{"error": "not found"})
				return
			}
			if c.Query("decoded") == "1" && f.Response != nil {
				// 动态构造一个带解码文本的响应
				type FlowView struct {
					ID         string      `json:"id"`
					Scheme     string      `json:"scheme"`
					RemoteAddr string      `json:"remote_addr"`
					StartAt    interface{} `json:"start_at"`
					EndAt      interface{} `json:"end_at"`
					LatencyMs  int64       `json:"latency_ms"`
					Request    interface{} `json:"request"`
					Response   interface{} `json:"response"`
				}
				rv := f.Response
				// 解码
				bodyText, _ := utils.DecodeBodyToText(rv.Body, rv.Headers)
				resp := gin.H{
					"status":      rv.Status,
					"status_code": rv.StatusCode,
					"headers":     rv.Headers,
					"proto":       rv.Proto,
					"length":      rv.Length,
					"body_text":   bodyText,
				}
				c.JSON(200, FlowView{ID: f.ID, Scheme: f.Scheme, RemoteAddr: f.RemoteAddr, StartAt: f.StartAt, EndAt: f.EndAt, LatencyMs: f.LatencyMs, Request: f.Request, Response: resp})
				return
			}
			c.JSON(200, f)
		})

		api.DELETE("/flows", func(c *gin.Context) {
			flowStore.Clear()
			c.JSON(200, gin.H{"ok": true})
		})

		// 新增统计信息API
		api.GET("/proxy/stats", func(c *gin.Context) {
			proxyMu.Lock()
			defer proxyMu.Unlock()
			if proxyInst == nil {
				c.JSON(200, gin.H{"network_stats": map[string]interface{}{}, "performance_stats": map[string]interface{}{}})
				return
			}

			networkStats := proxyInst.GetNetworkStats()
			performanceStats := proxyInst.GetPerformanceStats()

			c.JSON(200, gin.H{
				"network_stats":     networkStats,
				"performance_stats": performanceStats,
			})
		})
	}

	_ = r.Run(":8080")
}
