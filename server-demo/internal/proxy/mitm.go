package proxy

import (
	"bytes"
	"crypto/tls"
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

type ProxyServer struct {
	addr  string
	https bool
	srv   *http.Server
	mu    sync.Mutex
	store storage.FlowStorage
}

func NewProxyServer(addr string, https bool, store storage.FlowStorage) *ProxyServer {
	return &ProxyServer{addr: addr, https: https, store: store}
}

func (p *ProxyServer) Start() error {
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

		// 读取请求体
		var reqBody []byte
		if req.Body != nil {
			reqBody, _ = io.ReadAll(req.Body)
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		// 组装请求部分
		reqURL := req.URL.String()
		path := req.URL.Path
		query := req.URL.RawQuery
		host := req.Host

		flow := &models.Flow{
			ID:         flowID,
			Scheme:     req.URL.Scheme,
			RemoteAddr: req.RemoteAddr,
			StartAt:    start,
			Request: &models.HTTPRequest{
				Method:  req.Method,
				URL:     reqURL,
				Path:    path,
				Query:   query,
				Host:    host,
				Headers: models.CopyHeaders(req.Header),
				Body:    reqBody,
				Proto:   req.Proto,
				Length:  int(req.ContentLength),
			},
		}
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
		// 读取响应体
		var body []byte
		if resp != nil && resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(body))
		}
		end := time.Now()
		flow.EndAt = end
		flow.LatencyMs = end.Sub(flow.StartAt).Milliseconds()
		if resp != nil {
			flow.Response = &models.HTTPResponse{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				Headers:    models.CopyHeaders(resp.Header),
				Body:       body,
				Proto:      resp.Proto,
				Length:     int(resp.ContentLength),
			}
		}
		return resp
	})

	// HTTPS MITM（自签证书）
	if p.https {
		// 使用自定义根CA进行 MITM
		certPEM, keyPEM, _, err := EnsureCAExists()
		if err != nil {
			return err
		}
		caCert, caKey, err := ParseCA(certPEM, keyPEM)
		if err != nil {
			return err
		}
		// 为每个 host 动态签发含 SAN 的证书
		tlsConfigForHost := func(host string, _ *goproxy.ProxyCtx) (*tls.Config, error) {
			// 剥离端口，获取纯主机名
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
			return &tls.Config{Certificates: []tls.Certificate{pair}, MinVersion: tls.VersionTLS12, ServerName: h}, nil
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

func (p *ProxyServer) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.srv == nil {
		return nil
	}
	err := p.srv.Close()
	p.srv = nil
	return err
}

func (p *ProxyServer) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.srv != nil
}
