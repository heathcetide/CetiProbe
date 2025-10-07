package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"probe/internal/storage"
)

type Server struct {
	storage storage.Storage
	port    string
	server  *http.Server
}

func NewServer(storage storage.Storage, port string) *Server {
	return &Server{
		storage: storage,
		port:    port,
	}
}

func (s *Server) Start() error {
	router := s.setupRoutes()

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: router,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *Server) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// 静态文件服务
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// API路由
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/packets", s.handleGetPackets).Methods("GET")
	api.HandleFunc("/packets/filter", s.handleFilterPackets).Methods("POST")
	api.HandleFunc("/stats", s.handleGetStats).Methods("GET")
	api.HandleFunc("/clear", s.handleClear).Methods("POST")
	api.HandleFunc("/export", s.handleExport).Methods("GET")

	// WebSocket路由
	router.HandleFunc("/ws", s.handleWebSocket)

	// 主页面
	router.HandleFunc("/", s.handleIndex)

	return router
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

func (s *Server) handleGetPackets(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100 // 默认限制

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	packets := s.storage.GetPackets(limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"packets": packets,
		"count":   len(packets),
	})
}

func (s *Server) handleFilterPackets(w http.ResponseWriter, r *http.Request) {
	var filter storage.Filter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		http.Error(w, "无效的过滤条件", http.StatusBadRequest)
		return
	}

	packets := s.storage.GetPacketsByFilter(filter)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"packets": packets,
		"count":   len(packets),
	})
}

func (s *Server) handleGetStats(w http.ResponseWriter, r *http.Request) {
	stats := s.storage.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleClear(w http.ResponseWriter, r *http.Request) {
	s.storage.Clear()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "数据已清空",
	})
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	packets := s.storage.GetPackets(0) // 获取所有数据

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=packets.json")
		json.NewEncoder(w).Encode(packets)

	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=packets.csv")
		s.exportCSV(w, packets)

	default:
		http.Error(w, "不支持的导出格式", http.StatusBadRequest)
	}
}

func (s *Server) exportCSV(w http.ResponseWriter, packets []*storage.PacketInfo) {
	// CSV头部
	fmt.Fprintln(w, "时间戳,源IP,目标IP,源端口,目标端口,协议,长度,HTTP方法,HTTP URL,HTTP状态,User Agent,Content Type")

	// CSV数据
	for _, packet := range packets {
		fmt.Fprintf(w, "%s,%s,%s,%d,%d,%s,%d,%s,%s,%s,%s,%s\n",
			packet.Timestamp.Format(time.RFC3339),
			packet.SrcIP,
			packet.DstIP,
			packet.SrcPort,
			packet.DstPort,
			packet.Protocol,
			packet.Length,
			packet.HTTPMethod,
			packet.HTTPURL,
			packet.HTTPStatus,
			packet.UserAgent,
			packet.ContentType,
		)
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 这里可以实现WebSocket实时推送功能
	// 为了简化，暂时返回一个简单的响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "WebSocket功能待实现",
	})
}
